package fb

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/env"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/api/option"
)

type Fb struct {
	Store *firestore.Client
}

func CreateFirebaseClients(lc fx.Lifecycle, logger *logrus.Logger) (*Fb, error) {
	var app *firebase.App
	var store *firestore.Client

	config := option.WithCredentialsFile(env.FirebaseConfigPath)
	app, err := firebase.NewApp(context.Background(), nil, config)
	if err != nil {
		logger.Errorln("[Firebase] create client failed:", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			store, err = app.Firestore(ctx)
			if err != nil {
				logger.Errorln("[Firebase] create client failed:", err)
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := store.Close()
			if err != nil {
				logger.Errorln("[Firebase] close client failed:", err)
				return err
			}

			return nil
		},
	})

	return &Fb{
		Store: store,
	}, nil
}
