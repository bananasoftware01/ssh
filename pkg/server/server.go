package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	wishbubbletea "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/bananasoftware/banana-ssh/pkg/ui"
)

type Config struct {
	Host        string
	Port        int
	HostKeyPath string
	IdleTimeout time.Duration
	MaxTimeout  time.Duration
}

func DefaultConfig() Config {
	return Config{
		Host:        "0.0.0.0",
		Port:        2222,
		HostKeyPath: ".ssh/banana_ed25519",
		IdleTimeout: 15 * time.Minute,
		MaxTimeout:  1 * time.Hour,
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	var clientAddr string
	if remote := s.RemoteAddr(); remote != nil {
		clientAddr = remote.String()
	}
	sshUser := s.User()

	// Create new app model for this session
	m := ui.NewAppModel(clientAddr, sshUser)

	opts := wishbubbletea.MakeOptions(s)
	return m, opts
}

func Start(cfg Config) error {
	addr := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))

	// Ensure directory for host key exists
	if cfg.HostKeyPath != "" {
		dir := filepath.Dir(cfg.HostKeyPath)
		if dir != "." && dir != "" {
			_ = os.MkdirAll(dir, 0700)
		}
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "2006/01/02 15:04:05",
		Prefix:          "banana-ssh",
	})

	srv, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithIdleTimeout(cfg.IdleTimeout),
		wish.WithMaxTimeout(cfg.MaxTimeout),
		wish.WithMiddleware(
			wishbubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.MiddlewareWithLogger(logger),
		),
	)
	if err != nil {
		return fmt.Errorf("create wish server: %w", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Starting Banana Software SSH server", "addr", addr, "host_key", cfg.HostKeyPath)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			logger.Error("SSH server error", "err", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("Stopping SSH server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown ssh server: %w", err)
	}

	logger.Info("SSH server stopped")
	return nil
}
