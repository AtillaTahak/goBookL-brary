package book

import (
	"github.com/AtillaTahaK/gobooklibrary/pkg/db"
)

func GetAllBooks() ([]Book, error) {
	var books []Book
	if err := db.DB.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func GetBookByID(id uint) (*Book, error) {
	var book Book
	if err := db.DB.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func CreateBook(book *Book) error {
	if err := db.DB.Create(book).Error; err != nil {
		return err
	}
	return nil
}

func UpdateBook(id uint, updatedBook *Book) (*Book, error) {
	var book Book
	if err := db.DB.First(&book, id).Error; err != nil {
		return nil, err
	}

	// Update only non-zero fields
	if err := db.DB.Model(&book).Updates(updatedBook).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func DeleteBook(id uint) error {
	if err := db.DB.Delete(&Book{}, id).Error; err != nil {
		return err
	}
	return nil
}

func SearchBooks(query string) ([]Book, error) {
	var books []Book
	if err := db.DB.Where("title ILIKE ? OR author ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
