package user

import (
	"app/common"
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	db "github.com/dancannon/gorethink"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strconv"
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

func GetUserListHandler(user sessionauth.User, r render.Render, params martini.Params, sess *db.Session, req *http.Request) {
	arg_map := map[string]interface{}{"authuser": user}

	// Paginate, todo: to helpers
	page, ok := params["page"]
	if !ok {
		page = "1"
	}

	curr_page, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println(err)
		r.HTML(404, "404", arg_map)
		return
	}

	if curr_page <= 0 {
		r.HTML(404, "404", arg_map)
		return
	}

	row, err := db.Table("users").Count().RunRow(sess)
	if err != nil {
		fmt.Println(err)
		r.HTML(404, "404", arg_map)
		return
	}
	var total_items int
	err = row.Scan(&total_items)

	if err != nil {
		fmt.Println(err)
		r.HTML(404, "404", arg_map)
		return
	}

	var pager = common.GetPaginated(total_items, 10, curr_page)

	rows, _ := db.Table("users").WithFields("id", "email", "created").Skip(pager.StartPoint).Limit(pager.PerPage).Run(sess)
	var users []User
	for rows.Next() {
		var p User
		err := rows.Scan(&p)
		if err != nil {
			fmt.Println(err)
			return
		}
		users = append(users, p)
	}

	arg_map["users"] = users
	arg_map["pager"] = pager
	r.HTML(200, "users", arg_map)
}

// we'll just wrap a slice of ints
type items struct {
	Stuff []string
}

func GetUserProfileHandler(user sessionauth.User, r render.Render, params martini.Params, sess *db.Session) {

	id := params["id"]
	arg_map := map[string]interface{}{"authuser": user}

	row, _ := db.Table("users").Get(id).RunRow(sess)
	if row.IsNil() {
		fmt.Println("Not found!")
		r.HTML(404, "404", arg_map)
		return
	}

	var u User
	err := row.Scan(&u)
	if err != nil {
		fmt.Println(err)
		return
	}
	arg_map["currentuser"] = u
	r.HTML(200, "profile", arg_map)
}

func GetUserDeleteHandler(user sessionauth.User, r render.Render, params martini.Params, sess *db.Session) {
	//arg_map := map[string]interface{}{"authuser": user}
	query := db.Table("users").Get(params["id"]).Delete()
	row, err := query.RunWrite(sess)
	if err != nil {
		fmt.Println("Wrong query")
		fmt.Println(err)
		return
	}
	fmt.Println(row)
	r.Redirect("/users")

}
