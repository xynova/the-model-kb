package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/xynova/roster-pick/internal/config"
	derrors "github.com/xynova/roster-pick/internal/errors"
	"github.com/xynova/roster-pick/internal/observability"
	"github.com/xynova/roster-pick/internal/pick"
)

// Picker is the pick service surface for MCP and tests.
type Picker interface {
	Pick(ctx context.Context, req pick.Request) (pick.Result, error)
}

// ServerConfig configures the MCP stdio server.
type ServerConfig struct {
	ScoresPath string
	Version    string
}

// ServeStdio runs the roster_pick MCP tool over stdio.
func ServeStdio(cfg ServerConfig) error {
	pickerCfg := config.Config{ScoresPath: cfg.ScoresPath, Limit: 5}
	if pickerCfg.ScoresPath == "" {
		c, err := config.CreateConfig()
		if err != nil {
			return err
		}
		pickerCfg = c
	}
	p, err := pick.CreatePicker(pickerCfg, observability.NopLogger{})
	if err != nil {
		return err
	}
	return serveWithPicker(p, cfg.Version)
}

func serveWithPicker(p Picker, version string) error {
	if version == "" {
		version = "dev"
	}
	s := mcpserver.NewMCPServer("roster-pick", version)
	tool := mcp.NewTool("roster_pick",
		mcp.WithDescription("Rank Model Picks host offers for a natural-language task (fit, cost, host limits). Do not grep roster markdown."),
		mcp.WithString("task", mcp.Required(), mcp.Description("Natural-language task")),
	)
	s.AddTool(tool, rosterPickHandler(p))
	return mcpserver.ServeStdio(s)
}

// RosterPickHandlerForTest exposes the tool handler for unit tests.
func RosterPickHandlerForTest(p Picker) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return rosterPickHandler(p)
}

func rosterPickHandler(p Picker) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		task, err := req.RequireString("task")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("code=%s %v", derrors.CodeInvalidArgument, err)), nil
		}
		res, err := p.Pick(ctx, pick.Request{Task: task})
		if err != nil {
			if derrors.IsCode(err, derrors.CodeNotFound) {
				body, encErr := json.Marshal(res)
				if encErr != nil {
					return nil, encErr
				}
				return mcp.NewToolResultText(string(body)), nil
			}
			code := derrors.CodeFailed
			var de *derrors.Error
			if errors.As(err, &de) {
				code = de.Code
			}
			return mcp.NewToolResultError(fmt.Sprintf("code=%s %s", code, err.Error())), nil
		}
		body, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(string(body)), nil
	}
}
