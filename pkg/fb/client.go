package fb

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/NEO-TAT/tat_auto_roll_call_service/pkg/env"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/api/option"
)

type Fb struct {
	Store   *firestore.Client
	Message *messaging.Client
}

func CreateFirebaseClients(lc fx.Lifecycle, logger *logrus.Logger) (*Fb, error) {
	var app *firebase.App
	var store *firestore.Client
	var message *messaging.Client

	ctx := context.Background()

	config := option.WithCredentialsFile(env.FirebaseConfigPath)
	app, err := firebase.NewApp(ctx, nil, config)
	if err != nil {
		logger.Errorln("[Firebase] create client failed:", err)
		return nil, err
	}

	store, err = app.Firestore(ctx)
	if err != nil {
		logger.Errorln("[Firebase] create client failed:", err)
		return nil, err
	}

	message, err = app.Messaging(ctx)
	if err != nil {
		logger.Errorln("[Firebase] create client failed:", err)
		return nil, err
	}

	lc.Append(fx.Hook{
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
		Store:   store,
		Message: message,
	}, nil
}
