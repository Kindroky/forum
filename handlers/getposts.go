package handlers

import (
	"database/sql"
	"fmt"
	"forum/db"
	"forum/models"
	"log"
	"strings"
)

func GetPosts() ([]Post, error) {
	dbConn := db.GetDBConnection()
	// Query to join posts and users table
	rows, err := dbConn.Query(`
		SELECT 
			posts.id, 
			posts.title, 
			posts.content, 
			posts.category, 
			posts.likes_count, 
			posts.dislikes_count, 
			posts.user_id, 
			users.username, 
			CASE 
				WHEN users.LP >= 100 THEN 'Legend' 
				WHEN users.LP >= 50 THEN 'Pro' 
				ELSE 'Novice' 
			END AS rank, 
			users.LP
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.id DESC;
	`)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for rows.Next() {
		var post Post
		var user models.User
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.Rank, &user.LP)
		if err != nil {
			log.Printf("Error scanning post data: %v", err)
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostById(id string) (models.Post, error) {
	dbConn := db.GetDBConnection()
	post := &models.Post{}
	err := dbConn.QueryRow(`SELECT posts.id, posts.title, posts.content, posts.created_at, users.id, users.username 
		FROM posts
		JOIN Users ON posts.user_id = users.ID
		WHERE posts.id = ?`, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.User.ID, &post.User.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return *post, err
		}
		fmt.Printf("Error retrieving post: %v\n", err)
		return *post, err
	}
	return *post, nil
}

func GetPostsByCategory(categories []string) ([]Post, error) {
	dbConn := db.GetDBConnection()
	// Build the query dynamically based on the number of categories
	query := `
		SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP 
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE `
	conditions := make([]string, len(categories))
	args := []any{}
	for i, category := range categories {
		conditions[i] = "category LIKE ?"
		args = append(args, "%"+category+"%") // Add wildcard to match partial values
	}
	query += strings.Join(conditions, " OR ")
	query += " ORDER BY posts.id DESC;"
	// Execute the query
	rows, err := dbConn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		var user models.User
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.LP)
		if err != nil {
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}
	return posts, nil
}
