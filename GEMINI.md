# Gemini Instructions

## Project Overview

This is a terminal user interface (TUI) application built entirely in Go using the Charm ecosystem. The project focuses on creating interactive command-line interfaces with modern styling and component-based architecture.

## Technology Stack

- **Language**: Go
- **TUI Framework**: Bubbletea (BubbleTea)
- **UI Components**: Bubbles
- **Styling**: Lipgloss
- **Architecture**: Model-View-Update (MVU) pattern

## Code Style and Conventions

### Go Standards

- Use standard Go formatting with go fmt and go vet
- Follow effective Go naming conventions (camelCase for private, PascalCase for public)
- Use Go modules for dependency management
- Implement proper error handling with wrapped errors
- Use context.Context for cancellation and timeouts where appropriate
- Keep functions focused and single-purpose

### Bubbletea Architecture

- Follow the Model-View-Update (MVU) pattern strictly
- Implement tea.Model interface for all components
- Use tea.Cmd for asynchronous operations and side effects
- Handle all user input through the Update method
- Render UI completely in the View method
- Keep models immutable - return new instances from Update
- Use tea.Program for application lifecycle management

### Model Design

- Structure models with clear state representation
- Embed sub-models for complex components
- Use enums or constants for application states
- Implement proper state transitions
- Keep model fields exported only when necessary for embedding
- Use composition over inheritance for complex UIs

### Message Handling

- Define custom message types for application events
- Use tea.Msg interface for all messages
- Implement message routing patterns for complex applications
- Use tea.Batch for multiple commands
- Handle tea.KeyMsg and tea.WindowSizeMsg appropriately
- Create domain-specific messages for business logic

### Bubbles Components

- Leverage existing Bubbles components before building custom ones
- Customize Bubbles components through their configuration options
- Compose multiple Bubbles components in larger views
- Handle focus management between multiple input components
- Use appropriate Bubbles components: textinput, textarea, list, table, progress, spinner, etc.
- Implement proper key bindings for component interaction

### Lipgloss Styling

- Use Lipgloss for all styling and layout
- Create reusable style definitions as package-level variables
- Use adaptive colors that work in different terminal environments
- Implement consistent spacing and alignment
- Use Lipgloss layout functions for complex positioning
- Create style themes for consistent visual design
- Use borders, padding, and margins effectively

## Application Structure
- Organize models in separate files by feature or component
- Create a styles package for Lipgloss definitions
- Separate message types into dedicated files
- Use commands package for tea.Cmd implementations
- Implement proper initialization in Init methods
- Structure views with clear hierarchy and composition

## State Management

- Keep application state centralized in the main model
- Use sub-models for component-specific state
- Implement proper state synchronization between models
- Use channels for communication with goroutines
- Handle application lifecycle events properly
- Implement proper cleanup in Quit scenarios

## Input Handling

- Handle all keyboard input through tea.KeyMsg
- Implement proper key binding patterns
- Use key sequences for complex operations
- Handle special keys (Ctrl+C, Esc, Enter) appropriately
- Implement context-sensitive key bindings
- Provide clear feedback for user actions

## Asynchronous Operations

- Use tea.Cmd for all asynchronous operations
- Implement proper error handling in commands
- Use context for cancellable operations
- Handle network requests and file I/O asynchronously
- Implement loading states and progress indicators
- Use tea.Tick for periodic updates

## Layout and Responsive Design

- Use Lipgloss layout functions for responsive designs
- Handle terminal resizing with tea.WindowSizeMsg
- Implement proper text wrapping and truncation
- Use flexible layouts that adapt to terminal size
- Design for various terminal widths and heights
- Implement scrolling for content that exceeds viewport

## Error Handling

- Display user-friendly error messages in the TUI
- Log technical errors appropriately
- Implement graceful degradation for non-critical errors
- Use proper error recovery patterns
- Handle terminal-specific errors gracefully
- Provide clear error states in the UI

## Performance Considerations

- Minimize unnecessary re-renders in View methods
- Use efficient string building techniques
- Cache expensive computations where appropriate
- Implement lazy loading for large datasets
- Optimize Lipgloss style calculations
- Handle large lists and tables efficiently

## Testing Patterns

- Write unit tests for model Update methods
- Test message handling logic independently
- Mock external dependencies in commands
- Test view rendering with different model states
- Implement integration tests for complete workflows
- Use table-driven tests for message handling

## Accessibility Considerations

- Ensure keyboard navigation works for all functionality
- Use appropriate contrast ratios in color schemes
- Provide alternative text representations where needed
- Support screen readers through proper text output
- Implement clear focus indicators
- Design for users with different terminal capabilities

## Development Workflow

- Use go run for development and testing
- Implement proper logging for debugging
- Use build tags for development vs production features
- Test in different terminal environments
- Use proper signal handling for graceful shutdown
- Implement hot reloading where beneficial

## Code Generation Preferences

- Always provide complete, production-ready code without examples or samples
- Generate the final implementation directly, not demonstrations or tutorials
- Skip explanatory examples and provide only the working solution
- Deliver finished code that can be used immediately in the project
- Avoid showing "how to" examples - just implement the requested functionality
- Provide the actual code needed, not instructional or educational samples
- Focus on delivering final implementations rather than teaching patterns
- Include proper error handling and edge cases in all implementations
- Ensure code follows Go best practices and idiomatic patterns
- Generate complete model implementations with Init, Update, and View methods