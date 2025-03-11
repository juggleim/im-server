package models

type Prompt struct {
	ID          int64
	UserId      string
	Prompts     string
	CreatedTime int64
	UpdatedTime int64
	AppKey      string
}

type IPromptStorage interface {
	Create(prompt Prompt) (int64, error)
	UpdatePrompts(appkey, userId string, id int64, prompts string) error
	DelPrompts(appkey, userId string, id int64) error
	BatchDelPrompts(appkey, userId string, ids []int64) error
	FindPrompt(appkey, userId string, id int64) (*Prompt, error)
	QryPrompts(appkey, userId string, limit int64, startId int64) ([]*Prompt, error)
}
