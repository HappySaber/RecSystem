package main

import (
	"catalog-microservice/internal/importer/anilist"
	"catalog-microservice/internal/importer/tmdb"
	"catalog-microservice/internal/storage/postgresql"
	"fmt"
	"log"
	"os"
	"sync"

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
	log.Println(TMDBApiKey)
	client := tmdb.NewClient(TMDBApiKey)
	importer := tmdb.NewImporter(client, db, db, db)
	animeClient := anilist.NewClient()
	animeImporter := anilist.NewImporter(animeClient, db, db)

	wg := sync.WaitGroup{}
	errCh := make(chan error, 3)
	runImport := func(name string, f func() error) {
		defer wg.Done()
		if err := f(); err != nil {
			errCh <- fmt.Errorf("%s import failed: %w", name, err)
		}
	}

	wg.Add(3)

	go runImport("Movies import", func() error {
		return importer.ImportPopularMovies()
	})
	go runImport("TMDB Series", func() error {
		return importer.ImportPopularSeries()
	})

	go runImport("AniList Anime", func() error {
		return animeImporter.ImportPopularAnime()
	})

	wg.Wait()
	close(errCh)
	for err := range errCh {
		log.Println("ERROR:", err)
	}
	log.Println("All imports finished")
}
