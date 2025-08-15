# Geminii Instructions

1. **Context Usage**

Gemini **must always reference context7** for decisions and code completions.

2. **Formatting & Linting**

Follow **`gofmt` rules strictly. No lint errors or warnings allowed. Use idiomatic Go formatting and style conventions. Fix issues preemptively.

3. **Build/Run Restrictions**

Never build or run the app or trigger any compile/run processes. Only generate or edit source code.

4. **Dependencies**

Check `go.mod` for current dependencies. Suggest new dependencies **only if they:**

* Are well-maintained and popular.
* Improve type safety, maintainability, or developer experience.
* Have zero or minimal impact on binary size or performance.

5. **Type Safety & Files**

* Strictly use **typed structs and interfaces**.
* Use **dedicated `.go` files per domain of responsibility**, e.g. `dns_test.go`, `http_test.go`, `ui.go`, `models.go`.
* Avoid use of `interface{}` except when absolutely necessary; prefer concrete types.
* Use `context.Context` for all functions that involve network calls or cancellations.

6. **Accessibility & UX**

* Follow Bubbletea and Lipgloss best practices for terminal accessibility (keyboard navigation, focus handling).
* Use semantic and reusable UI components from Bubbles.
* Ensure color choices have sufficient contrast.

7. **Security**

Never hardcode secrets or sensitive information. Use config/environment variables for anything sensitive.

8. **Error Handling**

* Return errors explicitly.
* Wrap errors with context (`fmt.Errorf("...: %w", err)`).
* Validate inputs strictly before processing.
* Provide actionable error messages.

9. **Confirmation for Large Edits**

Before any major refactor or large code generation, prompt the user for confirmation.
