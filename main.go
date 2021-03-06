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
