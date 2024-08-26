package apis

import (
	"encoding/json"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"im-server/services/fileplugin/dbs"
	"im-server/services/fileplugin/fileengine"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type FileConf struct {
	AppKey  string                 `json:"app_key,omitempty"`
	Channel string                 `json:"channel,omitempty"`
	Enable  int                    `json:"enable"`
	Conf    map[string]interface{} `json:"conf,omitempty"`
}

func GetFileConf(c *gin.Context) {
	appkey := c.Query("app_key")
	channel := c.Query("channel")

	dao := dbs.FileConfDao{}
	fileConf, err := dao.FindFileConf(appkey, channel)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			services.SuccessHttpResp(c, nil)
			return
		}
		c.JSON(http.StatusNotFound, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_Default,
		})
		return
	}

	var resp FileConf
	resp.AppKey = appkey
	resp.Channel = channel

	_ = json.Unmarshal([]byte(fileConf.Conf), &resp.Conf)
	services.SuccessHttpResp(c, resp)

}

func SetFileConf(ctx *gin.Context) {
	var req FileConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	pushConf := ""
	var conf any
	if strings.EqualFold(req.Channel, fileengine.ChannelAws) {
		conf = tools.MapToStruct[fileengine.S3Config](req.Conf)
	} else if strings.EqualFold(req.Channel, fileengine.ChannelQiNiu) {
		conf = tools.MapToStruct[fileengine.QiNiuConfig](req.Conf)
	} else if strings.EqualFold(req.Channel, fileengine.ChannelOss) {
		conf = tools.MapToStruct[fileengine.OssConfig](req.Conf)
	} else if strings.EqualFold(req.Channel, fileengine.ChannelMinio) {
		conf = tools.MapToStruct[fileengine.MinioConfig](req.Conf)
	}
	pushConf = tools.ToJson(conf)

	//save to db
	dao := dbs.FileConfDao{}
	dao.Upsert(dbs.FileConfDao{
		AppKey:  req.AppKey,
		Channel: req.Channel,
		Conf:    pushConf,
	})

	services.SuccessHttpResp(ctx, nil)
}

func GetFileConfCurrSwitch(c *gin.Context) {
	appkey := c.Query("app_key")
	dao := dbs.FileConfDao{}
	fConf, err := dao.FindEnableFileConf(appkey)
	if err != nil {
		c.JSON(http.StatusOK, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_NoFileEngine,
		})
		return
	}
	services.SuccessHttpResp(c, gin.H{"channel": fConf.Channel})
}

func GetFileConfs(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	dao := dbs.FileConfDao{}
	confs, err := dao.FindFileConfs(appkey)

	fileConfs := []*FileConf{}
	if err == nil {
		for _, conf := range confs {
			fileConfs = append(fileConfs, &FileConf{
				Channel: conf.Channel,
				Enable:  conf.Enable,
			})
		}
	}
	services.SuccessHttpResp(ctx, map[string]interface{}{
		"file_confs": fileConfs,
	})
}

func SetFileConfSwitch(c *gin.Context) {
	var req struct {
		AppKey  string `json:"app_key"`
		Channel string `json:"channel"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}

	dao := dbs.FileConfDao{}
	err := dao.UpdateEnable(req.AppKey, req.Channel)
	if err != nil {
		services.FailHttpResp(c, services.AdminErrorCode_Default)
		return
	}
	services.SuccessHttpResp(c, nil)
}
