package entities

import (
	"fmt"
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Card struct {
	Id        int  `db:"iid"`
	ClientId  int  `db:"iclientid"`
	IsBlocked bool `db:"bblocked"`
	IsTest    bool `db:"btest"`
}

type GeneratedCard struct {
	CardNum string
	CVC     string
}

func GetCardById(tx *sqlx.Tx, id int, lock bool) (Card, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := Card{}

		forUpdStr := ""
		if lock {
			forUpdStr = " FOR UPDATE"
		}

		err := tx.Get(&crd, `select iid, iclientid, btest, bblocked from ls.tcards where iid = $1`+forUpdStr, id)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		return crd, err
	}, tx)
	return res.(Card), err
}

func (crd *Card) GetClient(tx *sqlx.Tx) (Client, error) {
	return GetClientById(tx, crd.ClientId)
}

func ExtractCardNum(fullNum string) (string, error) {
	if len(fullNum) < 13 || len(fullNum) > 16 {
		return "", errors.WrongFormatErr
	}
	return fullNum[0:13], nil
}

func getLastCardId(tx *sqlx.Tx) (int, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		id := ""
		rows, err := tx.Query("SELECT max(sCardNum) FROM ls.tcards")

		if err != nil {
			log.Println(err)
			return id, err
		}
		defer rows.Close()

		if rows.Next() {
			err = rows.Scan(&id)
			if err != nil {
				log.Println(err)
				return id, err
			}

			id = id[4:12]
		} else {
			return id, errors.NotFoundErr
		}

		return id, err
	}, tx)
	id, err := strconv.Atoi(res.(string))
	return id, err
}

func calcCheckSum(str string) (int, error) {
	sum := 0
	for i, char := range str {
		tmp, err := strconv.Atoi(string(char))
		if err != nil {
			return 0, err
		}

		if i%2 == 1 {
			tmp = tmp * 3
		}
		sum += tmp
	}
	return (10 - sum%10) % 10, nil
}

func getNewCardNum(id int) (string, error) {
	cardDefaultMask := "1248"

	sid := fmt.Sprintf("%08d", id)

	res := ""

	//конкатенация с маской карты
	res += cardDefaultMask
	//конкатенация с номером карты
	res += sid

	//конрольный разряд
	sum, err := calcCheckSum(res) //EAN13 checksum
	if err != nil {
		return "", err
	}
	res += strconv.Itoa(sum)
	return res, nil
}

func generateCVC() string {
	return strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
}

func generateNew() (GeneratedCard, error) {
	lastCardId, err := getLastCardId(nil)
	if err != nil {
		return GeneratedCard{}, err
	}

	num, err := getNewCardNum(lastCardId + 1)
	if err != nil {

	}

	cvc := generateCVC()
	return GeneratedCard{CardNum: num, CVC: cvc}, nil
}

func getExpDate(from time.Time) time.Time {
	return from.AddDate(2, 0, 0)
}

func createCard(tx *sqlx.Tx, num string, expDate time.Time, temp bool, test bool, cvc string, clientId int, virtual bool, blocked bool) error {
	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := Card{}
		acc := Account{}

		rows, err := tx.Query("INSERT INTO ls.tCardAccounts (nBonuses, nBlockedBonuses, btest) VALUES (0.0, 0.0, $1) RETURNING iID;", test)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&acc.Id)
			if err != nil {
				log.Println(err)
				return crd, err
			}
		}
		dtbound := pq.NullTime{Time: time.Now(), Valid: true}
		if !virtual {
			dtbound = pq.NullTime{Time: time.Now(), Valid: false}
		}
		rows, err = tx.Query(`INSERT INTO ls.tCards (sCardNum, dtExpired, bBlocked, bTemporary, btest, sCVC, 
									    	iclientid, bvirtual, dtbound) 
									  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING iID;`, num, expDate, blocked, temp,
			test, cvc, clientId, virtual, dtbound)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&crd.Id)
			if err != nil {
				log.Println(err)
				return crd, err
			}
		}

		_, err = tx.Exec(`INSERT INTO ls.tcardaccounts_cards (icardid, icardaccid) VALUES ($1, $2) RETURNING iid;`,
			crd.Id, clientId)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		return nil, err
	}, tx)

	return err
}

func GenerateCardOnline(tx *sqlx.Tx, virtual bool, clientId int) error {
	gc, err := generateNew()
	if err != nil {
		return err
	}
	expDate := getExpDate(time.Now())
	err = createCard(tx, gc.CardNum, expDate, false, false, gc.CVC, clientId, virtual, false)
	return err
}

func GetCardByNum(tx *sqlx.Tx, num string, blockForUpdate bool) (Card, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := Card{}
		forUpdStr := ""
		if blockForUpdate {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&crd, "select iid, iclientid, btest, bblocked from ls.tcards where scardnum = $1"+forUpdStr, num)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		return crd, err
	}, tx)
	return res.(Card), err
}

func GetCardlistByClient(tx *sqlx.Tx, clientId int) ([]Card, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := []Card{}
		err := tx.Select(&crd, "select iid, iclientid, btest, bblocked from ls.tcards where iclientid = $1", clientId)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		return crd, err
	}, tx)
	return res.([]Card), err
}
