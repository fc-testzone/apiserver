package v1

const (
	ApiVersion = "v1"
)

const (
	IndexRequest    = "/"
	RegisterRequest = "/register"
	LoginRequest    = "/login"
	ContentRequest  = "/posts"
)

type RegisterData struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
}
