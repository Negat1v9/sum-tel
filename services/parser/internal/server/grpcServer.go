package server

import (
	"context"

	parserv1 "github.com/Negat1v9/sum-tel/services/parser/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/services/parser/internal/parser/service"
	"google.golang.org/grpc"
)

type ParserGRPCServer struct {
	parserv1.UnimplementedTgParserServer

	parserService *service.ParserService
}

func Register(gRPCServer *grpc.Server, msgService *service.ParserService) {
	parserv1.RegisterTgParserServer(gRPCServer, &ParserGRPCServer{parserService: msgService})
}

func (s *ParserGRPCServer) ParseNewChannel(ctx context.Context, req *parserv1.NewChannelRequest) (*parserv1.NewChannelResponse, error) {
	return s.parserService.ParseNewChannel(ctx, req.ChannelID, req.Username)
}

func (s *ParserGRPCServer) ParseMessages(ctx context.Context, req *parserv1.ParseMessagesRequest) (*parserv1.ParseMessagesResponse, error) {
	return s.parserService.ParseMessages(ctx, req.ChannelID, req.Username)
}
