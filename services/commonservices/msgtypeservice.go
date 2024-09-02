package commonservices

const (
	MentionType_All        string = "mention_all"
	MentionType_Someone    string = "mention_someone"
	MentionType_AllSomeone string = "mention_all_someone"
)

type TextMsg struct {
	Content string `json:"content"`
	Extra   string `json:"extra"`
}
