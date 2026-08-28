package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"github.com/bananasoftware/banana-ssh/pkg/server"
	"github.com/bananasoftware/banana-ssh/pkg/ui"
)

var (
	Version = "1.0.0"
	Commit  = "head"
	Date    = "2026-08-28"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "tui", "local", "run":
			runLocalTUI()
			return
		case "version", "-v", "--version":
			fmt.Printf("Banana Software SSH v%s (%s, built %s)\n", Version, Commit, Date)
			return
		}
	}

	// Server mode
	cfg := server.DefaultConfig()

	// Environment variable overrides
	if envHost := os.Getenv("SSH_HOST"); envHost != "" {
		cfg.Host = envHost
	} else if envHost := os.Getenv("HOST"); envHost != "" {
		cfg.Host = envHost
	}

	if envPort := os.Getenv("SSH_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = p
		}
	} else if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = p
		}
	}

	if envKey := os.Getenv("SSH_KEY_PATH"); envKey != "" {
		cfg.HostKeyPath = envKey
	}

	// CLI flags
	var hostFlag string
	var portFlag int
	var keyFlag string
	var tuiFlag bool

	flag.StringVar(&hostFlag, "host", cfg.Host, "SSH server bind host")
	flag.IntVar(&portFlag, "port", cfg.Port, "SSH server listen port")
	flag.StringVar(&keyFlag, "key", cfg.HostKeyPath, "Path to SSH host private key")
	flag.BoolVar(&tuiFlag, "tui", false, "Run in interactive local TUI mode without SSH server")

	flag.Parse()

	if tuiFlag {
		runLocalTUI()
		return
	}

	cfg.Host = hostFlag
	cfg.Port = portFlag
	cfg.HostKeyPath = keyFlag

	fmt.Println("🍌 Banana Software SSH Terminal")
	fmt.Printf("Starting SSH server on %s:%d (Host key: %s)\n", cfg.Host, cfg.Port, cfg.HostKeyPath)
	fmt.Println("Connect via: ssh " + cfg.Host + " -p " + strconv.Itoa(cfg.Port))
	fmt.Println("--------------------------------------------------")

	if err := server.Start(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runLocalTUI() {
	m := ui.NewAppModel("127.0.0.1", os.Getenv("USER"))
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
