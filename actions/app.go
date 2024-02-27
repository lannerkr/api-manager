package actions

import (
	"net/http"
	"os"

	"api-manager/locales"
	"api-manager/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n/v2"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
// var ENV = envy.Get("GO_ENV", "development")
var ENV = envy.Get("GO_ENV", "production")

var (
	app *buffalo.App
	T   *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_ilmosslvpn_session",
			//SessionStore: sessions.NewCookieStore([]byte("SESSION_SECRET")),
			Logger: JSONLogger(logger.DebugLevel),
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		//app.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		app.Use(translations())
		app.Use(UserIPMiddleware)
		app.Use(CookieHandler)
		app.Middleware.Skip(UserIPMiddleware, AccessHandler, TestHandler)

		app.GET("/test", TestHandler)

		app.GET("/user", UserHandler)
		app.GET("/", AuthNew)
		app.GET("/access", AccessHandler)
		//app.GET("/", DashboardHandler)

		app.Use(SetCurrentUser)
		app.Use(Authorize)
		// app.GET("/users/new", UsersNew)
		// app.POST("/users", UsersCreate)
		app.GET("/signin", AuthNew)
		app.POST("/signin", AuthCreate)
		app.DELETE("/signout", AuthDestroy)
		app.GET("/signout", AuthDestroy)
		app.Middleware.Skip(Authorize, HomeHandler, UserHandler, AuthNew, AuthCreate, AccessHandler, TestHandler)

		app.GET("/admin", AdminHandler)
		app.POST("/adminpw", AdminpwHandler)

		app.GET("/home", DashboardHandler)
		app.GET("/active", ActiveHandler)
		app.GET("/active/{count}", ActiveHandler)

		app.GET("/user/create", CreateUserHandler)
		app.GET("/user/upper/users", UserUpperHandler)
		app.GET("/user/{realm}", UserTableHandler)
		app.GET("/user/{realm}/{user_id}", RealmUserHandler)

		app.GET("/user/approve/{realm}/{user_id}", ApproveHandler)
		app.GET("/user/unapprove/{realm}/{user_id}", UnapproveHandler)
		app.GET("/user/permit/{realm}/{user_id}", PermitHandler)
		app.GET("/user/protect/{realm}/{user_id}", ProtectHandler)

		app.GET("/user/status/{realm}/{user_id}", UserStatusHandler)
		app.GET("/user/resetPassword/{realm}/{user_id}", PwdResetHandler)
		app.GET("/user/unlockPwd/{realm}/{user_id}", PWDunlockHandler)

		app.GET("/user/delete/{realm}/{user_id}/{session_id}", DeleteSessionHandler)
		app.GET("/user/totp/{realm}/{user_id}", TOTPresetHandler)
		app.GET("/user/totpunlock/{realm}/{user_id}", TOTPunlockHandler)

		app.GET("/admin/staticip/{realm}", StaticIPHandler)

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

var logfile *os.File

func Logfile(fplog *os.File) error {
	logfile = fplog
	return nil
}

func JSONLogger(lvl logger.Level) logger.FieldLogger {
	l := logrus.New()
	l.Level = lvl
	//l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(logfile)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger.Logrus{FieldLogger: l}
}

func CookieHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Session().Session.Options.HttpOnly = true
		c.Session().Session.Options.Secure = true
		c.Session().Session.Options.SameSite = 3 //3:Strict, 1:None, 2:Lax
		c.Session().Session.Options.MaxAge = 3600

		return next(c)
	}
}
