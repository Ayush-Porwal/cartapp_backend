package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool;

func InitializeDBPool(connectionURL string) (*pgxpool.Pool, error) {
	dbconfig, err := pgxpool.ParseConfig(connectionURL);
	
	if err != nil {
		return nil, err;
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig);

	if err != nil {
		return nil, err;
	}

	// defer dbpool.Close();

	return dbpool, err;
}


// func GetDBPool() *pgxpool.Pool {
// 	return dbpool;
// }

// func ConnectDb() (*pgx.Conn, error) {

// 	err := godotenv.Load();
	
// 	if err != nil {
// 		return nil, err;
//     }
	
// 	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"));

// 	if err != nil {
// 		return nil, err;
//     }

//     defer conn.Close(context.Background());
	
// 	return conn, nil;
// }