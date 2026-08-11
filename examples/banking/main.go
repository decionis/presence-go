// A wire transfer with a verified Presence Record.
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

	// The Intent and Context this check protects. The sandbox seals a
	// deterministic fixture; a production Presence Check opens a Session with
	// this exact payload (docs/presence-check.md).
	wire := presence.ActionContext{
		Intent:           "banking.wire_transfer.execute",
		Surface:          "web_banking_portal",
		ActorID:          "usr_treasurer_014",
		TargetResourceID: "acc_operating_01",
		Amount:           250000,
		Currency:         "USD",
		RecipientIBAN:    "DE89370400440532013000",
	}
	fmt.Println("Intent:", wire.Intent)
	fmt.Printf("Amount: %.0f %s to ACME\n", wire.Amount, wire.Currency)

	ctx := context.Background()
	result, err := client.RunSandboxScenario(ctx, presence.ScenarioAllow, key("wire"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}
	fmt.Println("Policy verdict:", result.PolicyVerdict)

	bundle := result.Decision.VerificationBundle
	if bundle == nil {
		fmt.Println("No verification bundle on this deployment; nothing to verify.")
		return
	}

	verification, err := client.VerifyDossier(ctx, bundle.DossierID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}
	fmt.Println("Dossier:", verification.DossierID)
	fmt.Println("Signature valid:", verification.SignatureValid)
	fmt.Println("Chain valid:", verification.ChainValid)
	fmt.Println("Valid:", verification.Valid)

	if result.PolicyVerdict == presence.VerdictAuthorized && verification.Valid {
		fmt.Println("Execute: wire transfer may proceed.")
	} else {
		fmt.Println("Reject: do not execute the wire transfer.")
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
