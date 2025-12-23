package main

import (
	"catalog-microservice/internal/importer/tmdb"
	"catalog-microservice/internal/storage/postgresql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//cfg := config.MustLoad()

	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("failed to load .env: %v", err))
	}
	db, err := postgresql.New()
	if err != nil {
		log.Fatalf("Could not open database %s", err)
	}
	TMDBApiKey := os.Getenv("TMDB_APIKEY")
	//log.Println(TMDBApiKey)
	client := tmdb.NewClient(TMDBApiKey)
	importer := tmdb.NewImporter(client, db, db, db)

	if err := importer.ImportPopularMovies(); err != nil {
		log.Fatal(err)
	}
	if err := importer.ImportPopularTVShows(); err != nil {
		log.Fatal(err)
	}
}
