package session

type Session struct {
	ID       string `bson:"id"`
	Username string `bson:"username"`
}

type SessionSearchFilter struct {
	ID       string `bson:"id,omitempty"`
	Username string `bson:"username,omitempty"`
}
