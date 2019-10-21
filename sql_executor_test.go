package dbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertCountDelete(t *testing.T) {
	assert := assert.New(t)

	// create table
	tDatabase.DropTable(USER_TABLE)
	err := tDatabase.CreateTable(USER_TABLE)
	assert.Nil(err)

	// insert rows
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	assert.Nil(err)
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[1])
	assert.Nil(err)
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[2])
	assert.Nil(err)

	// count all
	n, err := tDatabase.T(USER_TABLE).CountAll()
	assert.Nil(err)
	assert.Equal(n, 3)

	// count with query
	n, err = tDatabase.T(USER_TABLE).Count("userid=?", TestUsers[0].Userid)
	assert.Nil(err)
	assert.Equal(n, 1)

	// count with empty
	n, err = tDatabase.T(USER_TABLE).Count("userid=?", "xx")
	assert.Nil(err)
	assert.Equal(n, 0)

	// select row
	user1 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().
		Filter("userid=?", TestUsers[0].Userid).One(&user1))
	TestUsers[0].Id = user1.Id
	assert.Equal(user1, TestUsers[0])

	// replace row
	user1.Nickname = "new_nickname"
	user1.Password = "new_password"
	user1.UpdateTime = "2019-02-01 00:00:00"
	_, err = tDatabase.T(USER_TABLE).Replace(&user1)
	assert.Nil(err)
	user2 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user1.Id).
		One(&user2))
	assert.Equal(user1, user2)

	// delete with query
	assert.Nil(tDatabase.T(USER_TABLE).Delete("userid=?", TestUsers[0].Userid))
	n, err = tDatabase.T(USER_TABLE).Count("userid=?", TestUsers[0].Userid)
	assert.Nil(err)
	assert.Equal(n, 0)

}

func TestInsertCountDeleteWithTx(t *testing.T) {
	assert := assert.New(t)

	// create table
	tDatabase.DropTable(USER_TABLE)
	err := tDatabase.CreateTable(USER_TABLE)
	assert.Nil(err)

	// insert rows
	tx, err := tDatabase.Begin()
	assert.Nil(err)
	_, err = tx.T(USER_TABLE).Insert(&TestUsers[0])
	assert.Nil(err)
	_, err = tx.T(USER_TABLE).Insert(&TestUsers[1])
	assert.Nil(err)
	_, err = tx.T(USER_TABLE).Insert(&TestUsers[2])
	assert.Nil(err)
	assert.Nil(tx.Commit())

	// count all
	n, err := tDatabase.T(USER_TABLE).CountAll()
	assert.Nil(err)
	assert.Equal(n, 3)

	// delete with query
	tx, err = tDatabase.Begin()
	assert.Nil(err)
	assert.Nil(tDatabase.T(USER_TABLE).Delete("userid=?", TestUsers[0].Userid))
	assert.Nil(tx.Commit())

	n, err = tDatabase.T(USER_TABLE).Count("userid=?", TestUsers[0].Userid)
	assert.Nil(err)
	assert.Equal(n, 0)
}
