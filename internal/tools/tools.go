package tools

import (
	"context"
	"fmt"

	"github.com/corani/mcp-human-go/internal/human"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Register(srv *server.MCPServer, ask *human.Ask) {
	tools := []Tool{
		newAskHuman(ask),
	}

	for _, tool := range tools {
		srv.AddTool(tool.Schema(), tool.Handler)
	}
}

type Tool interface {
	Schema() mcp.Tool
	Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type askHumanTool struct {
	human *human.Ask
}

func newAskHuman(ask *human.Ask) Tool {
	return &askHumanTool{ask}
}

func (c *askHumanTool) Schema() mcp.Tool {
	return mcp.NewTool("ask_human",
		mcp.WithDescription("Ask a human a question only if neither you nor the user can answer it, such as project-specific context, local environment details, or non-public information. Do NOT use this tool to clarify, disambiguate, or follow up on the user's intent or question. Always ask the user directly if you need clarification. Use this tool only for facts unavailable to both you and the user."),
		mcp.WithString("question",
			mcp.Description("The question to ask the human. Use markdown formatting for clarity. Do not use this for clarifications or follow-ups; ask the user directly instead."),
			mcp.Required()),
		mcp.WithString("context",
			mcp.Description("Context for the question, if any."),
			mcp.Required()),
	)
}

func (c *askHumanTool) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	question := request.GetString("question", "")
	if question == "" {
		return toError(fmt.Errorf("`question` is required"))
	}

	// optional context
	context := request.GetString("context", "")

	answer, err := c.human.Ask(question, context)
	if err != nil {
		return toError(err)
	}

	return mcp.NewToolResultText(answer), nil
}

func toError(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}
