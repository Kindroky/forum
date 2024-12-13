package db

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var database *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	var err error
	database, err = sql.Open("sqlite3", "db/forum.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	CreateTables()

	database.SetMaxOpenConns(10)   // Limit max open connections
	database.SetMaxIdleConns(5)    // Limit idle connections
	database.SetConnMaxLifetime(0) // No connection timeout
	ensureSessionIDColumn()
	log.Println("Database connection initialized with pooling")
	return database
}

func CloseDB() {
	if database != nil {
		log.Println("Closing database connection")
		database.Close()
	}
}

func GetDBConnection() *sql.DB {
	if database == nil {
		log.Fatal("Database connection is nil or closed")
	}
	return database
}

func ensureSessionIDColumn() {
	addColumnIfNotExists("users", "session_id", "TEXT")
	addColumnIfNotExists("posts", "likes_count", "INTEGER DEFAULT 0")
	addColumnIfNotExists("posts", "dislikes_count", "INTEGER DEFAULT 0")
}

func addColumnIfNotExists(tableName, columnName, columnType string) {
	query := fmt.Sprintf(`
        SELECT 1 FROM pragma_table_info('%s') WHERE name = '%s';`, tableName, columnName)

	var exists int
	err := database.QueryRow(query).Scan(&exists)
	if err == sql.ErrNoRows || exists == 0 {
		alterQuery := fmt.Sprintf(`
            ALTER TABLE %s ADD COLUMN %s %s;`, tableName, columnName, columnType)
		_, err = database.Exec(alterQuery)
	} else if err != nil {
		log.Printf("Error checking for column %s in table %s: %v", columnName, tableName, err)
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
		created_at TEXT NOT NULL,
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

func UpdateLikesAndDislikes(postID, likeChange, dislikeChange int) error {
	db := GetDBConnection() // Use the global connection
	likeC := ""
	dislikeC := ""
	err := db.QueryRow("SELECT likes_count, dislikes_count FROM posts WHERE id=?", postID).Scan(&likeC, &dislikeC)
	/*
		_, err := db.Exec(Que, likeChange, dislikeChange, postID)
		if err != nil {
			return fmt.Errorf("failed to update likes/dislikes: %w", err)
		}
		return nil
	*/
	if err != nil {
		fmt.Println(err)
		return err
	}
	countLike, _ := strconv.Atoi(likeC)
	countDislike, _ := strconv.Atoi(dislikeC)
	countLike += likeChange
	countDislike += dislikeChange
	_, err = db.Exec("UPDATE posts SET likes_count = ?, dislikes_count = ? WHERE id=?", countLike, countDislike, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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

func GetPostById(id string) (models.Post, error) {
	dbConn := GetDBConnection()
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
