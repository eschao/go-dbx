package dbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	assert := assert.New(t)

	// create table
	tDatabase.DropTable(USER_TABLE)
	tDatabase.DropTable(USER_LOGIN_TABLE)
	assert.Nil(tDatabase.CreateTable(USER_TABLE))
	assert.Nil(tDatabase.CreateTable(USER_LOGIN_TABLE))

	var err error
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[0])
	assert.Nil(err)
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[1])
	assert.Nil(err)
	_, err = tDatabase.T(USER_TABLE).Insert(&TestUsers[2])
	assert.Nil(err)
	_, err = tDatabase.T(USER_LOGIN_TABLE).Insert(&TestUserLogins[0])
	assert.Nil(err)
	_, err = tDatabase.T(USER_LOGIN_TABLE).Insert(&TestUserLogins[1])
	assert.Nil(err)
	_, err = tDatabase.T(USER_LOGIN_TABLE).Insert(&TestUserLogins[2])
	assert.Nil(err)

	user := User{}
	userLogin := UserLogin{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().
		InnerJoin(USER_LOGIN_TABLE, "userid", "userid").SelectAll().
		Filter("user.userid=?", TestUsers[1].Userid).One(&user, &userLogin))

	TestUsers[1].Id = user.Id
	TestUserLogins[1].Id = userLogin.Id
	assert.Equal(user, TestUsers[1])
	assert.Equal(userLogin, TestUserLogins[1])

	users := []User{}
	userLogins := []UserLogin{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().
		InnerJoin(USER_LOGIN_TABLE, "userid", "userid").SelectAll().
		All(&users, &userLogins))
	assert.Equal(len(users), 3)
	assert.Equal(len(userLogins), 3)
}
