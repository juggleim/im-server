package actorsystem

import (
	context "context"
	"fmt"
	"sync"
	"time"

	"im-server/commons/gmicro/actorsystem/rpc"
	"im-server/commons/gmicro/logs"

	grpc "google.golang.org/grpc"
)

type RpcClient struct {
	Address  string
	connPool sync.Pool
	// conn        *grpc.ClientConn
	// msgClient   rpc.RpcMessageClient
	// isConnected bool
}

func NewRpcClient(address string) *RpcClient {
	client := &RpcClient{
		Address: address,
		connPool: sync.Pool{
			New: func() interface{} {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientLogger))
				if err != nil {
					logs.Error("Can not create grpc connect. address:%v\terr:", address, err)
				}
				return conn
			},
		},
	}

	return client
}

// func (client *RpcClient) connect() {
// 	tmpConn, err := grpc.Dial(client.Address, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	} else {
// 		client.isConnected = true
// 	}
// 	client.conn = tmpConn
// 	client.msgClient = rpc.NewRpcMessageClient(tmpConn)
// }

func (client *RpcClient) DisConnect() {
	//client.conn.Close()
}

func (client *RpcClient) Send(req *rpc.RpcMessageRequest) {
	conn := client.connPool.Get().(*grpc.ClientConn)
	defer client.connPool.Put(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcClient := rpc.NewRpcMessageClient(conn)
	resp, err := grpcClient.Send(ctx, req)
	if err != nil {
		logs.Error("resp:%v\terr:", resp, err)
	}
}

func clientLogger(ctx context.Context, method string, req, reply interface{}, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	begin := time.Now()
	err := invoker(ctx, method, req, reply, conn, opts...)
	during := time.Since(begin).Milliseconds()
	logs.Debug(fmt.Sprintf("[%s] method=%v err=%v req=%v reqply=%v %v\n", begin.Format(time.RFC3339), method, err, req, reply, during))
	return err
}
