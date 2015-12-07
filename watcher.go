package main

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	Dir    string
	Ignore string
	w      *fsnotify.Watcher
}

func (w *Watcher) Setup() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	w.w = watcher

	err = w.addListeners()
	if err != nil {
		w.w.Close()
		return err
	}
	return nil
}

func (w *Watcher) Monitor(f func() error) error {
	for {
		select {
		case event := <-w.w.Events:
			log.Print(event)
			f()
		case err := <-w.w.Errors:
			log.Print("Error:", err)
		}
	}
}

func (w *Watcher) Close() error {
	return w.w.Close()
}

func (w *Watcher) addListeners() error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ignore := strings.Split(w.Ignore, ",")
	err := recursiveAdd(w.Dir, ignore, w.w, wg)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func recursiveAdd(dir string, ignore []string, watcher *fsnotify.Watcher, wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, i := range ignore {
		if dir == i {
			log.Println("Ignoring directry", dir)
			return nil
		}
	}

	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New(dir + "Is not a directory")
	}

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, entry := range files {
		name := entry.Name()
		if name == "." || name == ".." {
			//			log.Println("Found", name)
			continue
		}
		if entry.IsDir() {
			//			log.Println("Found subdirectory", name)
			wg.Add(1)
			go recursiveAdd(path.Join(dir, name), ignore, watcher, wg)
		}
	}
	return nil

}
