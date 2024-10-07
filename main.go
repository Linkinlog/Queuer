package main

import (
	"github.com/linkinlog/queuer/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
