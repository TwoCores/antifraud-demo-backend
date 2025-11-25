package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"antifraud-demo-backend/internal/auth"

	"github.com/google/uuid"
)

var (
	dbClient *DB
)

func SetDB(db *DB) { dbClient = db }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing id"})
		return
	}

	// Look up existing user
	user, err := dbClient.GetUserByID(req.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}

	pm := r.Header.Get("X-Phone-Model")
	os := r.Header.Get("X-OS")

	sess := &LoginSession{
		UserID:     user.ID,
		When:       time.Now().UTC(),
		PhoneModel: pm,
		OS:         os,
	}
	if err := dbClient.SaveSession(sess); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to save session"})
		return
	}

	tok, _ := auth.GenerateToken(user.ID, pm, os)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: tok})
}

func GetUsersMeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := jwtClaimsFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	user, err := dbClient.GetUserByID(claims.UserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}

	_ = json.NewEncoder(w).Encode(UserDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    string(user.Status),
	})
}

func ListCardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := jwtClaimsFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	cards, err := dbClient.ListCardsForUser(claims.UserId)
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

func GetCardByNumberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cardNumber := r.URL.Query().Get("n")
	if cardNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing card number"})
		return
	}

	card, err := dbClient.GetCardByNumber(cardNumber)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "card not found"})
		return
	}

	_ = json.NewEncoder(w).Encode(CardLookupResponse{
		ID:     card.ID.String(),
		Number: card.Number,
		UserID: card.UserID,
	})
}

func DoTransferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := jwtClaimsFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	sessions, err := dbClient.ListSessionsForUser(claims.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to get sessions"})
		return
	}

	feats := ComputeFeatures(claims.UserId, sessions, time.Now().UTC())

	pr, err := PredictFraud(feats)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to call predictor"})
		return
	}

	fromCardID, err := uuid.Parse(req.FromCardID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid from_card_id"})
		return
	}
	toCardID, err := uuid.Parse(req.ToCardID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid to_card_id"})
		return
	}

	t := &Transfer{
		FromUserID: claims.UserId,
		FromCardID: fromCardID,
		ToCardID:   toCardID,
		Amount:     req.Amount,
		When:       time.Now().UTC(),
		FraudScore: pr.FraudProbability,
		IsBlocked:  pr.BlockTransaction,
	}
	if err := dbClient.SaveTransfer(t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to persist transfer"})
		return
	}

	_ = json.NewEncoder(w).Encode(TransferResponse{
		ID:         t.ID.String(),
		FromUserID: t.FromUserID,
		FromCardID: t.FromCardID.String(),
		ToCardID:   t.ToCardID.String(),
		Amount:     t.Amount,
		When:       t.When.Format(time.RFC3339),
		FraudScore: t.FraudScore,
		IsBlocked:  t.IsBlocked,
	})
}

func ListTransfersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := jwtClaimsFromContext(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	list, err := dbClient.ListTransfersForUser(claims.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to list transfers"})
		return
	}

	transferDTOs := make([]TransferResponse, len(list))
	for i, t := range list {
		transferDTOs[i] = TransferResponse{
			ID:         t.ID.String(),
			FromUserID: t.FromUserID,
			FromCardID: t.FromCardID.String(),
			ToCardID:   t.ToCardID.String(),
			Amount:     t.Amount,
			When:       t.When.Format(time.RFC3339),
			FraudScore: t.FraudScore,
			IsBlocked:  t.IsBlocked,
		}
	}

	_ = json.NewEncoder(w).Encode(TransferListResponse{Transfers: transferDTOs})
}
