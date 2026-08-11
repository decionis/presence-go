// Your first Presence Check.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/decionis/presence-go/presence"
)

func main() {
	apiKey := os.Getenv("PRESENCE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Set PRESENCE_API_KEY to the one-time sandbox credential")
		os.Exit(1)
	}

	client := presence.NewClient(env("PRESENCE_API_URL", "https://presence.decionis.com"), apiKey)

	result, err := client.RunSandboxScenario(context.Background(), presence.ScenarioAllow, key("quickstart"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}

	fmt.Println("Policy verdict:", result.PolicyVerdict)
	fmt.Println("Effective verdict:", result.EffectiveVerdict, "(sandbox is SHADOW)")
	fmt.Println("Evidence class:", result.Decision.EvidenceClass)
	if bundle := result.Decision.VerificationBundle; bundle != nil {
		fmt.Println("Dossier:", bundle.DossierID)
	}
}

func env(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func key(prefix string) string {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}
	return prefix + "-" + hex.EncodeToString(bytes[:])
}
