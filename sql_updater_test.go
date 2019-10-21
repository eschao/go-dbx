package dbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	// create table
	tDatabase.DropTable(USER_TABLE)
	err := tDatabase.CreateTable(USER_TABLE)
	assert.Nil(err)

	// insert rows
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	assert.Nil(err)

	// select row
	user1 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().
		Filter("userid=?", TestUsers[0].Userid).One(&user1))

	// update all columns
	user1.Nickname = "nickname1"
	user1.Password = "password1"
	user1.UpdateTime = "2019-02-01 00:00:00"
	_, err = tDatabase.T(USER_TABLE).Update("id=?", user1.Id).With(&user1)
	assert.Nil(err)
	// check updated row
	user2 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user1.Id).
		One(&user2))
	assert.Equal(user1, user2)

	// update some columns
	user2.Password = "password2"
	user2.UpdateTime = "2019-03-01 00:00:00"
	_, err = tDatabase.T(USER_TABLE).Update("id=?", user2.Id).
		Columns("password", "update_time").With(&user2)
	assert.Nil(err)
	// check updated row
	user3 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user2.Id).
		One(&user3))
	assert.Equal(user2, user3)
}
