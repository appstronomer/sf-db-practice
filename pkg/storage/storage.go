package storage

// Post - публикация.
type Post struct {
	ID          int `bson:"ID"`
	Title       string `bson:"Title"`
	Content     string `bson:"Content"`
	AuthorID    int `bson:"AuthorID"`
	AuthorName  string `bson:"AuthorName"`
	CreatedAt   int64 `bson:"CreatedAt"`
	PublishedAt int64 `bson:"PublishedAt"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
