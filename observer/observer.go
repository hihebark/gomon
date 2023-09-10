package observer

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hihebark/gomon/engine"
	"github.com/hihebark/gomon/rule"
)

const (
	// DefaultCmd contant
	DefaultCmd string = "go run main.go"
)

// Observer struct
type Observer struct {
	sync.Mutex
	rule rule.Rule
	args []string
	cmd  *exec.Cmd
}

// NewObserver create new observer
func NewObserver(rule rule.Rule, args []string) *Observer {
	return &Observer{
		rule: rule,
		args: args,
	}
}

// Start the observer
func (o *Observer) Start() {
	userStdIn := make(chan string)
	go read(userStdIn)
	go func() {
		select {
		case cases := <-userStdIn:
			if cases == o.rule.Restartable {
				o.restart()
			}
			if cases == "exit" || cases == "stop" {
				o.exit()
			}
		}
	}()
	command := DefaultCmd
	if len(o.args) > 1 {
		command = strings.Join(o.args, " ")
	} else if o.rule.ExecCommand != "" {
		command = o.rule.ExecCommand
	}
	c := strings.Split(command, " ")
	fmt.Printf("\033[31m\033[1m[gomon]\033[0m starting `\033[1m%s\033[0m`\n", command)
	go func() {
		if err := o.observe(); err != nil {
			fmt.Printf("[ERROR] While watching for change:\n%v\n", err)
		}
	}()
	// Cant reach the cmd it still executing ....
	cmd, err := engine.ExecuteAndCapture(c[0], c[1:])
	if err != nil {
		fmt.Printf("[ERROR] While runnig this command %s, %v\n", command, err)
	}
	o.cmd = <-cmd
	o.cmd.Wait()
}

func (o *Observer) restart() {
	fmt.Println("\033[31m\033[1m[gomon]\033[0m restarting due to changes...")
	err := engine.KillCommand(o.cmd)
	if err != nil {
		fmt.Printf("[ERROR] App already finished %v\n", err)
	}
	// check if there is commande on restart.
	if o.rule.Events.OnRestart != "" {
		fmt.Println("OnRestart", o.rule.Events.OnRestart)
	}
}

func (o *Observer) exit() {
	// OnExit commande if none os.Exit(0)
	fmt.Printf("\033[31m\033[1m[gomon]\033[0m Exiting...")
	os.Exit(0)
}

func read(input chan<- string) {
	for {
		var str string
		_, err := fmt.Scanf("%s\n", &str)
		if err != nil {
			panic(err)
		}
		input <- str
	}
}

func (o *Observer) observe() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if string(event.Op) != "" && string(event.Op) != "Chmod" {
					o.restart()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}()
	err = o.addPaths(watcher)
	if err != nil {
		return err
	}
	<-done
	return nil
}

func contains(str string, arr []string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}

func (o *Observer) addPaths(watcher *fsnotify.Watcher) error {
	if o.rule.Watch == "*.*" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		return filepath.Walk(".", func(fp string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// add ignored folders and files
			if !contains(info.Name(), o.rule.Ignore) {
				if contains(filepath.Ext(info.Name()), o.rule.Ext) {
					err := watcher.Add(path.Join(wd, fp))
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
	for _, v := range strings.Split(o.rule.Watch, ",") {
		err := watcher.Add(v)
		if err != nil {
			return err
		}
	}
	return nil
}
