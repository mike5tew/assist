package models

// User represents a system user
type User struct {
	ID        string `json:"id" bson:"_id"`
	Username  string `json:"username" bson:"username"`
	Email     string `json:"email" bson:"email"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

// Project represents a user project
type Project struct {
	ID          string     `json:"id" bson:"_id"`
	Name        string     `json:"name" bson:"name"`
	Description string     `json:"description" bson:"description"`
	Materials   []Material `json:"materials" bson:"materials"`
	CreatedAt   int64      `json:"created_at" bson:"created_at"`
	UpdatedAt   int64      `json:"updated_at" bson:"updated_at"`
}
