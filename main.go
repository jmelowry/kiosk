package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/jmelowry/kiosk/menu"
)

func setupCleanupHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		resetTerminal()
		os.Exit(0)
	}()
}

func resetTerminal() {
	cmd := exec.Command("tput", "reset")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Also try stty sane which helps in many cases
	exec.Command("stty", "sane").Run()

	// Clear screen
	fmt.Print("\033[H\033[2J")
}

func main() {
	// Setup cleanup for unexpected termination
	setupCleanupHandler()

	// Ensure terminal is reset on exit
	defer resetTerminal()

	if err := menu.Start(); err != nil {
		resetTerminal()
		log.Fatal(err)
	}

	// Final reset to ensure clean exit
	resetTerminal()
}
