package main

import (
	"flag"
	"log"
)

func main() {
	watchDir := flag.String("d", ".", "Watch directory")
	inPort := flag.Int("p", 3000, "Listen port")
	outPort := flag.Int("o", 3001, "cmd listening port")
	ignore := flag.String("i", "log,tmp,.git", "Comman separated list of directories to ignore")
	command := flag.String("c", "rails s -p 3001", "Command to run")
	flag.Parse()

	cmd := &Cmd{
		Command: *command,
	}

	w := &Watcher{
		Dir:    *watchDir,
		Ignore: *ignore,
	}

	log.Print("Registering file observers")
	err := w.Setup()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting process")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening in localhost:", *inPort)
	http := &Http{*inPort, *outPort}
	go http.ListenAndServe()

	log.Print("Waiting for changes")
	w.Monitor(cmd.Restart)

}
