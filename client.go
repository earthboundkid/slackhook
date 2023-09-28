package slackhook

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/errorx"
	"github.com/carlmjohnson/requests"
)

// Logger is any function that behaves like slog.InfoContext, slog.DebugContext, etc.
type Logger = func(ctx context.Context, msg string, args ...any)

// NoOpLogger does nothing with log messages.
func NoOpLogger(ctx context.Context, msg string, args ...any) {}

// Magic URL to use a mock client in dev.
const MockClient = "slack://mock"

// Client posts messages to a Slack webhook URL.
type Client struct {
	hookURL string
}

// Set the hookURL for a Client.
// If hookURL is [MockClient],
// Client will use a mock client.
// Set is part of the flag.Value interface.
func (c *Client) Set(hookURL string) error {
	_, err := url.Parse(hookURL)
	if err != nil {
		return err
	}
	c.hookURL = hookURL
	return nil
}

// String is part of the flag.Value interface.
func (c *Client) String() string {
	return ""
}

// Post message to Slack.
// Logger must not be nil.
// Uses http.DefaultClient if cl is nil.
// Returns an error if response is not 200 OK.
func (sc *Client) Post(ctx context.Context, l Logger, cl *http.Client, msg Message) (err error) {
	defer errorx.Trace(&err)

	isMock := sc.hookURL == MockClient
	if isMock {
		cl = &http.Client{
			Transport: requests.ReplayString("HTTP/1.1 200 OK\r\n\r\n"),
		}
	}

	if isMock {
		b, _ := json.Marshal(&msg)
		l(ctx, "slackhook: PostCtx", "mock-client", isMock, "output", b)
	} else {
		l(ctx, "slackhook: PostCtx", "mock-client", isMock)
	}
	return requests.
		URL(sc.hookURL).
		Client(cl).
		BodyJSON(msg).
		Fetch(ctx)
}

// Message is the JSON object expected by Slack
type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Pretext   string  `json:"pretext,omitempty"`
	Fallback  string  `json:"fallback"`
	Color     string  `json:"color,omitempty"`
	Title     string  `json:"title,omitempty"`
	TitleLink string  `json:"title_link,omitempty"`
	Text      string  `json:"text,omitempty"`
	TimeStamp int64   `json:"ts,omitempty"`
	Fields    []Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
