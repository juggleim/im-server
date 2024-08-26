package tests

import "fmt"

const (
	NavURL    string = "http://127.0.0.1:8081/navigator/general"
	ApiURL    string = "http://127.0.0.1:8080"
	Appkey    string = "appkey"
	AppSecret string = "appsecret"
	WsAddr    string = "ws:127.0.0.1:9002"

	Token1 string = "CgZhcHBrZXkaIGrCchMUwd0haRDMl1PO8YZcrAlRe/s5KoYbbgS+/6a7"
	Token2 string = "CgZhcHBrZXkaIJU1zLmUbcP/rEWWrMfD4FRjvqhoRNe0S4630tK4dw2v"
	Token3 string = "CgZhcHBrZXkaILOhkwQ2bL6rSIHjt36xNHArIAXU+PDKfNBmtuH0gsi3"
	Token4 string = "CgZhcHBrZXkaIIJsXAQ6l3Snkch8hu+y5DCIziqkneaJ8aKjEWp4TR3T"
	Token5 string = "CgZhcHBrZXkaIOzq9LLCymFSeg50SMykuwh/F03ogefoopYNAHqkOZtx"

	User1 string = "t_userid1"
	User2 string = "t_userid2"
	User3 string = "t_userid3"
	User4 string = "t_userid4"
	User5 string = "t_userid5"

	Group1 string = "t_groupid1"
	Group2 string = "t_groupid2"
	Group3 string = "t_groupid3"

	TimeFormat string = "2006-01-02 15:04:05"
)

func Print(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
