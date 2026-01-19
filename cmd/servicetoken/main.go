package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

func main() {
	// Flags
	privateKeyHex := flag.String("private-key", "", "Hex-encoded Ed25519 private key (64 bytes)")
	serviceName := flag.String("name", "", "Service name (e.g., 'product-query-service')")
	role := flag.String("role", "service_account", "Role for the service account")
	duration := flag.Duration("duration", 8760*time.Hour, "Token validity duration (default: 1 year)")
	permissions := flag.String("permissions", "", "Comma-separated permissions (optional, uses role permissions by default)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Generate a long-lived PASETO token for service-to-service authentication.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -private-key=<hex> -name=product-query-service -duration=8760h\n", os.Args[0])
	}

	flag.Parse()

	if *privateKeyHex == "" {
		fmt.Fprintln(os.Stderr, "Error: -private-key is required")
		flag.Usage()
		os.Exit(1)
	}

	if *serviceName == "" {
		fmt.Fprintln(os.Stderr, "Error: -name is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse private key
	keyBytes, err := hex.DecodeString(*privateKeyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid private key hex: %v\n", err)
		os.Exit(1)
	}

	if len(keyBytes) != 64 {
		fmt.Fprintf(os.Stderr, "Error: private key must be 64 bytes (got %d)\n", len(keyBytes))
		os.Exit(1)
	}

	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(keyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse private key: %v\n", err)
		os.Exit(1)
	}

	// Build token
	now := time.Now()

	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(*duration))
	token.SetSubject(fmt.Sprintf("service:%s", *serviceName))
	token.SetString("role", *role)
	token.SetString("type", "access")
	token.SetString("service_name", *serviceName)

	// Set permissions
	var perms []string
	if *permissions != "" {
		// Use provided permissions
		for _, p := range splitAndTrim(*permissions, ",") {
			if p != "" {
				perms = append(perms, p)
			}
		}
	} else {
		// Default service account permissions
		perms = []string{
			"products:read",
			"categories:read",
			"attributes:read",
		}
	}
	token.Set("permissions", perms)

	// Sign token
	signedToken := token.V4Sign(privateKey, nil)

	// Output
	fmt.Println("=== Service Account Token ===")
	fmt.Println()
	fmt.Printf("Service:     %s\n", *serviceName)
	fmt.Printf("Role:        %s\n", *role)
	fmt.Printf("Permissions: %v\n", perms)
	fmt.Printf("Valid For:   %s\n", *duration)
	fmt.Printf("Expires:     %s\n", now.Add(*duration).Format(time.RFC3339))
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(signedToken)
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
