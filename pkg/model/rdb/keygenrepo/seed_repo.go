package keygenrepo

import (
	"time"

	"github.com/jmoiron/sqlx"
)

//Seed seedテーブル
type Seed struct {
	ID        uint8      `db:"id"`
	Seed      string     `db:"seed"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// GetSeedAll seedテーブル全体を返す(しかし、1行しかない想定)
func (r *KeygenRepository) GetSeedAll() ([]Seed, error) {
	var seeds []Seed
	err := r.db.Select(&seeds, "SELECT * FROM seed")

	return seeds, err
}

// GetSeedOne idが１のseedを返す
func (r *KeygenRepository) GetSeedOne() (Seed, error) {
	var seed Seed
	err := r.db.Get(&seed, "SELECT * FROM seed WHERE id=1")

	return seed, err
}

// GetSeedCount レコード数を返す
func (r *KeygenRepository) GetSeedCount() (int64, error) {
	var count int64
	err := r.db.Get(&count, "SELECT count(id) FROM seed")

	return count, err
}

// InsertSeed レコードをinsertする
func (r *KeygenRepository) InsertSeed(seed string, tx *sqlx.Tx, isCommit bool) (int64, error) {

	sql := `
INSERT INTO seed (seed, updated_at) 
VALUES (:seed, :updated_at)
`
	//logger.Debugf("sql: %s", sql)

	t := time.Now()
	seedRecord := Seed{
		Seed:      seed,
		UpdatedAt: &t,
	}

	if tx == nil {
		tx = r.db.MustBegin()
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
