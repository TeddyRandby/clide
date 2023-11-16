package main

import (
	"fmt"
	"os"
	"strings"

	clide "github.com/TeddyRandby/clide/app"
)

func main() {
	args := os.Args[1:]

	params := make(map[string]string)
	c := clide.New(params)

	if !c.Ok() {
		c.Run()
		return
	}

	if len(args) > 0 && args[0] == "@" {
		// Process the args as builtins, don't run
		switch args[1] {
		case "ls":

			leaves := c.Leaves()

			for _, leaf := range leaves {
				fmt.Printf("%s\t%s\t%s\n", leaf.Title(), leaf.Description(), leaf.Steps())
			}

			return
		default:
			m, _ := c.Error(fmt.Sprintf("Unknown builtin command '%s'", args[1]))
			m.Run()
			return
		}
	}

	steps := make([]string, 0)

	for _, arg := range args {
		if arg[0] == '-' {
			split := strings.Split(arg, "=")

			if len(split) == 1 {
				m, _ := c.Error(fmt.Sprintf("Argument '%s' has no value", arg))
				m.Run()
				return
			}

			params[split[0][1:]] = split[1]
		} else {
			steps = append(steps, strings.ToLower(arg))
		}
	}

	c = clide.New(params)

	if !c.Ok() {
		c.Run()
		return
	}

	for _, step := range steps {
		c, _ = c.SelectPath(step)
	}

	c.Run()
}
