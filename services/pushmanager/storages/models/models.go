package models

import "context"

type IUserTagStorage interface {
	AddUserTags(ctx context.Context, appKey string, userID string, tags ...string) error
	DeleteUserTags(ctx context.Context, appKey string, userID string, tags ...string) error
	ClearUserTag(ctx context.Context, appKey string, userID string) error
	GetUserWithTags(ctx context.Context, appKey string, condition Condition, page int, perPage int) (userIDs []string, err error)
	GetUserTags(ctx context.Context, appKey string, userID string) (tags []string, err error)
}

type Condition struct {
	TagsAnd []string `json:"tags_and"`
	TagsOr  []string `json:"tags_or"`
}
