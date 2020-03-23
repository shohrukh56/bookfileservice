package main

import (
	"flag"
	"github.com/shohrukh56/DI/pkg/di"
	"github.com/shohrukh56/bookFileService/cmd/app"
	"github.com/shohrukh56/bookFileService/pkg/core/file"
	"github.com/shohrukh56/jwt/pkg/jwt"
	"github.com/shohrukh56/mux/pkg/mux"
	"net"
	"net/http"
	"os"
)

var (
	host = flag.String("host", "0.0.0.0", "Server host")
	port = flag.String("port", "9999", "Server port")
)

const (
	envHost = "HOST"
	envPort = "PORT"
)

func main() {
	flag.Parse()
	serverHost := checkENV(envHost, *host)
	serverPort := checkENV(envPort, *port)
	addr := net.JoinHostPort(serverHost, serverPort)
	secret := jwt.Secret("secret")
	start(addr, secret)
}

func checkENV(env string, loc string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return loc
	}
	return str
}
func start(addr string, secret jwt.Secret) {
	container := di.NewContainer()
	container.Provide(
		func() string {return "files"},
		func() jwt.Secret { return secret },
		app.NewServer,
		file.NewService,
		mux.NewExactMux,
	)
	container.Start()
	var appServer *app.Server
	container.Component(&appServer)

	panic(http.ListenAndServe(addr, appServer))
}
