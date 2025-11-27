package main

import (
	"antifraud-demo-backend/internal"
	auth "antifraud-demo-backend/internal/auth"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var (
	port int

	dsn string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/antifraud_demo?sslmode=disable")

	port = viper.GetInt("PORT")
	dsn = viper.GetString("DATABASE_URL")
}

func main() {
	log.Println("Starting server...")

	db, err := internal.NewDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	internal.SetDB(db)

	mux := http.NewServeMux()

	// User login endpoints
	mux.HandleFunc("POST /login", internal.LoginHandler)
	mux.HandleFunc("POST /admin/login", internal.LoginSUHandler)

	// User endpoints
	mux.Handle("GET /users/me", auth.AuthMiddleware(http.HandlerFunc(internal.GetUsersMeHandler)))
	mux.Handle("GET /cards", auth.AuthMiddleware(http.HandlerFunc(internal.ListCardsHandler)))
	mux.Handle("GET /cards/lookup", auth.AuthMiddleware(http.HandlerFunc(internal.GetCardByNumberHandler)))
	mux.Handle("POST /transfer", auth.AuthMiddleware(http.HandlerFunc(internal.DoTransferHandler)))
	mux.Handle("GET /transfers", auth.AuthMiddleware(http.HandlerFunc(internal.ListTransfersHandler)))

	// Superuser endpoints
	mux.Handle("GET /admin/users", auth.AuthMiddleware(
		auth.RequireSuperuserMiddleware(
			http.HandlerFunc(internal.ListUsersHandler),
		),
	))
	mux.Handle("GET /admin/users/cards", auth.AuthMiddleware(
		auth.RequireSuperuserMiddleware(
			http.HandlerFunc(internal.ListCardsByUserHandler),
		),
	))
	mux.Handle("GET /admin/users/transfers", auth.AuthMiddleware(
		auth.RequireSuperuserMiddleware(
			http.HandlerFunc(internal.ListTransfersByUserHandler),
		),
	))
	mux.Handle("GET /admin/analytics/transfers", auth.AuthMiddleware(
		auth.RequireSuperuserMiddleware(
			http.HandlerFunc(internal.AnalyticsTransfersHandler),
		),
	))

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	host := fmt.Sprintf("0.0.0.0:%d", port)

	log.Printf("Serving %s", host)
	log.Fatal(http.ListenAndServe(host, corsHandler(mux)))
}
