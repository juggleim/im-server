package kvdbcommons

import "encoding/base64"

func Bytes2SafeKey(bs []byte) string {
	if len(bs) <= 0 {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bs)
}
