package dbx

const (
	USER_TABLE       = "user"
	USER_LOGIN_TABLE = "user_login"
	USER_OAUTH_TABLE = "user_oauth"
)

type User struct {
	Id         int64  `json:"id"          db:"id"          sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT" mysql:"int NOT NULL PRIMARY KEY AUTO_INCREMENT"`
	Userid     string `json:"userid"      db:"userid"      sqlite:"TEXT NOT NULL"                     mysql:"varchar(32) NOT NULL"`
	Nickname   string `json:"nickname"    db:"nickname"    sqlite:"TEXT"                              mysql:"varchar(64) NOT NULL DEFAULT ''"`
	Password   string `json:"password"    db:"password"    sqlite:"TEXT"                              mysql:"varchar(32) NOT NULL DEFAULT ''"`
	UpdateTime string `json:"update_time" db:"update_time" sqlite:"INTEGER"                           mysql:"datetime NOT NULL DEFAULT '2000-01-01 00:00:00'"`
}

type UserLogin struct {
	Id         int64  `json:"id"          column:"id"          sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT" mysql:"int NOT NULL PRIMARY KEY AUTO_INCREMENT"`
	Userid     string `json:"userid"      column:"userid"      sqlite:"TEXT UNIQUE NOT NULL"              mysql:"varchar(32) NOT NULL UNIQUE"`
	OAuthId    string `json:"oauth_id"    column:"oauth_id"    sqlite:"TEXT UNIQUE NOT NULL"              mysql:"varchar(64) NOT NULL UNIQUE DEFAULT ''"`
	LastLogin  string `json:"last_login"  column:"last_login"  sqlite:"INTEGER"                           mysql:"datetime NOT NULL DEFAULT '2000-01-01 00:00:00'"`
	LastIP     int64  `json:"last_ip"     column:"last_ip"     sqlite:"INTEGER"                           mysql:"int NOT NULL DEFAULT ''"`
	UpdateTime string `json:"update_time" column:"update_time" sqlite:"INTEGER"                           mysql:"datetime NOT NULL DEFAULT '2000-01-01 00:00:00'"`
}

type UserOAuth struct {
	Id         int64  `json:"id"          column:"id"          sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT" mysql:"int NOT NULL PRIMARY KEY AUTO_INCREMENT"`
	Userid     string `json:"userid"      column:"userid"      sqlite:"TEXT UNIQUE NOT NULL"              mysql:"varchar(32) NOT NULL UNIQUE"`
	OAuthId    string `json:"oauth_id"    column:"oauth_id"    sqlite:"TEXT UNIQUE NOT NULL"              mysql:"varchar(64) NOT NULL UNIQUE"`
	App        string `json:"app"         column:"app"         sqlite:"TEXT"                              mysql:"varchar(16) NOT NULL DAEFAULT ''"`
	Url        string `json:"url"         column:"url"         sqlite:"TEXT"                              mysql:"varchar(256) NOT NULL DEFAULT ''"`
	Token      string `json:"token"       column:"token"       sqlite:"TEXT"                              mysql:"varchar(64) NOT NULL DEFAULT ''"`
	ExpireTime string `json:"expire_time" column:"expire_time" sqlite:"INTEGER"                           mysql:"datetime NOT NULL DEFAULT '2000-01-01 00:00:00'"`
	UpdateTime string `json:"update_time" column:"update_time" sqlite:"INTEGER"                           mysql:"datetime NOT NULL DEFAULT '2000-01-01 00:00:00'"`
}

var TestUsers = []User{
	{
		Id:         -1,
		Userid:     "15600362000",
		Nickname:   "eschao",
		Password:   "15600362000",
		UpdateTime: "2019-01-01 00:00:00",
	},
	{
		Id:         -1,
		Userid:     "12520343000",
		Nickname:   "chaozh",
		Password:   "12520343000",
		UpdateTime: "2019-01-02 00:00:00",
	},
	{
		Id:         -1,
		Userid:     "12901060000",
		Nickname:   "zc",
		Password:   "12901060000",
		UpdateTime: "2019-01-03 00:00:00",
	},
}

var TestUserLogins = []UserLogin{
	{
		Id:         -1,
		Userid:     "15600362000",
		OAuthId:    "qq_15600362000",
		LastLogin:  "2019-07-01 00:00:00",
		LastIP:     1024,
		UpdateTime: "2019-07-01 00:00:00",
	},
	{
		Id:         -1,
		Userid:     "12520343000",
		OAuthId:    "wechat_12520343000",
		LastLogin:  "2019-07-02 00:00:00",
		LastIP:     2048,
		UpdateTime: "2019-07-02 00:00:00",
	},
	{
		Id:         -1,
		Userid:     "12901060000",
		OAuthId:    "weibo_12901060000",
		LastLogin:  "2019-07-03 00:00:00",
		LastIP:     3096,
		UpdateTime: "2019-07-03 00:00:00",
	},
}

var TestUserOAuths = []UserOAuth{
	{
		Id:         -1,
		Userid:     "15600362000",
		OAuthId:    "qq_15600362000",
		App:        "qq",
		Url:        "qq_url",
		Token:      "qq_token",
		ExpireTime: "2020-07-01 00:00:00",
		UpdateTime: "2019-07-01 00:00:00",
	},
	{
		Id:         -1,
		Userid:     "12520343000",
		OAuthId:    "wechat_12520343000",
		App:        "wechat",
		Url:        "wechat_url",
		Token:      "wechat_token",
		ExpireTime: "2020-07-02 00:00:00",
		UpdateTime: "2019-07-02 00:00:00",
	},
	{
		Id:         -1,
		Userid:     "12901060000",
		OAuthId:    "weibo_12901060000",
		App:        "weibo",
		Url:        "weibo_url",
		Token:      "weibo_token",
		ExpireTime: "2020-07-03 00:00:00",
		UpdateTime: "2019-07-03 00:00:00",
	},
}
