package handlers

import (
	"forum/db"
	"forum/models"
	"html/template"
	"net/http"
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
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &user.Username, &user.Rank, &user.LP)
		if err != nil {
			return nil, err
		}
		post.User = user
		posts = append(posts, post)
	}

	return posts, nil
}

func PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		Error(w, r, http.StatusBadRequest, "Post ID is required.")
		return
	}

	authenticated := false
	var user models.User

	cookie, err := r.Cookie("session_id")
	if err == nil {
		dbConn := db.GetDBConnection()
		err = dbConn.QueryRow(`
			SELECT id, username, LP, session_id 
			FROM users WHERE session_id = ?`, cookie.Value).Scan(&user.ID, &user.Username, &user.LP, &user.SessionID)
		if err == nil {
			authenticated = true
		}
	}

	dbConn := db.GetDBConnection()
	var post models.Post
	err = dbConn.QueryRow(`
		SELECT posts.id, posts.title, posts.content, posts.category, posts.likes_count, posts.dislikes_count, posts.user_id, users.username, users.LP 
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.LikesCount, &post.DislikesCount, &post.UserID, &post.User.Username, &post.User.LP)
	if err != nil {
		Error(w, r, http.StatusNotFound, "Post not found.")
		return
	}

	comments, err := FetchComments(post.ID)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "An error occurred while fetching comments.")
		return
	}

	data := struct {
		Authenticated bool
		User          models.User
		Post          models.Post
		Comments      []models.Comment
	}{
		Authenticated: authenticated,
		User:          user,
		Post:          post,
		Comments:      comments,
	}

	tmpl, err := template.ParseFiles("templates/postdetails.html")
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "An error occurred while loading the post page.")
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "An error occurred while rendering the post page.")
	}
}
