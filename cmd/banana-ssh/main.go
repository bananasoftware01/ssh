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
	// Help / Usage
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "-v", "--version":
			fmt.Printf("Banana Software CLI / SSH v%s (%s, built %s)\n", Version, Commit, Date)
			return
		case "help", "-h", "--help":
			printUsage()
			return
		case "tui", "local", "run":
			runLocalTUI()
			return
		}
	}

	// Determine whether to run in server mode or TUI mode
	isServer := false
	for _, arg := range os.Args[1:] {
		if arg == "serve" || arg == "server" || arg == "-server" || arg == "--server" ||
			arg == "-port" || arg == "--port" || arg == "-host" || arg == "--host" {
			isServer = true
			break
		}
	}

	// If explicitly asked for server or running server subcommand
	if isServer {
		runServer()
		return
	}

	// Default: Run interactive local TUI
	runLocalTUI()
}

func printUsage() {
	fmt.Printf("🍌 Banana Software CLI v%s\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  banana              Run interactive landing page in terminal (default)")
	fmt.Println("  banana serve        Start SSH server")
	fmt.Println("  banana version      Print version info")
	fmt.Println("  banana help         Show this help message")
	fmt.Println("\nServer options:")
	fmt.Println("  -host string        SSH server bind host (default: 0.0.0.0)")
	fmt.Println("  -port int           SSH server listen port (default: 2222)")
	fmt.Println("  -key string         Path to SSH host private key")
}

func runServer() {
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

	// Filter out "serve" / "server" subcommands before parsing flags
	var args []string
	for _, a := range os.Args[1:] {
		if a != "serve" && a != "server" {
			args = append(args, a)
		}
	}

	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.StringVar(&cfg.Host, "host", cfg.Host, "SSH server bind host")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "SSH server listen port")
	fs.StringVar(&cfg.HostKeyPath, "key", cfg.HostKeyPath, "Path to SSH host private key")
	_ = fs.Parse(args)

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
