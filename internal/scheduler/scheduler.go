package scheduler

import (
	"context"
	"fmt"
	"time"
)

func RunTicker(ctx context.Context, round func(ctx context.Context), interval int) {
	if interval == 0 {
		interval = 15
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			round(ctx)
		case <-ctx.Done():
			fmt.Println("interrupt, canceling")
			return
		}
	}
}
