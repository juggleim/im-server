package apis

import (
	"im-server/services/admingateway/services"
	"im-server/services/pushmanager/storages/dbs"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIosCer(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	//isProduct := ctx.Query("is_product")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	dao := dbs.IosCertificateDao{}
	cerDao, err := dao.Find(appkey)
	if err == nil {
		cerDao.Certificate = []byte{}
		cerDao.VoipCert = []byte{}
	}
	services.SuccessHttpResp(ctx, cerDao)
}

type IosPushReq struct {
	AppKey    string `json:"app_key"`
	Package   string `json:"package"`
	IsProduct int    `json:"is_product"`
	CertPath  string `json:"cert_path"`
	CertPwd   string `json:"cert_pwd"`
}

func SetIosPushConf(ctx *gin.Context) {
	var req IosPushReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	//save to db
	dao := dbs.IosCertificateDao{}
	err := dao.Upsert(dbs.IosCertificateDao{
		Package:   req.Package,
		CertPath:  req.CertPath,
		AppKey:    req.AppKey,
		CertPwd:   req.CertPwd,
		IsProduct: req.IsProduct,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	services.SuccessHttpResp(ctx, nil)
}

func UploadIosCer(ctx *gin.Context) {
	certPath := ctx.PostForm("cert_path")
	certPwd := ctx.PostForm("cert_pwd")
	voipCertPath := ctx.PostForm("voip_cert_path")
	voipCertPwd := ctx.PostForm("voip_cert_pwd")
	iosPackage := ctx.PostForm("package")
	appkey := ctx.PostForm("app_key")
	if appkey == "" || iosPackage == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
		})
		return
	}
	isProStr := ctx.PostForm("is_product")
	isProduct := 0
	if isProStr != "" {
		isProduct, _ = strconv.Atoi(isProStr)
	}

	iosCerBs := []byte{}
	fileHeader, err := ctx.FormFile("ioscer")
	if err == nil {
		file, err := fileHeader.Open()
		if err == nil {
			defer file.Close()
			bs, err := io.ReadAll(file)
			if err == nil {
				iosCerBs = append(iosCerBs, bs...)
			}
		}
	}
	voipIosCerBs := []byte{}
	voipFileHeader, err := ctx.FormFile("voip_ioscer")
	if err == nil {
		voipFile, err := voipFileHeader.Open()
		if err == nil {
			defer voipFile.Close()
			bs, err := io.ReadAll(voipFile)
			if err == nil {
				voipIosCerBs = append(voipIosCerBs, bs...)
			}
		}
	}
	//save to db
	dao := dbs.IosCertificateDao{}
	err = dao.Upsert(dbs.IosCertificateDao{
		AppKey:       appkey,
		Package:      iosPackage,
		IsProduct:    isProduct,
		CertPwd:      certPwd,
		VoipCertPwd:  voipCertPwd,
		Certificate:  iosCerBs,
		CertPath:     certPath,
		VoipCert:     voipIosCerBs,
		VoipCertPath: voipCertPath,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	services.SuccessHttpResp(ctx, nil)
}
