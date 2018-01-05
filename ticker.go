package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// TickerCallback callback function to run on ticker interval
type TickerCallback func(ctx context.Context)

// startTicker start the ticker and run the callback at each interval
func startTicker(ctx context.Context, tick TickerCallback) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		// cleanup the ticker on return
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				go tick(ctx)
			case <-ctx.Done():
				log.Info().Msg("Ticker stopping by context done")
				return
			}
		}
	}()
	log.Info().Msg("Ticker Started")
}
