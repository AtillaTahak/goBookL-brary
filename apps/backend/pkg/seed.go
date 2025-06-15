package pkg

import (
	"log"

	"github.com/AtillaTahaK/gobooklibrary/auth"
	"github.com/AtillaTahaK/gobooklibrary/book"
	"github.com/AtillaTahaK/gobooklibrary/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

func seedDatabase() {
	var userCount int64
	db.DB.Model(&auth.User{}).Count(&userCount)

	if userCount > 0 {
		log.Println("Database already has data, skipping seeding")
		return
	}

	log.Println("Seeding database with initial data...")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUser := auth.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@booklibrary.com",
		Role:     "admin",
	}

	if err := db.DB.Create(&adminUser).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
	} else {
		log.Println("Created admin user (username: admin, password: admin123)")
	}

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	regularUser := auth.User{
		Username: "user",
		Password: string(hashedPassword),
		Email:    "user@booklibrary.com",
		Role:     "user",
	}

	if err := db.DB.Create(&regularUser).Error; err != nil {
		log.Printf("Failed to create regular user: %v", err)
	} else {
		log.Println("Created regular user (username: user, password: user123)")
	}

	sampleBooks := []book.Book{
		{
			Title:  "1984",
			Author: "George Orwell",
			Year:   1949,
			Genre:  "Dystopian Fiction",
			ISBN:   "978-0-452-28423-4",
		},
		{
			Title:  "Brave New World",
			Author: "Aldous Huxley",
			Year:   1932,
			Genre:  "Science Fiction",
			ISBN:   "978-0-06-085052-4",
		},
		{
			Title:  "To Kill a Mockingbird",
			Author: "Harper Lee",
			Year:   1960,
			Genre:  "Fiction",
			ISBN:   "978-0-06-112008-4",
		},
		{
			Title:  "The Great Gatsby",
			Author: "F. Scott Fitzgerald",
			Year:   1925,
			Genre:  "Classic Literature",
			ISBN:   "978-0-7432-7356-5",
		},
		{
			Title:  "Pride and Prejudice",
			Author: "Jane Austen",
			Year:   1813,
			Genre:  "Romance",
			ISBN:   "978-0-14-143951-8",
		},
	}

	for _, bookData := range sampleBooks {
		if err := db.DB.Create(&bookData).Error; err != nil {
			log.Printf("Failed to create book %s: %v", bookData.Title, err)
		}
	}

	log.Printf("Created %d sample books", len(sampleBooks))
	log.Println("Database seeding completed!")
}
