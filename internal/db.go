package internal

import (
	"time"

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

func (db *DB) GetSuperuserByUsername(username string) (*Superuser, error) {
	var su Superuser
	err := db.conn.Get(&su, `SELECT id, username, password_hash FROM superusers WHERE username=$1`, username)
	if err != nil {
		return nil, err
	}
	return &su, nil
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

func (db *DB) ListAllUsers() ([]*User, error) {
	out := make([]*User, 0)
	err := db.conn.Select(&out, `SELECT id, first_name, last_name, status FROM users`)
	return out, err
}

func (db *DB) ListAllCardsByUser(userId string) ([]*Card, error) {
	out := make([]*Card, 0)
	err := db.conn.Select(&out, `SELECT id, user_id, number, balance, status FROM cards WHERE user_id=$1`, userId)
	return out, err
}

func (db *DB) ListAllTransfersByUser(userId string) ([]*Transfer, error) {
	out := make([]*Transfer, 0)
	err := db.conn.Select(&out, `SELECT id, from_user_id, from_card_id, to_card_id, amount, when_ts as when, fraud_score, is_blocked FROM transfers WHERE from_user_id=$1`, userId)
	return out, err
}

func (db *DB) ListAllLoginSessionsByUser(userId string) ([]*LoginSession, error) {
	out := make([]*LoginSession, 0)
	err := db.conn.Select(&out, `SELECT id, user_id, when_ts as when, phone_model, os FROM login_sessions WHERE user_id=$1`, userId)
	return out, err
}

// TODO: well, it works, but needs some refactoring
type TransferAnalytics struct {
	TotalTransfers   int                         `json:"total_transfers"`
	BlockedTransfers int                         `json:"blocked_transfers"`
	DailyStats       []TransferAnalyticsDayStats `json:"daily_stats"`
}

type TransferAnalyticsDayStats struct {
	Date       string `json:"date"`
	Total      int    `json:"total"`
	Blocked    int    `json:"blocked"`
	Successful int    `json:"successful"`
}

func (db *DB) GetTransferAnalytics(start, end *time.Time) (*TransferAnalytics, error) {
	var total, blocked int
	var args []interface{}
	var where string
	if start != nil && end != nil {
		where = "WHERE when_ts >= $1 AND when_ts <= $2"
		args = append(args, *start, *end)
	} else if start != nil {
		where = "WHERE when_ts >= $1"
		args = append(args, *start)
	} else if end != nil {
		where = "WHERE when_ts <= $1"
		args = append(args, *end)
	}

	totalQuery := "SELECT COUNT(*) FROM transfers " + where
	if err := db.conn.Get(&total, totalQuery, args...); err != nil {
		return nil, err
	}

	blockedQuery := "SELECT COUNT(*) FROM transfers " + where
	if where == "" {
		blockedQuery += " WHERE is_blocked=true"
	} else {
		blockedQuery += " AND is_blocked=true"
	}
	if err := db.conn.Get(&blocked, blockedQuery, args...); err != nil {
		return nil, err
	}

	dailyStats := []TransferAnalyticsDayStats{}
	dailyQuery := `SELECT DATE(when_ts) as date, COUNT(*) as total, SUM(CASE WHEN is_blocked THEN 1 ELSE 0 END) as blocked, SUM(CASE WHEN NOT is_blocked THEN 1 ELSE 0 END) as successful FROM transfers ` + where + ` GROUP BY DATE(when_ts) ORDER BY DATE(when_ts)`
	rows, err := db.conn.Queryx(dailyQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var stat TransferAnalyticsDayStats
		if err := rows.Scan(&stat.Date, &stat.Total, &stat.Blocked, &stat.Successful); err != nil {
			return nil, err
		}
		dailyStats = append(dailyStats, stat)
	}

	return &TransferAnalytics{
		TotalTransfers:   total,
		BlockedTransfers: blocked,
		DailyStats:       dailyStats,
	}, nil
}
