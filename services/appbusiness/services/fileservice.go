package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"

	"google.golang.org/protobuf/proto"
)

func GetFileCred(ctx context.Context, req *pbobjs.QryFileCredReq) (errs.IMErrorCode, *pbobjs.QryFileCredResp) {
	code, respObj, err := bases.SyncRpcCall(ctx, "file_cred", bases.GetRequesterIdFromCtx(ctx), req, func() proto.Message {
		return &pbobjs.QryFileCredResp{}
	})
	ret := &pbobjs.QryFileCredResp{}
	if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
		ret = respObj.(*pbobjs.QryFileCredResp)
	}
	return errs.IMErrorCode_SUCCESS, ret
}
