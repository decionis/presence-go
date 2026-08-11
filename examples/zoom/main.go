// An executive meeting instruction held for review.
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

	// The Intent and Context this check protects: an instruction given on a
	// video call, confirmed before anyone acts on it.
	meeting := presence.ActionContext{
		Intent:           "meeting.executive_instruction.confirm",
		Surface:          "zoom_meeting",
		ActorID:          "usr_ceo_0001",
		TargetResourceID: "mtg_board_q3",
	}
	fmt.Println("Intent:", meeting.Intent)

	result, err := client.RunSandboxScenario(context.Background(), presence.ScenarioEscalate, key("meeting"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}
	fmt.Println("Policy verdict:", result.PolicyVerdict)

	switch result.PolicyVerdict {
	case presence.VerdictAuthorized:
		fmt.Println("Execute: the instruction is trusted.")
	case presence.VerdictEscalate:
		fmt.Println("Hold: instruction awaits human review.")
	default:
		fmt.Println("Reject: do not act on the instruction.")
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
