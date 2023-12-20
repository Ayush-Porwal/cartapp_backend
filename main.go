package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ayush-Porwal/cartapp_backend/db"
	"github.com/Ayush-Porwal/cartapp_backend/routers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	os.Setenv("ENVIRONMENT", "development");
	mode := os.Getenv("ENVIRONMENT");

	var filename string;

	if mode == "development" {
		filename = ".env.local"
	} else {
		filename = ".env.prod"
	}

	err := godotenv.Load(filename);
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v", err);
    }
	
	dbpool, err := db.InitializeDBPool(os.Getenv("DATABASE_URL"));

	if err != nil  {
		fmt.Fprintf(os.Stderr, "Could not open a connection pool to database: %v\n", err);
	}

	router := chi.NewRouter();
	
	router.Use(middleware.Logger);
	
	// mounting the routes
	router.Mount("/cart", routers.Cart(dbpool));
	router.Mount("/user", routers.User(dbpool));
	router.Mount("/products", routers.Products(dbpool));

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":3000", router);
}