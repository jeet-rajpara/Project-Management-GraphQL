package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"project_management/api/constants"
	"project_management/api/dataloaders"
	"project_management/api/middlewares"
	"project_management/database"
	er "project_management/errors"
	"project_management/graph"
	"project_management/utils"
	"project_management/utils/socket"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

const defaultPort = "8010"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Error in establishing database connection: %v", err)
	}

	router := chi.NewRouter()
	router.Use(database.Middleware(db))
	router.Use(middlewares.AuthMiddleware)
	router.Use(dataloaders.LoaderMiddleware)

	resolver := &graph.Resolver{}

	resolverConfig := graph.Config{Resolvers: resolver}

	resolverConfig.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		authToken := ctx.Value(constants.AuthTokenCtxKey)
		if authToken == "" {
			return nil, er.TokenNotFound
		}
		//validate token
		userId, tokenError := utils.VerifyToken(ctx, fmt.Sprintf("%v", authToken))
		if tokenError != nil {
			return nil, tokenError
		}
		ctx = context.WithValue(ctx, constants.UserIDCtxKey, userId)
		return next(ctx)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(resolverConfig))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	// Serve static files from the public directory
	// fs := http.FileServer(http.Dir("./public"))
	// router.Handle("/*", fs)

	socketServer := socket.SocketConnection()
	socket.RegisterEvents(socketServer)
	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
		defer socketServer.Close()
	}()

	// Handle WebSocket connections
	router.Handle("/socket.io/", socketServer)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, router))
}
