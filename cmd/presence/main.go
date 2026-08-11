package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	presence "github.com/decionis/presence-go/presence"
)

var version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}
	if args[0] == "version" || args[0] == "--version" {
		fmt.Println(version)
		return nil
	}
	baseURL := env("PRESENCE_API_URL", "https://presence.decionis.com")
	client := presence.NewClient(baseURL, os.Getenv("PRESENCE_API_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	switch args[0] {
	case "quickstart":
		return runScenario(ctx, client, args[1:], "allow")
	case "sandbox":
		if len(args) < 3 || args[1] != "run" {
			return errors.New("usage: presence sandbox run <allow|block|escalate>")
		}
		return runScenario(ctx, client, args[2:], args[2])
	case "dossier":
		if len(args) != 3 || args[1] != "verify" {
			return errors.New("usage: presence dossier verify <dossier-id>")
		}
		if client.APIKey == "" {
			return errors.New("PRESENCE_API_KEY is required")
		}
		result, err := client.VerifyDossier(ctx, args[2])
		if err != nil {
			return err
		}
		return printJSON(result)
	case "doctor":
		if err := client.Health(ctx); err != nil {
			return err
		}
		fmt.Println("Presence API reachable at", baseURL)
		return nil
	default:
		return usage()
	}
}

func runScenario(ctx context.Context, client *presence.Client, args []string, fallback string) error {
	flags := flag.NewFlagSet("sandbox run", flag.ContinueOnError)
	idempotency := flags.String("idempotency-key", newKey(), "stable retry key")
	jsonOutput := flags.Bool("json", false, "print the full response")
	scenario := fallback
	if fallback == "allow" && len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		scenario = args[0]
		args = args[1:]
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	if client.APIKey == "" {
		return errors.New("PRESENCE_API_KEY is required")
	}
	result, err := client.RunSandboxScenario(ctx, presence.SandboxScenario(scenario), *idempotency)
	if err != nil {
		return err
	}
	if *jsonOutput {
		return printJSON(result)
	}
	fmt.Println("Policy verdict:", result.PolicyVerdict)
	fmt.Println("Effective verdict:", result.EffectiveVerdict, "(sandbox is SHADOW)")
	fmt.Println("Evidence class:", result.Decision.EvidenceClass)
	if result.Decision.VerificationBundle != nil {
		fmt.Println("Dossier:", result.Decision.VerificationBundle.DossierID)
	}
	return nil
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func newKey() string {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return fmt.Sprintf("quickstart-%d", time.Now().UnixNano())
	}
	return "quickstart-" + hex.EncodeToString(bytes[:])
}

func env(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func usage() error {
	return errors.New("usage: presence <quickstart|sandbox run|dossier verify|doctor|version>")
}
