package slackhook

import (
	"context"
	"net/http"

	"github.com/carlmjohnson/errutil"
	"github.com/carlmjohnson/requests"
)

// Client posts messages to a Slack webhook
type Client struct {
	rb *requests.Builder
}

// New returns mock client if hookURL and c are blank. Uses http.DefaultClient if c is nil.
func New(hookURL string, c *http.Client) *Client {
	if hookURL == "" && c == nil {
		hookURL = "protocol://nosuch.example"
		c = &http.Client{
			Transport: requests.ReplayString("HTTP/1.1 200 OK\r\n\r\n"),
		}
	}
	return &Client{
		requests.
			URL(hookURL).
			Client(c).
			CheckStatus(http.StatusOK),
	}
}

// PostCtx posts message to Slack with context.
// Noop if client is nil.
// Returns an error if response is not 200 OK.
func (sc *Client) PostCtx(ctx context.Context, msg Message) (err error) {
	defer errutil.Trace(&err)
	return sc.rb.Clone().
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
