package plugins

import (
	"context"
	"encoding/json"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"onah-go/database"
)

// Context wraps the environment of a command execution
type Context struct {
	Client    *whatsmeow.Client
	Event     *events.Message
	SenderJID string
	ChatJID   string
	PushName  string
	Command   string
	Args      []string
	Query     string
	User      *database.User
}

// Reply sends a text message replying to the user's incoming message
func (ctx *Context) Reply(text string) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(ctx.Event.Info.ID),
				Participant:   proto.String(ctx.Event.Info.Sender.ToNonAD().String()),
				QuotedMessage: ctx.Event.Message,
			},
		},
	}
	_, err := ctx.Client.SendMessage(context.Background(), ctx.Event.Info.Chat, msg)
	return err
}

// React sends an emoji reaction to the message
func (ctx *Context) React(emoji string) error {
	msg := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(ctx.Event.Info.Chat.String()),
				FromMe:    proto.Bool(ctx.Event.Info.IsFromMe),
				ID:        proto.String(ctx.Event.Info.ID),
			},
			Text: proto.String(emoji),
		},
	}
	_, err := ctx.Client.SendMessage(context.Background(), ctx.Event.Info.Chat, msg)
	return err
}

// Command defines a bot command and its rules
type Command struct {
	Name      string
	Tags      []string
	Help      string
	Limit     bool
	Premium   bool
	OwnerOnly bool
	Execute   func(ctx *Context) error
}

var commands []Command

// Register adds a command to the global command registry
func Register(cmd Command) {
	commands = append(commands, cmd)
}

// GetCommands retrieves all registered commands
func GetCommands() []Command {
	return commands
}

// Button defines a structure for WhatsApp Interactive Message buttons
type Button struct {
	Type string // "reply" (quick_reply) or "url" (cta_url)
	Text string
	ID   string // command keywords for "reply", or URL link for "url"
}

// SendButtons sends an interactive message containing action buttons
func (ctx *Context) SendButtons(bodyText string, footerText string, buttons []Button) error {
	var flowButtons []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton

	for _, btn := range buttons {
		var name string
		var params map[string]interface{}

		if btn.Type == "url" {
			name = "cta_url"
			params = map[string]interface{}{
				"display_text": btn.Text,
				"url":          btn.ID,
				"merchant_url": btn.ID,
			}
		} else {
			name = "quick_reply"
			params = map[string]interface{}{
				"display_text": btn.Text,
				"id":           btn.ID,
			}
		}

		bpBytes, err := json.Marshal(params)
		if err != nil {
			return err
		}
		bpJSON := string(bpBytes)

		flowButtons = append(flowButtons, &waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
			Name:             proto.String(name),
			ButtonParamsJSON: proto.String(bpJSON),
		})
	}

	msgVersion := int32(1)
	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String(bodyText),
			},
			Footer: &waE2E.InteractiveMessage_Footer{
				Text: proto.String(footerText),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons:        flowButtons,
					MessageVersion: &msgVersion,
				},
			},
		},
	}

	bizNode := binary.Node{
		Tag: "biz",
		Content: []binary.Node{{
			Tag: "interactive",
			Attrs: binary.Attrs{"type": "native_flow", "v": "1"},
			Content: []binary.Node{{
				Tag: "native_flow",
				Attrs: binary.Attrs{"v": "9", "name": "mixed"},
			}},
		}},
	}

	_, err := ctx.Client.SendMessage(context.Background(), ctx.Event.Info.Chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &[]binary.Node{bizNode},
	})
	return err
}

// ListRow defines a row in a WhatsApp Single Select List
type ListRow struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

// ListSection defines a section containing rows in a List
type ListSection struct {
	Title string    `json:"title"`
	Rows  []ListRow `json:"rows"`
}

// SendList sends a Single Select List Message
func (ctx *Context) SendList(bodyText, footerText, buttonText string, sections []ListSection) error {
	params := map[string]interface{}{
		"title":    buttonText,
		"sections": sections,
	}
	bpBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	bpJSON := string(bpBytes)

	buttons := []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
		{
			Name:             proto.String("single_select"),
			ButtonParamsJSON: proto.String(bpJSON),
		},
	}

	msgVersion := int32(1)
	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String(bodyText),
			},
			Footer: &waE2E.InteractiveMessage_Footer{
				Text: proto.String(footerText),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons:        buttons,
					MessageVersion: &msgVersion,
				},
			},
		},
	}

	bizNode := binary.Node{
		Tag: "biz",
		Content: []binary.Node{{
			Tag: "interactive",
			Attrs: binary.Attrs{"type": "native_flow", "v": "1"},
			Content: []binary.Node{{
				Tag: "native_flow",
				Attrs: binary.Attrs{"v": "9", "name": "mixed"},
			}},
		}},
	}

	_, err = ctx.Client.SendMessage(context.Background(), ctx.Event.Info.Chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &[]binary.Node{bizNode},
	})
	return err
}
