package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// структура для запроса создания, поля с временем указаны как string так как требуется format: MM-YYYY
type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"` // format: MM-YYYY
	EndDate     *string `json:"end_date,omitempty"`
}

// структура для запроса обновления, поля с временем указаны как string так как требуется format: MM-YYYY, а так же добалены omitempty
// для того чтобы в случае отсутсвия значения поля оно не передавалось пустым
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

// структура для запроса сумарной стоимсоти подписок, поля с временем указаны как string так как требуется format: MM-YYYY
type TotalCostRequest struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	StartDate   string `json:"start_date"` // format: MM-YYYY
	EndDate     string `json:"end_date"`   // format: MM-YYYY
}
