package db

import (
	"database/sql"
	"fmt"
	"log"

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
		comlikes_count INTEGER DEFAULT 0,
		comdislikes_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createComLikesTable := `
	CREATE TABLE IF NOT EXISTS comlikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		comlike_type INTEGER NOT NULL, -- 1 for like, -1 for dislike
		FOREIGN KEY(comment_id) REFERENCES comments(id),
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

	_, err = database.Exec(createComLikesTable)
	if err != nil {
		log.Fatalf("Error creating likes table: %v", err)
	}

	log.Println("Database tables initialized successfully.")
}

func UpdateLikesAndDislikes(postID, likeChange, dislikeChange int) error {
	db := GetDBConnection()
	var likeCount, dislikeCount int

	// Fetch current counts
	err := db.QueryRow("SELECT likes_count, dislikes_count FROM posts WHERE id=?", postID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Update counts
	likeCount += likeChange
	dislikeCount += dislikeChange
	_, err = db.Exec("UPDATE posts SET likes_count = ?, dislikes_count = ? WHERE id=?", likeCount, dislikeCount, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func UpdateComLikesAndDislikes(commentID, comlikeChange, comdislikeChange int) error {
	db := GetDBConnection()

	// Fetch current counts
	var comLikesCount, comDislikesCount int
	err := db.QueryRow("SELECT comlikes_count, comdislikes_count FROM comments WHERE id=?", commentID).
		Scan(&comLikesCount, &comDislikesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("comment not found")
		}
		fmt.Println("Error fetching comment counts:", err)
		return err
	}

	// Increment the counts
	comLikesCount += comlikeChange
	comDislikesCount += comdislikeChange

	// Update the counts
	_, err = db.Exec("UPDATE comments SET comlikes_count = ?, comdislikes_count = ? WHERE id=?", comLikesCount, comDislikesCount, commentID)
	if err != nil {
		fmt.Println("Error updating comment counts:", err)
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
	db := GetDBConnection()
	stmt, err := db.Prepare("UPDATE posts SET comments_count = comments_count + 1 WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID)
	return err
}
