package internal

type LoginRequest struct {
	ID string `json:"id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CardDTO struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
	Status  string  `json:"status"`
}

type CardLookupResponse struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	UserID string `json:"user_id"`
}

type TransferRequest struct {
	FromCardID string  `json:"from_card_id"`
	ToCardID   string  `json:"to_card_id"`
	Amount     float64 `json:"amount"`
}

type TransferResponse struct {
	ID         string  `json:"id"`
	FromUserID string  `json:"from_user_id"`
	FromCardID string  `json:"from_card_id"`
	ToCardID   string  `json:"to_card_id"`
	Amount     float64 `json:"amount"`
	When       string  `json:"when"`
	FraudScore float64 `json:"fraud_score"`
	IsBlocked  bool    `json:"is_blocked"`
}

type TransferListResponse struct {
	Transfers []TransferResponse `json:"transfers"`
}

type CardListResponse struct {
	Cards []CardDTO `json:"cards"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
