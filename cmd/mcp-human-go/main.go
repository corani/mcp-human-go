package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/corani/mcp-human-go/internal/config"
	"github.com/corani/mcp-human-go/internal/human"
	"github.com/corani/mcp-human-go/internal/memory"
	"github.com/corani/mcp-human-go/internal/tools"
	"github.com/corani/mcp-human-go/internal/web"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed system-prompt.txt
var INSTRUCTIONS string

func main() {
	logfile, err := os.OpenFile("mcpserver.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()

	out := io.MultiWriter(logfile, os.Stdout)

	handler := slog.NewTextHandler(out, &slog.HandlerOptions{})
	logger := slog.New(handler)

	conf := config.MustLoad(logger)
	_ = conf

	hooks := new(server.Hooks)

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		logger.Error("Error in MCP method",
			slog.String("method", string(method)),
			slog.Any("id", id),
			slog.Any("message", message),
			slog.String("error", err.Error()),
		)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		logger.Info("Success in MCP method",
			slog.String("method", string(method)),
			slog.Any("id", id),
			slog.Any("message", message),
		)
	})

	instructions := INSTRUCTIONS +
		fmt.Sprintf("\n\nThe current date is: %v", time.Now().Format("2006-01-02"))

	srv := server.NewMCPServer(
		"mcp-human-go", "1.0.0",
		server.WithLogging(),
		server.WithInstructions(instructions),
		server.WithRecovery(),
		server.WithHooks(hooks),
	)

	mem := memory.NewMemoryDB()
	api := web.NewAPI(conf, mem)
	ask := human.NewAsk(conf, mem)

	tools.Register(srv, ask)

	// TODO(daniel): probably shouldn't use a lambda here, and we should check the request params.
	srv.AddPrompt(mcp.NewPrompt("instructions"),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return mcp.NewGetPromptResult("instructions", []mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(instructions)),
			}), nil
		})

	// TODO(daniel): probably shouldn't use a lambda here, and we should check the request params.
	srv.AddResource(mcp.NewResource("file:///mcpserver.log", "server log"),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			logfile.Sync()

			// TODO(daniel): reading a file that's open for writing is a bad idea, but this is just a demo.
			bs, err := os.ReadFile(logfile.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to read log file: %w", err)
			}

			contents := mcp.TextResourceContents{
				Text:     string(bs),
				URI:      request.Params.URI,
				MIMEType: "text/plain",
			}

			return []mcp.ResourceContents{contents}, nil
		})

	go func() {
		logger.Info("Starting web server",
			slog.String("address", fmt.Sprintf("http://localhost:%d", conf.WebPort)))
		if err := api.Start(); err != nil {
			logger.Error("Failed to start web server",
				slog.String("error", err.Error()),
			)
		}
		logger.Info("Web server stopped")
	}()
	defer api.Shutdown()

	sse := server.NewSSEServer(srv,
		server.WithSSEEndpoint("/mcp"),
	)
	defer sse.Shutdown(context.Background())

	go func() {
		logger.Info("Starting SSE server",
			slog.String("address", fmt.Sprintf("http://localhost:%d/mcp", conf.SsePort)))

		if err := sse.Start(fmt.Sprintf("0.0.0.0:%d", conf.SsePort)); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start SSE server",
				slog.String("error", err.Error()),
			)
		}

		logger.Info("SSE server stopped")
	}()

	if err := server.ServeStdio(srv); !errors.Is(err, context.Canceled) {
		logger.Error("Failed to serve MCP server",
			slog.String("error", err.Error()),
		)
	}

	logger.Info("MCP server stopped")
}
