package services

import (
	"errors"
	"testing"

	"im-server/commons/pbdefines/pbobjs"
)

func TestSetGlobalConverTagsIfChanged(t *testing.T) {
	conf := &ConverConfItem{
		AppKey:           "app",
		ConverId:         "conver",
		ChannelType:      pbobjs.ChannelType_Group,
		GlobalConverTags: map[string]bool{"tag_a": true},
	}

	persistCalls := 0
	err := conf.SetGlobalConverTagsIfChanged(map[string]bool{"tag_a": true}, func() error {
		persistCalls++
		return nil
	})
	if err != nil {
		t.Fatalf("SetGlobalConverTagsIfChanged() error = %v", err)
	}
	if persistCalls != 0 {
		t.Fatalf("persist calls = %d, want 0 for unchanged tags", persistCalls)
	}

	err = conf.SetGlobalConverTagsIfChanged(map[string]bool{"tag_b": true}, func() error {
		persistCalls++
		return nil
	})
	if err != nil {
		t.Fatalf("SetGlobalConverTagsIfChanged() error = %v", err)
	}
	if persistCalls != 1 {
		t.Fatalf("persist calls = %d, want 1 for changed tags", persistCalls)
	}
	if !globalConverTagsEqual(conf.GlobalConverTags, map[string]bool{"tag_b": true}) {
		t.Fatalf("GlobalConverTags = %v, want tag_b=true", conf.GlobalConverTags)
	}
}

func TestSetGlobalConverTagsIfChangedKeepsCacheOnPersistFailure(t *testing.T) {
	conf := &ConverConfItem{
		AppKey:           "app",
		ConverId:         "conver",
		ChannelType:      pbobjs.ChannelType_Private,
		GlobalConverTags: map[string]bool{"tag_a": true},
	}
	wantErr := errors.New("persist failed")

	err := conf.SetGlobalConverTagsIfChanged(map[string]bool{"tag_b": true}, func() error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("SetGlobalConverTagsIfChanged() error = %v, want %v", err, wantErr)
	}
	if !globalConverTagsEqual(conf.GlobalConverTags, map[string]bool{"tag_a": true}) {
		t.Fatalf("GlobalConverTags = %v, want original tag_a=true", conf.GlobalConverTags)
	}
}
