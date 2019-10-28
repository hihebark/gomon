package main

import (
	"fmt"

	"github.com/hihebark/gomon/core"
	"github.com/hihebark/gomon/rules"
)

func main() {
	fmt.Println("gomon - 0.0.1")
	rules := rules.NewRules()
	fmt.Printf("rules: %v\n", rules)
	fmt.Printf("cmd: %s\n", core.Execute("ls", []string{"-la"}))
	core.Watcher()
}
