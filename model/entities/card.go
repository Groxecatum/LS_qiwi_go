package entities

import (
	"fmt"
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Card struct {
	Id       int `db:"iid"`
	ClientId int `db:"iclientid"`
}

type GeneratedCard struct {
	CardNum string
	CVC     string
}

func GetCardById(tx *sqlx.Tx, id int) (Card, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := Actor{}
		err := tx.Get(&crd, `select * from ls.tcards where iid = $1`, id)
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
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
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
	_, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := Card{}
		acc := Account{}

		rows, err := tx.Query("INSERT INTO ls.tCardAccounts (nBonuses, nBlockedBonuses, btest) VALUES (0.0, 0.0, $1) RETURNING iID;", num)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		if rows.Next() {
			err := rows.Scan(&acc.Id)
			if err != nil {
				log.Println(err)
				return crd, err
			}
		}

		return nil, err
	}, tx)

	return err
	//new VoidSqlCallbackExecutor((conn, stmt, rs) -> {

	//	Pattern pattern = Pattern.compile("^[0-9]{13}$");
	//	Matcher matcher = pattern.matcher(cardNumber);
	//	if (!matcher.matches()) {
	//		throw new CardNumberFormatException(cardNumber);
	//	}
	//
	//	stmt = conn.prepareStatement("INSERT INTO ls.tCardAccounts (nBonuses, nBlockedBonuses, btest) VALUES (0.0, 0.0, ?) RETURNING iID;");
	//	stmt.setBoolean(1, test);
	//	rs = stmt.executeQuery();
	//	if (!rs.next()) {
	//		throw new SimpleException(cardNumber.concat("card account has not created!"));
	//	}
	//	int cardAccID = rs.getInt(1);
	//	stmt.close();
	//
	//	stmt = conn.prepareStatement("INSERT INTO ls.tCards (sCardNum, dtExpired, bBlocked, bTemporary, btest, sCVC, iclientid, bvirtual, dtbound) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING iID;");
	//	stmt.setString(1, cardNumber);
	//	java.sql.Date dt = new java.sql.Date(expDate.getTime());
	//	stmt.setDate(2, dt);
	//	stmt.setBoolean(3, blocked);
	//	stmt.setBoolean(4, temporary);
	//	stmt.setBoolean(5, test);
	//	stmt.setString(6, CVC);
	//	stmt.setInt(7, clientId);
	//	stmt.setBoolean(8, virtual == null ? false : virtual);
	//	if (virtual == null ? false : virtual){
	//		stmt.setTimestamp(9,new Timestamp(new Date().getTime()));
	//	} else {
	//		stmt.setNull(9,Types.TIMESTAMP);
	//	}
	//	rs = stmt.executeQuery();
	//	if (!rs.next()) {
	//		throw new SimpleException(cardNumber.concat("Card has not been bound!"));
	//	}
	//	int cardID = rs.getInt(1);
	//	stmt.close();
	//
	//	stmt = conn.prepareStatement("INSERT INTO ls.tcardaccounts_cards (icardid, icardaccid) VALUES (?, ?) RETURNING iid;");
	//	stmt.setInt(1, cardID);
	//	stmt.setInt(2, cardAccID);
	//	rs = stmt.executeQuery();
	//	if (!rs.next()) {
	//		throw new SimpleException(cardNumber.concat("Realtions Card <-> Account has not been created!"));
	//	}
	//}, extConn).execute();
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
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := Card{}
		forUpdStr := ""
		if blockForUpdate {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&crd, "select * from ls.tcards where iid = $1"+forUpdStr, num)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		return crd, err
	}, tx)
	return res.(Card), err
}
