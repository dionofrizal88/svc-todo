package main

import (
	"crypto/tls"
	"log"
	"net"
	pbAuth "svc-todo/pb/auth"
	pbTodo "svc-todo/pb/todo"
	pbUser "svc-todo/pb/user"
	"svc-todo/server/db"
	"svc-todo/server/services"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
    secretKey     = "secret"
    tokenDuration = 15 * time.Minute
)

var addr string = "0.0.0.0:50051"

func main(){

	// implement TLS Credentials
	// tlsCredentials, err := loadTLSCredentials()
    // if err != nil {
    //     log.Fatal("cannot load TLS credentials: ", err)
    // }

	lis, err := net.Listen("tcp", addr)
	if err != nil{
		log.Fatalf("Failed to listen on; %v\n", err)
	}

	log.Printf("Listening on %s\n", addr)

	// connect db mysql and redis
	db.Init()
	db.RedisClient()

	// create seed user
	userLoginStore := services.NewInMemoryUserLoginStore()
    err = seedUsers(userLoginStore)
    if err != nil {
        log.Fatal("cannot seed users: ", err)
    }

	// implement jwt in server
	jwtManager := services.NewJWTManager(secretKey, tokenDuration)

	// inisiasi auth interceptor
	interceptor := services.NewAuthInterceptor(jwtManager)

	// add interceptor to grpc server
	s := grpc.NewServer(
        grpc.UnaryInterceptor(interceptor.Unary()),
    )

	// s := grpc.NewServer(grpc.Creds(tlsCredentials),
    //     grpc.UnaryInterceptor(unaryInterceptor))

	// s := grpc.NewServer()

    authServer := services.NewAuthServer(userLoginStore, jwtManager)
    pbAuth.RegisterAuthServiceServer(s, authServer)

	userService := services.UserService{}
	pbUser.RegisterUserServiceServer(s, &userService)

	todoService := services.TodoService{}
	pbTodo.RegisterTodoServiceServer(s, &todoService)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	
	if err = s.Serve(lis); err != nil{
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
    // Load server's certificate and private key
    serverCert, err := tls.LoadX509KeyPair("cert/ca-cert.pem", "cert/ca-key.pem")
    if err != nil {
        return nil, err
    }

    // Create the credentials and return it
    config := &tls.Config{
        Certificates: []tls.Certificate{serverCert},
        ClientAuth:   tls.NoClientCert,
    }

    return credentials.NewTLS(config), nil
}



// User login
func createUser(userStore services.UserLoginStore, username, password, role string) error {
    user, err := services.NewUser(username, password, role)
    if err != nil {
        return err
    }
    return userStore.Save(user)
}

func seedUsers(userStore services.UserLoginStore) error {
    err := createUser(userStore, "admin", "secret", "admin")
    if err != nil {
        return err
    }
    return createUser(userStore, "user1", "secret", "user")
}