package entity

type Guest struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required,alphaunicode,max=255"`
	Surname string `json:"surname" validate:"required,alphaunicode,max=255"`
	Phone   string `json:"phone" validate:"required,e164,max=255"`
	Email   string `json:"email" validate:"max=255"`
	Country string `json:"country" validate:"max=255"`
}
