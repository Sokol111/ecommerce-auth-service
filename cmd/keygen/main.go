package main

import (
	"encoding/hex"
	"fmt"

	"aidanwoods.dev/go-paseto"
)

func main() {
	privateKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := privateKey.Public()

	privateKeyHex := hex.EncodeToString(privateKey.ExportBytes())
	publicKeyHex := hex.EncodeToString(publicKey.ExportBytes())

	fmt.Println("=== PASETO V4 Key Pair ===")
	fmt.Println()
	fmt.Println("Private Key (keep secret, only for auth-service):")
	fmt.Println(privateKeyHex)
	fmt.Println()
	fmt.Println("Public Key (share with all services):")
	fmt.Println(publicKeyHex)
}
