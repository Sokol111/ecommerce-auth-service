package http

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-auth-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/query"
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	"github.com/Sokol111/ecommerce-commons/pkg/security/token"
)

type authHandler struct {
	loginHandler        command.LoginHandler
	logoutHandler       command.LogoutHandler
	createUserHandler   command.CreateAdminUserHandler
	disableUserHandler  command.DisableAdminUserHandler
	enableUserHandler   command.EnableAdminUserHandler
	refreshTokenHandler command.RefreshTokenHandler
	getUserByIDHandler  query.GetAdminUserByIDHandler
	listUsersHandler    query.ListAdminUsersHandler
	getRolesHandler     query.GetRolesHandler
	permissionProvider  adminuser.RolePermissionProvider
	loginRateLimiter    *LoginRateLimiter
}

func newAuthHandler(
	loginHandler command.LoginHandler,
	logoutHandler command.LogoutHandler,
	createUserHandler command.CreateAdminUserHandler,
	disableUserHandler command.DisableAdminUserHandler,
	enableUserHandler command.EnableAdminUserHandler,
	refreshTokenHandler command.RefreshTokenHandler,
	getUserByIDHandler query.GetAdminUserByIDHandler,
	listUsersHandler query.ListAdminUsersHandler,
	getRolesHandler query.GetRolesHandler,
	permissionProvider adminuser.RolePermissionProvider,
	loginRateLimiter *LoginRateLimiter,
) httpapi.Handler {
	return &authHandler{
		loginHandler:        loginHandler,
		logoutHandler:       logoutHandler,
		createUserHandler:   createUserHandler,
		disableUserHandler:  disableUserHandler,
		enableUserHandler:   enableUserHandler,
		refreshTokenHandler: refreshTokenHandler,
		getUserByIDHandler:  getUserByIDHandler,
		listUsersHandler:    listUsersHandler,
		getRolesHandler:     getRolesHandler,
		permissionProvider:  permissionProvider,
		loginRateLimiter:    loginRateLimiter,
	}
}

var aboutBlankURL, _ = url.Parse("about:blank")

// Helper functions
func toOptDateTime(t *time.Time) httpapi.OptDateTime {
	if t == nil {
		return httpapi.OptDateTime{}
	}
	return httpapi.NewOptDateTime(*t)
}

func toPermissions(perms []adminuser.Permission) []httpapi.Permission {
	return lo.Map(perms, func(p adminuser.Permission, _ int) httpapi.Permission {
		return httpapi.Permission(p)
	})
}

func toAdminUserProfile(user *adminuser.AdminUser, permissions []adminuser.Permission) *httpapi.AdminUserProfile {
	return &httpapi.AdminUserProfile{
		ID:          uuid.MustParse(user.ID),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        httpapi.AdminRole(user.Role),
		Permissions: toPermissions(permissions),
		Enabled:     user.Enabled,
		CreatedAt:   user.CreatedAt,
		LastLoginAt: toOptDateTime(user.LastLoginAt),
	}
}

func toAdminUserResponse(user *adminuser.AdminUser) *httpapi.AdminUserResponse {
	return &httpapi.AdminUserResponse{
		ID:          uuid.MustParse(user.ID),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        httpapi.AdminRole(user.Role),
		Enabled:     user.Enabled,
		CreatedAt:   user.CreatedAt,
		ModifiedAt:  httpapi.NewOptDateTime(user.ModifiedAt),
		LastLoginAt: toOptDateTime(user.LastLoginAt),
	}
}

// AdminLogin implements adminLogin operation.
func (h *authHandler) AdminLogin(ctx context.Context, req *httpapi.LoginRequest) (httpapi.AdminLoginRes, error) {
	// Rate limit by email
	if !h.loginRateLimiter.Allow(ctx, req.Email) {
		return &httpapi.AdminLoginTooManyRequests{
			Status: 429,
			Type:   *aboutBlankURL,
			Title:  "Too many login attempts",
			Detail: httpapi.NewOptString("Please try again later"),
		}, nil
	}

	result, err := h.loginHandler.Handle(ctx, command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, adminuser.ErrInvalidCredentials) {
		return &httpapi.AdminLoginUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Invalid credentials",
		}, nil
	}
	if errors.Is(err, adminuser.ErrAdminUserDisabled) {
		return &httpapi.AdminLoginForbidden{
			Status: 403,
			Type:   *aboutBlankURL,
			Title:  "Account disabled",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &httpapi.AdminAuthResponse{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		ExpiresIn:        result.ExpiresIn,
		RefreshExpiresIn: result.RefreshExpiresIn,
		TokenType:        "Bearer",
		User:             *toAdminUserProfile(result.User, h.permissionProvider.GetPermissionsForRole(result.User.Role)),
	}, nil
}

// AdminLogout implements adminLogout operation.
func (h *authHandler) AdminLogout(ctx context.Context) (httpapi.AdminLogoutRes, error) {
	claims, err := h.getCurrentUserClaims(ctx)
	if err != nil {
		return &httpapi.AdminLogoutUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Unauthorized",
		}, nil
	}

	if err := h.logoutHandler.Handle(ctx, command.LogoutCommand{UserID: claims.UserID}); err != nil {
		return nil, err
	}

	return &httpapi.AdminLogoutNoContent{}, nil
}

// AdminGetProfile implements adminGetProfile operation.
func (h *authHandler) AdminGetProfile(ctx context.Context) (httpapi.AdminGetProfileRes, error) {
	claims, err := h.getCurrentUserClaims(ctx)
	if err != nil {
		return &httpapi.AdminGetProfileUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Unauthorized",
		}, nil
	}

	user, err := h.getUserByIDHandler.Handle(ctx, query.GetAdminUserByIDQuery{ID: claims.UserID})
	if err != nil {
		return nil, err
	}

	return toAdminUserProfile(user, h.permissionProvider.GetPermissionsForRole(user.Role)), nil
}

