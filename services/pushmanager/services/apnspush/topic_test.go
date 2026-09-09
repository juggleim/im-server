package apnspush

import (
	"strings"
	"testing"

	"github.com/sideshow/apns2"
)

func TestBundleValidationMatchesConsole(t *testing.T) {
	for _, topic := range []string{"", "com.example app", "com/example", "com_example", "应用", strings.Repeat("a", 101)} {
		if validTopic(topic) {
			t.Errorf("accepted invalid bundle %q", topic)
		}
	}
	base := strings.Repeat("a", 100)
	clients := fakeClients(&fakeTransport{})
	clients.ApnsVoipClient = clients.ApnsClient
	_, n, err := Route(clients, base, &apns2.Notification{}, true, "ab", "cd")
	if err != nil || n.Topic != base+".voip" {
		t.Fatalf("valid 100-byte bundle failed VoIP routing: %v", err)
	}
}
