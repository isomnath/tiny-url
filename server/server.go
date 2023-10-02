package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"
)

func recoverFromPanic() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				log.Log.HTTPErrorf(r, "Recovered from panic: %v", err)
				return
			}
		}()
		next(rw, r)
	}
}

func listenAndServe(ctx context.Context, apiServer *http.Server) {
	err := apiServer.ListenAndServe()
	if err != nil {
		log.Log.Fatalf("failed to start web router: %v", err)
	}
}

func waitForShutdown(apiServer *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	log.Log.Info("web server shutting down")
	// Finish all API calls being served and shutdown gracefully
	_ = apiServer.Shutdown(context.Background())
	log.Log.Info("web server shutting down")
}

func httpStatLogger() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		startTime := time.Now()
		next(rw, r)

		response := rw.(negroni.ResponseWriter)
		responseTime := time.Now()

		if r.URL.Path != config.GetAppHealthCheckAPIPath() {
			log.Log.HTTPStatInfo(r, startTime, responseTime, response.Status())
		}
	}
}

// StartServer - Starts the web server using a fully qualified mux router
func StartServer(router *mux.Router) {
	ctx := context.Background()
	log.Log.Infof("starting %s ... on port %d", config.GetAppName(), config.GetAppWebPort())

	handlerFunc := router.ServeHTTP

	n := negroni.New()
	n.Use(httpStatLogger())
	n.Use(recoverFromPanic())
	n.UseHandlerFunc(handlerFunc)

	portInfo := ":" + strconv.Itoa(config.GetAppWebPort())
	server := &http.Server{Addr: portInfo, Handler: n}

	go listenAndServe(ctx, server)
	waitForShutdown(server)
}
