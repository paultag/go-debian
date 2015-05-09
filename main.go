package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"pault.ag/x/go-debian/dependency"
	"pault.ag/x/go-debian/version"
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
	case "version":
		versionTool()
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
	version       | parse a version
	dependency    | parse dependency relations to json
		`,
	)
}

func dependencyTool() {
	if len(os.Args) <= 2 {
		fmt.Printf("Error! Give me a dependency to parse!\n")
		return
	}

	dep, err := dependency.Parse(os.Args[2])
	if err != nil {
		log.Fatalf("Oh no! %s", err)
		return
	}
	data, err := json.MarshalIndent(&dep, "", "  ")
	fmt.Printf("%s\n", data)
}

func versionTool() {
	if len(os.Args) <= 2 {
		fmt.Printf("Error! Give me a version to parse!\n")
		return
	}

	ver, err := version.Parse(os.Args[2])
	if err != nil {
		log.Fatalf("Oh no! %s", err)
		return
	}

	if ver.Revision == nil {
		fmt.Printf("[native] %d:%s\n", ver.Epoch, ver.Version)
	} else {
		fmt.Printf("         %d:%s-%s\n", ver.Epoch, ver.Version, *ver.Revision)
	}
}
