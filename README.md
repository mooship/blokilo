[![Go](https://github.com/mooship/blokilo/actions/workflows/go.yml/badge.svg)](https://github.com/mooship/blokilo/actions/workflows/go.yml) [![Gosec](https://github.com/mooship/blokilo/actions/workflows/gosec.yml/badge.svg)](https://github.com/mooship/blokilo/actions/workflows/gosec.yml)

# Blokilo — Ad Block Test TUI 🚫🛡️ 

Blokilo is a terminal-based tool for testing ad-blocking setups (hosts file, DNS filtering, Pi-hole, etc). It checks known ad/tracker domains to verify if they are blocked at the DNS or HTTP level, and presents results in a modern, accessible TUI.

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/T6T31HRCAR)

## ✨ Features

- 🚦 Test if ad/tracker domains are blocked (DNS/HTTP)
- 🕵️ Identify blocking via hosts file or DNS
- 📋 Built-in curated domain list (190+ verified ad/tracker domains)
- 📊 Live progress bar, results table, and summary view
- 🎨 Color-coded, accessible UI (Bubbletea, Bubbles, Lipgloss)
- 🛠️ Custom DNS server support (for Pi-hole, etc)
- ⚡ Parallel/concurrent test engine

## 🛠️ Installation

### Prerequisites 📦

- Go 1.24+
- Internet access for HTTP/DNS tests

### Build from Source 🏗️

```sh
go build -o blokilo ./cmd/blokilo
```

## 🚀 Usage

1. **Run the application:**
   ```sh
   ./blokilo
   ```

2. **Navigate the interface:**
   - Use arrow keys/Enter to select menu options
   - Start Test, Settings, Exit
   - View progress, results, and summary

## 🌐 Domain List

The app uses a built-in curated list of 190+ verified ad/tracker domains covering all major advertising networks. This list is maintained and updated by the developers to ensure optimal testing coverage and accuracy.

The domains are selected to represent:
- Major advertising networks (Google Ads, Facebook, etc.)
- Common tracking services
- Analytics platforms
- Ad servers and CDNs
- Known malware/phishing domains

This curated approach ensures consistent and reliable testing across all installations without requiring external dependencies or manual list management.

## ⚙️ Configuration

- **Custom DNS server:** Enter IP (optionally with :port, default 53) in Settings

## 🧪 Testing

Blokilo includes comprehensive test coverage across all components:

### Running Tests 🏃‍♂️

```sh
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Test Structure 🧩

- **Unit Tests:** Individual component testing (dns, http, models, ui)
- **Coverage:** All major functions and edge cases covered

The test suite ensures reliability and helps maintain code quality as the project evolves.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:

- Code of conduct
- Development setup
- Pull request process
- Coding standards

For bug reports and feature requests, please use the GitHub issue tracker.

## 🗂️ Project Structure

- `cmd/blokilo/main.go` — Entry point
- `internal/models/` — Domain, config, worker, results
- `internal/dns/` — DNS test logic
- `internal/http/` — HTTP test logic
- `internal/ui/` — TUI components (menu, progress, table, summary, settings)

## 📄 License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0). See the LICENSE file for details.

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/T6T31HRCAR)
