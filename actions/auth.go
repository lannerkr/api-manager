package actions

import (
	"api-manager/models"
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	c.Set("a_user", "")

	return c.Render(200, r.HTML("html/auth_new.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) (err error) {
	u := &models.User{}
	var adminusers map[string]string
	if err = c.Bind(u); err != nil {
		return errors.WithStack(err)
	}
	if adminusers, err = getUser(u.Email); err != nil {
		log.Fatal(err)
	}
	u.ID = u.Email
	u.PasswordHash = adminusers[u.ID]

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		c.Set("a_user", "")
		verrs := validate.NewErrors()
		verrs.Add("email", "사용자 계정/비밀번호가 잘못 입력되었습니다.")
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("html/auth_new.html"))
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}

	c.Session().Set("current_user_id", u.ID)

	welcome := "Welcome Back " + u.ID
	c.Flash().Add("success", welcome)

	return c.Redirect(302, "/home")
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	// c.Cookies().Delete("user_id")
	c.Flash().Add("success", "정상적으로 로그아웃 되었습니다!")
	return c.Redirect(302, "/signin")
}
