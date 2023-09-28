package slackhook_test

import (
	"context"
	"testing"

	"github.com/carlmjohnson/slackhook"
)

func TestClient(t *testing.T) {
	var c slackhook.Client
	if err := c.Set(slackhook.MockClient); err != nil {
		t.Fatal(err)
	}
	err := c.Post(context.Background(), slackhook.NoOpLogger, nil, slackhook.Message{
		Text: "Hello",
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if err := c.Set(""); err != nil {
		t.Fatal(err)
	}
	err = c.Post(context.Background(), slackhook.NoOpLogger, nil, slackhook.Message{
		Text: "Hello",
	})
	if err == nil {
		t.Fatal("want error; got nil")
	}
}
