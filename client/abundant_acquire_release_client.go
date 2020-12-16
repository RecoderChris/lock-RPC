package main

import (
	"../lockrpc"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"net"
	"os"
	"time"
)

func main() {
	// 先建立和服务器的连接的Transport
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9090"))
	if err != nil {
		fmt.Println("Error opening socket:", err)
		os.Exit(1)
	}

	// 创建二进制协议
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()


	// 接口需要context，以便在长操作时用户可以取消RPC调用
	ctx := context.Background()

	client := lockrpc.NewLockServeClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Println("Create client error!")
	}
	defer transport.Close()

	var clientId int32
	clientId = 0

	for i:=0;i<5;i++ {
		// acquire lock
		res, err := client.AcquireLock(ctx,clientId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		clientId = res.Cid
		if res.RetValue == 1{
			fmt.Println("[Client ",clientId, "  ]: Get lock, ready to enter CS!" )
			time.Sleep(time.Second*5)
		} else {
			fmt.Println("[Client ",clientId, "  ]: Get lock failed, CS is busy! Please wait..." )
		}
		// release lock
		res, err = client.ReleaseLock(ctx,clientId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		clientId = res.Cid
		if res.RetValue == 1{
			fmt.Println("[Client ",clientId, "  ]: Free lock, ready to quit CS!" )
		} else if res.RetValue == -2{
			fmt.Println("[Client ",clientId, "  ]: Free lock failed, haven't entered CS!" )
		} else {
			fmt.Println("[Client ",clientId, "  ]: Free lock failed, CS is already free!")
		}
	}
}