package deb_test

import (
	"fmt"
	"io"
	"log"

	"pault.ag/go/debian/deb"
)

func ExampleExamples() {
	debFile, err := deb.Load("/home/paultag/tmp/docker.io_1.8.3~ds1-2_amd64.deb")
	if err != nil {
		panic(err)
	}
	defer debFile.Close()

	for {
		entry, err := debFile.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if !entry.IsTarfile() {
			continue
		}

		tr, err := entry.Tarfile()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s:\n", entry.Name)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf(" -> %s:\n", hdr.Name)
		}

	}

}
