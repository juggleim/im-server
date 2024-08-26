package fcmpush

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FcmPush struct {
	credential string // json string of firebase credentials
}

func (f *FcmPush) client() (client *messaging.Client, err error) {
	var app *firebase.App

	opts := []option.ClientOption{option.WithCredentialsJSON([]byte(f.credential))}

	// Initialize firebase app
	app, err = firebase.NewApp(context.Background(), nil, opts...)

	if err != nil {
		err = errors.Wrap(err, "error in initializing firebase app")
		return
	}
	client, err = app.Messaging(context.Background())

	return
}

func (f *FcmPush) SendPush(title, body, token string) (err error) {
	client, err := f.client()
	if err != nil {
		return
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	_, err = client.Send(context.Background(), message)
	if err != nil {
		err = errors.Wrap(err, "error in sending push notification")
		return
	}

	return
}

func (f *FcmPush) MultipleSendPush(
	title string, body string, tokens []string) (err error) {
	client, err := f.client()
	if err != nil {
		return
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	}

	response, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		err = errors.Wrap(err, "error in sending push notification")
		return
	}

	for _, result := range response.Responses {
		if result.Error != nil {
			err = errors.Wrap(err, "error in sending push notification")
			return
		}
	}

	return
}
