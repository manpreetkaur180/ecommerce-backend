package user

type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
type LoginRequest struct {
    Email string `json:"email"`
}