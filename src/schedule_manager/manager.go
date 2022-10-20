package schedule_manager

import (
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/fb"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type Manager interface {
	UpdateSchedules() (err error)
}

type manager struct {
	logger      *logrus.Logger
	fb          *fb.Fb
	cacheClient *redis.UniversalClient
}

func CreateAndInitManager(
	logger *logrus.Logger,
	fb *fb.Fb,
	cacheClient *redis.UniversalClient,
) {
	mgr := &manager{
		logger:      logger,
		fb:          fb,
		cacheClient: cacheClient,
	}

	for {
		err := mgr.UpdateSchedules()
		if err != nil {
			logger.Error(err)
			return
		}

		time.Sleep(20 * time.Second)
	}
}
