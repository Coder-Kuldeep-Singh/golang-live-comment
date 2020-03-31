package main

import (
	"database/sql"
	"golang-live-comment/handlers"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func initDB() *sql.DB {
	dbhost := os.Getenv("DBHOST")
	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	dbport := os.Getenv("DBPORT")
	dbname := os.Getenv("DB")
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp("+dbhost+":"+dbport+")/"+dbname)
	if err != nil {
		log.Println("Connection String failed", err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

// Here we create a function to migrate the database and insert the first rows for the votes
func migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS comments(
			id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			comment TEXT NOT NULL
	);
   `
	_, err := db.Exec(sql)

	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize the database
	db := initDB()
	migrate(db)

	// Define the HTTP routes
	e.File("/", "public/index.html")
	e.GET("/comments", handlers.GetComments(db))
	e.POST("/comment", handlers.PushComment(db))

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
