package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	dbUrl := os.Getenv("DB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status": "down", "error": "%v"}`, err), 503)
		return
	}
	defer conn.Close(context.Background())

	var version string
	err = conn.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		http.Error(w, `{"status": "error"}`, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "version": "%s"}`, version)
}
