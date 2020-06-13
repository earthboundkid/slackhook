package slackhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return StatusErr(rsp.StatusCode)
	}

	// Drain connection
	_, err = io.Copy(ioutil.Discard, rsp.Body)
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

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
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

// Message is the JSON object expected by Slack
type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}
