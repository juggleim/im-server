package services

import (
	"testing"

	"im-server/commons/pbdefines/pbobjs"
)

func TestLargeGroupDispatchPolicy(t *testing.T) {
	tests := []struct {
		name                 string
		forceCloseOfflineMsg bool
		isOnline             bool
		wantCloseOffline     bool
		wantSkipUserInfo     bool
	}{
		{
			name:     "online member",
			isOnline: true,
		},
		{
			name:             "offline member",
			wantCloseOffline: true,
			wantSkipUserInfo: true,
		},
		{
			name:                 "online member with forced offline storage",
			forceCloseOfflineMsg: true,
			isOnline:             true,
			wantCloseOffline:     true,
		},
		{
			name:                 "offline member with forced offline storage",
			forceCloseOfflineMsg: true,
			wantCloseOffline:     true,
			wantSkipUserInfo:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCloseOffline, gotSkipUserInfo := largeGroupDispatchPolicy(tt.forceCloseOfflineMsg, tt.isOnline)
			if gotCloseOffline != tt.wantCloseOffline {
				t.Fatalf("closeOffline = %v, want %v", gotCloseOffline, tt.wantCloseOffline)
			}
			if gotSkipUserInfo != tt.wantSkipUserInfo {
				t.Fatalf("skipUserInfo = %v, want %v", gotSkipUserInfo, tt.wantSkipUserInfo)
			}
		})
	}
}

func TestSkipUserInfoForDispatch(t *testing.T) {
	tests := []struct {
		name        string
		channelType pbobjs.ChannelType
		isOffline   bool
		want        bool
	}{
		{name: "online group member", channelType: pbobjs.ChannelType_Group},
		{name: "offline group member", channelType: pbobjs.ChannelType_Group, isOffline: true, want: true},
		{name: "private member", channelType: pbobjs.ChannelType_Private, isOffline: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipUserInfoForDispatch(tt.channelType, tt.isOffline); got != tt.want {
				t.Fatalf("skipUserInfoForDispatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
