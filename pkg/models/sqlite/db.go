package sqlite

import (
	"context"
	"database/sql"
	"log"
	"os"

	turso "turso.tech/database/tursogo"
)

type Store struct {
	Db      *sql.DB
	Testing bool
}

func NewSQLiteDB(t bool) *sql.DB {

	dbStore := Store{Testing: t}

	if err := dbStore.getConnection(); err != nil {
		log.Fatalf("failed to connect to the database... Error: %s", err)
	}

	return dbStore.Db
}

func (dbStore *Store) getConnection() error {
	if dbStore.Db != nil {
		return nil
	}

	ctx := context.Background()

	syncDb, _ := turso.NewTursoSyncDb(ctx, turso.TursoSyncDbConfig{
		Path:      "pkg/models/sqlite/turso/app.db",
		RemoteUrl: os.Getenv("TURSO_DATABASE_URL"),
		AuthToken: os.Getenv("TURSO_AUTH_TOKEN"),
	})

	db, _ := syncDb.Connect(ctx)

	dbStore.Db = db
	log.Printf("Connected successfully to the database | TURSO: %s \n", os.Getenv("TURSO_DATABASE_URL"))

	return nil
}
