package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"nugu.dev/basement/pkg/models/postgres"
)

type application struct {
	AuthToken            string
	ReadBucketURL        string
	RWBucketURL          string
	activitiesRepository *postgres.ActivityRepository
	tagsRepository       *postgres.TagRepository
}

func main() {

	setEnvVars()

	db := postgres.NewPostgresDB(false) // true for tests database

	app := &application{
		AuthToken:            "123",
		ReadBucketURL:        os.Getenv("PUBLIC_SINGLE_READ_BUCKET"),
		RWBucketURL:          os.Getenv("PRIVATE_LIST_RW_BUCKET"),
		activitiesRepository: &postgres.ActivityRepository{Db: db},
		tagsRepository:       &postgres.TagRepository{Db: db},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: app.routes(),
	}

	fmt.Printf("Starting server on port %s...\n", os.Getenv("APP_PORT"))
	err := srv.ListenAndServe()
	log.Fatalln(err)
}

func setEnvVars() {
	env := os.Getenv("ENV")
	envPath, _ := os.Getwd()

	// In production the binary is controlled by systemd, so the path is "/".
	// In this case we need to get the executable path and set the service
	// environment to start as "ENV=PROD ./main".

	if env == "PROD" {
		envPath, _ = os.Executable()
		envPath = filepath.Dir(envPath)
	}

	if err := godotenv.Load(filepath.Join(envPath, ".env")); err != nil {
		log.Fatalln("No .env file found. Path:", envPath)
	} else {
		log.Println("Environment variables found. Path:", envPath)
	}
}
