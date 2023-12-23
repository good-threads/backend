package user

type User struct {
	Name         string `bson:"name"`
	PasswordHash []byte `bson:"password_hash"`
}
