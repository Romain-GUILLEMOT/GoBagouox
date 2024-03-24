package web

import (
	"GoBagouox/utils"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
)

type Route struct {
	Method     string
	Path       string
	Handler    gin.HandlerFunc
	Middleware []gin.HandlerFunc
}

func getRoutes() []Route {
	return []Route{
		{
			Method:     "GET",
			Path:       "/",
			Handler:    Hello,
			Middleware: nil,
		},
	}
}
func joinMiddlewares(middlewares []gin.HandlerFunc) string {
	if len(middlewares) == 0 {
		return ""
	}

	names := make([]string, len(middlewares))
	for i, middleware := range middlewares {
		names[i] = runtime.FuncForPC(reflect.ValueOf(middleware).Pointer()).Name()
	}

	return " [MIDDLEWARES: " + strings.Join(names, ", ") + "]"
}
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	routes := getRoutes()
	webport := os.Getenv("WEBSERVER_PORT")

	for _, route := range routes {
		middlewares := joinMiddlewares(route.Middleware)
		switch route.Method {
		case "GET":
			router.GET(route.Path, append(route.Middleware, route.Handler)...)
			utils.Debug("Route "+utils.Bold(utils.Purple(route.Path+middlewares))+" GET loaded.", 0)
		case "POST":
			router.POST(route.Path, append(route.Middleware, route.Handler)...)
			utils.Debug("Route "+utils.Bold(utils.Purple(route.Path+middlewares))+" POST loaded.", 0)
		case "PUT":
			router.PUT(route.Path, append(route.Middleware, route.Handler)...)
			utils.Debug("Route "+utils.Bold(utils.Purple(route.Path+middlewares))+" PUT loaded.", 0)
		case "DELETE":
			router.DELETE(route.Path, append(route.Middleware, route.Handler)...)
			utils.Debug("Route "+utils.Bold(utils.Purple(route.Path+middlewares))+" DELETE loaded.", 0)
		default:
			utils.Fatal("Unknown method for route "+route.Path+": "+route.Method, errors.New("Unknow type: "+route.Method), 0)
		}
	}
	server := &http.Server{
		Addr:    ":" + webport,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.Fatal("ERR-001: Can't start the webserver on port "+webport, err, 0)
		}
	}()
	utils.Info("Web server is now running on port "+webport+".", 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	if err := server.Shutdown(context.Background()); err != nil {
		utils.Fatal("ERR-002: Webserver Shutdown failed", err, 0)
	}
	utils.Info("Web server bot has now stopped.", 0)

}

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Accueil"})
}