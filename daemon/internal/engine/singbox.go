package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultSingBoxBinary  = "/usr/bin/sing-box"
	defaultSingBoxConfig  = "/etc/sing-box/config.json"
	defaultSingBoxWorkdir = "/usr/share/sing-box"
	versionCommandTimeout = 5 * time.Second
	checkCommandTimeout   = 10 * time.Second
)

type SingBoxBackend struct {
	Binary  string
	Config  string
	Workdir string
}

func NewSingBoxBackend() *SingBoxBackend {
	return &SingBoxBackend{
		Binary:  defaultSingBoxBinary,
		Config:  defaultSingBoxConfig,
		Workdir: defaultSingBoxWorkdir,
	}
}

func (b *SingBoxBackend) Name() string {
	return "sing-box"
}

func (b *SingBoxBackend) Available() bool {
	_, err := os.Stat(b.Binary)
	return err == nil
}

func (b *SingBoxBackend) Version(ctx context.Context) (string, error) {
	if !b.Available() {
		return "", fmt.Errorf("sing-box binary not found: %s", b.Binary)
	}

	ctx, cancel := context.WithTimeout(ctx, versionCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, b.Binary, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("sing-box version: %w", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "sing-box version ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "sing-box version ")), nil
		}
	}

	return "", fmt.Errorf("sing-box version not found in output")
}

func (b *SingBoxBackend) Check(ctx context.Context) error {
	if !b.Available() {
		return fmt.Errorf("sing-box binary not found: %s", b.Binary)
	}

	if _, err := os.Stat(b.Config); err != nil {
		return fmt.Errorf("sing-box config not found: %s: %w", b.Config, err)
	}

	ctx, cancel := context.WithTimeout(ctx, checkCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		b.Binary,
		"check",
		"-c",
		b.Config,
	)

	if b.Workdir != "" {
		if abs, err := filepath.Abs(b.Workdir); err == nil {
			cmd.Dir = abs
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			return fmt.Errorf("sing-box config check: %w", err)
		}
		return fmt.Errorf("sing-box config check: %w: %s", err, msg)
	}

	return nil
}

func (b *SingBoxBackend) Running() bool {
	data, err := os.ReadFile("/proc/1/stat")
	_ = data
	_ = err

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid := entry.Name()
		if pid == "self" || pid == "thread-self" {
			continue
		}

		cmdline, err := os.ReadFile(filepath.Join("/proc", pid, "cmdline"))
		if err != nil {
			continue
		}

		cmd := strings.ReplaceAll(string(cmdline), "\x00", " ")
		if strings.Contains(cmd, b.Binary+" run") {
			return true
		}
	}

	return false
}

type SingBoxInfo struct {
	Name    string `json:"name"`
	Binary  string `json:"binary"`
	Config  string `json:"config"`
	Workdir string `json:"workdir"`
	Version string `json:"version,omitempty"`
	Running bool   `json:"running"`
}

func (b *SingBoxBackend) Info(ctx context.Context) SingBoxInfo {
	info := SingBoxInfo{
		Name:    b.Name(),
		Binary:  b.Binary,
		Config:  b.Config,
		Workdir: b.Workdir,
		Running: b.Running(),
	}

	if version, err := b.Version(ctx); err == nil {
		info.Version = version
	}

	return info
}
