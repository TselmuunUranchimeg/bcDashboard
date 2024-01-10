package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"bcDashboard/middlewares"
	"bcDashboard/routes"
	"bcDashboard/services"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	godotenv.Load()
	dbAddress := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
	)
	db, err := services.InitDb(dbAddress)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if _, err = db.Exec(`DELETE FROM "check";`); err != nil {
		log.Fatal(err)
	}
	networks, err := services.InitNetworks(db)
	if err != nil {
		log.Fatal(err)
	}
	hashmaps, err := services.InitHashmaps(db, networks)
	if err != nil {
		log.Fatal(err)
	}
	for ind := range hashmaps {
		go services.BackgroundTask(db, hashmaps[ind], 7)
	}

	r := chi.NewRouter()

	r.Use(httprate.LimitByIP(100, time.Duration(1)*time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("AUDIENCE")},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(middlewares.AuthMiddleware)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", routes.SignUp(db))
		r.Post("/login", routes.Login(db))
		r.Get("/verify", routes.Verify())
	})
	r.Route("/wallet", func(r chi.Router) {
		r.Post("/create", routes.CreateWallet(db))
		r.Get("/fetch/{network}", routes.FetchAddresses(db))
		r.Get("/{id}/{address}", routes.FetchBalance(db, hashmaps))
	})
	r.Route("/transfer", func(r chi.Router) {
		r.Post("/eth", routes.TransferEthereum(db, hashmaps))
		r.Post("/tokens", routes.TransferTokens(db, hashmaps))
	})
	r.Get("/address/{address}", routes.FetchDetails(db))
	r.Post("/contract", routes.RegisterContract(db, hashmaps))
	r.Get("/networks", routes.FetchNetworks(db))

	fmt.Println("Server has started")
	log.Fatal(http.ListenAndServe("localhost:4000", r))
}
