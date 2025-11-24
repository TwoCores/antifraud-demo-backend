package internal

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	StatusActive  UserStatus = "active"
	StatusBlocked UserStatus = "blocked"
)

type User struct {
	ID        string     `json:"id" db:"id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Status    UserStatus `json:"status" db:"status"`
}

type CardStatus string

const (
	CardActive  CardStatus = "active"
	CardBlocked CardStatus = "blocked"
)

type Card struct {
	ID      uuid.UUID  `json:"id" db:"id"`
	UserID  string     `json:"user_id" db:"user_id"`
	Number  string     `json:"number" db:"number"`
	Balance float64    `json:"balance" db:"balance"`
	Status  CardStatus `json:"status" db:"status"`
}

type LoginSession struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	When       time.Time `json:"when" db:"when"`
	PhoneModel string    `json:"phone_model" db:"phone_model"`
	OS         string    `json:"os" db:"os"`
}

type Transfer struct {
	ID         uuid.UUID `json:"id" db:"id"`
	FromUserID string    `json:"from_user_id" db:"from_user_id"`
	FromCardID uuid.UUID `json:"from_card_id" db:"from_card_id"`
	ToCardID   uuid.UUID `json:"to_card_id" db:"to_card_id"`
	Amount     float64   `json:"amount" db:"amount"`
	When       time.Time `json:"when" db:"when"`
	FraudScore float64   `json:"fraud_score" db:"fraud_score"`
	IsBlocked  bool      `json:"is_blocked" db:"is_blocked"`
}
