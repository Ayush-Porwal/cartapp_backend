package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Cart() *chi.Mux {
	CartRouter := chi.NewRouter();

	CartRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, from inside cart!"));
	})

	return CartRouter;
}