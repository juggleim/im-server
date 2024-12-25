package apis

import (
	"bufio"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/services"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

func QrySensitiveWords(ctx *gin.Context) {
	sizeStr := ctx.Query("size")
	var size int64 = 50
	if sizeStr != "" {
		intVal, err := tools.String2Int64(sizeStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			size = intVal
		}
	}

	var page int64
	pageStr := ctx.Query("page")
	if pageStr != "" {
		intVal, err := tools.String2Int64(pageStr)
		if err == nil {
			page = intVal
		}
	}

	word := ctx.Query("word")

	var wordType int64
	wordTypeStr := ctx.Query("word_type")
	if wordTypeStr != "" {
		wordType, _ = tools.String2Int64(wordTypeStr)
	}

	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_sensitive_words", tools.RandStr(8), &pbobjs.QrySensitiveWordsReq{
		Size:     int32(size),
		Page:     int32(page),
		Word:     word,
		WordType: int32(wordType),
	}, func() proto.Message {
		return &pbobjs.QrySensitiveWordsResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	wordsResp := resp.(*pbobjs.QrySensitiveWordsResp)
	res := &SensitiveWords{
		Items:      []*SensitiveWord{},
		IsFinished: false,
		Total:      wordsResp.Total,
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
	tools.SuccessHttpResp(ctx, res)
}

type SensitiveWords struct {
	Items      []*SensitiveWord `json:"items"`
	Total      int32            `json:"total"`
	IsFinished bool             `json:"is_finished"`
}
type SensitiveWord struct {
	Id       string `json:"id"`
	Word     string `json:"word"`
	WordType int    `json:"word_type"`
}

func ImportSensitiveWords(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_MSG_PARAM_ILLEGAL)
		return
	}
	f, err := file.Open()
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	defer f.Close()

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
			tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
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
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "add_sensitive_words", tools.RandStr(8), rpcReq)
	tools.SuccessHttpResp(ctx, nil)
}

func AddSensitiveWords(ctx *gin.Context) {
	var req SensitiveWords
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_MSG_PARAM_ILLEGAL)
		return
	}
	rpcReq := &pbobjs.AddSensitiveWordsReq{
		Words: []*pbobjs.SensitiveWord{},
	}
	for _, word := range req.Items {
		rpcReq.Words = append(rpcReq.Words, &pbobjs.SensitiveWord{
			Word:     word.Word,
			WordType: pbobjs.SensitiveWordType(word.WordType),
		})
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "add_sensitive_words", tools.RandStr(8), rpcReq)
	tools.SuccessHttpResp(ctx, nil)
}

func DeleteSensitiveWords(ctx *gin.Context) {
	var req DelSensitiveWordsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_MSG_PARAM_ILLEGAL)
		return
	}

	rpcReq := &pbobjs.DelSensitiveWordsReq{
		Words: req.Words,
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "del_sensitive_words", tools.RandStr(8), rpcReq)
	tools.SuccessHttpResp(ctx, nil)
}

type DelSensitiveWordsReq struct {
	Words []string `json:"words"`
}
