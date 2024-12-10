package handlers

import (
	"database/sql"
	"fmt"
	"forum/db"
	"net/http"
	"strconv"
)

type Comment struct {
    ID        int
    UserID    int
    Content   string
    Author    string // Ajout de l'auteur
    CreatedAt string // Ajout de la date de création
}

func detailPostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID du post depuis les paramètres de l'URL
	postIDStr := r.URL.Query().Get("postID")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	// Initialiser le modèle du post
	var post Post

	// Récupérer les détails du post depuis la base de données
	query := `SELECT posts.id, posts.title, posts.content, users.username, posts.created_at
              FROM posts
              JOIN users ON posts.user_id = users.id
              WHERE posts.id = ?`
	dbConn := db.GetDBConnection()
	err = dbConn.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Récupérer les commentaires associés au post
	rows, err := dbConn.Query(`
    SELECT comments.id, comments.user_id, comments.content, users.username, comments.created_at
    FROM comments
    JOIN users ON comments.user_id = users.id
    WHERE comments.post_id = ?
    ORDER BY comments.created_at DESC
`, postID)
if err != nil {
    http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
    return
}
defer rows.Close()

	defer rows.Close()

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.Author, &comment.CreatedAt); err != nil {
			http.Error(w, "Error reading comments", http.StatusInternalServerError)
			return
		}
		post.Comments = append(post.Comments, comment)
	}	

	// Générer la réponse HTML (ou JSON)
	// Si vous préférez générer du JSON, décommentez cette section :
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(post)

	// Sinon, générer une réponse HTML simple :
	fmt.Fprintf(w, "<h1>%s</h1><p>%s</p><p><i>By %s on %s</i></p>",
		post.Title, post.Content, post.Author, post.CreatedAt)
	fmt.Fprintf(w, "<h2>Comments</h2>")
	for _, comment := range post.Comments {
		fmt.Fprintf(w, "<p><b>User %d:</b> %s <i>on %s</i></p>", comment.UserID, comment.Content, comment.CreatedAt)
	}
}
