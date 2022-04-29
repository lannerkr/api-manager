package actions

import (
	"fmt"

	htp "github.com/foomo/htpasswd"
	"golang.org/x/crypto/bcrypt"
)

func getUser(name string) (u map[string]string, err error) {
	file := configuration.Htpasswd
	var adminusers map[string]string
	if adminusers, err = htp.ParseHtpasswdFile(file); err != nil {
		fmt.Println(err)
	}
	//fmt.Println(adminusers)
	return adminusers, nil
}

func Htp_change(name, cpwd, npwd string) (err error) {
	file := configuration.Htpasswd
	var adminusers map[string]string
	if adminusers, err = htp.ParseHtpasswdFile(file); err != nil {
		fmt.Println(err)
	}
	for u, p := range adminusers {
		if name == u {
			if CheckPasswordHash(cpwd, p) {
				// passwdUpdate := "htpasswd -iB -C 10 " + file + " " + name + " " + npwd
				// cmd, output := exec.Command("sh", "-c", passwdUpdate), new(strings.Builder)
				// cmd.Stdout = output
				// cmd.Run()
				if err := htp.SetPassword(file, name, npwd, htp.HashBCrypt); err != nil {
					return err
				}
				return nil
			} else {
				err = fmt.Errorf("현재 비밀번호가 맞지 않습니다")
				return err
			}
		}
	}
	err = fmt.Errorf("사용자 " + name + "(이)가 존재하지 않습니다")
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
