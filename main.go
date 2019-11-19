package main

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-logrusutil/logrusutil/errfield"
	"github.com/go-logrusutil/logrusutil/logctx"
	log "github.com/sirupsen/logrus"
)

func main() {
	addr := os.Getenv("ADDR")

	// setup logging
	log.SetOutput(os.Stdout)
	// set errfield.Formatter to log the error fields
	log.SetFormatter(&errfield.Formatter{
		Formatter: &log.JSONFormatter{},
	})
	logctx.Default.WithField("addr", addr).Info("server starting")

	// bootstrap HTTP server
	mux := http.NewServeMux()
	registerHandler(mux)

	// create a server
	server := http.Server{
		Addr:    addr,
		Handler: logMiddleware(mux),
	}
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logctx.Default.WithError(err).Fatal("server unexpectedly closed")
	}
	logctx.Default.Info("server closed")
}

func registerHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := try(); err != nil {
			logctx.From(r.Context()).WithError(err).Error("try failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logctx.From(r.Context()).Info("try succeded")
		io.WriteString(w, "Hello World")
	})
}

func try() error {
	p := struct {
		X int
		Y int
	}{
		rand.Intn(3),
		rand.Intn(3),
	}
	if p.X != p.Y {
		// return error with a structured field
		return errfield.Add(errors.New("failed to generate an excelent point"), "point", p)
	}
	return nil
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add a random reqID field for each request
		reqID := rand.Int()
		logEntry := logctx.Default.WithField("reqID", reqID)
		logEntry.Info("request started")

		// setting contextual log entry for the handler
		ctx := logctx.New(r.Context(), logEntry)
		next.ServeHTTP(w, r.WithContext(ctx))

		logEntry.Info("request finished")
	})
}
