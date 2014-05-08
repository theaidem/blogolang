package user

import (
	"app/common"
	db "github.com/dancannon/gorethink"
	"github.com/martini-contrib/sessionauth"
	"time"
)

type User struct {
	Id            string `form:"-" gorethink:"id,omitempty"`
	Email         string `form:"email" gorethink:"email"`
	Password      string `form:"password" gorethink:"password"`
	Created       time.Time
	authenticated bool `form:"-" gorethink:"-"`
}

var sess *db.Session

func init() {
	sess = common.GetSess()
}

func GenerateAnonymousUser() sessionauth.User {
	return &User{}
}

func (u *User) Login() {
	u.authenticated = true
}

func (u *User) Logout() {
	u.authenticated = false
}

func (u *User) IsAuthenticated() bool {
	return u.authenticated
}

func (u *User) UniqueId() interface{} {
	return u.Id
}

func (u *User) GetById(id interface{}) error {
	row, err := db.Table("users").Get(id).RunRow(sess)
	if err != nil {
		return err
	}
	if !row.IsNil() {
		if err := row.Scan(&u); err != nil {
			return err
		}
	}
	return nil
}
