package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// This script generates an Ed25519 key pair for the commercial license system.
// For setup instructions, see docs/INTERNAL_LICENSE_SETUP.md

func main() {
	// Generate Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Encode as base64 for easy storage in env files
	publicKeyB64 := base64.StdEncoding.EncodeToString(publicKey)
	privateKeyB64 := base64.StdEncoding.EncodeToString(privateKey)

	fmt.Println("=== Ed25519 License Key Pair Generated ===")
	fmt.Println()
	fmt.Println("PUBLIC KEY (for self-hosted instances):")
	fmt.Println(publicKeyB64)
	fmt.Println()
	fmt.Println("PRIVATE KEY (for hosted server only - KEEP SECRET!):")
	fmt.Println(privateKeyB64)
	fmt.Println()
	fmt.Println("Add to production .env:")
	fmt.Println("LICENSE_PRIVATE_KEY=" + privateKeyB64)
	fmt.Println()
	fmt.Println("Distribute to self-hosted customers:")
	fmt.Println("LICENSE_PUBLIC_KEY=" + publicKeyB64)
	fmt.Println()
	fmt.Println("For full setup instructions, see docs/INTERNAL_LICENSE_SETUP.md")
}
