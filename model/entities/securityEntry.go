package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons"
	"github.com/jmoiron/sqlx"
)

type SecurityEntry struct {
	Id       int    `db:"iid"         xml:"id"`
	Hash     string `db:"shash"       xml:"shash"`
	Salt     string `db:"ssalt"       xml:"ssalt"`
	HashType int    `db:"sihashtype"  xml:"sihashtype"`
}

const (
	NONE   = 0
	BCRYPT = 1
	MD5    = 2
)

func CreateNewSecurityEntry(tx *sqlx.Tx, secureValue *string, hashType int) (SecurityEntry, error) {

	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		return SecurityEntry{}, nil
	}, tx)

	return res.(SecurityEntry), err
	//return (Security) new ObjectSqlCallbackExecutor((conn, stmt, rs) -> {
	//	String salt = generateSalt(hashType);
	//	String hash;
	//
	//	if (hashType == HashType.BCRYPT) {
	//		hash = BCrypt.hashpw(secureValue + salt, BCrypt.gensalt(13));
	//	} else if (hashType == HashType.MD5) {
	//		hash = getMD5(secureValue + salt);
	//	} else if (hashType == HashType.NONE) {
	//		hash = "00000000000000000000000000000000";
	//	} else {
	//		throw new SimpleException("Security.setNewSecureValue wrong hashType: " + hashType);
	//	}
	//
	//	if (hashType != 0) {
	//		stmt = conn.prepareStatement("INSERT INTO ls.tsecurity (shash, ssalt, siHashType) VALUES (?, ?, ?) RETURNING iid;");
	//	} else {
	//		stmt = conn.prepareStatement("INSERT INTO ls.tsecurity (shash, ssalt) VALUES (?, ?) RETURNING iid;");
	//	}
	//
	//	stmt.setString(1, hash);
	//	stmt.setString(2, salt);
	//
	//	if (hashType != 0) {
	//		stmt.setShort(3, hashType);
	//	}
	//
	//	rs = stmt.executeQuery();
	//
	//	Security security;
	//
	//	if (rs.next()) {
	//		security = new Security(rs.getInt(1), hash, salt, (hashType != 0) ? hashType : 1);
	//	} else {
	//		throw new SimpleException("Error during creating security field");
	//	}
	//
	//	rs.close();
	//	stmt.close();
	//
	//	return security;
	//},extConn).execute();

}
