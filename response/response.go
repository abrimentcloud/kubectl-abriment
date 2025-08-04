package response

type Response struct {
	Success    bool          `json:"success"`
	StatusCode int           `json:"status_code"`
	ErrorCode  int           `json:"error_code"`
	Message    string        `json:"message"`
	Data       LoginResponse `json:"data"`
}

type LoginResponse struct {
	Domain       Domain   `json:"domain"`
	Project      Project  `json:"project"`
	User         User     `json:"user"`
	Roles        []string `json:"roles"`
	IsAdmin      bool     `json:"is_admin"`
	IsIdentified bool     `json:"is_identified"`
	Token        Token    `json:"token"`
}

type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Enabled bool   `json:"enabled"`
}

type Token struct {
	ID        string `json:"id"`
	Unscoped  string `json:"unscoped"`
	Expires   string `json:"expires"`
	Life      int    `json:"life"`
	HeaderKey string `json:"header_key"`
}
