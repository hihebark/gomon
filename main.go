package main

import (
	"fmt"

	"github.com/hihebark/gomon/engine"
	"github.com/hihebark/gomon/observer"
	"github.com/hihebark/gomon/rule"
)

func main() {
	fmt.Println("\033[31m\033[1m[gomon]\033[0m 0.0.2\033[0m")
	mrule := &rule.Rule{}
	if engine.FileExists("gomon.conf") {
		body, err := engine.ReadFile("gomon.conf")
		if err != nil {
			fmt.Printf("[ERROR] While reading file:\n%v\n", err)
		}
		mrule, err = rule.NewRuleFromFile(body)
		if err != nil {
			fmt.Printf("[ERROR] While reading rules from file:\n%v\n", err)
		}
	} else {
		mrule = rule.NewRule()
	}
	fmt.Printf("\033[31m\033[1m[gomon]\033[0m to restart at any time, enter `%s`\n", mrule.Restartable)
	fmt.Printf("\033[31m\033[1m[gomon]\033[0m watching dir(s): %s\n", mrule.Watch)
	observer.NewObserver(*mrule).Start()
}
