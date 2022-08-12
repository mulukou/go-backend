package main

import (
	"fmt"
	"log"
	"os"
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

var postgresLock = &sync.Mutex{}

var postgresInstance *PostgresHandler

func getPostgres() *PostgresHandler {
	if postgresInstance == nil {
		postgresLock.Lock()
		defer postgresLock.Unlock()
		if postgresInstance == nil {
			fmt.Println("Creating connection to PostgreSQL.")
			db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_CONNECTION")), &gorm.Config{})
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
