package apis

import (
	"im-server/services/admingateway/services"
	"im-server/services/pushmanager/storages/dbs"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IosCert struct {
	Package   string `json:"package"`
	CertLen   int    `json:"cert_len"`
	AppKey    string `json:"app_key"`
	CertPwd   string `json:"cert_pwd"`
	IsProduct int    `json:"is_product"`
}

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
	cerDao, _ := dao.Find(appkey)
	cerDao.Certificate = []byte{}
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
	iosPackage := ctx.PostForm("package")
	appkey := ctx.PostForm("app_key")
	isProStr := ctx.PostForm("is_product")
	isProduct := 0
	if isProStr != "" {
		isProduct, _ = strconv.Atoi(isProStr)
	}

	fileHeader, err := ctx.FormFile("ioscer")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	defer file.Close()
	iosCerBs, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	//save to db
	dao := dbs.IosCertificateDao{}
	err = dao.Upsert(dbs.IosCertificateDao{
		Package:     iosPackage,
		Certificate: iosCerBs,
		AppKey:      appkey,
		CertPwd:     certPwd,
		IsProduct:   isProduct,
		CertPath:    certPath,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}
	services.SuccessHttpResp(ctx, nil)
}
