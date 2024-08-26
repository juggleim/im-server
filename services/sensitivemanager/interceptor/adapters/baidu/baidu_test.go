package baidu

import (
	"fmt"
	"testing"
)

var interceptor = NewInterceptor("oCaDB7oTHPehR7sHn7q7WP8u", "qlQD7aaVon8bDvTZIUvcIvBqxYtyDW5C", false)

func TestInterceptor_InterceptText(t *testing.T) {
	result, _, err := interceptor.InterceptText("我是你爹, 你好， 免费翻墙,找小姐, 傻逼")
	fmt.Println(result, err)
}

func TestInterceptor_InterceptImage(t *testing.T) {
	fmt.Println(interceptor.InterceptImage("https://pic.rmb.bdstatic.com/bjh/240423/dump/bc7bb919742e0e31cf9381dfcb5b0e1e.jpeg"))
}
