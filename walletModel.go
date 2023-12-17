package golang_gorm

import "time"

type Wallet struct {
	ID        string    `gorm:"primary_key;column:id"`
	UserId    string    `gorm:"column:user_id"`
	Balance   int64     `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreatedTime"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreatedTime;autoUpdatedTime"`
	User      *User     `gorm:"foreignKey:user_id;references:id"`
}

func (w *Wallet) TableName() string {
	return "wallets"
}
