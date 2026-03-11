package main

import (
	"catalog-microservice/internal/config"
	"catalog-microservice/internal/importer/tmdb"
	producer "catalog-microservice/internal/kafka"
	"catalog-microservice/internal/storage/postgresql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	cfg := config.MustLoad()

	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("failed to load .env: %v", err))
	}
	db, err := postgresql.New()
	if err != nil {
		log.Fatalf("Could not open database %s", err)
	}
	kafkaProducer := producer.NewKafkaProducer(
		cfg.KafkaConfig,
	)

	TMDBApiKey := os.Getenv("TMDB_APIKEY")
	//log.Println(TMDBApiKey)
	log.Println("HTTP_PROXY =", os.Getenv("HTTP_PROXY"))
	log.Println("HTTPS_PROXY =", os.Getenv("HTTPS_PROXY"))
	log.Println("http_proxy =", os.Getenv("http_proxy"))
	log.Println("https_proxy =", os.Getenv("https_proxy"))

	client := tmdb.NewClient(TMDBApiKey)
	importer := tmdb.NewImporter(client, db, db, db, kafkaProducer)
	// animeClient := anilist.NewClient()
	// animeImporter := anilist.NewImporter(animeClient, db, db)

	wg := sync.WaitGroup{}
	errCh := make(chan error, 3)
	runImport := func(name string, f func() error) {
		defer wg.Done()
		if err := f(); err != nil {
			errCh <- fmt.Errorf("%s import failed: %w", name, err)
		}
	}

	wg.Add(1)

	go runImport("Movies import", func() error {
		return importer.ImportPopularMovies()
	})
	// go runImport("TMDB Series", func() error {
	// 	return importer.ImportPopularSeries()
	// })

	// go runImport("AniList Anime", func() error {
	// 	return animeImporter.ImportPopularAnime()
	// })

	wg.Wait()
	close(errCh)
	for err := range errCh {
		log.Println("ERROR:", err)
	}
	log.Println("All imports finished")
}
