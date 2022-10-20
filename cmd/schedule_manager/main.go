package main

import (
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/fb"
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/goredis"
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/logrux"
	"github.com/NEO-TAT/tat_auto_roll_call_service/src/schedule_manager"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			logrux.NewLogger,
			goredis.NewGoRedisClient,
			fb.CreateFirebaseClients,
		),
		fx.Invoke(
			schedule_manager.CreateAndInitManager,
		),
	).Run()
}
