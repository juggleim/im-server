package services

import (
	"context"
	"im-server/commons/pbdefines/pbobjs"
	"testing"
)

func TestCancelKickUserPush(t *testing.T) {
	for _, tc := range []struct {
		name string
		req  *pbobjs.KickUserReq
		want int
	}{
		{"enabled", &pbobjs.KickUserReq{UserId: "user", CancelPush: true}, 1},
		{"disabled", &pbobjs.KickUserReq{UserId: "user"}, 0},
		{"empty user", &pbobjs.KickUserReq{CancelPush: true}, 0},
		{"filtered connections", &pbobjs.KickUserReq{UserId: "user", CancelPush: true, Platforms: []string{"iOS"}, DeviceIds: []string{"device"}}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			calls := 0
			cancelKickUserPush(ctx, tc.req, func(gotCtx context.Context, userId string) {
				calls++
				if gotCtx != ctx || userId != tc.req.UserId {
					t.Fatalf("incorrect removal context or user: %s", userId)
				}
			})
			if calls != tc.want {
				t.Fatalf("removal calls=%d, want %d", calls, tc.want)
			}
		})
	}
}
