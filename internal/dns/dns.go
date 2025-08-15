package dns

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/mooship/blokilo/internal/models"
)

const (
	MaxRetries       = 2
	DefaultTimeout   = 5 * time.Second
	RetryDelay       = 100 * time.Millisecond
	SystemDNSTimeout = 3 * time.Second
)

func GetSystemDNS() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), SystemDNSTimeout)
	defer cancel()

	var dnsServer string
	var err error

	switch runtime.GOOS {
	case "windows":
		dnsServer, err = getWindowsDNS(ctx)
	case "darwin":
		dnsServer, err = getMacDNS(ctx)
	case "linux":
		dnsServer, err = getLinuxDNS(ctx)
	default:
		dnsServer, err = getLinuxDNS(ctx)
	}

	if err != nil {
		return "System", fmt.Errorf("failed to get system DNS: %w", err)
	}

	return dnsServer, nil
}

func getWindowsDNS(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "ipconfig", "/all")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run ipconfig: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "DNS Servers") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				dns := strings.TrimSpace(parts[1])
				if dns != "" {
					return dns + ":53", nil
				}
			}
		}
	}

	return "", fmt.Errorf("failed to find DNS server in ipconfig output")
}

func getMacDNS(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "scutil", "--dns")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run scutil: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver[") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				dns := strings.TrimSpace(parts[1])
				if dns != "" {
					return dns + ":53", nil
				}
			}
		}
	}

	return "", fmt.Errorf("failed to find DNS server in scutil output")
}

func getLinuxDNS(ctx context.Context) (string, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return "", fmt.Errorf("failed to open /etc/resolv.conf: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				dns := parts[1]
				if dns != "" {
					return dns + ":53", nil
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan /etc/resolv.conf: %w", err)
	}

	return "", fmt.Errorf("failed to find nameserver in /etc/resolv.conf")
}

func TestDomainDNS(ctx context.Context, domain string, dnsServer string) models.TestResult {
	if domain == "" {
		return models.TestResult{
			Domain:       domain,
			Status:       models.StatusError,
			ResponseTime: 0,
			Err:          fmt.Errorf("domain cannot be empty"),
		}
	}

	domain = strings.TrimSpace(domain)
	if strings.Contains(domain, " ") {
		return models.TestResult{
			Domain:       domain,
			Status:       models.StatusError,
			ResponseTime: 0,
			Err:          fmt.Errorf("invalid domain name: contains spaces"),
		}
	}

	const maxRetries = MaxRetries
	const timeout = DefaultTimeout

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result := testDomainDNSOnce(ctx, domain, dnsServer, timeout)

		if result.Status != models.StatusError {
			return result
		}

		if attempt == maxRetries {
			return result
		}

		select {
		case <-ctx.Done():
			return models.TestResult{
				Domain:       domain,
				Status:       models.StatusError,
				ResponseTime: 0,
				Err:          ctx.Err(),
			}
		case <-time.After(RetryDelay):
		}
	}

	return models.TestResult{
		Domain:       domain,
		Status:       models.StatusError,
		ResponseTime: 0,
		Err:          fmt.Errorf("max retries exceeded"),
	}
}

func testDomainDNSOnce(ctx context.Context, domain string, dnsServer string, timeout time.Duration) models.TestResult {
	c := new(dns.Msg)
	c.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	client := new(dns.Client)
	client.Timeout = timeout

	if dnsServer == "" {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		start := time.Now()
		resolver := &net.Resolver{}
		ips, err := resolver.LookupIPAddr(timeoutCtx, domain)
		elapsed := time.Since(start)

		if err != nil {
			return models.TestResult{
				Domain:       domain,
				Status:       models.StatusError,
				ResponseTime: elapsed,
				Err:          fmt.Errorf("system DNS lookup failed: %w", err),
			}
		}

		if len(ips) == 0 {
			return models.TestResult{
				Domain:       domain,
				Status:       models.StatusBlocked,
				ResponseTime: elapsed,
				Err:          nil,
			}
		}

		for _, ip := range ips {
			if ip.IP.Equal(net.IPv4(0, 0, 0, 0)) {
				return models.TestResult{
					Domain:       domain,
					Status:       models.StatusBlocked,
					ResponseTime: elapsed,
					Err:          nil,
				}
			}
		}

		return models.TestResult{
			Domain:       domain,
			Status:       models.StatusResolved,
			ResponseTime: elapsed,
			Err:          nil,
		}
	}

	in, rtt, err := client.ExchangeContext(ctx, c, dnsServer)

	if err != nil {
		return models.TestResult{
			Domain:       domain,
			Status:       models.StatusError,
			ResponseTime: rtt,
			Err:          fmt.Errorf("DNS query failed: %w", err),
		}
	}

	if in == nil {
		return models.TestResult{
			Domain:       domain,
			Status:       models.StatusError,
			ResponseTime: rtt,
			Err:          fmt.Errorf("received nil DNS response"),
		}
	}

	if in.Rcode != dns.RcodeSuccess || len(in.Answer) == 0 {
		return models.TestResult{
			Domain:       domain,
			Status:       models.StatusBlocked,
			ResponseTime: rtt,
			Err:          nil,
		}
	}

	for _, ans := range in.Answer {
		if a, ok := ans.(*dns.A); ok {
			if a.A.Equal(net.IPv4(0, 0, 0, 0)) {
				return models.TestResult{
					Domain:       domain,
					Status:       models.StatusBlocked,
					ResponseTime: rtt,
					Err:          nil,
				}
			}
		}
	}

	return models.TestResult{
		Domain:       domain,
		Status:       models.StatusResolved,
		ResponseTime: rtt,
		Err:          nil,
	}
}
