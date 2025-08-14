# 🤝 Contributing to Blokilo

Thank you for your interest in contributing to Blokilo! This guide will help you understand our development process and standards.

## 🛠️ Development Environment Setup

### Prerequisites 📦

- **Go 1.25** - Required for building and testing
- **Git** - For version control
- **Terminal** - For testing the TUI interface

### Getting Started 🚀

1. **Fork and Clone:**
   ```sh
   git clone https://github.com/your-username/blokilo.git
   cd blokilo
   ```

2. **Install Dependencies:**
   ```sh
   go mod download
   ```

3. **Verify Setup:**
   ```sh
   go build -o blokilo ./cmd/blokilo
   ./blokilo
   ```

4. **Run Tests:**
   ```sh
   go test ./...
   ```

## 📝 Coding Standards

### Go Formatting and Linting 🧹

- **All code must pass `gofmt`** - Use `go fmt ./...`
- **All code must pass `golangci-lint`** - No warnings allowed
- **Run linting before committing:**
  ```sh
  golangci-lint run
  ```

### Code Style Guidelines 🎨

- **Type Safety:** Use typed structs and interfaces, avoid `interface{}` except when absolutely necessary
- **Context Usage:** Use `context.Context` for functions involving network calls or cancellations
- **Error Handling:** 
  - Return errors explicitly
  - Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
  - Provide actionable error messages
- **File Organization:** Use dedicated `.go` files per domain (e.g., `dns.go`, `http.go`, `ui.go`)

### Security Requirements 🔒

- **Never hardcode secrets** - Use config/environment variables
- **Validate all inputs** before processing
- **Follow secure coding practices**

## 🧪 Testing Requirements

### Test Coverage 📈

- **Unit tests required** for all new functions
- **Integration tests** for end-to-end functionality
- **Aim for high coverage** - run `go test -cover ./...`

### Running Tests 🏃‍♂️

```sh
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/dns/
```

### Test Structure 🧩

- Place tests in `*_test.go` files alongside source code
- Use table-driven tests where appropriate
- Test both success and error cases
- Mock external dependencies

## 🖥️ TUI Development Guidelines

### Accessibility Requirements ♿

- **Keyboard Navigation:** All UI must be navigable with arrow keys/Enter
- **Focus Handling:** Clear visual focus indicators
- **Color Contrast:** Ensure sufficient contrast for readability
- **Screen Reader Friendly:** Use semantic UI components

### UI Framework Standards 🧩

- **Bubbletea:** Follow Bubbletea patterns for components
- **Bubbles:** Use reusable components from Bubbles library
- **Lipgloss:** Consistent styling and theming

## 🏗️ Architecture Overview

### Project Structure 🗂️

```
cmd/blokilo/          # Application entry point
internal/models/      # Data structures and business logic
internal/dns/         # DNS testing functionality
internal/http/        # HTTP testing functionality
internal/ui/          # TUI components and views
```

### Component Responsibilities 🧩

- **Models:** Domain objects, configuration, results
- **DNS/HTTP:** Core testing logic, network operations
- **UI:** User interface, navigation, display logic

## 🔄 Contribution Workflow

### 1. Create an Issue First 📝

- Describe your proposed change, bug fix, or feature
- Wait for feedback before starting work
- Reference existing issues when possible

### 2. Development Process 🛠️

1. **Fork the repository** to your GitHub account
2. **Create a feature branch:** `git checkout -b feature/your-feature`
3. **Make your changes** following the coding standards
4. **Add/update tests** for your changes
5. **Run the full test suite:** `go test ./...`
6. **Run linting:** `golangci-lint run`
7. **Test the TUI manually** to ensure functionality

### 3. Commit Guidelines 🗒️

- **Write clear commit messages**
- **Use conventional commits format:**
  ```
  type(scope): description
  
  feat(dns): add custom DNS server support
  fix(ui): resolve navigation issue in settings
  test(http): add timeout test cases
  ```

### 4. Pull Request Process 🔀

- **Reference the related issue:** "Closes #123"
- **Describe your changes** in detail
- **Include screenshots** for UI changes
- **Ensure all checks pass** (tests, linting)
- **Be responsive to feedback**

## 📦 Dependencies Policy

### Adding New Dependencies ➕

New dependencies should:
- Be well-maintained and popular
- Improve type safety, maintainability, or developer experience
- Have minimal impact on binary size and performance
- Be discussed in an issue before addition

### Updating Dependencies ♻️

- Keep dependencies up to date
- Test thoroughly after updates
- Document breaking changes

## 🤗 Code of Conduct

Please be respectful and considerate in all interactions. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for full details.

## 🆘 Getting Help

- **Questions:** Open an issue with the "question" label
- **Bug Reports:** Use the bug report template
- **Feature Requests:** Use the feature request template
- **Security Issues:** Email maintainers privately

## 🌟 Recognition

Contributors will be recognized in:
- Project README
- Release notes for significant contributions
- GitHub contributors page

---

Thank you for helping make Blokilo better! Your contributions help improve ad-blocking testing for everyone.
