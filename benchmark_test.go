package dbx

import (
	"fmt"
	"testing"

	sqlx "github.com/jmoiron/sqlx"
)

// go test -run=XXX -bench=.

func BenchmarkRawInsert(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)

	u := TestUsers[0]
	db := tDatabase.DB()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		q := "INSERT INTO user(userid,nickname,password,update_time) VALUES(?,?,?,?)"
		stmt, _ := db.Prepare(q)
		stmt.Exec(&u.Userid, &u.Nickname, &u.Password, &u.UpdateTime)
		stmt.Close()
	}
}

func BenchmarkDbxInsert(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)

	t := tDatabase.T(USER_TABLE)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.Insert(&TestUsers[0])
	}
}

func BenchmarkRawSelectRow(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])

	u := User{}
	db := tDatabase.DB()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		q := "SELECT id, userid, nickname, password, update_time FROM user WHERE userid=?"
		r := db.QueryRow(q, TestUsers[0].Userid)
		r.Scan(&u.Id, &u.Userid, &u.Nickname, &u.Password, &u.UpdateTime)
	}
}

func BenchmarkDbxSelectRow(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])

	u := User{}
	t := tDatabase.T(USER_TABLE)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.SelectAll().Filter("userid=?", TestUsers[0].Userid).One(&u)
	}
}

func BenchmarkSqlxSelectRow(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])

	db, err := sqlx.Connect("sqlite3", TEST_DB_FILE)
	if err != nil {
		fmt.Printf("Can't open db for sqlx: %s\n", err.Error())
		return
	}

	u := User{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		//q := "SELECT id, userid, nickname, password, update_time FROM user WHERE userid=$1"
		q := "SELECT * FROM user WHERE userid=$1"
		db.Get(&u, q, TestUsers[0].Userid)
	}
}

func BenchmarkRawSelectRows(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	for i := 0; i < 20; i++ {
		tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	}

	users := []User{}
	db := tDatabase.DB()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		q := "SELECT id, userid, nickname, update_time FROM user WHERE userid=?"
		r, _ := db.Query(q, TestUsers[0].Userid)
		for r.Next() {
			u := User{}
			r.Scan(&u.Id, &u.Userid, &u.Nickname, &u.UpdateTime)
			users = append(users, u)
		}
		r.Close()
	}
}

func BenchmarkDbxSelectRows(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	for i := 0; i < 20; i++ {
		tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	}

	users := []User{}
	t := tDatabase.T(USER_TABLE)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.SelectAll().Filter("userid=?", TestUsers[0].Userid).All(&users)
	}
}

func BenchmarkSqlxSelectRows(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	for i := 0; i < 20; i++ {
		tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	}

	db, err := sqlx.Connect("sqlite3", TEST_DB_FILE)
	if err != nil {
		fmt.Printf("Can't open db for sqlx: %s\n", err.Error())
		return
	}

	users := []User{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		q := "SELECT id, userid, nickname, update_time FROM user WHERE userid=$1"
		db.Select(&users, q, TestUsers[0].Userid)
	}
}

func BenchmarkRawUpdate(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])

	u := TestUsers[1]
	u.Userid = TestUsers[0].Userid
	db := tDatabase.DB()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		q := "UPDATE user SET userid=?,nickname=?,password=?,update_time=? WHERE userid=?"
		db.Exec(q, &u.Userid, &u.Nickname, &u.Password, &u.UpdateTime, &u.Userid)
	}
}

func BenchmarkDbxUpdate(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])

	u := TestUsers[1]
	u.Userid = TestUsers[0].Userid
	t := tDatabase.T(USER_TABLE)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.Update("userid=?", u.Userid).
			Columns("userid", "nickname", "password", "update_time").With(&u)
	}
}

func BenchmarkRawInnerJoin(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.CreateTable(USER_LOGIN_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	tDatabase.T(USER_LOGIN_TABLE).Insert(&TestUserLogins[0])

	u := User{}
	uLogin := UserLogin{}
	db := tDatabase.DB()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		q := `
			SELECT user.id,user.userid,user.nickname,user.password,user.update_time, 
			user_login.id,user_login.userid,user_login.oauth_id,user_login.last_login, 
			user_login.last_ip,user_login.update_time FROM user INNER JOIN user_login 
			ON user.userid=user_login.userid WHERE user.userid=?
		`
		r := db.QueryRow(q, TestUsers[0].Userid)
		r.Scan(&u.Id, &u.Userid, &u.Nickname, &u.Password, &u.UpdateTime,
			&uLogin.Id, &uLogin.Userid, &uLogin.OAuthId, &uLogin.LastLogin,
			&uLogin.LastIP, &uLogin.UpdateTime)
	}
}

func BenchmarkDbxInnerJoin(b *testing.B) {
	DBLogger = nil
	tDatabase.DropTable(USER_TABLE)
	tDatabase.CreateTable(USER_TABLE)
	tDatabase.CreateTable(USER_LOGIN_TABLE)
	tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	tDatabase.T(USER_LOGIN_TABLE).Insert(&TestUserLogins[0])

	u := User{}
	uLogin := UserLogin{}
	t := tDatabase.T(USER_TABLE)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t.SelectAll().
			InnerJoin(USER_LOGIN_TABLE, "userid", "userid").SelectAll().
			Filter("user.userid=?", TestUsers[0].Userid).One(&u, &uLogin)
	}
}
