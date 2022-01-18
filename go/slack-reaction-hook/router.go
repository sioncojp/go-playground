package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// Server ...httpサーバを立てる。shutdownはgraceful restart
func Server(port string) error {
	_, cancel := context.WithCancel(context.Background())

	slack := NewSlack()
	go slack.ListenAndResponse()
	slack.client.Run()

	// router
	router := http.NewServeMux()
	router.Handle("/healthz", http.HandlerFunc(healthz))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      logging()(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	// graceful shutdown
	go func(cancel context.CancelFunc) {
		defer close(done)
		<-sigCh
		cancel()

		// dispatcher, workerを安全に停止してからserverをshutdown
		log.sugar.Info("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelTimeout()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctxTimeout); err != nil {
			log.sugar.Warnf("Could not gracefully shutdown the server: %v\n", err)
		}
	}(cancel)

	log.sugar.Infof("Server is ready to handle requests at %s", port)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Could not listen on %s : %v\n", port, err)
	}

	<-done
	log.sugar.Info("Server stopped")
	return nil
}
