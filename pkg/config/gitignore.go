package config

import (
	"fmt"
	"os"

	ignore "github.com/sabhiram/go-gitignore"
)

func (config TalhelperConfig) GenerateGitignore(outputDir string) error {
	for _, node := range config.Nodes {
		fileName := config.ClusterName + "-" + node.Hostname + ".yaml"
		err := createGitIgnore(outputDir, fileName)
		if err != nil {
			return err
		}
	}
	fileName := "talosconfig"
	err := createGitIgnore(outputDir, fileName)
	if err != nil {
		return err
	}
	fmt.Printf("generated .gitignore file in %s/.gitignore\n", outputDir)
	return nil
}

func createGitIgnore(path, line string) error {
	ignorefPath := path + "/.gitignore"

	file, err := os.OpenFile(ignorefPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	ignoreObject, err := ignore.CompileIgnoreFile(ignorefPath)
	if err != nil {
		file.Close()
		return err
	}

	if !ignoreObject.MatchesPath(line) {
		if _, err := file.Write([]byte(line + "\n")); err != nil {
			file.Close()
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}

		return nil
	}

	return nil
}
