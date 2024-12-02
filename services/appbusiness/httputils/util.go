package httputils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/models"
	"io"
	"net/http"
	"net/url"
)

type HttpContext struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	QueryParams url.Values

	AppKey        string
	CurrentUserId string
}

func (ctx *HttpContext) BindJson(req interface{}) error {
	return Body2Obj(ctx.Request.Body, req)
}

func (ctx *HttpContext) Query(key string) string {
	if ctx.QueryParams != nil {
		return ctx.QueryParams.Get(key)
	}
	return ""
}

func (ctx *HttpContext) ResponseErr(code errs.IMErrorCode) {
	appErr := errs.GetApiErrorByCode(code)
	ctx.Writer.WriteHeader(appErr.HttpCode)
	bs, _ := tools.JsonMarshal(appErr)
	ctx.Writer.Write(bs)
}

func (ctx *HttpContext) ResponseSucc(resp interface{}) {
	connonResp := &models.CommonResp{
		CommonError: models.CommonError{
			ErrorMsg: "success",
		},
		Data: resp,
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	bs, _ := tools.JsonMarshal(connonResp)
	ctx.Writer.Write(bs)
}

func (ctx *HttpContext) ToRpcCtx(userId string) context.Context {
	rpcCtx := context.Background()
	rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_AppKey, ctx.AppKey)
	rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_Session, fmt.Sprintf("app_%s", tools.GenerateUUIDShort11()))
	rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_IsFromApp, true)
	if userId != "" {
		rpcCtx = context.WithValue(rpcCtx, bases.CtxKey_RequesterId, userId)
	}
	return rpcCtx
}

func Read2String(read io.ReadCloser) string {
	buf := bytes.NewBuffer([]byte{})
	for {
		bs := make([]byte, 1024)
		c, err := read.Read(bs)
		buf.Write(bs)
		if err != nil || c < 1024 {
			break
		}
	}
	return buf.String()
}

func Read2Bytes(read io.ReadCloser) []byte {
	buf := bytes.NewBuffer([]byte{})
	for {
		bs := make([]byte, 1024)
		c, err := read.Read(bs)
		if err != nil || c < 1024 {
			if c > 0 {
				buf.Write(bs[:c])
			}
			break
		}
		buf.Write(bs)
	}
	return buf.Bytes()
}

func Body2Obj(read io.ReadCloser, obj interface{}) error {
	bs := Read2Bytes(read)
	if len(bs) <= 0 {
		return errors.New("no value")
	}
	return tools.JsonUnMarshal(bs, obj)
}
