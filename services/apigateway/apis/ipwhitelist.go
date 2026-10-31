package apis

import (
	"fmt"
	"net"
	"net/netip"
	"strings"

	"im-server/commons/errs"
	"im-server/commons/tools"

	"github.com/gin-gonic/gin"
)

type ipWhitelist struct {
	addresses map[netip.Addr]struct{}
	prefixes  []netip.Prefix
}

func NewIPWhitelistMiddleware(entries []string) (gin.HandlerFunc, error) {
	whitelist, err := parseIPWhitelist(entries)
	if err != nil {
		return nil, err
	}
	return func(ctx *gin.Context) {
		if whitelist == nil {
			ctx.Next()
			return
		}

		clientIP, ok := directClientIP(ctx.Request.RemoteAddr)
		if !ok || !whitelist.contains(clientIP) {
			tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_IP_FORBIDDEN)
			ctx.Abort()
			return
		}
		ctx.Next()
	}, nil
}

func parseIPWhitelist(entries []string) (*ipWhitelist, error) {
	if len(entries) == 0 {
		return nil, nil
	}
	whitelist := &ipWhitelist{addresses: make(map[netip.Addr]struct{})}
	for _, rawEntry := range entries {
		entry := strings.TrimSpace(rawEntry)
		if entry == "" {
			return nil, fmt.Errorf("empty entry")
		}
		if address, err := netip.ParseAddr(entry); err == nil {
			whitelist.addresses[address.Unmap()] = struct{}{}
			continue
		}
		prefix, err := netip.ParsePrefix(entry)
		if err != nil {
			return nil, fmt.Errorf("invalid IP or CIDR %q", rawEntry)
		}
		whitelist.prefixes = append(whitelist.prefixes, prefix.Masked())
	}
	return whitelist, nil
}

func directClientIP(remoteAddress string) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(remoteAddress)
	if err != nil {
		host = remoteAddress
	}
	address, err := netip.ParseAddr(strings.TrimSpace(host))
	if err != nil {
		return netip.Addr{}, false
	}
	return address.Unmap(), true
}

func (whitelist *ipWhitelist) contains(address netip.Addr) bool {
	if _, ok := whitelist.addresses[address.Unmap()]; ok {
		return true
	}
	for _, prefix := range whitelist.prefixes {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}
