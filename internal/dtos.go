package internal

type UserDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
}

type CardDTO struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
	Status  string  `json:"status"`
}

type TransferDTO struct {
	ID         string  `json:"id"`
	FromUserID string  `json:"from_user_id"`
	FromCardID string  `json:"from_card_id"`
	ToCardID   string  `json:"to_card_id"`
	Amount     float64 `json:"amount"`
	When       string  `json:"when"`
	FraudScore float64 `json:"fraud_score"`
	IsBlocked  bool    `json:"is_blocked"`
}

type LoginRequest struct {
	ID string `json:"id"`
}

type SuperuserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	IsSuperuser bool   `json:"is_superuser"`
}

type TransferRequest struct {
	FromCardID string  `json:"from_card_id"`
	ToCardID   string  `json:"to_card_id"`
	Amount     float64 `json:"amount"`
}

type UserListResponse struct {
	Users []UserDTO `json:"users"`
}

type CardLookupResponse struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	UserID string `json:"user_id"`
}

type CardListResponse struct {
	Cards []CardDTO `json:"cards"`
}

type TransferListResponse struct {
	Transfers []TransferDTO `json:"transfers"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
