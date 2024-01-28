package grpcconn

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	favPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite/favPb"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/jwtauth"
	"context"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type (
	GrpcClientFactoryHandler interface {
		Users() usersPb.UsersGrpcServiceClient
		FavAndCom() favPb.FavGrpcServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
		cfg       *config.Jwt
	}
)

func (g *grpcClientFactory) Users() usersPb.UsersGrpcServiceClient {
	return usersPb.NewUsersGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) FavAndCom() favPb.FavGrpcServiceClient {
	return favPb.NewFavGrpcServiceClient(g.client)
}

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
	opts := make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Printf("Error: Grpc Clinet Connection Failed %s", err.Error())
		return nil, errors.New("error: grpc client connection failed")
	}

	return &grpcClientFactory{
		client: clientConn,
	}, err
}

func (g *grpcAuth) unaryAuthorization(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		log.Printf("Error: Metadata Not Found")
		return nil, errors.New("error: metadata is not found")
	}

	authHeader, ok := md["auth"]
	if !ok {
		log.Printf("Error: Metadata Not Found")
		return nil, errors.New("error: metadata is not found")
	}

	if len(authHeader) == 0 {
		log.Printf("Error: Metadata Not Found")
		return nil, errors.New("error: metadata is not found")
	}

	cliams, err := jwtauth.ParseToken(string(authHeader[0]), g.cfg)
	if err != nil {
		log.Printf("Error: Parse Token Failed %s", err.Error())
		return nil, errors.New("error: token is invalid")
	}

	log.Printf("Cliams %v", cliams)
	return handler(ctx, req)
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {

	opts := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{
		secretKey: cfg.ApiSecretKey,
		cfg:       cfg,
	}

	opts = append(opts, grpc.UnaryInterceptor(grpcAuth.unaryAuthorization))

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error: failed to listen %v", err)
	}

	return grpcServer, lis
}