// AdminUserCreate implements adminUserCreate operation.
func (h *authHandler) AdminUserCreate(ctx context.Context, req *httpapi.AdminUserCreateRequest) (httpapi.AdminUserCreateRes, error) {
	user, err := h.createUserHandler.Handle(ctx, command.CreateAdminUserCommand{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      adminuser.Role(req.Role),
	})

	if errors.Is(err, adminuser.ErrEmailAlreadyExists) {
		return &httpapi.AdminUserCreateConflict{
			Status: 409,
			Type:   *aboutBlankURL,
			Title:  "Email already exists",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return toAdminUserResponse(user), nil
}

// AdminUserGetById implements adminUserGetById operation.
func (h *authHandler) AdminUserGetById(ctx context.Context, params httpapi.AdminUserGetByIdParams) (httpapi.AdminUserGetByIdRes, error) {
	user, err := h.getUserByIDHandler.Handle(ctx, query.GetAdminUserByIDQuery{ID: params.ID.String()})
	if errors.Is(err, persistence.ErrEntityNotFound) {
		return &httpapi.AdminUserGetByIdNotFound{
			Status: 404,
			Type:   *aboutBlankURL,
			Title:  "User not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return toAdminUserResponse(user), nil
}

// AdminUserList implements adminUserList operation.
func (h *authHandler) AdminUserList(ctx context.Context, params httpapi.AdminUserListParams) (httpapi.AdminUserListRes, error) {
	q := query.ListAdminUsersQuery{
		Page: params.Page,
		Size: params.Size,
	}

	if params.Role.IsSet() {
		role := adminuser.Role(params.Role.Value)
		q.Role = &role
	}
	if params.Enabled.IsSet() {
		q.Enabled = &params.Enabled.Value
	}
	if params.Search.IsSet() {
		q.Search = &params.Search.Value
	}

	result, err := h.listUsersHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	items := lo.Map(result.Items, func(user *adminuser.AdminUser, _ int) httpapi.AdminUserResponse {
		return *toAdminUserResponse(user)
	})

	return &httpapi.AdminUserListResponse{
		Items: items,
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}, nil
}

// AdminUserDisable implements adminUserDisable operation.
func (h *authHandler) AdminUserDisable(ctx context.Context, params httpapi.AdminUserDisableParams) (httpapi.AdminUserDisableRes, error) {
	claims, err := h.getCurrentUserClaims(ctx)
	if err != nil {
		return &httpapi.AdminUserDisableUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Unauthorized",
		}, nil
	}

	err = h.disableUserHandler.Handle(ctx, command.DisableAdminUserCommand{
		ID:            params.ID.String(),
		RequestUserID: claims.UserID,
	})

	if errors.Is(err, persistence.ErrEntityNotFound) {
		return &httpapi.AdminUserDisableNotFound{
			Status: 404,
			Type:   *aboutBlankURL,
			Title:  "User not found",
		}, nil
	}
	if errors.Is(err, adminuser.ErrCannotDisableSelf) || errors.Is(err, adminuser.ErrCannotDisableSuperAdmin) {
		return &httpapi.AdminUserDisableForbidden{
			Status: 403,
			Type:   *aboutBlankURL,
			Title:  err.Error(),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &httpapi.AdminUserDisableNoContent{}, nil
}

// AdminUserEnable implements adminUserEnable operation.
func (h *authHandler) AdminUserEnable(ctx context.Context, params httpapi.AdminUserEnableParams) (httpapi.AdminUserEnableRes, error) {
	err := h.enableUserHandler.Handle(ctx, command.EnableAdminUserCommand{
		ID: params.ID.String(),
	})

	if errors.Is(err, persistence.ErrEntityNotFound) {
		return &httpapi.AdminUserEnableNotFound{
			Status: 404,
			Type:   *aboutBlankURL,
			Title:  "User not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &httpapi.AdminUserEnableNoContent{}, nil
}

// TokenRefresh implements tokenRefresh operation.
func (h *authHandler) TokenRefresh(ctx context.Context, req *httpapi.TokenRefreshRequest) (httpapi.TokenRefreshRes, error) {
	result, err := h.refreshTokenHandler.Handle(ctx, command.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})

	if errors.Is(err, command.ErrRefreshTokenReused) {
		return &httpapi.TokenRefreshUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Token reused",
			Detail: httpapi.NewOptString("This refresh token has already been used. All sessions have been invalidated for security. Please log in again."),
		}, nil
	}
	if err != nil {
		return &httpapi.TokenRefreshUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Invalid or expired refresh token",
		}, nil
	}

	return &httpapi.TokenRefreshResponse{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		ExpiresIn:        result.ExpiresIn,
		RefreshExpiresIn: result.RefreshExpiresIn,
	}, nil
}

// GetRoles implements getRoles operation.
func (h *authHandler) GetRoles(ctx context.Context) (httpapi.GetRolesRes, error) {
	roles := h.getRolesHandler.Handle(ctx)

	items := lo.Map(roles, func(r query.RoleInfo, _ int) httpapi.RoleInfo {
		return httpapi.RoleInfo{
			Name:        httpapi.AdminRole(r.Name),
			Description: r.Description,
			Permissions: toPermissions(r.Permissions),
		}
	})

	return &httpapi.RolesResponse{
		Roles: items,
	}, nil
}

// getCurrentUserClaims extracts current user claims from context
func (h *authHandler) getCurrentUserClaims(ctx context.Context) (*token.Claims, error) {
	// Claims are validated by security handler and stored in context
	claims := token.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, errors.New("no claims in context")
	}
	return claims, nil
}
