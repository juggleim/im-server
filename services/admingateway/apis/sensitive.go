package apis

import (
	"bufio"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

func SensitiveWords(ctx *gin.Context) {
	sizeStr := ctx.Query("size")
	var size int64 = 50
	if sizeStr != "" {
		intVal, err := tools.String2Int64(sizeStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			size = intVal
		}
	}
	pageStr := ctx.Query("page")
	var page int64 = 1
	if pageStr != "" {
		intVal, err := tools.String2Int64(pageStr)
		if err == nil && intVal > 0 {
			page = intVal
		}
	}
	appkey := ctx.Query("app_key")
	services.SetCtxString(ctx, services.CtxKey_AppKey, appkey)
	code, resp, err := services.SyncApiCall(ctx, "qry_sensitive_words", "", tools.RandStr(8), &pbobjs.QrySensitiveWordsReq{
		Page: int32(page),
		Size: int32(size),
	}, func() proto.Message {
		return &pbobjs.QrySensitiveWordsResp{}
	})
	if err != nil {
		services.FailHttpResp(ctx, services.AdminErrorCode_ServerErr, err.Error())
		return
	}
	if code != services.AdminErrorCode_Success {
		services.FailHttpResp(ctx, services.AdminErrorCode(code), "")
		return
	}
	wordsResp := resp.(*pbobjs.QrySensitiveWordsResp)
	res := &QrySensitiveWordsResp{
		Items:      []*SensitiveWord{},
		IsFinished: false,
	}
	for _, senWord := range wordsResp.Words {
		res.Items = append(res.Items, &SensitiveWord{
			Id:       senWord.Id,
			Word:     senWord.Word,
			WordType: int(senWord.WordType),
		})
	}
	if len(res.Items) < int(size) {
		res.IsFinished = true
	}
	res.Total = wordsResp.Total
	tools.SuccessHttpResp(ctx, res)
}

type QrySensitiveWordsResp struct {
	Items      []*SensitiveWord `json:"items"`
	IsFinished bool             `json:"is_finished"`
	Total      int32            `json:"total"`
}
type SensitiveWord struct {
	AppKey   string `json:"app_key,omitempty"`
	Id       string `json:"id,omitempty"`
	Word     string `json:"word"`
	WordType int    `json:"word_type"`
}

func ImportSensitiveWords(ctx *gin.Context) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, err.Error())
		return
	}
	f, err := file.Open()
	defer f.Close()
	if err != nil {
		services.FailHttpResp(ctx, services.AdminErrorCode_ServerErr, err.Error())
		return
	}

	var (
		allWords []*pbobjs.SensitiveWord
	)
	reader := bufio.NewReader(f)
	for {
		bs, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			services.FailHttpResp(ctx, services.AdminErrorCode_ServerErr, err.Error())
			return
		}
		allWords = append(allWords, &pbobjs.SensitiveWord{
			Word:     string(bs),
			WordType: pbobjs.SensitiveWordType_replace_word,
		})
	}

	rpcReq := &pbobjs.AddSensitiveWordsReq{
		Words: allWords,
	}

	services.SyncApiCall(ctx, "add_sensitive_words", "", appKey, rpcReq, nil)

	services.SuccessHttpResp(ctx, nil)
}

func AddSensitiveWord(ctx *gin.Context) {
	var req SensitiveWord
	if err := ctx.ShouldBindJSON(&req); err != nil {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}

	rpcReq := &pbobjs.AddSensitiveWordsReq{
		Words: []*pbobjs.SensitiveWord{
			{
				Word:     req.Word,
				WordType: pbobjs.SensitiveWordType(req.WordType),
			},
		},
	}
	appKey := req.AppKey
	services.SetCtxString(ctx, services.CtxKey_AppKey, appKey)
	services.SyncApiCall(ctx, "add_sensitive_words", "", appKey, rpcReq, nil)

	services.SuccessHttpResp(ctx, nil)
}

func DeleteSensitiveWord(ctx *gin.Context) {
	var req SensitiveWord
	if err := ctx.ShouldBindJSON(&req); err != nil {
		services.FailHttpResp(ctx, services.AdminErrorCode_ParamError, "param illegal")
		return
	}
	appKey := req.AppKey
	rpcReq := &pbobjs.DelSensitiveWordsReq{
		Words: []string{req.Word},
	}
	services.SetCtxString(ctx, services.CtxKey_AppKey, appKey)
	services.SyncApiCall(ctx, "del_sensitive_words", "", appKey, rpcReq, nil)

	services.SuccessHttpResp(ctx, nil)
}
