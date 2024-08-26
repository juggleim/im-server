package jpush

import "testing"

var (
	appKey       = "0380e1068c6d9e3959f32c8c"
	masterSecret = "b2b7d868ffeda29c9da35072"
)

func TestJPushClient(t *testing.T) {
	client := NewJpushClient(appKey, masterSecret)
	msgId, err := client.Push(&Payload{
		Platform: NewPlatform().All(),
		Audience: NewAudience().SetRegistrationId("1104a89793a130a90cf"),
		Notification: &Notification{
			Alert: "JPush test",
			Android: &AndroidNotification{
				Alert: "hello, JPush!",
				Title: "JPush test",
			},
			Ios:      nil,
			WinPhone: nil,
		},
		Message: &Message{
			Title:   "Hello",
			Content: "Hello, JPush!",
		},
		SmsMessage: nil,
		Options:    nil,
		//Cid:        "1104a89793a130a90cf",
	})
	t.Log(msgId, err)
}
