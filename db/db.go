package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var database *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	var err error
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	CreateTables()
	ensureSessionIDColumn()
	return database
}

func GetDBConnection() *sql.DB {
	return database
}

func ensureSessionIDColumn() {
	_, err := database.Exec("ALTER TABLE users ADD COLUMN session_id TEXT;")
	if err != nil && err.Error() != "duplicate column name: session_id" {
		log.Printf("Error adding session_id column: %v", err)
	}
}

func CreateTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		session_id TEXT,
		LP INTEGER DEFAULT 0
	);`

	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		comments_count INTEGER DEFAULT 0,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		like_type INTEGER NOT NULL, -- 1 for like, -1 for dislike
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := database.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = database.Exec(createPostsTable)
	if err != nil {
		log.Fatalf("Error creating posts table: %v", err)
	}

	_, err = database.Exec(createLikesTable)
	if err != nil {
		log.Fatalf("Error creating likes table: %v", err)
	}

	_, err = database.Exec(createCommentsTable)
	if err != nil {
		log.Fatalf("Error creating comments table: %v", err)
	}

	log.Println("Database tables initialized successfully.")
}

func UpdateLikesAndDislikes(postID int, likeChange, dislikeChange int) error {
	_, err := database.Exec(`
		UPDATE posts
		SET likes_count = likes_count + ?, dislikes_count = dislikes_count + ?
		WHERE id = ?`, likeChange, dislikeChange, postID)
	return err
}

func CreatePost(title, content, category string, userID int) error {
	_, err := database.Exec(`
		INSERT INTO posts (title, content, category, user_id, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))`, title, content, category, userID)
	return err
}

func UpdateCommentsCount(postID int) error {
	_, err := database.Exec(`
		UPDATE posts
		SET comments_count = comments_count + 1
		WHERE id = ?`, postID)
	return err
}
