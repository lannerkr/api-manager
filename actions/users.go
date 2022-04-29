package actions

import (
	"api-manager/models"
	"net"

	"github.com/gobuffalo/buffalo"
)

func UserHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("html/signin.html"))
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			c.Set("current_user", u)
			c.Set("a_user", uid)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "사용자 로그인이 필요합니다.!")
			return c.Redirect(302, "/")
		}
		c.LogField("adminUser", c.Session().Get("current_user_id"))
		// logrus.WithField("user", (c.Session().Get("current_user_id"))).Info("request is called by ")
		return next(c)
	}
}

func UserIPMiddleware(next buffalo.Handler) buffalo.Handler {
	ipaccess := configuration.IPaccess
	return func(c buffalo.Context) error {
		if xRealIP := c.Request().Header.Get("X-Real-Ip"); len(xRealIP) > 0 {
			//log.Println("xRealIP : " + xRealIP)
			for _, ip := range ipaccess {
				if xRealIP == ip {
					c.Set("user_ip", xRealIP)
					c.LogField("adminUserIP", xRealIP)
					return next(c)
				}
			}
			return c.Redirect(302, "/access")
		} else {
			h, _, err := net.SplitHostPort(c.Request().RemoteAddr)
			if err != nil {
				return err
			}
			//log.Println("H : " + h)
			for _, ip := range ipaccess {
				if h == ip {
					c.Set("user_ip", h)
					c.LogField("adminUserIP", h)
					// logrus.WithField("ip", h).Info("access is accepted for ")
					return next(c)
				}
			}
		}

		return c.Redirect(302, "/access")
	}
}
