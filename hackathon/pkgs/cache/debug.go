package cache

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/redis/go-redis/v9"
	"time"
)

type DebugHook struct {
	log log.Logger
}

func NewDebugHook(log log.Logger) DebugHook {
	return DebugHook{
		log: log,
	}
}

func (h DebugHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h DebugHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		startAt := time.Now()
		err := next(ctx, cmd)
		logger := log.With(h.log, "args", cmd.Args())
		logger = log.With(logger, "err", err)
		logger = log.With(logger, "cost", time.Since(startAt).String())
		_ = logger.Log("[redis]%s", cmd.Name())
		return err
	}
}

func (h DebugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
