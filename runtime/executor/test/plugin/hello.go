package main

import (
	"fmt"
	task "github.com/devfeel/dottask"
)

func Exec(ctx *task.TaskContext) error {
	fmt.Println("print from plugin")
	return nil
}
