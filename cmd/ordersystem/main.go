package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	configs "github.com/soares-t-o/clean-arch/configs"
	"github.com/soares-t-o/clean-arch/internal/event/handler"
	"github.com/soares-t-o/clean-arch/internal/infra/graph"
	"github.com/soares-t-o/clean-arch/internal/infra/grpc/pb"
	"github.com/soares-t-o/clean-arch/internal/infra/grpc/service"
	"github.com/soares-t-o/clean-arch/internal/infra/web/webserver"
	"github.com/soares-t-o/clean-arch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(
		configs.RabbitMQUser,
		configs.RabbitMQPassword,
		configs.RabbitMQURL,
		configs.RabbitMQPort,
	)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	webserver.AddHandler("GET", "/order", webOrderHandler.List)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(
		*createOrderUseCase,
		*listOrderUseCase,
	)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(user, password, host, port string) *amqp.Channel {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
