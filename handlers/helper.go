package helper

import (
	"forum/db"
	"log"
)

type Post struct {
	ID      int
	Title   string
	Content string
	UserID  int
}

// GetPosts récupère tous les posts
func GetPosts() ([]Post, error) {
	dbConn := db.GetDBConnection()
	rows, err := dbConn.Query("SELECT * FROM posts")
	if err != nil {
		log.Printf("Erreur lors de la récupération des posts : %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			log.Printf("Erreur lors du scan des posts : %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
