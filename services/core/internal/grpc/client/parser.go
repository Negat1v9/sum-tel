package grpcclient

import (
	"context"
	"time"

	parserv1 "github.com/Negat1v9/sum-tel/services/core/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type TgParserClient struct {
	address string
	client  parserv1.TgParserClient
	log     *logger.Logger
	conn    *grpc.ClientConn

	closed bool
}

func NewTgParserClient(grpcURL string, log *logger.Logger) (*TgParserClient, error) {
	p := &TgParserClient{
		address: grpcURL,
		log:     log,
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	// start healthcheck in background
	go p.healthcheck()

	return p, nil
}

func (c *TgParserClient) ParseNewChannel(ctx context.Context, req *parserv1.NewChannelRequest) (*parserv1.NewChannelResponse, error) {
	return c.client.ParseNewChannel(ctx, req)
}

func (c *TgParserClient) ParseMessages(ctx context.Context, req *parserv1.ParseMessagesRequest) (*parserv1.ParseMessagesResponse, error) {
	return c.client.ParseMessages(ctx, req)
}

func (c *TgParserClient) GetNewsSources(ctx context.Context, req *parserv1.NewsSourcesRequest) (*parserv1.NewsSourcesResponse, error) {
	return c.client.NewsSources(ctx, req)

}

// Close closes the gRPC connection
func (c *TgParserClient) Close() error {
	c.closed = true
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// connect establishes a gRPC connection to the server
func (c *TgParserClient) connect() error {
	conn, err := grpc.NewClient(
		c.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = parserv1.NewTgParserClient(conn)
	c.closed = false

	return nil
}

// healthcheck periodically checks the health of the gRPC connection and reconnects if necessary
func (c *TgParserClient) healthcheck() {
	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {

		if c.closed {
			return
		}

		if !c.isConnectionHealthy() {
			c.log.Warnf("grpc_client.healthcheck: connection is not healthy, reconnecting")
			if err := c.connect(); err != nil {
				c.log.Errorf("grpc_client.connect: failed to reconnect: %v", err)
			}
		}
	}
}

// isConnectionHealthy checks if the gRPC connection is healthy
func (rc *TgParserClient) isConnectionHealthy() bool {

	if rc.conn == nil {
		return false
	}

	state := rc.conn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}
