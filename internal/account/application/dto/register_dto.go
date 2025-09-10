package dto

// RegisterDTO is used from client requests to the application, for registering new user account
type RegisterDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
	Role     string `json:"role" binding:"required"`
}
