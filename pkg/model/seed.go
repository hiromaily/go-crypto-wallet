package model

import (
	"github.com/jmoiron/sqlx"
	"time"
)

//Seed seedテーブル
type Seed struct {
	ID        uint8      `db:"id"`
	Seed      string     `db:"seed"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// GetSeedAll seedテーブル全体を返す(しかし、1行しかない想定)
func (m *DB) GetSeedAll() ([]Seed, error) {
	var seeds []Seed
	err := m.RDB.Select(&seeds, "SELECT * FROM seed")

	return seeds, err
}

// GetSeedOne idが１のseedを返す
func (m *DB) GetSeedOne() (Seed, error) {
	var seed Seed
	err := m.RDB.Get(&seed, "SELECT * FROM seed WHERE id=1")

	return seed, err
}

// GetSeedCount レコード数を返す
func (m *DB) GetSeedCount() (int64, error) {
	var count int64
	err := m.RDB.Get(&count, "SELECT count(id) FROM seed")

	return count, err
}

// InsertSeed レコードをinsertする
func (m *DB) InsertSeed(seed string, tx *sqlx.Tx, isCommit bool) (int64, error) {

	sql := `
INSERT INTO seed (seed, updated_at) 
VALUES (:seed, :updated_at)
`

	t := time.Now()
	seedRecord := Seed{
		Seed:      seed,
		UpdatedAt: &t,
	}

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	res, err := tx.NamedExec(sql, seedRecord)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if isCommit {
		tx.Commit()
	}
	id, _ := res.LastInsertId()

	return id, err
}
