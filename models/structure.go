package models

type Post struct {
	ID            int
	Title         string
	Content       string
	Author        string
	Category      string
	UserID        string
	CreatedAt     string
	LikesCount    int
	DislikesCount int
	Comments      []Comment
	User          User
}

type Comment struct {
	ID               int
	UserID           int
	Content          string
	Author           string
	CreatedAt        string
	ComLikesCount    int
	ComDislikesCount int
}

type HomepageData struct {
	Authenticated bool
	User          User
	Posts         []Post
}

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	Rank      string
	LP        int
	SessionID string
}
