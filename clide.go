package main

import (
	"fmt"
    "os"
	"clide/path"
    "clide/model"
)

func main() {

	root, err := path.FindRoot()

	if err != nil {
		fmt.Println(err)
	}

	if root == "" {
		fmt.Println("No clide project found.")
        os.Exit(1)
	}

    model.Run(root)
}
