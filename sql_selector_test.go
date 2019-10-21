package dbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
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

	// select all users
	users := []User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().All(&users))
	assert.Equal(len(users), 3)

	// select one row
	user1 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().
		Filter("userid=?", TestUsers[0].Userid).One(&user1))
	TestUsers[0].Id = user1.Id
	assert.Equal(user1, TestUsers[0])

	// order by id asc
	usersAsc := []User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Asc("id").All(&usersAsc))
	assert.Equal(usersAsc[0].Userid, TestUsers[0].Userid)
	assert.Equal(usersAsc[1].Userid, TestUsers[1].Userid)
	assert.Equal(usersAsc[2].Userid, TestUsers[2].Userid)

	// order by id desc
	usersDesc := []User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Desc("id").All(&usersDesc))
	assert.Equal(usersDesc[0].Userid, TestUsers[2].Userid)
	assert.Equal(usersDesc[1].Userid, TestUsers[1].Userid)
	assert.Equal(usersDesc[2].Userid, TestUsers[0].Userid)

	// select some columns
	user2 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).Select("nickname", "password").
		Filter("userid=?", TestUsers[0].Userid).One(&user2))
	assert.Equal(user2.Nickname, TestUsers[0].Nickname)
	assert.Equal(user2.Password, TestUsers[0].Password)
	assert.NotEqual(user2.Userid, TestUsers[0].Userid)
	assert.NotEqual(user2.UpdateTime, TestUsers[0].UpdateTime)

	// select by page
	users1 := []User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Desc("id").Offset(1).
		Limit(10).All(&users1))
	assert.Equal(len(users1), 2)
	assert.Equal(users1[0].Userid, TestUsers[1].Userid)
	assert.Equal(users1[1].Userid, TestUsers[0].Userid)
}
