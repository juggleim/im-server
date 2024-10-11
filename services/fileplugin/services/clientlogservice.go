package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
)

func ReportClientLogState(ctx context.Context, req *pbobjs.UploadLogStatusReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.ClientLogDao{}
	err := dao.UpdateLogUrl(appkey, req.MsgId, req.LogUrl, dbs.ClientLogState(req.State))
	if err != nil {
		logs.WithContext(ctx).Errorf("report client log state failed. err:%v", err)
	}
	return errs.IMErrorCode_SUCCESS
}
