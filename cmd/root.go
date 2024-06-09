package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"log"
	"os/exec"
	"strings"
)

var (
	address string
	model   string
	style   string
	locale  string
	RootCmd = &cobra.Command{
		Use:   "gcai",
		Short: "hello world",
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			if !isValidModel() {
				log.Fatalf("invalid locale: %s", model)
			}
			if !isValidStyle() {
				log.Fatalf("invalid style: %s", style)
			}
			if !isValidLocale() {
				log.Fatalf("invalid locale: %s", locale)
			}
			return nil
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&address, "address", "a", "http://127.0.0.1:11434", "schema and root path of the Ollama server")
	RootCmd.PersistentFlags().StringVarP(&model, "model", "m", "llama3", "name of the AI model to use for generating the commit message")
	RootCmd.PersistentFlags().StringVarP(&style, "style", "s", "simple", "style of the commit message, either 'simple', 'conventional' or 'gitmoji'")
	RootCmd.PersistentFlags().StringVarP(&locale, "locale", "l", "en", "language in which to generate the commit message, as a locale code")

	RootCmd.AddCommand(GenerateCmd)
}

func isValidModel() bool {
	if !strings.Contains(model, ":") {
		model = model + ":latest"
	}

	cmd := exec.Command("ollama", "list")
	out, _ := cmd.CombinedOutput()
	return strings.Contains(string(out), model)
}

func isValidStyle() bool {
	return style == "simple" || style == "conventional" || style == "gitmoji"
}

func isValidLocale() bool {
	_, err := language.Parse(locale)
	return err == nil
}
