package main

// version 1.01
// readlog grep function patch
// version 1.1
// status-page add
// version 1.2
// API response message customize
// Redirection customize
// version 1.2.1
// Redirection custiomize continue...
// default active-users to 200
// version 1.2.2
// start mode to production mode
// use TLS and add certificates
// customize HTML header and footer process
// version 2.0
// version update
// version 2.01
// image update adding apache-util
// version 2.02
// image update removing apache-util
// htpasswd update customize
// version 2.1
// sign-in process add
// ip access control add
// access user ip logging
// OTP-unlock add
// hide "delete session" if user's not logged in
// static ip config display
// coloring user log javascript myFunction() in user.plush.html
// access user and ip logging add
// PWD-unlock add ( delete user and create user )
// logging to file
// version 2.1.1
// embeded source customize
// version 2.2
// admin menu name add [AuthNew Handler, SetCurrentUser Handler]
// CookiHandler add
// remove /templates from embed
// ----------------------------
// git upload
// ----------------------------
// version 2.2.1
// admin menu name customize
// USB-Permit/USB-Protect/MAC-Unapprove function add
// make userlog include MAC status change PPS-[ADM31591]
// version 2.2.2
// button click comfirm massage [user.plush.html customize]
// version 2.2.3
// emp user create add [createuser.html customize]
// emp staticip view add [dashboard.html customize]

// version 3.0
// emp user modify
// version 3.1
// realm page Redesign , user page Redesign
// version 3.1.1
// user create bug fix
// version 3.1.2
// <div container> style customize (max width 1000 -> 1500) at [application.plush.html]
// userTable page optimize
// version 3.2
// userTable ReDesign
// version 3.2.1
// singleuserTable bug fix
// version 3.2.2
// bug fix and optimize
// change userHistory Days type int64 -> int [userTables.go]
// add timeNow at [UserTableHandler]
// last_login customize at [usertable.html]
// change usertable to order by last_login at [application.plush.html]
// emp user create modify [createUser.go]
// EMPmacApprove customize [ApproveHandler] [EMPmacApprove()]
// version 3.2.3
// userTable performance optimize [UserTableHandler] [userTables.go] [usertable.html]

import (
	"log"
	"net/http"
	"os"

	"api-manager/actions"

	"github.com/gobuffalo/buffalo/servers"
)

// main is the starting point for your Buffalo application.
// You can feel free and add to this `main` method, change
// what it does, etc...
// All we ask is that, at some point, you make sure to
// call `app.Serve()`, unless you don't want to start your
// application that is. :)

func main() {
	conf := os.Args[1]
	cert, certKey, logpath := actions.Config(conf)

	fpLog, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	if actions.Logfile(fpLog); err != nil {
		panic(err)
	}
	defer fpLog.Close()

	app := actions.App()

	serv := servers.Simple{Server: &http.Server{}}
	s := servers.WrapTLS(serv.Server, cert, certKey)
	if err := app.Serve(s); err != nil {
		log.Fatal(err)
	}

	// if err := app.Serve(); err != nil {
	// 	log.Fatal(err)
	// }
}

/*
# Notes about `main.go`

## SSL Support

We recommend placing your application behind a proxy, such as
Apache or Nginx and letting them do the SSL heavy lifting
for you. https://gobuffalo.io/en/docs/proxy

## Buffalo Build

When `buffalo build` is run to compile your binary, this `main`
function will be at the heart of that binary. It is expected
that your `main` function will start your application using
the `app.Serve()` method.

*/
