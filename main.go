package main

import (
	"fmt"
	"os"

	"github.com/budimanjojo/talhelper/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
