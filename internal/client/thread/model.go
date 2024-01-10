package thread

type Thread struct {
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Knots    []Knot `json:"knots" bson:"knots"`
	Username string `json:"-" bson:"username"`
}

type Knot struct {
	ID   string `json:"id" bson:"id"`
	Body string `json:"body" bson:"body"`
}
