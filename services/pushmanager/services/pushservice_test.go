package services

import (
	"im-server/services/commonservices"
	"testing"
)

func TestHandleJPushOptionsDoesNotMutateChannelTemplates(t *testing.T) {
	xiaomiTemplate := `{"keywords1":"${senderName}","keywords2":"${pushText}"}`
	oppoContentTemplate := map[string]string{
		"sender":  "${senderName}",
		"content": "${pushText}",
	}
	oppoTitleTemplate := map[string]string{
		"group": "${groupName}",
	}
	config := &commonservices.JPushOptions{
		Classification: 1,
		ThirdPartyChannel: &commonservices.JPushThirdPartyChannel{
			Xiaomi: &commonservices.JPushXiaomiChannel{
				MiTemplateId:    "xiaomi-template",
				MiTemplateParam: xiaomiTemplate,
			},
			Oppo: &commonservices.JPushOppoChannel{
				PrivateContentParameters: oppoContentTemplate,
				PrivateTitleParameters:   oppoTitleTemplate,
			},
		},
	}

	first := handleJPushOptions(config, map[string]string{
		"senderName": "sender-1",
		"pushText":   "message-1",
		"groupName":  "group-1",
	})
	second := handleJPushOptions(config, map[string]string{
		"senderName": "sender-2",
		"pushText":   "message-2",
		"groupName":  "group-2",
	})

	if got, want := first.ThirdPartyChannel.Xiaomi.MiTemplateParam, `{"keywords1":"sender-1","keywords2":"message-1"}`; got != want {
		t.Fatalf("first Xiaomi template params = %q, want %q", got, want)
	}
	if got, want := second.ThirdPartyChannel.Xiaomi.MiTemplateParam, `{"keywords1":"sender-2","keywords2":"message-2"}`; got != want {
		t.Fatalf("second Xiaomi template params = %q, want %q", got, want)
	}
	if got := config.ThirdPartyChannel.Xiaomi.MiTemplateParam; got != xiaomiTemplate {
		t.Fatalf("Xiaomi config template was mutated: got %q, want %q", got, xiaomiTemplate)
	}
	if first.ThirdPartyChannel.Xiaomi == second.ThirdPartyChannel.Xiaomi ||
		first.ThirdPartyChannel.Xiaomi == config.ThirdPartyChannel.Xiaomi {
		t.Fatal("Xiaomi channel was not deep-copied per request")
	}

	assertJPushStringMapValue(t, first.ThirdPartyChannel.Oppo.PrivateContentParameters, "sender", "sender-1")
	assertJPushStringMapValue(t, first.ThirdPartyChannel.Oppo.PrivateContentParameters, "content", "message-1")
	assertJPushStringMapValue(t, first.ThirdPartyChannel.Oppo.PrivateTitleParameters, "group", "group-1")
	assertJPushStringMapValue(t, second.ThirdPartyChannel.Oppo.PrivateContentParameters, "sender", "sender-2")
	assertJPushStringMapValue(t, second.ThirdPartyChannel.Oppo.PrivateContentParameters, "content", "message-2")
	assertJPushStringMapValue(t, second.ThirdPartyChannel.Oppo.PrivateTitleParameters, "group", "group-2")
	assertJPushStringMapValue(t, config.ThirdPartyChannel.Oppo.PrivateContentParameters, "sender", "${senderName}")
	assertJPushStringMapValue(t, config.ThirdPartyChannel.Oppo.PrivateContentParameters, "content", "${pushText}")
	assertJPushStringMapValue(t, config.ThirdPartyChannel.Oppo.PrivateTitleParameters, "group", "${groupName}")
	if first.ThirdPartyChannel.Oppo == second.ThirdPartyChannel.Oppo ||
		first.ThirdPartyChannel.Oppo == config.ThirdPartyChannel.Oppo {
		t.Fatal("OPPO channel was not deep-copied per request")
	}

	first.ThirdPartyChannel.Oppo.PrivateContentParameters["content"] = "changed"
	assertJPushStringMapValue(t, second.ThirdPartyChannel.Oppo.PrivateContentParameters, "content", "message-2")
	assertJPushStringMapValue(t, config.ThirdPartyChannel.Oppo.PrivateContentParameters, "content", "${pushText}")
}

func assertJPushStringMapValue(t *testing.T, values map[string]string, key, want string) {
	t.Helper()
	if got := values[key]; got != want {
		t.Fatalf("map value %q = %q, want %q", key, got, want)
	}
}
