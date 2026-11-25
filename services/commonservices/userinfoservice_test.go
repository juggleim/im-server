package commonservices

import (
	"testing"

	"im-server/commons/pbdefines/pbobjs"
)

func TestTargetUserInfoToDisplayUserInfoIncludesUserType(t *testing.T) {
	got := targetUserInfoToDisplayUserInfo(&TargetUserInfo{
		UserId:   "user-1",
		UserType: pbobjs.UserType_Bot,
	})
	if got == nil {
		t.Fatal("targetUserInfoToDisplayUserInfo() returned nil")
	}
	if got.UserId != "user-1" {
		t.Fatalf("UserId = %q, want user-1", got.UserId)
	}
	if got.UserType != pbobjs.UserType_Bot {
		t.Fatalf("UserType = %v, want %v", got.UserType, pbobjs.UserType_Bot)
	}
}

func TestTargetUserInfoToDisplayUserInfoHandlesNil(t *testing.T) {
	if got := targetUserInfoToDisplayUserInfo(nil); got != nil {
		t.Fatalf("targetUserInfoToDisplayUserInfo(nil) = %#v, want nil", got)
	}
}
