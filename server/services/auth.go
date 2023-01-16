package services

import (
	"context"
	"log"

	pbAuth "svc-todo/pb/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer is the server for authentication
type AuthServer struct {
	pbAuth.UnimplementedAuthServiceServer
	userStore  UserLoginStore
	jwtManager *JWTManager
}

// NewAuthServer returns a new auth server
func NewAuthServer(userStore UserLoginStore, jwtManager *JWTManager) pbAuth.AuthServiceServer {
	return &AuthServer{userStore: userStore, jwtManager: jwtManager}
}

func (server *AuthServer) Login(ctx context.Context, req *pbAuth.LoginRequest) (*pbAuth.LoginResponse, error) {
    log.Println("User login was invoked")

    user, err := server.userStore.Find(req.GetUsername())
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
    }

    if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
        return nil, status.Errorf(codes.NotFound, "incorrect username/password")
    }

    token, err := server.jwtManager.Generate(user)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot generate access token")
    }

    res := &pbAuth.LoginResponse{AccessToken: token}
    return res, nil
}