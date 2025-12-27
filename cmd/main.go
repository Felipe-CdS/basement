package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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
		Addr:    "5432",
		Handler: app.routes(),
	}

	fmt.Printf("Starting server...\n")
	err := srv.ListenAndServe()
	log.Fatalln(err)
}

func setEnvVars() {
	env := os.Getenv("ENV")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	envContent := ""

	switch env {
	case "local":
		envContent = "/basement/config/local"
	case "stage":
		envContent = "/basement/config/stage"
	case "prod":
		envContent = "/basement/config/prod"
	}

	if envContent == "" {
		log.Fatalln("No env declared. Check docker env vars.")
	}

	ssmClient := ssm.NewFromConfig(cfg)

	// TODO: use new() when go 1.26 releases
	boolHolder := true
	input := &ssm.GetParameterInput{
		Name:           &envContent,
		WithDecryption: &boolHolder,
	}

	result, err := ssmClient.GetParameter(ctx, input)
	if err != nil {
		log.Fatalf("No .env found. Check SSM | Err: %v", err)
	}

	envMap, err := godotenv.Unmarshal(*result.Parameter.Value)
	if err != nil {
		log.Fatalf("Failed to load env vars. Err: %v", err)
	}

	log.Println("Environment variables found. | ENV:", env)

	for key, value := range envMap {
		os.Setenv(key, value)
	}
}
