package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/xynova/roster-pick/internal/pick"
	"github.com/xynova/roster-pick/internal/scores"
)

type fakePicker struct{}

func (fakePicker) Pick(ctx context.Context, req pick.Request) (pick.Result, error) {
	return pick.Result{
		MatchedJobs:     []string{"cheap", "vision"},
		MatchedCapacity: "industrial",
		Offers: []scores.Offer{
			{Name: "gemma-4-26b-a4b", HostID: "deepinfra"},
		},
	}, nil
}

func TestMCPHandlerRoundTrip(t *testing.T) {
	t.Parallel()
	h := RosterPickHandlerForTest(fakePicker{})
	req := mcp.CallToolRequest{}
	req.Params.Name = "roster_pick"
	req.Params.Arguments = map[string]any{"task": "industrial cheap vision"}
	out, err := h(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if out.IsError {
		t.Fatalf("tool error: %v", out.Content)
	}
	text := out.Content[0].(mcp.TextContent).Text
	var body pick.Result
	if err := json.Unmarshal([]byte(text), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.MatchedJobs) == 0 || len(body.Offers) == 0 {
		t.Fatalf("body=%+v", body)
	}
}
