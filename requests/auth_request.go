package requests

type RegisterRequest struct {
	Name          string `json:"name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required|min=8,max=20"`
	FirebaseToken string `json:"firebase_token" validate:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
