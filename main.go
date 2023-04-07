package main

import (
	"log"

	"github.com/budimanjojo/talhelper/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
