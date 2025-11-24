package internal

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sqlx.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetUserByID(id string) (*User, error) {
	var user User
	err := db.conn.Get(&user, `SELECT id, first_name, last_name, status FROM users WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) EnsureUser(u *User) error {
	_, err := db.conn.Exec(`INSERT INTO users (id, first_name, last_name, status) VALUES ($1,$2,$3,$4) ON CONFLICT (id) DO UPDATE SET first_name=EXCLUDED.first_name,last_name=EXCLUDED.last_name,status=EXCLUDED.status`, u.ID, u.FirstName, u.LastName, u.Status)
	return err
}

func (db *DB) SaveSession(s *LoginSession) error {
	_, err := db.conn.Exec(`INSERT INTO login_sessions (user_id, when_ts, phone_model, os) VALUES ($1,$2,$3,$4)`, s.UserID, s.When, s.PhoneModel, s.OS)
	return err
}

func (db *DB) ListSessionsForUser(userId string) ([]*LoginSession, error) {
	out := make([]*LoginSession, 0)
	err := db.conn.Select(&out, `SELECT id, user_id, when_ts as when, phone_model, os FROM login_sessions WHERE user_id=$1 ORDER BY when_ts ASC`, userId)
	return out, err
}

func (db *DB) ListCardsForUser(userId string) ([]*Card, error) {
	out := make([]*Card, 0)
	err := db.conn.Select(&out, `SELECT id, user_id, number, balance, status FROM cards WHERE user_id=$1`, userId)
	return out, err
}

func (db *DB) GetCardByNumber(number string) (*Card, error) {
	var card Card
	err := db.conn.Get(&card, `SELECT id, user_id, number, balance, status FROM cards WHERE number=$1`, number)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (db *DB) SaveTransfer(t *Transfer) error {
	_, err := db.conn.Exec(`INSERT INTO transfers (from_user_id, from_card_id, to_card_id, amount, when_ts, fraud_score, is_blocked) VALUES ($1,$2,$3,$4,$5,$6,$7)`, t.FromUserID, t.FromCardID, t.ToCardID, t.Amount, t.When, t.FraudScore, t.IsBlocked)
	return err
}

func (db *DB) ListTransfersForUser(userId string) ([]*Transfer, error) {
	out := make([]*Transfer, 0)
	err := db.conn.Select(&out, `SELECT id, from_user_id, from_card_id, to_card_id, amount, when_ts as when, fraud_score, is_blocked FROM transfers WHERE from_user_id=$1 ORDER BY when_ts DESC`, userId)
	return out, err
}
