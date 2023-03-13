package goredis

import (
	"context"
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/env"
	"github.com/redis/go-redis/extra/rediscmd/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net"
)

type redisHook interface {
	redis.Hook
}

type logHook struct {
	logger *logrus.Logger
}

func createRedisHookWith(logger *logrus.Logger) redisHook {
	return &logHook{
		logger: logger,
	}
}

func (logHook) DialHook(_ redis.DialHook) redis.DialHook {
	return func(_ context.Context, network, addr string) (net.Conn, error) {
		// TODO: implementations.
		return nil, nil
	}
}

func (h logHook) ProcessHook(_ redis.ProcessHook) redis.ProcessHook {
	return func(_ context.Context, cmd redis.Cmder) error {
		if err := cmd.Err(); err != nil {
			h.logger.Errorf("[GoRedis] cmd: %s, err: %s", cmd.FullName(), err)
			return err
		}

		if env.IsDebugMode {
			h.logger.Infof("[GoRedis] cmd: %s", cmd.FullName())
		}

		return nil
	}
}

func (h logHook) ProcessPipelineHook(_ redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(_ context.Context, cmdList []redis.Cmder) error {
		summary, cmdStrings := rediscmd.CmdsString(cmdList)

		if err := cmdList[0].Err(); err != nil {
			h.logger.Errorf("[GoRedis] cmd: %s, cmdNum: %d, pipeline: %s, err: %s", cmdStrings, len(cmdList), summary, err)
			return err
		}

		if env.IsDebugMode {
			h.logger.Infof("[GoRedis] cmd: %s, cmdNum: %d, pipeline: %s", cmdStrings, len(cmdList), summary)
		}

		return nil
	}
}

// NewGoRedisClient return new redis.UniversalClient with fx.Lifecycle
func NewGoRedisClient(lc fx.Lifecycle, logger *logrus.Logger) (*redis.UniversalClient, error) {
	var c redis.UniversalClient

	if len(env.RedisHosts) <= 1 {
		c = redis.NewClient(&redis.Options{
			Addr:     env.RedisHosts[0],
			Password: env.RedisPassword,
			DB:       0,
		})
	} else {
		c = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    env.RedisHosts,
			Password: env.RedisPassword,
		})
	}

	hook := createRedisHookWith(logger)
	c.AddHook(hook)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := c.Ping(ctx).Err()
			if err != nil {
				logger.Errorln("Failed to establish connection to redis service")
			}

			return nil
		},
		OnStop: func(context.Context) error {
			return c.Close()
		},
	})

	return &c, nil
}
