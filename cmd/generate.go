package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	GenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "A brief description of your command",
		Run: func(cobraCmd *cobra.Command, args []string) {
			diff, err := getDiff()
			if err != nil {
				log.Fatalf("error: %s", err)
				return
			}

			reqData := map[string]interface{}{
				"model":  model,
				"stream": false,
				"system": "Act as a Linux server running a REST API. You reply with only valid JSON.",
				"prompt": getStylePrompt() + getPrePrompt() + diff,
			}

			reqBytes, err := json.Marshal(reqData)
			if err != nil {
				log.Fatalf("error: failed to marshal JSON: %v", err)
			}

			url := fmt.Sprintf("%s/api/generate", address)
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
			if err != nil {
				log.Fatalf("error: failed to make POST request: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("error: failed to read response body: %v", err)
			}

			if !json.Valid(body) {
				log.Fatalf("error: response body is not valid JSON: %s", body)
			}

			fmt.Println(string(body))
		},
	}
)

func getDiff() (string, error) {
	// Abort if ollama isn't ready
	cmd := exec.Command("ollama", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New("ollama seems to not be running")
	}
	// Abort if git isn't ready
	cmd = exec.Command("git", "status")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return "", errors.New("directory seems to not be a git repository")
	}
	// Get the diff
	cmd = exec.Command("git", "diff", "--staged")
	out, err = cmd.CombinedOutput()
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return "", errors.New("there are no staged changes")
	}
	return diff, err
}

func getPrePrompt() string {
	languages := display.English.Languages()
	tag, _ := language.Parse(locale)
	return fmt.Sprintf(`
Based on the diff below, generate a commit title and message, both in %s. Make sure the title is a good, concise commit
title (preferably less than 73 characters). The commit message can be longer and more detailed. Your entire response
will be parsed as JSON, so make sure it is valid with a title and message key: don't include anything but JSON in your
reply'.
`, languages.Name(tag))
}

func getStylePrompt() string {
	if style == "conventional" {
		return `
Choose the most applicable prefix your commit title from the following list:
- "feat:": adds or removes a new feature
- "fix:": fixes a bug
- "refactor:": rewrites/restructures code, however doesn't' change any API behaviour
- "perf:": special refactor that improves performance
- "style:": not affecting any meaning (white-space, formatting, missing semicolons, etc.)
- "test:": adds missing tests or corrects existing tests
- "docs:": affects documentation only
- "build:": affects build components like build tool, ci pipeline, dependencies, project version, etc.
- "ops:": affects operational components like infrastructure, deployment, backup, recovery, etc.
- "chore:": miscellaneous changes e.g. modifying .gitignore
`
	}
	if style == "gitmoji" {
		return `
Choose the most applicable prefix your commit title from the following list:
- "🎨": improves structure/format of the code
- "⚡️": improves performance
- "🔥": removes code or files
- "🐛": fixes a bug
- "🚑️": critical hotfix
- "✨": introduces new features
- "📝": adds or updates documentation
- "🚀": deploys stuff
- "💄": adds or updates the UI and style files
- "🎉": begins a project
- "✅": adds, updates, or passes tests
- "🔒️": fixes security or privacy issues
- "🔐": adds or updates secrets
- "🔖": release/version tags
- "🚨": fixes compiler/linter warnings
- "🚧": work in progress
- "💚": fixes CI build
- "⬇️": downgrades dependencies
- "⬆️": upgrades dependencies
- "📌": pins dependencies to specific versions
- "👷": adds or updates CI build system
- "📈": adds or updates analytics or track code
- "♻️": refactors code
- "➕": adds a dependency
- "➖": removes a dependency
- "🔧": adds or updates configuration files
- "🔨": adds or updates development scripts
- "🌐": internationalization and localization
- "✏️": fixes typos
- "💩": adds bad code that needs to be improved
- "⏪️": reverts changes
- "🔀": merges branches
- "📦️": adds or updates compiled files or packages
- "👽️": updates code due to external API changes
- "🚚": moves or renames resources (e.g.: files, paths, routes)
- "📄": adds or update license
- "💥": introduces breaking changes
- "🍱": adds or updates assets
- "♿️": improves accessibility
- "💡": adds or updates comments in source code
- "🍻": writes code drunkenly
- "💬": adds or updates text and literals
- "🗃️": performs database related changes
- "🔊": adds or updates logs
- "🔇": removes logs
- "👥": adds or updates contributor(s)
- "🚸": improves user experience / usability
- "🏗️": makes architectural changes
- "📱": works on responsive design
- "🤡": mocks things
- "🥚": adds or updates an easter egg
- "🙈": adds or updates a .gitignore file
- "📸": adds or updates snapshots
- "⚗️": performs experiments
- "🔍️": improves SEO
- "🏷️": adds or updates types
- "🌱": adds or updates seed files
- "🚩": adds, updates, or removes feature flags
- "🥅": catches errors
- "💫": adds or updates animations and transitions
- "🗑️": deprecates code that needs to be cleaned up
- "🛂": works on code related to authorization, roles and permissions
- "🩹": simple fix for a non-critical issue
- "🧐": data exploration/inspection
- "⚰️": removes dead code
- "🧪": adds a failing test
- "👔": adds or updates business logic
- "🩺": adds or updates healthcheck
- "🧱": infrastructure related changes
- "🧑‍💻": improves developer experience
- "💸": adds sponsorships or money related infrastructure
- "🧵": adds or updates code related to multithreading or concurrency
- "🦺": adds or updates code related to validation
`
	}
	return ""
}
