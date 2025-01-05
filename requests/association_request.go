package requests

type JoinGroupRequest struct {
	Code string `json:"code" gorm:"unique;not null" validate:"required,min=5,max=10"`
}
