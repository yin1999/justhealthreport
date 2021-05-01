package httpclient

import (
	"context"
	"testing"
)

func TestLogin(t *testing.T) {
	// ctx, cc := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cc()
	ctx := context.Background()
	_, err := login(ctx, [2]string{"182211901233", "304813"})
	if err != nil {
		t.Fatal(err)
	}
}
