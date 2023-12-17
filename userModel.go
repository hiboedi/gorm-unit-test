package golang_gorm

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"primaryKey;column:id;<-:create"`
	Password     string    `gorm:"column:password"`
	Name         Name      `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreatedTime;<-:create"`
	UpdatedAt    time.Time `gorm:"column:created_at;autoCreatedTime;autoUpdatedTime"`
	Information  string    `gorm:"-"`
	Wallet       Wallet    `gorm:"foreignKey:user_id;references:id"`
	Addresses    []Address `gorm:"foreignKey:user_id;references:id"`
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_id"`
}

type UserLog struct {
	ID        int    `gorm:"primaryKey;column:id;autoIncrement"`
	UserId    string `gorm:"column:user_id"`
	Action    string `gorm:"column:action"`
	CreatedAt int64  `gorm:"column:created_at;autoCreatedTime:milli"`
	UpdatedAt int64  `gorm:"column:created_at;autoCreatedTime:milli;autoUpdatedTime:milli"`
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == "" {
		u.ID = "user-" + time.Now().Format("2006102150405")
	}

	return nil
}

func (u *User) TableName() string {
	return "users"
}

func (l *UserLog) TableName() string {
	return "user_logs"
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`
}
