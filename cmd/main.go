package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"nugu.dev/basement/pkg/models/sqlite"
)

// These come from build ldflags
// Check dockerfile or air.toml for more info
var GitTag string
var GitShortSHA string
var BuildTime string

type application struct {
	AuthToken            string
	ReadBucketURL        string
	RWBucketURL          string
	activitiesRepository *sqlite.ActivityRepository
	tagsRepository       *sqlite.TagRepository

	GitTag      string
	GitShortSHA string
	BuildTime   string
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Fail .env")
	}

	db := sqlite.NewSQLiteDB(false) // true for tests database

	app := &application{
		AuthToken:   "123",
		GitTag:      GitTag,
		GitShortSHA: GitShortSHA,
		BuildTime:   BuildTime,

		ReadBucketURL: os.Getenv("PUBLIC_SINGLE_READ_BUCKET"),
		RWBucketURL:   os.Getenv("PRIVATE_LIST_RW_BUCKET"),

		activitiesRepository: &sqlite.ActivityRepository{Db: db},
		tagsRepository:       &sqlite.TagRepository{Db: db},
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: app.routes(),
	}

	fmt.Printf("> Starting server...\n")
	fmt.Printf("> Version: %s | Build: %s | Date: %s\n", GitTag, GitShortSHA, BuildTime)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
