package models

import (
	"context"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	CustomDNSServer string
}

func EnsureDNSPort(addr string) string {
	if addr == "" {
		return "8.8.8.8:53"
	}

	addr = strings.TrimSpace(addr)

	if host, port, err := net.SplitHostPort(addr); err == nil {
		if portNum, err := strconv.Atoi(port); err == nil && portNum > 0 && portNum <= 65535 {
			if net.ParseIP(host) != nil || isValidHostname(host) {
				return net.JoinHostPort(host, port)
			}
		}
	}

	if ip := net.ParseIP(addr); ip != nil {
		return net.JoinHostPort(addr, "53")
	}

	if isValidHostname(addr) && !strings.Contains(addr, ":") {
		return addr + ":53"
	}

	return "8.8.8.8:53"
}

func isValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}
	if strings.Contains(hostname, " ") ||
		strings.HasPrefix(hostname, "-") ||
		strings.HasSuffix(hostname, "-") ||
		strings.Contains(hostname, "..") {
		return false
	}

	labels := strings.Split(hostname, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if len(label) > 0 {
			first := label[0]
			last := label[len(label)-1]
			if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || (first >= '0' && first <= '9')) ||
				!((last >= 'a' && last <= 'z') || (last >= 'A' && last <= 'Z') || (last >= '0' && last <= '9')) {
				return false
			}
		}
		for _, r := range label {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-') {
				return false
			}
		}
	}

	if strings.Count(hostname, ".") == 0 && strings.Contains(hostname, "-") {
		if len(hostname) > 15 || strings.Count(hostname, "-") > 2 {
			return false
		}
	}

	return true
}

type ConfigStore struct {
	mu  sync.RWMutex
	cfg Config
}

func NewConfigStore(cfg Config) *ConfigStore {
	return &ConfigStore{cfg: cfg}
}

func (s *ConfigStore) Get() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *ConfigStore) Set(cfg Config) {
	s.mu.Lock()
	s.cfg = cfg
	s.mu.Unlock()
}

type configKey struct{}

func ContextWithConfig(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configKey{}, cfg)
}

func ConfigFromContext(ctx context.Context) (Config, bool) {
	cfg, ok := ctx.Value(configKey{}).(Config)
	return cfg, ok
}
