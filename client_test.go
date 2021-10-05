package slackhook_test

import (
	"context"
	"testing"

	"github.com/carlmjohnson/slackhook"
)

func TestClient(t *testing.T) {
	c := slackhook.New("", nil)
	err := c.PostCtx(context.Background(), slackhook.Message{
		Text: "Hello",
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
