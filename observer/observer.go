package observer

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hihebark/gomon/rule"
)

type Observer struct {
	sync.Mutex
	rule rule.Rule
}

func NewObserver(rule rule.Rule) *Observer {
	return &Observer{
		rule: rule,
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
	fmt.Printf("\033[31m\033[1m[gomon]\033[0m starting `%s`\n", o.rule.ExecCommand)
	if err := o.observe(); err != nil {
		fmt.Printf("[ERROR] While watching for change:\n%v\n", err)
	}
}

func (o *Observer) restart() {
	fmt.Println("\033[31m\033[1m[gomon]\033[0m restarting due to changes...")
	// check if there is commande on restart.
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
	o.addPaths(watcher)
	err = watcher.Add(o.rule.Watch)
	if err != nil {
		return err
	}
	<-done
	return nil
}

func (o *Observer) addPaths(watcher *fsnotify.Watcher) {
	for _, v := range o.rule.Ignore {
		fmt.Printf("Ignoring: %s\n", v)
	}
}
