package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"pault.ag/x/go-debian/dependency"
)

func main() {
	if len(os.Args) <= 1 {
		helpTool()
		return
	}

	switch os.Args[1] {
	case "help":
		helpTool()
		return
	case "dependency":
		dependencyTool()
		return
	}

	helpTool()

}

func helpTool() {
	fmt.Printf(
		"%s\n",
		`
go-debian
=========

Commands:

	help          | show this help
	dependency    | parse dependency relations to json
		`,
	)
}

func dependencyTool() {
	if len(os.Args) <= 2 {
		fmt.Printf("Error! Give me a version to parse!\n")
		return
	}

	dep, err := dependency.Parse(os.Args[2])
	if err != nil {
		log.Fatalf("Oh no! %s", err)
	}
	data, err := json.MarshalIndent(&dep, "", "  ")
	fmt.Printf("%s\n", data)
}
