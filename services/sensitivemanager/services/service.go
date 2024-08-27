package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/sensitivemanager/dbs"
)

func QrySensitiveWords(ctx context.Context, req *pbobjs.QrySensitiveWordsReq) (errs.IMErrorCode, *pbobjs.QrySensitiveWordsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.SensitiveWordDao{}
	resp := &pbobjs.QrySensitiveWordsResp{
		Words: []*pbobjs.SensitiveWord{},
	}
	if req.Size > 0 { //page size
		list, total, err := dao.QrySensitiveWordsWithPage(appkey, int64(req.Page), int64(req.Size), req.Word, req.WordType)
		if err != nil {
			return errs.IMErrorCode_API_INTERNAL_RESP_FAIL, resp
		}
		for _, word := range list {
			idStr, _ := tools.EncodeInt(word.ID)
			resp.Words = append(resp.Words, &pbobjs.SensitiveWord{
				Id:       idStr,
				Word:     word.Word,
				WordType: pbobjs.SensitiveWordType(word.WordType),
			})
		}
		resp.Total = int32(total)
	} else { //limit offset
		var startId int64 = 0
		if req.Offset != "" {
			intVal, err := tools.DecodeInt(req.Offset)
			if err == nil {
				startId = intVal
			}
		}
		list, err := dao.QrySensitiveWords(appkey, int64(req.Limit), startId)
		if err != nil {
			return errs.IMErrorCode_API_INTERNAL_RESP_FAIL, resp
		}
		for _, word := range list {
			idStr, _ := tools.EncodeInt(word.ID)
			resp.Words = append(resp.Words, &pbobjs.SensitiveWord{
				Id:       idStr,
				Word:     word.Word,
				WordType: pbobjs.SensitiveWordType(word.WordType),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}
