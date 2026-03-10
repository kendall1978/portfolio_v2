package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
	"google.golang.org/api/option"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/kroberts_personal?sslmode=disable"
	}

	// Connect to Postgres
	if err := database.Connect(dbURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Connect to Firebase
	ctx := context.Background()
	opt := option.WithCredentialsFile("firebase-service-account.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to init Firebase: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to init Firestore: %v", err)
	}
	defer client.Close()

	// Migrate blog posts
	migrateBlogPosts(ctx, client)

	// Seed home sections from old hardcoded data
	seedHomeSections()

	// Seed social links from old Home.vue
	seedSocialLinks()

	// Migrate profile image from Firebase assets
	migrateProfileImage(ctx, client)

	fmt.Println("Migration complete!")
}

func migrateBlogPosts(ctx context.Context, client *firestore.Client) {
	docs, err := client.Collection("blog-posts").OrderBy("date", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed to read blog-posts: %v", err)
	}

	for _, doc := range docs {
		data := doc.Data()
		post := models.BlogPost{
			Title:     fmt.Sprintf("%v", data["title"]),
			Content:   fmt.Sprintf("%v", data["content"]),
			ImagePath: fmt.Sprintf("%v", data["imageUrl"]),
		}
		// Parse date - the old app stored dates as strings
		if dateStr, ok := data["date"].(string); ok {
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				post.Date = t
			} else if t, err := time.Parse("01/02/2006", dateStr); err == nil {
				post.Date = t
			} else if t, err := time.Parse("January 2, 2006", dateStr); err == nil {
				post.Date = t
			} else {
				post.Date = time.Now()
				log.Printf("Could not parse date '%s' for post '%s', using now", dateStr, post.Title)
			}
		} else {
			post.Date = time.Now()
		}

		if err := database.DB.Create(&post).Error; err != nil {
			log.Printf("Failed to insert post '%s': %v", post.Title, err)
		} else {
			log.Printf("Migrated post: %s", post.Title)
		}
	}
}

func seedHomeSections() {
	sections := []models.HomeSection{
		{
			Title:     "My Education",
			Content:   "I've been interested in technology and code since I was in high school. My first project-setting me up to love this field-was putting RetroPie on a Raspberry Pi. Ever since then I've loved this field. So I went to Ozarks Technical Community College and finished a degree in Computer Information Science",
			SortOrder: 1,
		},
		{
			Title:     "My Life",
			Content:   "Nerd for life, I love work on a computer and write code. I also find enjoyment in activities like fishing, hunting, hiking, and boating. Put me in a kayak on a river somewhere and I am one happy guy.",
			SortOrder: 2,
		},
		{
			Title:     "My Future",
			Content:   "My long term goal is to be a well respected full stack software engineer with a project manager role in mind. My short term goals are always to continue to challege myself with new technologies and write quality code. At the end of the day I just want to learn all that I can about web development.",
			SortOrder: 3,
		},
	}

	for _, s := range sections {
		if err := database.DB.Create(&s).Error; err != nil {
			log.Printf("Failed to insert section '%s': %v", s.Title, err)
		} else {
			log.Printf("Seeded section: %s", s.Title)
		}
	}
}

func seedSocialLinks() {
	links := []models.SocialLink{
		{Platform: "GitHub", URL: "https://github.com/kendall1978?tab=repositories", Icon: "github", SortOrder: 1},
		{Platform: "LinkedIn", URL: "https://www.linkedin.com/in/kendall-roberts-1b8b3a171/", Icon: "linkedin", SortOrder: 2},
		{Platform: "Twitter", URL: "https://twitter.com/kendertsrondall", Icon: "twitter", SortOrder: 3},
		{Platform: "Instagram", URL: "https://www.instagram.com/kendall_roberts49/", Icon: "instagram", SortOrder: 4},
	}

	for _, l := range links {
		if err := database.DB.Create(&l).Error; err != nil {
			log.Printf("Failed to insert link '%s': %v", l.Platform, err)
		} else {
			log.Printf("Seeded link: %s", l.Platform)
		}
	}
}

func migrateProfileImage(ctx context.Context, client *firestore.Client) {
	doc, err := client.Collection("assets").Doc("homePic").Get(ctx)
	if err != nil {
		log.Printf("Could not read assets/homePic: %v (skipping)", err)
		return
	}

	data := doc.Data()
	if imgURL, ok := data["img"].(string); ok {
		setting := models.SiteSetting{
			Key:   "profile_image_url",
			Value: imgURL,
		}
		database.DB.Where("key = ?", "profile_image_url").FirstOrCreate(&setting)
		log.Printf("Set profile image URL: %s", imgURL)
	}

	// Also set default site title
	titleSetting := models.SiteSetting{Key: "site_title", Value: "Kendall Roberts"}
	database.DB.Where("key = ?", "site_title").FirstOrCreate(&titleSetting)

	subtitleSetting := models.SiteSetting{Key: "site_subtitle", Value: "Web-Developer, Student, and Part-Time Dad Joke Comedian"}
	database.DB.Where("key = ?", "site_subtitle").FirstOrCreate(&subtitleSetting)
}
