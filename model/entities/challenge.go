package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Challenge struct {
	Id        int       `db:""`
	DataValue string    `db:""`
	ExpDate   time.Time `db:""`
}

type ChallengeData struct {
	VerificationCode string `xml:"verificationCode,attr"`
	TrnRequestId     int64  `xml:"trnRequestId,attr"`
}

const (
	EXPIRED  = 0
	ACTIVE   = 1
	FINISHED = 2
)

func (challenge *Challenge) Check(tx *sqlx.Tx) bool {
	expired := challenge.ExpDate.After(time.Now())
	if expired {
		challenge.FinishChallenge(tx)
	}
	return !expired
}

func GetChallengeByHash(tx *sqlx.Tx, hash string) (Challenge, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		c := Challenge{}

		err := tx.Get(&c, `select iid, sdata, dtexpired from ls.tchallenges where shash = $1 AND sistate = ?`,
			hash, ACTIVE)
		if err != nil {
			log.Println(err)
			return c, err
		}

		return c, err
	}, tx)
	return res.(Challenge), err
}

func (challenge *Challenge) FinishChallenge(tx *sqlx.Tx) error {
	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		_, err := tx.Exec(`UPDATE ls.tchallenges SET sistate = ? WHERE iid = ?;`, EXPIRED, challenge.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return nil, err
	}, tx)

	return err
}
