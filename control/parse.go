package control

import (
	"bufio"
	"fmt"
	"strings"
)

type Deb822 struct {
	Values map[string]string
	Order  []string
}

func ParseDeb822(reader *bufio.Reader) (ret *Deb822, ohshit error) {

	ret = &Deb822{
		Values: map[string]string{},
		Order:  []string{},
	}

	var key = ""
	var value = ""
	var noop = " \n\r\t"

	for {
		line, err := reader.ReadString('\n')

		if err != nil || line == "\n" {
			break
		}

		if line[0] == ' ' {
			line = line[1:]
			ret.Values[key] += "\n" + strings.Trim(line, noop)
			continue
		}

		els := strings.SplitN(line, ":", 2)

		switch len(els) {
		case 2:
			key = strings.Trim(els[0], noop)
			value = strings.Trim(els[1], noop)

			ret.Values[key] = value
			ret.Order = append(ret.Order, key)
			continue
		default:
			return nil, fmt.Errorf("The shit.")
		}
	}

	return
}
