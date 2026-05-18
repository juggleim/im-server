package server

import (
	"errors"
	"fmt"
	"im-server/commons/errs"
	"im-server/commons/gmicro/utils"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

func (server *ImWebsocketServer) ImBotServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connect upgrade:", err)
		return
	}
	clientIp := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if clientIp == "" {
		clientIp = conn.RemoteAddr().String()
	}
	clientHost := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if clientHost == "" {
		clientHost = strings.TrimSpace(r.Host)
	}
	child := &ImBotWebsocketChild{
		stopChan:         make(chan bool, 1),
		wsConn:           conn,
		isActive:         true,
		messageListener:  server.BotMsgListener,
		latestActiveTime: time.Now().UnixMilli(),
	}
	utils.SafeGo(func() {
		child.startWsListener(clientIp, clientHost)
	})
}

type ImBotWebsocketChild struct {
	stopChan         chan bool
	wsConn           *websocket.Conn
	isActive         bool
	messageListener  ImListener
	latestActiveTime int64
	ticker           *time.Ticker
}

func (child *ImBotWebsocketChild) startWsListener(clientIp, clientHost string) {
	handler := ImBotWebsocketMsgHandler{child.messageListener}
	ctx := &BotWsHandleContextImpl{
		conn:       child.wsConn,
		wsChild:    child,
		lock:       &sync.RWMutex{},
		attachment: &sync.Map{},
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ConnectSession, tools.GenerateUUIDShort11())
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ConnectCreateTime, time.Now().UnixMilli())
	imcontext.SetContextAttr(ctx, imcontext.StateKey_CtxLocker, &sync.Mutex{})
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Limiter, rate.NewLimiter(100, 100))
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientIp, clientIp)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientHost, clientHost)

	child.startTicker(ctx, handler)

	for child.isActive {
		_, message, err := child.wsConn.ReadMessage()
		child.latestActiveTime = time.Now().UnixMilli()

		if err != nil {
			if child.isActive {
				child.Stop()
				handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_NET_ERR, err)
			}
			break
		}

		wsMsg := &codec.ImWebsocketMsg{}
		err = tools.PbUnMarshal(message, wsMsg)
		if err != nil {
			fmt.Println("failed to decode pb data:", err)
			child.Stop()
			handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_PB_DECODE_FAIL, err)
			break
		}

		handler.HandleRead(ctx, wsMsg)
	}
}

func (child *ImBotWebsocketChild) startTicker(ctx imcontext.WsHandleContext, handler ImBotWebsocketMsgHandler) {
	if child.ticker == nil {
		child.ticker = time.NewTicker(5 * time.Second)
	} else {
		child.ticker.Reset(5 * time.Second)
	}
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				current := time.Now().UnixMilli()
				interval := current - child.latestActiveTime
				if interval > 300*1000 {
					child.Stop()
					handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_HEARTBEAT_TIMEOUT, errors.New("bot inactive more than 5min"))
					return
				}
			case <-child.stopChan:
				return
			}
		}
	}(child.ticker)
}

func (child *ImBotWebsocketChild) Stop() {
	if child.isActive {
		child.isActive = false
		child.stopChan <- true
		if child.wsConn != nil {
			child.wsConn.Close()
		}
		close(child.stopChan)
	}
}

type BotWsHandleContextImpl struct {
	wsChild    *ImBotWebsocketChild
	conn       *websocket.Conn
	attachment interface{}
	lock       *sync.RWMutex
}

func (ctx *BotWsHandleContextImpl) Write(message interface{}) {
	imMsg, ok := message.(codec.IMessage)
	if ok {
		wsImMsg := imMsg.ToImWebsocketMsg()
		bs, err := tools.PbMarshal(wsImMsg)
		if err == nil {
			ctx.lock.Lock()
			defer ctx.lock.Unlock()
			_ = ctx.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			ctx.conn.WriteMessage(websocket.BinaryMessage, bs)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("No IMessage to transfer to WebsocketMsg.")
	}
}

func (ctx *BotWsHandleContextImpl) Close(err error) {
	if ctx.wsChild != nil {
		ctx.wsChild.Stop()
	}
}

func (ctx *BotWsHandleContextImpl) Attachment() imcontext.Attachment {
	return ctx.attachment
}

func (ctx *BotWsHandleContextImpl) SetAttachment(attachment imcontext.Attachment) {
	ctx.attachment = attachment
}

func (ctx *BotWsHandleContextImpl) IsActive() bool {
	return ctx.wsChild.isActive
}

func (ctx *BotWsHandleContextImpl) RemoteAddr() string {
	if ctx.conn != nil {
		return ctx.conn.RemoteAddr().String()
	}
	return ""
}
