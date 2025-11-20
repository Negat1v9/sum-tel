package server

import (
	"context"
	"log/slog"

	parserv1 "github.com/Negat1v9/sum-tel/services/parser/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/services/parser/internal/parser/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ParserGRPCServer struct {
	parserv1.UnimplementedTgParserServer

	parserService *service.ParserService
}

func Register(gRPCServer *grpc.Server, msgService *service.ParserService) {
	parserv1.RegisterTgParserServer(gRPCServer, &ParserGRPCServer{parserService: msgService})
}

// handler parse new telegram channel information
func (s *ParserGRPCServer) ParseNewChannel(ctx context.Context, req *parserv1.NewChannelRequest) (*parserv1.NewChannelResponse, error) {
	r, err := s.parserService.ParseNewChannel(ctx, req.GetChannelID(), req.GetUsername())
	if err != nil {
		slog.Error("ParseNewChannel:", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "failed parse new channel")
	}
	return r, nil
}

// handler parse new telegram channel latest messages
func (s *ParserGRPCServer) ParseMessages(ctx context.Context, req *parserv1.ParseMessagesRequest) (*parserv1.ParseMessagesResponse, error) {
	r, err := s.parserService.ParseMessages(ctx, req.GetChannelID(), req.GetUsername())
	if err != nil {
		slog.Error("ParseMessages:", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, "failed parse channel messages")
	}
	return r, nil
}
