package services

type ServerLogType string

const (
	ServerLogType_UserConnect ServerLogType = "userconnect"
	ServerLogType_Connect     ServerLogType = "connect"
	ServerLogType_Business    ServerLogType = "business"
)
