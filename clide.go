package main

import (
	. "clide/model"
	"fmt"
	"os"
	"strings"
)

func main() {

	args := os.Args[1:]

	if len(args) > 0 && args[0] == "@" {
		// Process the args as builtins, don't run
		switch args[1] {
		case "ls":
			c := New(nil)

            if !c.Ok() {
                c.Run()
                return
            }

            leaves := c.Leaves()

			for _, leaf := range leaves {
				fmt.Printf("%s %s\n", leaf.Title(), leaf.Description())
			}

            return
		default:
            m, _ := Clide{}.Error(fmt.Sprintf("Unknown builtin command '%s'", args[1]))
            m.Run()
			return
		}
	}

	params := make(map[string]string)
	steps := make([]string, 0)

	clide := New(params)

	for _, arg := range args {
		if arg[0] == '-' {
			split := strings.Split(arg, "=")

			if len(split) == 1 {
                m, _ := clide.Error(fmt.Sprintf("Argument '%s' has no value", arg))
                m.Run()
				return
			}

			params[split[0][1:]] = split[1]
		} else {
			steps = append(steps, strings.ToLower(arg))
		}
	}

	clide = New(params)

    if !clide.Ok() {
        clide.Run()
        return
    }

	for _, step := range steps {
		i := clide.Index(step)

		if i == -1 {
            m, _ := clide.Error(fmt.Sprintf("Unknown command '%s'", step))
            m.Run()
			return
		}

        clide, _ = clide.SelectPath(i)
	}

	clide.Run()
}
