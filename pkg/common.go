package pkg

const (
	CREATEADDRROUTE = "/users/create"
	USERADDROUTE    = "/users/"
	REGULARURL      = "https://localhost:8080"
)

type HttpData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
