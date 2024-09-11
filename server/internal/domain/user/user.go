package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"_email"`
	Username     string `json:"username" bson:"_username"`
	PasswordHash string `json:"-" bson:"_password"`
}
