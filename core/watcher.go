package core

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func Watcher(path string) error {
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					if string(event.Op) != "Chmod" {
						fmt.Printf("New Event %s, %v\n", event.Name, event.Op)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}()
	err = watcher.Add(path)
	if err != nil {
		return err
	}
	<-done
	return nil
}
