# jchTerm

jchTerm is a custom terminal application built in Go, using the Charm libraries for the text-based user interface (TUI). It provides a familiar terminal experience while offering built-in commands and the potential for easy extension.

## Features

- Command-line interface similar to traditional terminals
- Built-in commands for enhanced functionality
- Command history navigation
- Integration with system shell for executing standard commands
- Custom command processing
- Logging for debugging and auditing
- API integration (currently with Anthropic's Claude AI)

## Built-in Commands

- `clear`: Clears the terminal screen
- `rebuild`: Rebuilds the application
- `edit`: Opens the main.go file in the default editor
- `ask`: Sends a query to Claude AI (using Anthropic's API)

## Project Structure

- `main.go`: Initializes the application, sets up logging, and runs the TUI
- `model.go`: Implements the terminal user interface, handles user input and command execution
- `commands.go`: Processes both built-in and shell commands
- `builtin.go`: Implements built-in commands
- `claude.go`: Handles integration with Claude AI
- `config.go`: Stores application-wide constants
- `logger.go`: Provides logging functionality

## Building and Running

To build and run jchTerm:

1. Ensure you have Go installed on your system.
2. Clone the repository:
   ```
   git clone https://github.com/yourusername/jchterm.git
   cd jchterm
   ```
3. Build the application:
   ```
   go build -o jchterm cmd/jchterm/main.go
   ```
4. Run jchTerm:
   ```
   ./jchterm
   ```

## Configuration

Update the `config.go` file to modify application-wide settings such as window dimensions and API keys.

## Contributing

We welcome contributions to jchTerm! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to get started.

## License

[Insert your chosen license here]

## Future Development

jchTerm is designed with extensibility in mind. We plan to continue adding new features and improving existing ones through a "dog fooding" approach, using jchTerm to build itself.

## Contact

jchackert@me.com
