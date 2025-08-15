# GEMINI.md

## 1. Context Usage

* Gemini **must always reference Context7 MCP** for:

  * All decisions
  * Code completions
  * Dependency suggestions
  * UX/UI choices
* Always **verify the latest state of Context7** before acting.

---

## 2. Formatting & Linting

* Follow **`gofmt` rules strictly**.
* Use **idiomatic Go style**.
* Fix linting issues **preemptively**; no warnings or errors allowed.
* Example:

  ```go
  // Good
  type User struct {
      ID   int
      Name string
  }
  ```

---

## 3. Build & Run Restrictions

* **Never build or run** the app in code generation.
* Gemini only generates or edits **source code**.

---

## 4. Dependencies

* Check `go.mod` for current dependencies **before suggesting new ones**.
* Suggest **new dependencies only if they**:

  * Are popular and well-maintained.
  * Improve type safety, maintainability, or developer experience.
  * Have minimal or zero impact on binary size/performance.
* Always justify new dependencies with **Context7 MCP**.

---

## 5. Type Safety & File Structure

* Strictly use **typed structs and interfaces**.
* Avoid `interface{}` unless absolutely necessary. Prefer concrete types.
* Use `context.Context` for all network calls or cancellable operations.
* Use **dedicated `.go` files per domain**:

  * `dns_test.go` — DNS tests
  * `http_test.go` — HTTP tests
  * `ui.go` — Terminal UI logic
  * `models.go` — Data models

---

## 6. Accessibility & UX

* Follow **Bubbletea** and **Lipgloss** best practices for:

  * Terminal accessibility
  * Keyboard navigation
  * Focus management
* Use **semantic, reusable components** from Bubbles.
* Ensure **color contrast is sufficient** for readability.
* Confirm all UX decisions with **Context7 MCP**.

---

## 7. Security

* **Never hardcode secrets** or sensitive information.
* Use **config/environment variables** for credentials, tokens, and secrets.
* Validate all external input strictly.

---

## 8. Error Handling

* Return errors explicitly; never ignore them.
* Wrap errors with context:

  ```go
  if err != nil {
      return fmt.Errorf("failed to load config: %w", err)
  }
  ```
* Validate inputs before processing.
* Provide **actionable and clear error messages**.

---

## 9. Confirmation for Large Edits

* Gemini must **prompt the user** before:

  * Major refactors
  * Large code generation
* Always document what will change and why **Context7 MCP** requires it.
