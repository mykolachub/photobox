package response

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"usename"`
	Password string `json:"password"`
}
