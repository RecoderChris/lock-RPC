package main

import (
	"../lockrpc"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"net"
	"os"
)

var lock int32 = 1
var lockOwner int32 = 0
var cid int32 = 0
type rpcService struct{}

func (s *rpcService) AcquireLock(_ context.Context,clientId int32) (r *lockrpc.RetType, err error) {
	var ret lockrpc.RetType
	if clientId == 0{
		cid = cid + 1
		clientId = cid
		ret.Cid = cid
	}
	if lock == 1{
		fmt.Println("[Server  ]: Client ", clientId, " get lock...")
		lock = 0
		lockOwner = clientId
		ret.RetValue = 1
		ret.Cid = clientId
		return &ret,nil
	} else{
		fmt.Println("[Server  ]: Client ", clientId, " get lock failed! The lockOwner is Client ",lockOwner,"...")
		ret.RetValue = -1
		ret.Cid = clientId
		return &ret, nil
	}
}

func (s *rpcService) ReleaseLock(_ context.Context, clientId int32) (r *lockrpc.RetType, err error) {
	var ret lockrpc.RetType
	if clientId == 0{
		cid = cid + 1
		clientId = cid
		ret.Cid = cid
	}
	if lock == 0 && lockOwner == clientId {
		fmt.Println("[Server  ]: Client ", clientId, " free lock...")
		lock = 1
		lockOwner = 0
		ret.RetValue = 1
		ret.Cid = clientId
		return &ret,nil
	}else if lock == 1{
		fmt.Println("[Server  ]: Client ", clientId, " try to do a dummy free...")
		ret.RetValue = -1
		ret.Cid = clientId
		return &ret,nil
	}else{
		fmt.Println("[Server  ]: Client ", clientId,  " free lock failed! The lockOwner is Client ",lockOwner,"...")
		ret.RetValue = -2
		ret.Cid = clientId
		return &ret,nil
	}
}

func main() {
	// 创建服务器
	serverTransport, err := thrift.NewTServerSocket(net.JoinHostPort("127.0.0.1", "9090"))
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	// 创建Processor，用一个端口处理多个服务
	handler := &rpcService{}
	processor := lockrpc.NewLockServeProcessor(handler)
	server := thrift.NewTSimpleServer2(processor, serverTransport)

	fmt.Println("[Server  ]: Copyright@UCAS@ICT-Xinmiao Zhang")
	fmt.Println("[Server  ]: Hello world! ")
	if err := server.Serve(); err != nil {
		panic(err)
	}

}