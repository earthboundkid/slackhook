package slackhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

// Posts message to Slack. Noop if client is nil.
// Returns StatusErr if response is not 200 OK.
func (sc *Client) Post(msg Message) error {
	if sc == nil {
		return nil
	}
	blob, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	r := bytes.NewReader(blob)
	rsp, err := sc.c.Post(sc.hookURL, "application/json", r)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return StatusErr{rsp.StatusCode}
	}
	return nil
}

// StatusErr is an unexpected status
type StatusErr struct {
	Code int
}

func (se StatusErr) Error() string {
	return fmt.Sprintf("unexpected status: %q", se.Code)
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Fallback  string  `json:"fallback"`
	Color     string  `json:"color"`
	Title     string  `json:"title"`
	TitleLink string  `json:"title_link"`
	Text      string  `json:"text"`
	TimeStamp int64   `json:"ts"`
	Fields    []Field `json:"fields"`
}

// Message is the JSON object expected by Slack
type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}
