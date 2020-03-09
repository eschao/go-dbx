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
	_, err = tDatabase.T(USER_TABLE).Update("id=?", user1.Id).Value(&user1)
	assert.Nil(err)
	// check updated row
	user2 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user1.Id).
		One(&user2))
	assert.Equal(user1, user2)

	// update some columns with row
	user2.Password = "password2"
	user2.UpdateTime = "2019-03-01 00:00:00"
	_, err = tDatabase.T(USER_TABLE).Update("id=?", user2.Id).
		Set("password", "update_time").Value(&user2)
	assert.Nil(err)
	// check updated row
	user3 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user2.Id).
		One(&user3))
	assert.Equal(user2, user3)

	// update some columns with values
	_, err = tDatabase.T(USER_TABLE).Update("id=?", user2.Id).
		Set("password", "update_time").Values("password3", "2019-04-01 00:00:00")
	assert.Nil(err)
	user4 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user2.Id).
		One(&user4))
	assert.Equal(user4.Password, "password3")
	assert.Equal(user4.UpdateTime, "2019-04-01 00:00:00")

	// update with value map
	valueMap := map[string]interface{}{
		"password":    "password4",
		"update_time": "2019-05-01 00:00:00",
	}
	_, err = tDatabase.T(USER_TABLE).Update("id=?", user2.Id).ValueMap(valueMap)
	assert.Nil(err)
	user5 := User{}
	assert.Nil(tDatabase.T(USER_TABLE).SelectAll().Filter("id=?", user2.Id).
		One(&user5))
	assert.Equal(user5.Password, "password4")
	assert.Equal(user5.UpdateTime, "2019-05-01 00:00:00")
}
