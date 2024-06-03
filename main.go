package main

import (
	"log"

	"github.com/budimanjojo/talhelper/v3/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
