package main

import (
	"fmt"

	"github.com/hihebark/gomon/core"
	"github.com/hihebark/gomon/rule"
)

func main() {
	fmt.Println("\033[31m\033[1m[gomon]\033[0m 0.0.2\033[0m")
	fmt.Println("\033[31m\033[1m[gomon]\033[0m to restart at any time, enter `rs`")
	fmt.Println("\033[31m\033[1m[gomon]\033[0m watching dir(s): *.*")
	fmt.Println("\033[31m\033[1m[gomon]\033[0m watching extensions: go")
	fmt.Println("\033[31m\033[1m[gomon]\033[0m starting `go run main.go`")
	fmt.Println("")
	body, err := core.ReadFile("./.gomon.conf")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	rule, err := rule.NewRuleFromFile(body)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	err = core.Watcher(rule.Watch)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
