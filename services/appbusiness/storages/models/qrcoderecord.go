package models

type QrCodeRecordStatus int

var (
	QrCodeRecordStatus_Default QrCodeRecordStatus = 0
	QrCodeRecordStatus_OK      QrCodeRecordStatus = 1
)

type QrCodeRecord struct {
	CodeId      string
	Status      QrCodeRecordStatus
	CreatedTime int64
	UserId      string
	AppKey      string
}

type IQrCodeRecordStorage interface {
	Create(item QrCodeRecord) error
	FindById(appkey, codeId string) (*QrCodeRecord, error)
	UpdateStatus(appkey, codeId string, status QrCodeRecordStatus, userId string) error
}
