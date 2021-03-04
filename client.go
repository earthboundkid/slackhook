package slackhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

// Client posts messages to a Slack webhook
type Client struct {
	hookURL string
	c       *http.Client
}

// New returns nil if hookURL is blank. Uses http.DefaultClient if c is nil.
func New(hookURL string, c *http.Client) *Client {
	if hookURL == "" {
		return nil
	}
	if c == nil {
		c = http.DefaultClient
	}
	return &Client{hookURL, c}
}

// Set implements flag.Value
func (sc *Client) Set(s string) error {
	if s == "" {
		return nil
	}
	_, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("bad Slack incoming webbook hook URL: %v", err)
	}
	sc.hookURL = s
	return nil
}

// String implements flag.Value
func (sc *Client) String() string {
	return ""
}

// Post message to Slack. Noop if client is nil.
// Returns StatusErr if response is not 200 OK.
func (sc *Client) Post(msg Message) error {
	ctx := context.Background()
	return sc.PostCtx(ctx, msg)
}

// PostCtx posts message to Slack with context.
// Noop if client is nil.
// Returns StatusErr if response is not 200 OK.
func (sc *Client) PostCtx(ctx context.Context, msg Message) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("problem sending message to Slack: %w", err)
		}
	}()
	if sc == nil || sc.hookURL == "" {
		return nil
	}
	blob, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	r := bytes.NewReader(blob)
	rsp, err := ctxhttp.Post(ctx, sc.c, sc.hookURL, "application/json", r)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// Drain connection
	// Slack shouldn't be giving us anything to discard,
	// but if they do, read some of it to try to reuse connections
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, rsp.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}

	if rsp.StatusCode != http.StatusOK {
		return StatusErr(rsp.StatusCode)
	}

	return err
}

// StatusErr is an unexpected status
type StatusErr int

func (se StatusErr) String() string {
	return http.StatusText(int(se))
}

func (se StatusErr) Error() string {
	return fmt.Sprintf("unexpected status: %d %s",
		int(se), se.String())
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
