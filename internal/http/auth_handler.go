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
	loginHandler          command.LoginHandler
	changePasswordHandler command.ChangePasswordHandler
	createUserHandler     command.CreateAdminUserHandler
	updateUserHandler     command.UpdateAdminUserHandler
	disableUserHandler    command.DisableAdminUserHandler
	enableUserHandler     command.EnableAdminUserHandler
	refreshTokenHandler   command.RefreshTokenHandler
	getUserByIDHandler    query.GetAdminUserByIDHandler
	listUsersHandler      query.ListAdminUsersHandler
	getRolesHandler       query.GetRolesHandler
	getPermissionsHandler query.GetPermissionsHandler
}

func newAuthHandler(
	loginHandler command.LoginHandler,
	changePasswordHandler command.ChangePasswordHandler,
	createUserHandler command.CreateAdminUserHandler,
	updateUserHandler command.UpdateAdminUserHandler,
	disableUserHandler command.DisableAdminUserHandler,
	enableUserHandler command.EnableAdminUserHandler,
	refreshTokenHandler command.RefreshTokenHandler,
	getUserByIDHandler query.GetAdminUserByIDHandler,
	listUsersHandler query.ListAdminUsersHandler,
	getRolesHandler query.GetRolesHandler,
	getPermissionsHandler query.GetPermissionsHandler,
) httpapi.Handler {
	return &authHandler{
		loginHandler:          loginHandler,
		changePasswordHandler: changePasswordHandler,
		createUserHandler:     createUserHandler,
		updateUserHandler:     updateUserHandler,
		disableUserHandler:    disableUserHandler,
		enableUserHandler:     enableUserHandler,
		refreshTokenHandler:   refreshTokenHandler,
		getUserByIDHandler:    getUserByIDHandler,
		listUsersHandler:      listUsersHandler,
		getRolesHandler:       getRolesHandler,
		getPermissionsHandler: getPermissionsHandler,
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

func toAdminRole(role adminuser.Role) httpapi.AdminRole {
	switch role {
	case adminuser.RoleSuperAdmin:
		return httpapi.AdminRoleSuperAdmin
	case adminuser.RoleCatalogManager:
		return httpapi.AdminRoleCatalogManager
	case adminuser.RoleViewer:
		return httpapi.AdminRoleViewer
	default:
		return httpapi.AdminRoleViewer
	}
}

func toDomainRole(role httpapi.AdminRole) adminuser.Role {
	switch role {
	case httpapi.AdminRoleSuperAdmin:
		return adminuser.RoleSuperAdmin
	case httpapi.AdminRoleCatalogManager:
		return adminuser.RoleCatalogManager
	case httpapi.AdminRoleViewer:
		return adminuser.RoleViewer
	default:
		return adminuser.RoleViewer
	}
}

func toPermissions(perms []adminuser.Permission) []httpapi.Permission {
	return lo.Map(perms, func(p adminuser.Permission, _ int) httpapi.Permission {
		return httpapi.Permission(p)
	})
}

func toAdminUserProfile(user *adminuser.AdminUser) *httpapi.AdminUserProfile {
	return &httpapi.AdminUserProfile{
		ID:          uuid.MustParse(user.ID),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Role:        toAdminRole(user.Role),
		Permissions: toPermissions(user.GetPermissions()),
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
		Role:        toAdminRole(user.Role),
		Enabled:     user.Enabled,
		CreatedAt:   user.CreatedAt,
		ModifiedAt:  httpapi.NewOptDateTime(user.ModifiedAt),
		LastLoginAt: toOptDateTime(user.LastLoginAt),
	}
}

// AdminLogin implements adminLogin operation.
func (h *authHandler) AdminLogin(ctx context.Context, req *httpapi.LoginRequest) (httpapi.AdminLoginRes, error) {
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
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
		User:         *toAdminUserProfile(result.User),
	}, nil
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

	return toAdminUserProfile(user), nil
}

// AdminChangePassword implements adminChangePassword operation.
func (h *authHandler) AdminChangePassword(ctx context.Context, req *httpapi.ChangePasswordRequest) (httpapi.AdminChangePasswordRes, error) {
	claims, err := h.getCurrentUserClaims(ctx)
	if err != nil {
		return &httpapi.AdminChangePasswordUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Unauthorized",
		}, nil
	}

	err = h.changePasswordHandler.Handle(ctx, command.ChangePasswordCommand{
		UserID:          claims.UserID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})

	if errors.Is(err, adminuser.ErrInvalidCredentials) {
		return &httpapi.AdminChangePasswordUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Wrong current password",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &httpapi.AdminChangePasswordNoContent{}, nil
}

// AdminUserCreate implements adminUserCreate operation.
func (h *authHandler) AdminUserCreate(ctx context.Context, req *httpapi.AdminUserCreateRequest) (httpapi.AdminUserCreateRes, error) {
	user, err := h.createUserHandler.Handle(ctx, command.CreateAdminUserCommand{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      toDomainRole(req.Role),
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

// AdminUserUpdate implements adminUserUpdate operation.
func (h *authHandler) AdminUserUpdate(ctx context.Context, req *httpapi.AdminUserUpdateRequest) (httpapi.AdminUserUpdateRes, error) {
	cmd := command.UpdateAdminUserCommand{
		ID: req.ID.String(),
	}

	if req.FirstName.IsSet() {
		cmd.FirstName = &req.FirstName.Value
	}
	if req.LastName.IsSet() {
		cmd.LastName = &req.LastName.Value
	}
	if req.Role.IsSet() {
		role := toDomainRole(req.Role.Value)
		cmd.Role = &role
	}

	user, err := h.updateUserHandler.Handle(ctx, cmd)
	if errors.Is(err, persistence.ErrEntityNotFound) {
		return &httpapi.AdminUserUpdateNotFound{
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
		role := toDomainRole(params.Role.Value)
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

	items := make([]httpapi.AdminUserResponse, len(result.Items))
	for i, user := range result.Items {
		items[i] = *toAdminUserResponse(user)
	}

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

	if err != nil {
		return &httpapi.TokenRefreshUnauthorized{
			Status: 401,
			Type:   *aboutBlankURL,
			Title:  "Invalid or expired refresh token",
		}, nil
	}

	return &httpapi.TokenRefreshResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// GetRoles implements getRoles operation.
func (h *authHandler) GetRoles(ctx context.Context) (httpapi.GetRolesRes, error) {
	roles := h.getRolesHandler.Handle(ctx)

	items := lo.Map(roles, func(r query.RoleInfo, _ int) httpapi.RoleInfo {
		return httpapi.RoleInfo{
			Name:        toAdminRole(r.Name),
			Description: r.Description,
			Permissions: toPermissions(r.Permissions),
		}
	})

	return &httpapi.RolesResponse{
		Roles: items,
	}, nil
}

// GetPermissions implements getPermissions operation.
func (h *authHandler) GetPermissions(ctx context.Context) (httpapi.GetPermissionsRes, error) {
	permissions := h.getPermissionsHandler.Handle(ctx)

	items := lo.Map(permissions, func(p query.PermissionInfo, _ int) httpapi.PermissionInfo {
		return httpapi.PermissionInfo{
			Name:        httpapi.Permission(p.Name),
			Description: p.Description,
			Resource:    p.Resource,
		}
	})

	return &httpapi.PermissionsResponse{
		Permissions: items,
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
