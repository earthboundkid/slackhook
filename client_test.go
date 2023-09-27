package slackhook_test

import (
	"context"
	"testing"

	"github.com/carlmjohnson/slackhook"
)

func TestClient(t *testing.T) {
	c := slackhook.New(slackhook.MockClient, nil)
	err := c.PostCtx(context.Background(), slackhook.NoOpLogger, slackhook.Message{
		Text: "Hello",
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	c = slackhook.New("", nil)
	err = c.PostCtx(context.Background(), slackhook.NoOpLogger, slackhook.Message{
		Text: "Hello",
	})
	if err == nil {
		t.Fatal("want error; got nil")
	}
}
