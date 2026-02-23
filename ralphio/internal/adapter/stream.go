package adapter

import (
	"encoding/json"
	"strings"
)

// streamMsg is the top-level structure for a single NDJSON line emitted by
// an AI agent. Fields are union-typed across agent formats; only relevant
// fields are populated for any given message type.
type streamMsg struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype,omitempty"`

	// Claude / Cursor: type=="assistant"
	Message *assistantMessage `json:"message,omitempty"`

	// Claude result: type=="result", subtype=="success"
	Result string `json:"result,omitempty"`

	// opencode / kilo: type=="text"
	Part *partContent `json:"part,omitempty"`

	// opencode / kilo streaming delta: type=="message_update"
	AssistantMessageEvent *assistantEvent `json:"assistantMessageEvent,omitempty"`
}

type assistantMessage struct {
	Content []contentBlock `json:"content"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type partContent struct {
	Text string `json:"text"`
}

type assistantEvent struct {
	Type  string `json:"type"`
	Delta string `json:"delta,omitempty"`
}

// ParseStreamLine parses one line of NDJSON agent output and returns the
// displayable text it carries, or "" if the line carries no text. Non-JSON
// lines are returned as-is. Ported from ralph/src/lib/agent-stream.ts.
func ParseStreamLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	var msg streamMsg
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		// Not JSON — return raw line (plain-text agent output).
		return line
	}

	switch msg.Type {
	case "assistant":
		// Claude / Cursor format: extract text from content blocks.
		if msg.Message == nil {
			return ""
		}
		var parts []string
		for _, block := range msg.Message.Content {
			if block.Type == "text" && block.Text != "" {
				parts = append(parts, block.Text)
			}
		}
		return strings.Join(parts, "")

	case "result":
		// Final result message (Claude / Cursor).
		if msg.Subtype == "success" {
			return msg.Result
		}
		return ""

	case "text":
		// opencode / kilo per-part text.
		if msg.Part != nil {
			return msg.Part.Text
		}
		return ""

	case "message_update":
		// opencode / kilo streaming delta.
		if msg.AssistantMessageEvent != nil &&
			msg.AssistantMessageEvent.Type == "text_delta" {
			return msg.AssistantMessageEvent.Delta
		}
		return ""

	case "step_finish":
		// Lifecycle marker — no displayable text.
		return ""

	default:
		// Unknown message type — return raw line for visibility.
		return line
	}
}
