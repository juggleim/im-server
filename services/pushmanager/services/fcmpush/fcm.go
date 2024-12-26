package fcmpush

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FcmPushClient struct {
	fcmClient *messaging.Client
}

func NewFcmPushClient(jsonBs []byte) (*FcmPushClient, error) {
	opts := []option.ClientOption{option.WithCredentialsJSON(jsonBs)}
	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "error in initializing firebase app")
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}
	return &FcmPushClient{
		fcmClient: client,
	}, nil
}

func (f *FcmPushClient) SendPush(title, body, token string, jcExts map[string]interface{}) error {
	client := f.fcmClient
	if client == nil {
		return errors.New("not initialized")
	}
	extMap := map[string]string{}
	for k, v := range jcExts {
		extMap[k] = fmt.Sprintf("%v", v)
	}
	message := &messaging.Message{
		Data: extMap,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Title:       title,
				Body:        body,
				ClickAction: "com.j.im.intent.MESSAGE_CLICK",
			},
		},
		Token: token,
	}

	_, err := client.Send(context.Background(), message)
	if err != nil {
		return errors.Wrap(err, "error in sending push notification")
	}

	return nil
}

func (f *FcmPushClient) MultipleSendPush(title string, body string, tokens []string) error {
	client := f.fcmClient
	if client == nil {
		return errors.New("not initialized")
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
		return errors.Wrap(err, "error in sending push notification")
	}

	for _, result := range response.Responses {
		if result.Error != nil {
			return errors.Wrap(err, "error in sending push notification")
		}
	}

	return nil
}
