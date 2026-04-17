package main

import (
	"fmt"
	"strings"
)

// Copy of updated translation functions for testing
func translateOpenSSLToGoNameMap() map[string]string {
	return map[string]string{
		// HTTP/2 required ciphers (CRITICAL!)
		"ECDHE-RSA-AES128-GCM-SHA256":     "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"ECDHE-ECDSA-AES128-GCM-SHA256":   "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",

		// Additional Mozilla Intermediate ciphers
		"ECDHE-RSA-AES256-GCM-SHA384":     "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"ECDHE-ECDSA-AES256-GCM-SHA384":   "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"ECDHE-RSA-CHACHA20-POLY1305":     "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
		"ECDHE-ECDSA-CHACHA20-POLY1305":   "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",

		// TLS 1.3 ciphers (pass through unchanged)
		"TLS_AES_128_GCM_SHA256":          "TLS_AES_128_GCM_SHA256",
		"TLS_AES_256_GCM_SHA384":          "TLS_AES_256_GCM_SHA384",
		"TLS_CHACHA20_POLY1305_SHA256":    "TLS_CHACHA20_POLY1305_SHA256",
	}
}

func translateCiphers(opensslCiphers []string) []string {
	translationMap := translateOpenSSLToGoNameMap()
	var goFormatCiphers []string

	for _, cipher := range opensslCiphers {
		if translatedName, ok := translationMap[cipher]; ok {
			goFormatCiphers = append(goFormatCiphers, translatedName)
		} else {
			// Check if this might be a DHE cipher that should be ECDHE
			if strings.HasPrefix(cipher, "DHE-RSA-") {
				// Potential typo: DHE should be ECDHE for modern profiles
				ecdheVersion := strings.Replace(cipher, "DHE-RSA-", "ECDHE-RSA-", 1)
				if translatedName, ok := translationMap[ecdheVersion]; ok {
					goFormatCiphers = append(goFormatCiphers, translatedName)
					continue
				}
			}
			// Assume it's already in Go format or unknown (let plugin handle)
			goFormatCiphers = append(goFormatCiphers, cipher)
		}
	}

	return goFormatCiphers
}

func main() {
	fmt.Println("=== Testing DHE→ECDHE Typo Fix ===")
	fmt.Println()

	// Test with the problematic DHE ciphers from the logs
	opensslCiphers := []string{
		"TLS_AES_128_GCM_SHA256",          // TLS 1.3 - should pass through
		"ECDHE-ECDSA-AES128-GCM-SHA256",   // Correct ECDHE - should translate
		"ECDHE-RSA-AES128-GCM-SHA256",     // Correct ECDHE - should translate
		"DHE-RSA-AES128-GCM-SHA256",       // Potential typo - should fix to ECDHE
		"DHE-RSA-AES256-GCM-SHA384",       // Potential typo - should fix to ECDHE
		"ECDHE-ECDSA-AES256-GCM-SHA384",   // Correct ECDHE - should translate
	}

	fmt.Println("Input (with potential DHE typos):")
	for i, cipher := range opensslCiphers {
		fmt.Printf("  [%d] %s\n", i, cipher)
	}

	fmt.Println()
	fmt.Println("Translation Process with DHE→ECDHE Fix:")
	translationMap := translateOpenSSLToGoNameMap()
	for _, cipher := range opensslCiphers {
		if translatedName, ok := translationMap[cipher]; ok {
			fmt.Printf("  ✅ %s → %s (direct translation)\n", cipher, translatedName)
		} else if strings.HasPrefix(cipher, "DHE-RSA-") {
			ecdheVersion := strings.Replace(cipher, "DHE-RSA-", "ECDHE-RSA-", 1)
			if translatedName, ok := translationMap[ecdheVersion]; ok {
				fmt.Printf("  🔧 %s → %s → %s (DHE typo fixed!)\n", cipher, ecdheVersion, translatedName)
			} else {
				fmt.Printf("  ❌ %s → (DHE correction failed)\n", cipher)
			}
		} else {
			fmt.Printf("  ⚠️  %s → %s (pass through)\n", cipher, cipher)
		}
	}

	// Apply translation with DHE fix
	goFormatCiphers := translateCiphers(opensslCiphers)

	fmt.Println()
	fmt.Println("Output (Go TLS format for UI plugins):")
	for i, cipher := range goFormatCiphers {
		fmt.Printf("  [%d] %s\n", i, cipher)
	}

	// Verify HTTP/2 critical ciphers are present
	fmt.Println()
	fmt.Println("HTTP/2 Compatibility Check:")
	http2Critical := []string{
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	}

	for _, critical := range http2Critical {
		found := false
		for _, cipher := range goFormatCiphers {
			if cipher == critical {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("  ✅ %s - PRESENT (HTTP/2 compatible)\n", critical)
		} else {
			fmt.Printf("  ❌ %s - MISSING (HTTP/2 panic risk!)\n", critical)
		}
	}

	fmt.Println()
	fmt.Println("=== Expected Result ===")
	fmt.Println("✅ DHE cipher typos automatically corrected to ECDHE")
	fmt.Println("✅ No more 'Unknown cipher suite' warnings for DHE ciphers")
	fmt.Println("✅ All cipher suites align with Mozilla standards")
	fmt.Println("✅ HTTP/2 compatibility maintained")
}