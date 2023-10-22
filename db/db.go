package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDb() error {

	err := godotenv.Load();
	
	if err != nil {
		return err;
    }
	
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"));

	if err != nil {
		return err;
    }

    defer conn.Close(context.Background());
	
	return nil;
}