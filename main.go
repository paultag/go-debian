/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"pault.ag/x/go-debian/control"
	"pault.ag/x/go-debian/dependency"
	"pault.ag/x/go-debian/version"
)

// Main {{{

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
	case "control":
		controlTool()
		return
	case "dsc":
		dscTool()
		return
	case "dscSort":
		topsortDscTool()
		return
	}

	helpTool()

}

// }}}

// Help {{{

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
	control       | parse debian/control relations to json
	dsc           | parse a .dsc into json
		`,
	)
}

// }}}

// Dependency Tool {{{

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

// }}}

// Version Tool {{{

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

	if ver.Native {
		fmt.Printf("[native] %d:%s\n", ver.Epoch, ver.Version)
	} else {
		fmt.Printf("         %d:%s-%s\n", ver.Epoch, ver.Version, ver.Revision)
	}
}

// }}}

// Control Tool {{{

func controlTool() {
	if len(os.Args) <= 2 {
		fmt.Printf("Error! Give me a file to parse!\n")
		return
	}
	file, err := os.Open(os.Args[2])
	dep, err := control.ParseControl(bufio.NewReader(file))
	if err != nil {
		log.Fatalf("Oh no! %s", err)
		return
	}
	data, err := json.MarshalIndent(&dep, "", "  ")
	fmt.Printf("%s\n", data)
}

// }}}

// DSC Tool {{{

func dscTool() {
	if len(os.Args) <= 2 {
		fmt.Printf("Error! Give me a file to parse!\n")
		return
	}
	file, err := os.Open(os.Args[2])
	dep, err := control.ParseDsc(bufio.NewReader(file))
	if err != nil {
		log.Fatalf("Oh no! %s", err)
		return
	}
	data, err := json.MarshalIndent(&dep, "", "  ")
	fmt.Printf("%s\n", data)
}

// }}}

// Topsort DSC Tool {{{

func topsortDscTool() {
	if len(os.Args) <= 2 {
		fmt.Printf("Error! Give me a bunch of dsc files!\n")
		return
	}

	possies := map[string][]dependency.Possibility{}
	for _, dscFile := range os.Args[2:] {
		file, err := os.Open(dscFile)
		dep, err := control.ParseDsc(bufio.NewReader(file))
		if err != nil {
			log.Fatalf("Oh no! %s", err)
			return
		}

		arch, err := dependency.ParseArch("amd64")
		possies[dep.Source] = dep.BuildDepends.GetPossibilities(*arch)
	}

	order, err := dependency.SortDependencies(possies)
	if err != nil {
		fmt.Printf("Error!: %s", err)
		return
	}

	for _, el := range order {
		fmt.Printf(" - %s\n", el)
	}
}

// }}}

// vim: foldmethod=marker
