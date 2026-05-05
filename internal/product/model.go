package product



type Product struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Price     float64   `gorm:"not null"`
    Stock     int       `gorm:"not null"`

}