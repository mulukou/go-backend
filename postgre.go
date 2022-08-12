package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string
	Description string
	Picture     string
	Price       float64
}

type PostgresHandler struct {
	postgresInstance *gorm.DB
}

var postgreLock = &sync.Mutex{}

var postgresInstance *PostgresHandler

func getPostgre() *PostgresHandler {
	if postgresInstance == nil {
		postgreLock.Lock()
		defer postgreLock.Unlock()
		if postgresInstance == nil {
			fmt.Println("Creating connection to PostgreSQL.")
			db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=extramemepassword dbname=postgres port=5432 sslmode=disable"), &gorm.Config{})
			if err != nil {
				log.Fatalln(err)
			}
			db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
			db.AutoMigrate(&Item{})
			postgresInstance = &PostgresHandler{db}
		} else {
			fmt.Println("PostgreSQL instance already created.")
		}
	} else {
		fmt.Println("PostgreSQL instance already created.")
	}

	return postgresInstance
}
