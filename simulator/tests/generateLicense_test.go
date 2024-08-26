package tests

import (
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/commons/logs"
	"im-server/services/commonservices"
	"testing"
)

func TestGenerateLicense(t *testing.T) {
	//str := commonservices.GenerateLicenseStr(&pbobjs.LicenseConf{
	//	Appkey:            "test",
	//	Secret:            "test123",
	//	SecureKey:         "test123",
	//	CreatedAt:         time.Now().UnixMilli(),
	//	EndedAt:           time.Now().Add(time.Hour * 24 * 365 * 10).UnixMilli(),
	//	AppName:           "test",
	//	RegistedUserCount: 0,
	//	GroupMemberCount:  0,
	//	OpenPush:          true,
	//})
	//fmt.Println(str)

	if err := configures.InitConfigures(); err != nil {
		//logs.Error("Init Configures failed.", err)
		fmt.Println("Init Configures failed.", err)
		return
	}
	//init logs
	logs.InitLogs()
	//init mysql
	if err := dbcommons.InitMysql(); err != nil {
		logs.Error("Init Mysql failed.", err)
		return
	}
	fmt.Println(commonservices.GetAppInfo("test"))
}
