package routers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
    ProductID   int     `json:"productid"`
    ProductName string  `json:"productname"`
    Category    string  `json:"category"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
    InStock     bool    `json:"instock"`
}



func formSQLQuery (r *http.Request) string {
	productId := chi.URLParam(r, "id");
	category := r.URL.Query().Get("category");

	var sqlQuery string;

	if productId == "" {
		sqlQuery = "SELECT * FROM products";

		if category != "" {
			sqlQuery += fmt.Sprintf(" WHERE category = '%s'", category);
		}
	} else {
		sqlQuery = fmt.Sprintf("SELECT * FROM products WHERE productid = %s", productId);
	}

	return sqlQuery;
}

func queryDB(query string, dbpool *pgxpool.Pool) (pgx.Rows, error) {
	return dbpool.Query(context.Background(), query)
}

func getJsonData(r *http.Request, dbpool *pgxpool.Pool) ([]byte, error) {
	sqlQuery := formSQLQuery(r);

	rows, err := queryDB(sqlQuery, dbpool);

	if err != nil {
		return nil, err
	}

	defer rows.Close();

	var products []Product

	// Iterate over the query results and append them to the products slice.
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductID, &product.ProductName, &product.Category, &product.Price, &product.Description, &product.InStock ); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Convert the products to JSON.
	productsJSON, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	return productsJSON, nil
}

func Products(dbpool *pgxpool.Pool) *chi.Mux {
	ProductsRouter := chi.NewRouter();

	ProductsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {

		allProducts, err := getJsonData(r, dbpool)

		if err != nil {
			http.Error(w, "Error getting products.", http.StatusInternalServerError)
		}

		// Set response headers.
		w.Header().Set("Content-Type", "application/json")
	
		// Send the JSON response.
		w.WriteHeader(http.StatusOK)
		w.Write(allProducts)
	})

	ProductsRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		
		product, err := getJsonData(r, dbpool)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error getting products.", http.StatusInternalServerError)
		}

		// Set response headers.
		w.Header().Set("Content-Type", "application/json")
	
		// Send the JSON response.
		w.WriteHeader(http.StatusOK)
		w.Write(product)
	})

	return ProductsRouter;
}