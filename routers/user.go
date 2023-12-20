package routers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserBodyData struct {
	UserId int  `json:"userid"`
	Username string `json:"username"`
	Email string  `json:"email"`
	Password string  `json:"password"`
	CreatedAt time.Time  `json:"createdat"`
}

type UserSignInData struct {
	Username string `json:"username"`
	Password string  `json:"password"`
}

func createJWTToken (signinData UserSignInData) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": signinData.Username, "Password": signinData.Password})
	
	return tokenString, err
}

func CreateUser(w http.ResponseWriter, r *http.Request,  dbpool *pgxpool.Pool) {
	var userBody UserBodyData;

	decoder := json.NewDecoder(r.Body);

	if err := decoder.Decode(&userBody); err != nil {
        http.Error(w, "Failed to decode user body: " + err.Error(), http.StatusBadRequest);
        return
    }

	var UserId int;
	var CreatedAt time.Time;
	sqlQuery := "INSERT INTO Users (Username, Email, Password) VALUES ( $1, $2, $3) RETURNING userid, createdat;"

	err := dbpool.QueryRow(context.Background(), sqlQuery, userBody.Username, userBody.Email, userBody.Password).Scan(&UserId, &CreatedAt);

	if err != nil {
		http.Error(w, "Failed to signup the user: " + err.Error(), http.StatusBadRequest);
		return
	}

	userBody.UserId = UserId
	userBody.CreatedAt = CreatedAt

	userJson, err := json.Marshal(userBody);

	if err != nil {
		http.Error(w, "Failed to send the user info back to client: " + err.Error(), http.StatusBadRequest);
		return
	}

	// Set response headers.
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response.
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func User(dbpool *pgxpool.Pool) *chi.Mux {
	UserRouter := chi.NewRouter();

	UserRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var tokenAuth *jwtauth.JWTAuth;
		tokenAuth = jwtauth.New("HS256", []byte("my_suprise_secret"), nil) // replace with secret key
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
		fmt.Sprintf("%s", tokenString)
		fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
		w.Write([]byte("HEllo"));
	})

	UserRouter.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		CreateUser(w, r, dbpool);
	})

	UserRouter.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		var signinData UserSignInData;

		decoder := json.NewDecoder(r.Body);
	
		if err := decoder.Decode(&signinData); err != nil {
			http.Error(w, "Failed to decode user body: " + err.Error(), http.StatusBadRequest);
			return
		}

		jwtToken := createJWTToken(signinData)

		fmt.Println("Outside: ", jwtToken)
		
	})

	return UserRouter
}