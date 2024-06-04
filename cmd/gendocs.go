package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var GendocsCmd = &cobra.Command{
	Use:    "gendocs <output-dir>",
	Short:  "Generate documentation for the CLI",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		if err := os.MkdirAll(dir, 0o777); err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(filepath.Join(dir, "cli.md"))
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		data := &bytes.Buffer{}
		if _, err := io.WriteString(file, "# CLI\n\n"); err != nil {
			log.Fatal(err)
		}

		if err := genMarkdownReference(rootCmd, data, linkHandler); err != nil {
			log.Fatal(err)
		}

		if _, err := io.WriteString(file, editData(data.String())); err != nil {
			log.Fatal(err)
		}
	},
}

// linkHandler will change default generated link to markdown section link
func linkHandler(name string) string {
	base := strings.TrimSuffix(name, ".md")
	base = strings.ReplaceAll(base, "_", "-")
	return "#" + strings.ToLower(base)
}

// genMarkdownReference is the same as GenMarkDownTree, but
// with custom filePrepender and linkHandler
func genMarkdownReference(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.DisableAutoGenTag = true
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		if err := genMarkdownReference(c, w, linkHandler); err != nil {
			return err
		}
	}

	return doc.GenMarkdownCustom(cmd, w, linkHandler)
}

// editData will take the generated doc data and do something to it
func editData(data string) string {
	// Trim whitespaces
	data = strings.TrimSpace(data)

	// Replace "-----------" with "\n```"
	re := regexp.MustCompile(`^\s*-{5,}$`)
	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if re.MatchString(line) {
			lines[i] = "\n```"
		}
	}
	return strings.Join(lines, "\n")
}

func init() {
	RootCmd.AddCommand(gendocsCmd)
}
