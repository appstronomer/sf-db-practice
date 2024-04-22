package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "sf"
	collectionName = "posts"
)

// Хранилище данных.
type Store struct {
	client mongo.Client
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Store{client: *client}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.client.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var data []storage.Post
	for cur.Next(context.Background()) {
		var post storage.Post
		err := cur.Decode(&post)
		if err != nil {
			return nil, err
		}
		data = append(data, post)
	}
	return data, cur.Err()
}

func (s *Store) AddPost(post storage.Post) error {
	collection := s.client.Database(databaseName).Collection(collectionName)
	curTimestamp := time.Now().UnixMilli()
	post.CreatedAt = curTimestamp
	post.PublishedAt = curTimestamp
	_, err := collection.InsertOne(context.Background(), post)
	return err
}

func (s *Store) UpdatePost(post storage.Post) error {
	collection := s.client.Database(databaseName).Collection(collectionName)
	filter := bson.M{"ID": post.ID}
	update := bson.M{"$set": bson.M{
		"Title":       post.Title,
		"Content":     post.Content,
		"AuthorID":    post.AuthorID,
		"AuthorName":  post.AuthorName,
		"PublishedAt": time.Now().UnixMilli(),
	}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("update post id=%v : no documents modified", post.ID)
	}
	return nil
}

func (s *Store) DeletePost(post storage.Post) error {
	collection := s.client.Database(databaseName).Collection(collectionName)
	filter := bson.M{"ID": post.ID}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("delete post id=%v : no documents deleted", post.ID)
	}
	return nil
}
