package dbs

import (
	"im-server/commons/dbcommons"
)

type IosCertificateDao struct {
	Package     string `gorm:"package" json:"package"`
	Certificate []byte `gorm:"certificate" json:"certificate"`
	CertPath    string `gorm:"cert_path" json:"cert_path"`
	AppKey      string `gorm:"app_key" json:"app_key"`
	CertPwd     string `gorm:"cert_pwd" json:"cert_pwd"`
	IsProduct   int    `gorm:"is_product" json:"is_product"`
	// CreatedTime time.Time `gorm:"created_time"`
}

func (cer IosCertificateDao) TableName() string {
	return "ioscertificates"
}

func (cer IosCertificateDao) FindByPackage(appkey, packageName string) (*IosCertificateDao, error) {
	var item IosCertificateDao
	err := dbcommons.GetDb().Where("app_key=? and package=?", appkey, packageName).Take(&item).Error
	return &item, err
}

func (cer IosCertificateDao) Upsert(item IosCertificateDao) error {
	if len(item.Certificate) <= 0 {
		err := dbcommons.GetDb().Exec("INSERT INTO ioscertificates (app_key,package,is_product,cert_pwd,cert_path) "+
			"VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE "+
			"is_product=?,cert_pwd=?,cert_path=?,package=?",
			item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.CertPath,
			item.IsProduct, item.CertPwd, item.CertPath, item.Package,
		).Error
		return err
	}
	err := dbcommons.GetDb().Exec("INSERT INTO ioscertificates (app_key,package,is_product,certificate,cert_pwd, cert_path) "+
		"VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "+
		"is_product=?,certificate=?,cert_pwd=?,cert_path=?,package=?",
		item.AppKey, item.Package, item.IsProduct, item.Certificate, item.CertPwd, item.CertPath,
		item.IsProduct, item.Certificate, item.CertPwd, item.CertPath, item.Package,
	).Error
	return err
}

func (cer IosCertificateDao) Create(item IosCertificateDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (cer IosCertificateDao) Find(appkey string) (*IosCertificateDao, error) {
	var item IosCertificateDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
