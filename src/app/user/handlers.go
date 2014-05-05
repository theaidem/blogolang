package user

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	db "github.com/dancannon/gorethink"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"net/http"
	"time"
)

func GetSignUpHandler(user sessionauth.User, r render.Render) {
	if user.IsAuthenticated() {
		r.Redirect("/")
		return
	}
	r.HTML(200, "signup", nil)
}

func PostSignUpHandler(session sessions.Session, newUser User, r render.Render, req *http.Request, sess *db.Session) {
	if session.Get(sessionauth.SessionKey) != nil {
		fmt.Println("Logged in already! Logout first.")
		r.Redirect("/")
		return
	}
	var userInDb User
	query := db.Table("users").Filter(db.Row.Field("email").Eq(newUser.Email))
	row, err := query.RunRow(sess)
	if err == nil && !row.IsNil() {
		// Register, error case.
		if err := row.Scan(&userInDb); err != nil {
			fmt.Println("Error reading DB")
		} else {
			fmt.Println("User already exists. Redirecting to login.")
		}

		r.Redirect(sessionauth.RedirectUrl)
		return
	} else { // User doesn't exist, continue with registration.
		if row.IsNil() {
			fmt.Println("User doesn't exist. Registering...")
		} else {
			fmt.Println(err)
		}
	}

	// Try to compare passwords
	pass1Hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	pass2String := req.FormValue("confirmpassword")
	passErr := bcrypt.CompareHashAndPassword(pass1Hash, []byte(pass2String))

	if passErr != nil {
		fmt.Println("Error, passwords don't match.", passErr)
	} else { // passwords are the same, insert user to db
		newUser.Password = string(pass1Hash)
		newUser.Created = time.Now().Local()
		db.Table("users").Insert(newUser).RunWrite(sess)
		fmt.Println("Register done. Try to login.")
	}
	r.Redirect(sessionauth.RedirectUrl)
}

func GetSignInHandler(user sessionauth.User, r render.Render) {
	if user.IsAuthenticated() {
		r.Redirect("/")
		return
	}
	r.HTML(200, "signin", nil)
}

func PostSignInHandler(session sessions.Session, userLoggingIn User, r render.Render, req *http.Request, sess *db.Session) {
	if session.Get(sessionauth.SessionKey) != nil {
		r.Redirect("/")
		return
	}

	var userInDb User
	query := db.Table("users").Filter(db.Row.Field("email").Eq(userLoggingIn.Email))
	row, err := query.RunRow(sess)

	if err == nil && !row.IsNil() {
		if err := row.Scan(&userInDb); err != nil {
			fmt.Println("Error scanning user in DB")
			r.Redirect(sessionauth.RedirectUrl)
			return
		}
	} else {
		if row.IsNil() {
			fmt.Println("User doesn't exist")
		} else {
			fmt.Println(err)
		}
		r.Redirect(sessionauth.RedirectUrl)
		return
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(userInDb.Password), []byte(userLoggingIn.Password))
	if passErr != nil {
		fmt.Println("Wrong Password")
		r.Redirect(sessionauth.RedirectUrl)
	} else {
		err := sessionauth.AuthenticateSession(session, &userInDb)
		if err != nil {
			fmt.Println("Wrong Auth")
			r.JSON(500, err)
		}
		params := req.URL.Query()
		redirect := params.Get(sessionauth.RedirectParam)
		r.Redirect(redirect)
	}

	r.HTML(200, "signin", nil)
}

func GetSignOutHandler(session sessions.Session, user sessionauth.User, r render.Render) {
	sessionauth.Logout(session, user)
	r.Redirect("/")
}
