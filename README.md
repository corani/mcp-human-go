# 🤖 Ask a Human MCP

A Model Context Protocol (MCP) server and web UI for human-in-the-loop question answering.

- 🧑‍💻 Serve questions to a human operator via a modern web interface
- 🗃️ Store and manage (question, context) → answer pairs in an in-memory database
- 🔗 Integrate with the MCP protocol for use in agent workflows
- 💬 View, filter, and answer questions in a responsive, Bootstrap-powered UI

> [!note]
> Heavily inspired by [Masony817/ask-human-mcp](https://github.com/Masony817/ask-human-mcp) (written in Python) and [KOBA789/human-in-the-loop](https://github.com/KOBA789/human-in-the-loop) (written in Rust).

## ✨ Features

- **🤝 MCP Server**: Implements the Model Context Protocol for agent/human collaboration
- **🌐 Web UI**: Single-page application for listing, filtering, and answering questions
- **⚡ In-Memory Database**: Fast, thread-safe storage of questions and answers
- **🔄 Live Updates**: UI auto-refreshes and shows connection status
- **🖥️ Modern UX**: Responsive design, keyboard shortcuts, and modal dialogs

## 🚀 Getting Started

### Prerequisites

- Go 1.24+

### Running the Server

1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd mcp-human-go
   ```
2. Build and run:
   ```sh
   go run ./cmd/mcp-human-go
   ```
3. Open the web UI:
   - Visit [http://localhost:8990/ui](http://localhost:8990/ui) (or the port configured in your settings)

### ⚙️ Configuration

You can configure the server using environment variables or a `.env` file in the project root (or in your XDG config directory). The following options are supported (see `internal/config/config.go`):

| Variable    | Default | Description                                 |
|-------------|---------|---------------------------------------------|
| `WEB_PORT`  | `8990`  | Port for the web UI/API                     |
| `SSE_PORT`  | `8989`  | Port for the SSE (event) server             |
| `MAX_WAIT`  | `60`    | Max wait time (in seconds) for human answer |

You can create a `.env` file like this:

```
WEB_PORT=8990
SSE_PORT=8989
MAX_WAIT=60
```

Or set variables in your shell before running:

```sh
export WEB_PORT=8990
export MAX_WAIT=60
```

## 🛠️ MCP Client Configuration

You can configure your MCP client using a `mcp.json` file. This allows you to connect to the human server via either SSE (HTTP) or stdio. Example:

```jsonc
{
    "servers": {
        "mcp-human-sse": {
            "url": "http://localhost:8989/mcp"
        },
        "mcp-human-stdio": {
            "type": "stdio",
            "command": "go",
            "args": [
                "run",
                "./cmd/mcp-human-go/"
            ]
        }
    }
}
```

- **mcp-human-sse**: Connects to the server using HTTP/SSE at the given URL.
- **mcp-human-stdio**: Launches the server as a subprocess using stdio (useful for local development or integration).

You can select which server to use in your MCP client configuration or tooling.

### 📝 Usage

- ✍️ New questions are added to the in-memory database and appear in the web UI.
- 👁️ Click a question to view details and answer it in a modal dialog.
- 🕵️ Use the filter to hide answered questions.
- 🚦 The UI shows when the backend is unreachable and resumes when reconnected.
- ⌨️ Use Ctrl+Enter to submit answers quickly.

## 📁 Project Structure

- `cmd/mcp-human-go/` — Main entry point
- `internal/memory/` — In-memory database logic
- `internal/web/` — Web server and UI
- `internal/human/` — Human interaction logic
- `internal/tools/` — MCP tool registration

## 🙏 Main dependencies

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) — Model Context Protocol Go implementation
- [gofiber/fiber](https://github.com/gofiber/fiber) — Web framework for Go

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.