package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ErrAlreadyRunning = errors.New("Process is already running")

type Cmd struct {
	Command string
	command string
	args    []string
	cmd     *exec.Cmd
	die     chan (error)
}

func (c *Cmd) Start() error {
	if c.cmd != nil {
		return nil
	}

	if c.die == nil {
		c.die = make(chan (error))
	}

	parts := strings.Split(c.Command, " ")
	command := parts[0]
	args := parts[1:]

	c.cmd = exec.Command(command, args...)
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	err := c.cmd.Start()
	if err != nil {
		c.cmd = nil
		return err
	}

	go func() {
		c.die <- c.cmd.Wait()
	}()

	select {
	case e := <-c.die:
		return e
	case <-time.After(time.Second * 5):
		return nil
	}
}

func (c *Cmd) Restart() error {
	if c.cmd == nil {
		return c.Start()
	}

	c.cmd.Process.Kill()
	// Ignore the error

	// Wait until dies
	<-c.die
	c.cmd = nil
	return c.Start()
}
