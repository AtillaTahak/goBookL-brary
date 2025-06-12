package book


var books = []Book{
	{ID: 1, Title: "1984", Author: "George Orwell", Year: 1949},
	{ID: 2, Title: "Brave New World", Author: "Aldous Huxley", Year: 1932},
}

var nextID = 3

func GetAllBooks() ([]Book) {
	return books
}

func GetBookByID(id int) *Book {
	for _, b := range books {
		if b.ID == id {
			return &b
		}
	}
	return nil
}

func AddBook(b Book) Book {
	b.ID = nextID
	nextID++
	books = append(books, b)
	return b
}

func UpdateBook(id int, data Book) *Book {
	for i, b := range books {
		if b.ID == id {
			data.ID = id
			books[i] = data
			return &books[i]
		}
	}
	return nil
}

func DeleteBook(id int) bool {
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			return true
		}
	}
	return false
}
