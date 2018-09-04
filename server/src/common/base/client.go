package base

import (
	"google.golang.org/grpc"
	"sync"
	"time"
	"math/rand"
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	maxConnNum        = 3
	clientDialTimeout = 30 * time.Second
)

var (
	connPool map[string][]*grpc.ClientConn
	lock     sync.RWMutex
)

//noinspection ALL
type BaseClient struct {
	ServiceName string
	ip          string
	port        string
	grpcClient  interface{}
	conn        *grpc.ClientConn
}

func NewBaseClent(name, ip, port string, grpcClientConstructor func(conn *grpc.ClientConn) interface{}) *BaseClient {

	client := &BaseClient{
		ServiceName: name,
		ip:          ip,
		port:        port,
	}

	client.conn = client.getConn()
	client.grpcClient = grpcClientConstructor(client.conn)

	return client

}
func (c *BaseClient) WithClientFunc(ctx context.Context,clientFunc func(
	clientInterface interface{},
	ctx context.Context,
	opts ...grpc.CallOption,
) (interface{}, error)) (interface{}, error) {

	return c.WithClientFuncTimeout(ctx,clientFunc, clientDialTimeout)
}

func (c *BaseClient) WithClientFuncTimeout(ctx context.Context,clientFunc func(
	clientInterface interface{},
	ctx context.Context,
	opts ...grpc.CallOption,
) (interface{}, error), timeout time.Duration) (interface{}, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var trailer metadata.MD
	return c.unaryCallWithContext(ctx, clientFunc, grpc.Trailer(&trailer))

}

func (c *BaseClient) unaryCallWithContext(ctx context.Context, clientFunc func(clientInterface interface{}, ctx context.Context, opts ...grpc.CallOption) (interface{}, error), opts ...grpc.CallOption) (interface{}, error) {
	resp, err := clientFunc(c.grpcClient, ctx, opts...)
	return resp, err
}

func (c *BaseClient) getConn() *grpc.ClientConn {
	lock.Lock()
	defer lock.Unlock()

	target := c.getTargetName()

	_, ok := connPool[target]
	if !ok {
		connPool[target] = make([]*grpc.ClientConn, 0)
	}
	pool := connPool[target]

	index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(maxConnNum)

	if len(pool) < maxConnNum {

		ctx, _ := context.WithTimeout(context.Background(), clientDialTimeout)

		dialOpts := []grpc.DialOption{
			grpc.WithBlock(),
		}

		//todo etcd

		conn, err := grpc.DialContext(ctx, target, dialOpts...)
		if err != nil {
			return nil
		}
		pool = append(pool, conn)
		return conn
	}

	return pool[index]

}

func (c *BaseClient) getTargetName() string {
	//todo etcd
	return c.ip + ":" + c.port
}
