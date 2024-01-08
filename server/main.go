package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
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
	client, err := ethclient.Dial(os.Getenv("NETWORK_URL"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	hm, err := services.AddWalletAddresses(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go func() {
		for {
			services.BackgroundTask(client, db, hm)
			time.Sleep(time.Second * 20)
		}
	}()
	r := chi.NewRouter()
	r.Use(httprate.LimitByIP(50, time.Duration(1)*time.Minute))
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
		r.Get("/create", routes.CreateWallet(db, hm))
		r.Get("/fetch", routes.FetchAddresses(db))
		r.Get("/{address}", routes.FetchBalance(client, db))
	})
	r.Route("/transfer", func(r chi.Router) {
		r.Post("/eth", routes.TransferEthereum(client))
		r.Post("/tokens", routes.TransferTokens(client))
	})
	r.Get("/address/{address}", routes.FetchDetails(db))
	fmt.Println("Server has started")
	log.Fatal(http.ListenAndServe("localhost:4000", r))
}
