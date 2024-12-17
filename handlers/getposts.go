package handlers

import (
	"database/sql"
	"fmt"
	"forum/db"
	"forum/models"
	"strings"
)

func GetPosts() ([]models.Post, error) {
	dbConn := db.GetDBConnection()
	rows, err := dbConn.Query(`
		SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP 
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.id DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var user models.User
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.LP)
		if err != nil {
			return nil, err
		}
		// Rank calculation integrated
		user.Rank = getRank(user.LP)
		post.User = user
		post.Rank = user.Rank
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

func GetPostsByCategory(categories []string) ([]models.Post, error) {
	dbConn := db.GetDBConnection()
	query := `
        SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE `
	conditions := make([]string, len(categories))
	args := []any{}
	for i, category := range categories {
		conditions[i] = "category LIKE ?"
		args = append(args, "%"+category+"%")
	}
	query += strings.Join(conditions, " OR ")
	query += " ORDER BY posts.id DESC;"

	rows, err := dbConn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var user models.User
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.LP)
		if err != nil {
			return nil, err
		}
		// Rank calculation integrated
		user.Rank = getRank(user.LP)
		post.User = user
		post.Rank = user.Rank
		posts = append(posts, post)
	}
	return posts, nil
}

func GetLikedPosts(userID int) ([]models.Post, error) {
	dbConn := db.GetDBConnection()
	query := `
        SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP, posts.created_at
        FROM likes
        JOIN posts ON likes.post_id = posts.id
        JOIN users ON posts.user_id = users.id
        WHERE likes.user_id = ? AND likes.like_type = 1
        ORDER BY posts.created_at DESC;
    `
	rows, err := dbConn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &post.User.Username, &post.User.LP, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetCreatedPosts retrieves posts created by the user.
func GetCreatedPosts(userID int) ([]models.Post, error) {
	dbConn := db.GetDBConnection()
	query := `
        SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.user_id = ?
        ORDER BY posts.created_at DESC;
    `
	rows, err := dbConn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &post.User.Username, &post.User.LP, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
