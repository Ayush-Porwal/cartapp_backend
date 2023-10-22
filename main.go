package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ayush-Porwal/cartapp_backend/db"
	"github.com/Ayush-Porwal/cartapp_backend/routers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	
	err := db.ConnectDb();

	if err != nil  {
		fmt.Fprintf(os.Stderr, "Could not connect to database: %v\n", err);
	}

	router := chi.NewRouter();
	
	router.Use(middleware.Logger);
	
	// mounting the routes
	router.Mount("/cart", routers.Cart());
	router.Mount("/products", routers.Products());

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":3000", router);
}