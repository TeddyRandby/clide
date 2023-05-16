package main

import (
	. "clide/model"
	"fmt"
	"os"
	"strings"
)

func main() {

	args := os.Args[1:]

	if args[0] == "$" {
		// Process the args as builtins, don't run
	}

	params := make(map[string]string)
	steps := make([]string, 0)

    for _, arg := range args {
        if arg[0] == '-' {
            split := strings.Split(arg, "=")

            if len(split) == 1 {
                fmt.Println("Error: Parameter", arg, "has no value")
                os.Exit(1)
            }

            params[split[0][1:]] = split[1]
        } else {
            steps = append(steps, strings.ToLower(arg))
        }
    }

	clide := New(&params)

    for _, step := range steps {
        i := clide.Index(step)

        if i == -1 {
            fmt.Println("Error: Command", step, "not found")
            os.Exit(1)
        }

        clide = clide.SelectPath(i)
    }

	clide.Run()
}
