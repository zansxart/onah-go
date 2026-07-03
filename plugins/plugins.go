package plugins

import (
	"context"

	"go.mau.fi/whatsmeow"
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
