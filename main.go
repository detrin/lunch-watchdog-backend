package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/detrin/lunch-watchdog-backend/watchdog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	var menus []watchdog.Menu
	menuKolkovna, err := watchdog.ScrapeMenuKolkovna()
	if err != nil {
		log.Fatalf("Error scraping menu: %v", err)
	}
	fmt.Printf("Menu for %s - %s\n", menuKolkovna.Name, menuKolkovna.Date)
	for i, item := range menuKolkovna.MenuItems {
		fmt.Printf("Item %d: %#v\n", i+1, item)
	}
	menus = append(menus, *menuKolkovna)

	menuMerkur, err := watchdog.ScrapeMenuMerkur()
	if err != nil {
		log.Fatalf("Error scraping menu: %v", err)
	}
	fmt.Printf("Menu for %s - %s\n", menuMerkur.Name, menuMerkur.Date)
	for i, item := range menuMerkur.MenuItems {
		fmt.Printf("Item %d: %#v\n", i+1, item)
	}
	menus = append(menus, *menuMerkur)

	jsonData, err := json.MarshalIndent(menus, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling menu to JSON: %v", err)
	}

	fmt.Println(string(jsonData))

	endpoint := "eu2.contabostorage.com"
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := "lunch-watchdog"
	objectName := "menus.json"
	content := bytes.NewReader(jsonData)
	opts := minio.PutObjectOptions{ContentType: "application/json"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	info, err := s3Client.PutObject(ctx, bucketName, objectName, content, int64(content.Len()), opts)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
