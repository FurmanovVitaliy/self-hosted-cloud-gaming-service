package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"_email"`
	Username     string `json:"username" bson:"_username"`
	PasswordHash string `json:"-" bson:"_password"`
}

// ??This User create DTO that we wil get from client
type CreateUserReq struct {
	Email    string `json:"email" bson:"_email"`
	Username string `json:"username" bson:"_username"`
	Password string `json:"password" bson:"_password"`
}

type CreateUserRes struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Email    string `json:"email" bson:"_email"`
	Username string `json:"username" bson:"_username"`
}

type LogingUserReq struct {
	Email    string `json:"email" bson:"_email"`
	Password string `json:"password" bson:"_password"`
}

type LogingUserRes struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Username    string `json:"username" bson:"_username"`
	AccessToken string `json:"access_token" bson:"_access_token"`
}
