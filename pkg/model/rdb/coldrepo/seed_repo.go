package coldrepo

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

// GetSeedOne idが１のseedを返す
func (r *ColdRepository) GetSeedOne() (Seed, error) {
	var seed Seed
	err := r.db.Get(&seed, "SELECT * FROM seed WHERE id=1")

	return seed, err
}

// InsertSeed レコードをinsertする
func (r *ColdRepository) InsertSeed(seed string, tx *sqlx.Tx, isCommit bool) (int64, error) {

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
