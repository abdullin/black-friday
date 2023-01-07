package node

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// Cancel creates a new cancelable context and registers a function to
// cancel the context when the user sends a SIGINT signal (e.g. by pressing
// Ctrl+C).
func Cancel() (context.Context, func()) {
	ctx := context.Background()
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			fmt.Println("Ctrl-C")
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, func() {
		signal.Stop(c)
		cancel()
	}
}
