package db

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func InitDB(dataSourceName string) *sql.DB {
    var err error
    database, err = sql.Open("sqlite3", dataSourceName)
    if err != nil {
        log.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
    }
    return database
}

// GetDBConnection retourne la connexion à la base de données
func GetDBConnection() *sql.DB {
    return database
}

// CreateTables initialise les tables de la base de données
func CreateTables() {
	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	_, err := database.Exec(createPostsTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table posts : %v", err)
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`
	_, err = database.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table users : %v", err)
	}

	log.Println("Tables créées ou déjà existantes.")
}