package main

import (
	"fmt"
	"os"
	"strings"

	clide "github.com/TeddyRandby/clide/app"
)

func is_builtin(args []string) bool {
  return len(args) > 0 && args[0][0] == '@'
}

func get_builtin(args []string) string {
  return args[0][1:]
}

func main() {
	args := os.Args[1:]

	params := make(map[string]string)
	c := clide.New(params)

	if !c.Ok() {
		c.Run()
		return
	}

	if is_builtin(args) {
    c.Builtin(get_builtin(args))
    return
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
