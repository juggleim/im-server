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

type BaseStreamMsg interface {
	GetStreamId() string
	GetMsgType() string
}

var InnerMsgType_StreamMsg string = "jg:streamtext"

type StreamMsg struct {
	msgType string `json:"-"`

	StreamId   string `json:"stream_id"`
	Content    string `json:"content"`
	IsFinished bool   `json:"is_finished"`
	Seq        int    `json:"seq"`
}

func NewStreamMsg(streamId, content string, isFinished bool, seq int) *StreamMsg {
	return &StreamMsg{
		msgType: InnerMsgType_StreamMsg,

		StreamId:   streamId,
		Content:    content,
		IsFinished: isFinished,
		Seq:        seq,
	}
}

func (msg *StreamMsg) GetStreamId() string {
	return msg.StreamId
}

func (msg *StreamMsg) GetMsgType() string {
	return msg.msgType
}

var InnerMsgType_StreamAppendMsg string = "jg:streamappend"

type StreamAppendMsg struct {
	msgType string `json:"-"`

	StreamId   string `json:"stream_id"`
	Content    string `json:"content"`
	Seq        int    `json:"seq"`
	IsFinished bool   `json:"is_finished"`
}

func NewStreamAppendMsg(streamId, content string, isFinished bool, seq int) *StreamAppendMsg {
	return &StreamAppendMsg{
		msgType: InnerMsgType_StreamAppendMsg,

		StreamId:   streamId,
		Content:    content,
		IsFinished: isFinished,
		Seq:        seq,
	}
}

func (msg *StreamAppendMsg) GetStreamId() string {
	return msg.StreamId
}

func (msg *StreamAppendMsg) GetMsgType() string {
	return msg.msgType
}
