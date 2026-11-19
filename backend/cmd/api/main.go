package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jkusuda/orbital/backend/internal/api"
	"github.com/jkusuda/orbital/backend/internal/ingest"
	"github.com/jkusuda/orbital/backend/internal/nasa"
	"github.com/jkusuda/orbital/backend/internal/store"
)

const refreshEvery = 6 * time.Hour

func main() {
	key := os.Getenv("NASA_API_KEY")
	if key == "" {
		log.Fatal("NASA_API_KEY is not set")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	go refreshLoop(ctx, ingest.New(nasa.NewClient(key), st), refreshEvery)

	httpSrv := &http.Server{Addr: ":" + port, Handler: api.NewServer(st).Routes()}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpSrv.Shutdown(shutdownCtx)
	}()

	log.Printf("listening on :%s", port)
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	log.Println("shut down cleanly")
}

func refreshLoop(ctx context.Context, ing *ingest.Ingester, every time.Duration) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()

	for {
		if err := ing.Refresh(ctx); err != nil && ctx.Err() == nil {
			log.Printf("refresh: %v", err)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}
