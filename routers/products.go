package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Products() *chi.Mux {
	ProductsRouter := chi.NewRouter();

	ProductsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, from inside products!"));
	})

	return ProductsRouter;
}