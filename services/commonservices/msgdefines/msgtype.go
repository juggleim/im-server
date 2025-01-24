package msgdefines

const (
	MentionType_All        string = "mention_all"
	MentionType_Someone    string = "mention_someone"
	MentionType_AllSomeone string = "mention_all_someone"
)

var InnerMsgType_Text string = "jg:text"

type TextMsg struct {
	Content string `json:"content"`
	Extra   string `json:"extra"`
}

var InnerMsgType_Img string = "jg:img"
var InnerMsgType_Voice string = "jg:voice"
var InnerMsgType_File string = "jg:file"
var InnerMsgType_Video string = "jg:video"
var InnerMsgType_Merge string = "jg:merge"
var InnerMsgType_VoiceCall string = "jg:voicecall"
var InnerMsgType_CallFinishNtf string = "jg:callfinishntf"

var InnerMsgType_StreamText string = "jgs:text"

type StreamMsg struct {
	Content string `json:"content"`
}
