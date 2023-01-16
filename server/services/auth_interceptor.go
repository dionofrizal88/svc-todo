package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"svc-todo/server/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
    jwtManager      *JWTManager
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
    return &AuthInterceptor{jwtManager}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        log.Println("Unary interceptor in endpoint:", info.FullMethod)

        err := interceptor.authorize(ctx, info.FullMethod)
        if err != nil {
            return nil, err
        }

        // add limiter in interceptor
        err = redisLimitter()
        if err != nil {
            return nil, err
        }

        return handler(ctx, req)
    }
}

func redisLimitter() error{
    limiterStatus, msg, err := db.CreateLimitter()
    if err != nil {
        return err
    }
    if limiterStatus == false{
        newMsg := fmt.Sprintf("To many request: %s", msg)
        return status.Errorf(codes.Unauthenticated, newMsg)
    }
    return nil
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {

    // exclue method not use authorize
    const authServicePath = "/svc_todo.AuthService/"

    excludeMethod := []string{
        authServicePath + "Login",
    }

    for _, excludeMethod := range excludeMethod {
        if method == excludeMethod {
            return nil
        }
    }

    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return status.Errorf(codes.Unauthenticated, "metadata is not provided")
    }

    values := md["authorization"]
    if len(values) == 0 {
        return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
    }

    accessToken := strings.Replace(values[0], "Bearer ", "", 1)
    _, err := interceptor.jwtManager.Verify(accessToken)
    if err != nil {
        return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
    }

    return nil
}
