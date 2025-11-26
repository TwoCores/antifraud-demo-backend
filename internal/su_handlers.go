package internal

import (
	"encoding/json"
	"net/http"
)

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := dbClient.ListAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to list users"})
		return
	}

	userDTOs := make([]UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = UserDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    string(user.Status),
		}
	}

	_ = json.NewEncoder(w).Encode(UserListResponse{Users: userDTOs})
}

func ListCardsByUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing user id"})
		return
	}

	cards, err := dbClient.ListAllCardsByUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to list cards"})
		return
	}

	cardDTOs := make([]CardDTO, len(cards))
	for i, card := range cards {
		cardDTOs[i] = CardDTO{
			ID:      card.ID.String(),
			UserID:  card.UserID,
			Number:  card.Number,
			Balance: card.Balance,
			Status:  string(card.Status),
		}
	}

	_ = json.NewEncoder(w).Encode(CardListResponse{Cards: cardDTOs})
}

func ListTransfersByUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing user id"})
		return
	}

	transfers, err := dbClient.ListAllTransfersByUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to list transfers"})
		return
	}

	transferDTOs := make([]TransferDTO, len(transfers))
	for i, t := range transfers {
		transferDTOs[i] = TransferDTO{
			ID:         t.ID.String(),
			FromUserID: t.FromUserID,
			FromCardID: t.FromCardID.String(),
			ToCardID:   t.ToCardID.String(),
			Amount:     t.Amount,
			When:       t.When.Format("2006-01-02T15:04:05Z07:00"),
			FraudScore: t.FraudScore,
			IsBlocked:  t.IsBlocked,
		}
	}

	_ = json.NewEncoder(w).Encode(TransferListResponse{Transfers: transferDTOs})
}
