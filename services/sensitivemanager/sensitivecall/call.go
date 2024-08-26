package sensitivecall

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"

	"google.golang.org/protobuf/proto"
)

func FilterCall(ctx context.Context, text string) (res *pbobjs.SensitiveFilterResp, errCode errs.IMErrorCode, err error) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	var result proto.Message
	errCode, result, err = bases.SyncRpcCall(ctx, "sensitive_filter_text", appKey, &pbobjs.SensitiveFilterReq{
		Text: text,
	}, func() proto.Message {
		return &pbobjs.SensitiveFilterResp{}
	})
	if err == nil && errCode == errs.IMErrorCode_SUCCESS {
		res = result.(*pbobjs.SensitiveFilterResp)
	}

	return
}
