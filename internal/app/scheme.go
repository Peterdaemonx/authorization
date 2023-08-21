package app

import (
	"context"
	"fmt"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	mastercardScheme "gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/mastercard"
	visaScheme "gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/visa"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/sequences"
	"gitlab.cmpayments.local/creditcard/platform"
)

func Mip(ctx context.Context, logger platform.Logger, pool *connection.Pool, mip *mastercardScheme.Mip, ss sequences.Store, shutdown func()) *mastercardScheme.Mip {
	if mip == nil {
		stanGen := sequences.NewDaily(100, 100000, 999999, ss)

		go func() {
			err := stanGen.Fill(ctx)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("Cannot keep STAN list filled: %s", err.Error()))
				shutdown()
			}
		}()

		mip := mastercardScheme.NewMip(pool, &stanGen)
		return &mip
	}

	return mip
}

func Eas(ctx context.Context, logger platform.Logger, pool *connection.Pool, eas *visaScheme.Eas, ss sequences.Store, sourceID string, tickDelay time.Duration, shutdown func()) *visaScheme.Eas {
	if eas == nil {
		stanGen := sequences.NewDaily(100, 100000, 999999, ss)

		go func() {
			err := stanGen.Fill(ctx)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("Cannot keep STAN list filled: %s", err.Error()))
				shutdown()
			}
		}()

		eas := visaScheme.NewEas(pool, &stanGen, sourceID)
		if tickDelay != 0 {
			go func() {
				// receive time from config.
				ticker := time.NewTicker(tickDelay)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						err := eas.Echo(ctx)
						if err != nil {
							logger.Error(ctx, err.Error())
						}
					case <-ctx.Done():
						return
					}
				}
			}()

		}
		return &eas
	}

	return eas
}
