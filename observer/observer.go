package observer

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hihebark/gomon/engine"
	"github.com/hihebark/gomon/rule"
)

type Observer struct {
	sync.Mutex
	rule rule.Rule
	args []string
	cmd  *exec.Cmd
}

func NewObserver(rule rule.Rule, args []string) *Observer {
	return &Observer{
		rule: rule,
		args: args,
	}
}

func (o *Observer) Start() {
	userStdIn := make(chan string)
	go read(userStdIn)
	go func() {
		select {
		case restart := <-userStdIn:
			if restart == o.rule.Restartable {
				o.restart()
			}
		}
	}()
	///.command := o.rule.ExecCommand
	///.if len(o.args) != 0 {
	///.	command = strings.Join(o.args, " ")
	///.}
	///.fmt.Printf("\033[31m\033[1m[gomon]\033[0m starting `%s`\n", command)
	///.cmd, err := engine.ExecuteAndCapture("go", []string{"run", "main.go"})
	///.if err != nil {
	///.	fmt.Printf("[ERROR] While runnig this command %s\n", command)
	///.}
	///.o.cmd = cmd
	if err := o.observe(); err != nil {
		fmt.Printf("[ERROR] While watching for change:\n%v\n", err)
	}
}

func (o *Observer) restart() {
	fmt.Println("\033[31m\033[1m[gomon]\033[0m restarting due to changes...")
	err := engine.KillCommand(o.cmd)
	if err != nil {
		fmt.Printf("[ERROR] While killing the %v %v\n", o.cmd.Args, err)
	}
	// check if there is commande on restart.
	//if o.rule.Events.OnRestart != "" {
	//
	//}
}

func (o *Observer) exit() {
	// OnExit commande if none os.Exit(0)
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
					fmt.Printf("New Event %s, %v\n", event.Name, event.Op)
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

func (o *Observer) addPaths(watcher *fsnotify.Watcher) error {
	if o.rule.Watch == "*.*" {
		err := watcher.Add(".")
		return err
	} else {
		for _, v := range strings.Split(o.rule.Watch, ",") {
			err := watcher.Add(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
