package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

type AdminUser struct {
	Username   string `form:"name"`
	CurrPwd    string `form:"cpwd"`
	Password   string `form:"npwd"`
	ConfirmPwd string `form:"cfpwd"`
}

func AdminHandler(c buffalo.Context) error {

	return c.Render(http.StatusOK, r.HTML("html/adminuserpw.html"))
}

func AdminpwHandler(c buffalo.Context) error {

	user := &AdminUser{}
	if err := c.Bind(user); err != nil {
		return err
	}

	//fmt.Println(user)
	if user.Username != "" && user.Password == user.ConfirmPwd {
		if err := Htp_change(user.Username, user.CurrPwd, user.Password); err != nil {
			warning := err.Error()
			c.Flash().Add("warning", warning)
		} else {
			inform := user.Username + "의 비밀번호가 변경되었습니다. 새로운 비밀번호를 사용해 로그인 하십시오"
			c.Flash().Add("success", inform)
			c.Redirect(302, "/signout")
		}
	} else if user.Username != "" && user.Password != "" {
		warning := "새로운 비밀번호가 일치하지 않습니다"
		c.Flash().Add("warning", warning)
	}

	return c.Render(http.StatusOK, r.HTML("html/adminuserpw.html"))
}
