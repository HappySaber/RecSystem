package main

import (
	"catalog-microservice/internal/importer/anilist"
	"catalog-microservice/internal/storage/postgresql"
	"fmt"
	"log"

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
	//TMDBApiKey := os.Getenv("TMDB_APIKEY")
	//log.Println(TMDBApiKey)
	// client := tmdb.NewClient(TMDBApiKey)
	// importer := tmdb.NewImporter(client, db, db, db)
	animeClient := anilist.NewClient()
	animeImporter := anilist.NewImporter(animeClient, db, db)

	// if err := importer.ImportPopularMovies(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := importer.ImportPopularSeries(); err != nil {
	// 	log.Fatal(err)
	// }

	if err := animeImporter.ImportPopularAnime(); err != nil {
		log.Fatal(err)
	}
}
