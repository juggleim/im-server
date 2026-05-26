package server

import (
	"errors"
	"fmt"
	"im-server/commons/errs"
	"im-server/commons/gmicro/utils"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"im-server/services/connectmanager/server/imhttpmsghandlers"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type ImWebsocketServer struct {
	MessageListener ImListener
	BotMsgListener  ImListener
}

func (server *ImWebsocketServer) AsyncStart(port int) {
	var mux *http.ServeMux = commonservices.GetDefaultHttpServeMux()
	mux.HandleFunc("/im", server.ImWsServer)
	mux.HandleFunc("/imbot", server.ImBotServer)
	mux.HandleFunc("/im/publish", imhttpmsghandlers.ImHttpPubHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})
	go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4196,
	WriteBufferSize: 1124,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (server *ImWebsocketServer) ImWsServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connect upgrade:", err)
		return
	}
	referer := strings.TrimSpace(r.Header.Get("Origin"))
	if referer == "" {
		referer = strings.TrimSpace(r.Header.Get("Referer"))
	}
	clientIp := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if clientIp == "" {
		clientIp = conn.RemoteAddr().String()
	}
	clientHost := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if clientHost == "" {
		clientHost = strings.TrimSpace(r.Host)
	}
	child := &ImWebsocketChild{
		stopChan:        make(chan struct{}),
		wsConn:          conn,
		messageListener: server.MessageListener,
	}
	child.isActive.Store(true)
	child.latestActiveTime.Store(time.Now().UnixMilli())

	utils.SafeGo(func() {
		child.startWsListener(referer, clientIp, clientHost)
	})
}

func (server *ImWebsocketServer) Stop() {}

type ImWebsocketChild struct {
	stopChan         chan struct{}
	wsConn           *websocket.Conn
	isActive         atomic.Bool
	messageListener  ImListener
	latestActiveTime atomic.Int64
	ticker           *time.Ticker
	stopOnce         sync.Once
}

func (child *ImWebsocketChild) startWsListener(referer, clientIp, clientHost string) {
	handler := IMWebsocketMsgHandler{child.messageListener}
	ctx := &WsHandleContextImpl{
		conn:       child.wsConn,
		wsChild:    child,
		lock:       &sync.RWMutex{},
		attachment: &sync.Map{},
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ConnectSession, tools.GenerateUUIDShort11())
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ConnectCreateTime, time.Now().UnixMilli())
	imcontext.SetContextAttr(ctx, imcontext.StateKey_CtxLocker, &sync.Mutex{})
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Limiter, rate.NewLimiter(100, 10))
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Referer, referer)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientIp, clientIp)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientHost, clientHost)

	// 设置读超时，防止半开连接永久阻塞
	child.wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
	child.wsConn.SetPongHandler(func(string) error {
		child.wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
		child.latestActiveTime.Store(time.Now().UnixMilli())
		return nil
	})

	// start ticker
	child.startTicker(ctx, handler)

	for child.isActive.Load() {
		_, message, err := child.wsConn.ReadMessage()
		// record
		child.latestActiveTime.Store(time.Now().UnixMilli())

		if err != nil {
			if child.isActive.Load() {
				child.Stop()
				handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_NET_ERR, err)
			}
			break
		}

		// 重置读超时，防止下一次 ReadMessage 因 deadline 过期立即报错
		child.wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// decode
		wsMsg := &codec.ImWebsocketMsg{}
		err = tools.PbUnMarshal(message, wsMsg)
		if err != nil {
			fmt.Println("failed to decode pb data:", err)
			child.Stop()
			handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_PB_DECODE_FAIL, err)
			break
		}

		// decrypt
		wsMsg.Decrypt(ctx)

		handler.HandleRead(ctx, wsMsg)
	}
}

func (child *ImWebsocketChild) startTicker(ctx imcontext.WsHandleContext, handler IMWebsocketMsgHandler) {
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
				if !child.isActive.Load() {
					return
				}
				// 发送 Ping 检测连接是否存活
				if err := child.wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
					child.Stop()
					handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_HEARTBEAT_TIMEOUT, err)
					return
				}
				current := time.Now().UnixMilli()
				interval := current - child.latestActiveTime.Load()
				if interval > 300*1000 {
					child.Stop()
					handler.HandleException(ctx, errs.IMErrorCode_CONNECT_CLOSE_HEARTBEAT_TIMEOUT, errors.New("user inactive more than 5min"))
					return
				}
			case <-child.stopChan:
				return
			}
		}
	}(child.ticker)
}

func (child *ImWebsocketChild) Stop() {
	child.stopOnce.Do(func() {
		child.isActive.Store(false)
		if child.wsConn != nil {
			child.wsConn.Close()
		}
		close(child.stopChan)
	})
}

type WsHandleContextImpl struct {
	wsChild    *ImWebsocketChild
	conn       *websocket.Conn
	attachment interface{}
	lock       *sync.RWMutex
}

func (ctx *WsHandleContextImpl) Write(message interface{}) {
	imMsg, ok := message.(codec.IMessage)
	if ok {
		wsImMsg := imMsg.ToImWebsocketMsg()
		// encrypt
		wsImMsg.Encrypt(ctx)
		bs, err := tools.PbMarshal(wsImMsg)
		if err == nil {
			ctx.lock.Lock()
			defer ctx.lock.Unlock()
			_ = ctx.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err = ctx.conn.WriteMessage(websocket.BinaryMessage, bs)
			if err != nil {
				fmt.Println("write result:", err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("No IMessage to transfer to WebsocketMsg.")
	}
}

func (ctx *WsHandleContextImpl) Close(err error) {
	if ctx.wsChild != nil {
		ctx.wsChild.Stop()
	}
}
func (ctx *WsHandleContextImpl) Attachment() imcontext.Attachment {
	return ctx.attachment
}
func (ctx *WsHandleContextImpl) SetAttachment(attachment imcontext.Attachment) {
	ctx.attachment = attachment
}
func (ctx *WsHandleContextImpl) IsActive() bool {
	return ctx.wsChild.isActive.Load()
}
func (ctx *WsHandleContextImpl) RemoteAddr() string {
	if ctx.conn != nil {
		return ctx.conn.RemoteAddr().String()
	}
	return ""
}
