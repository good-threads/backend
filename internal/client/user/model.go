package user

type User struct {
	Name         string   `bson:"name"`
	PasswordHash []byte   `bson:"passwordHash"`
	Threads      []string `bson:"threads"`
}
