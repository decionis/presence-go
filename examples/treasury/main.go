// A treasury wire read through Shadow Mode enforcement.
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

	// The Intent and Context this check protects: a high-value treasury
	// movement behind the high-value trigger floor.
	treasury := presence.ActionContext{
		Intent:           "corporate_treasury.wire_transfer.execute",
		Surface:          "treasury_workstation",
		ActorID:          "usr_treasury_ops_07",
		TargetResourceID: "acc_treasury_global_01",
		Amount:           25000000.00,
		Currency:         "USD",
	}
	fmt.Println("Intent:", treasury.Intent)

	result, err := client.RunSandboxScenario(context.Background(), presence.ScenarioBlock, key("treasury"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}
	fmt.Println("Policy verdict:", result.PolicyVerdict)

	enforcement := result.Decision.Enforcement
	if enforcement == nil {
		fmt.Println("No enforcement outcome on this deployment.")
		return
	}

	fmt.Println("Mode:", enforcement.Mode)
	fmt.Println("Effective verdict:", enforcement.EffectiveVerdict)
	fmt.Println("True verdict:", enforcement.TrueVerdict)
	fmt.Println("Downgraded:", enforcement.Downgraded)
	fmt.Println("Kill switch:", enforcement.KillSwitchActive)

	if enforcement.Mode == "SHADOW" && enforcement.Downgraded {
		fmt.Println("Shadow Mode recorded the block; execution proceeded unchanged.")
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
