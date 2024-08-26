package hwpush

import (
	"context"
	"testing"
)

func TestPush(t *testing.T) {
	client, err := NewHwPushClient("110722347", "fc053abca92e8e626a3419ac75f34a1ceb594df0c523b35b5ef564c224b5169b")
	if err != nil {
		t.Error(err)
		return
	}

	res, err := client.SendMessage(context.Background(), &MessageRequest{
		Message: &Message{
			Notification: &Notification{
				Title: "title test",
				Body:  "body test",
			},
			Android: &AndroidConfig{
				Notification: &AndroidNotification{
					Title:        "wahu",
					Body:         "body test",
					DefaultSound: true,
					ClickAction: &ClickAction{
						Type: 3,
					},
				},
			},
			Token: []string{"IQAAAACy06JCAACMzefr-DO_VkRF6ToByRDloyM2nP-8YO1VQnbbXNf2eNveQoOnynPxFLsmJEnnSeq4-teEM-gGmVqfDfSsSTVBV9ljK8OgQJFfRQ"},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
