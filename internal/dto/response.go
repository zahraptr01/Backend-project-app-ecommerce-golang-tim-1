package dto

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type ResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResponValidatePhone struct {
	Status string `json:"status"`
	Phone  string `json:"phone"`
	Vendor string `json:"vendor"`
}

type RegisterResponse struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Token    string `json:"token,omitempty"`
}
