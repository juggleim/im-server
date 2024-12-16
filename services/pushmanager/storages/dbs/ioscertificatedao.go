package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type IosCertificateDao struct {
	Package     string `gorm:"package" json:"package"`
	Certificate []byte `gorm:"certificate" json:"certificate"`
	CertPath    string `gorm:"cert_path" json:"cert_path"`
	AppKey      string `gorm:"app_key" json:"app_key"`
	CertPwd     string `gorm:"cert_pwd" json:"cert_pwd"`
	IsProduct   int    `gorm:"is_product" json:"is_product"`

	VoipCert     []byte `gorm:"voip_cert" json:"voip_cert"`
	VoipCertPwd  string `gorm:"voip_cert_pwd" json:"voip_cert_pwd"`
	VoipCertPath string `gorm:"voip_cert_path" json:"voip_cert_path"`
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
	var sql string = ""
	if len(item.Certificate) > 0 && len(item.VoipCert) > 0 {
		sql = fmt.Sprintf("INSERT INTO %s (app_key,package,is_product,cert_pwd,voip_cert_pwd,certificate,cert_path,voip_cert,voip_cert_path)VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE package=VALUES(package),is_product=VALUES(is_product),cert_pwd=VALUES(cert_pwd),voip_cert_pwd=VALUES(voip_cert_pwd),certificate=VALUES(certificate),cert_path=VALUES(cert_path),voip_cert=VALUES(voip_cert),voip_cert_path=VALUES(voip_cert_path)", cer.TableName())
		return dbcommons.GetDb().Exec(sql, item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.VoipCertPwd, item.Certificate, item.CertPath, item.VoipCert, item.VoipCertPath).Error
	} else if len(item.Certificate) > 0 {
		sql = fmt.Sprintf("INSERT INTO %s (app_key,package,is_product,cert_pwd,voip_cert_pwd,certificate,cert_path)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE package=VALUES(package),is_product=VALUES(is_product),cert_pwd=VALUES(cert_pwd),voip_cert_pwd=VALUES(voip_cert_pwd),certificate=VALUES(certificate),cert_path=VALUES(cert_path)", cer.TableName())
		return dbcommons.GetDb().Exec(sql, item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.VoipCertPwd, item.Certificate, item.CertPath).Error
	} else if len(item.VoipCert) > 0 {
		sql = fmt.Sprintf("INSERT INTO %s (app_key,package,is_product,cert_pwd,voip_cert_pwd,voip_cert,voip_cert_path)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE package=VALUES(package),is_product=VALUES(is_product),cert_pwd=VALUES(cert_pwd),voip_cert_pwd=VALUES(voip_cert_pwd),voip_cert=VALUES(voip_cert),voip_cert_path=VALUES(voip_cert_path)", cer.TableName())
		return dbcommons.GetDb().Exec(sql, item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.VoipCertPwd, item.VoipCert, item.VoipCertPath).Error
	} else {
		sql = fmt.Sprintf("INSERT INTO %s (app_key,package,is_product,cert_pwd,voip_cert_pwd)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE package=VALUES(package),is_product=VALUES(is_product),cert_pwd=VALUES(cert_pwd),voip_cert_pwd=VALUES(voip_cert_pwd)", cer.TableName())
		return dbcommons.GetDb().Exec(sql, item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.VoipCertPwd).Error
	}
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
