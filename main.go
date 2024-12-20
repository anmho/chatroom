package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var port = flag.String("port", "8080", "http service port")
var logger = slog.Default()

func serveHome(w http.ResponseWriter, r *http.Request) {
	logger = logger.With("handler", "serveHome")
	logger.Info(r.URL.String())

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)

	addr := fmt.Sprintf(":%s", *port)
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logger = logger.With("addr", addr, "port", *port)
	logger.Info("server listening")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("unexpected server error", slog.Any("error", err))
		os.Exit(1)
	}
}
