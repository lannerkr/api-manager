package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}

func TestHandler(c buffalo.Context) error {
	c.Set("a_user", "test")

	// NewSessionHandler(c.Response(), c.Request())
	return c.Render(http.StatusOK, r.HTML("html/test.html"))
}

func AccessHandler(c buffalo.Context) error {
	c.Set("a_user", "")
	warning := "페이지 접근 권한이 없습니다!"
	c.Flash().Add("warning", warning)

	return c.Render(http.StatusOK, r.HTML("html/access.html"))
}

type SystemStatus struct {
	Sysdate    string
	Uptime     string
	Configdate string
	Licensed   string
	Current    string
	Maxlast    string
	Cpu        string
	Swap       string
	Disk       string
}

func DashboardHandler(c buffalo.Context) error {

	var cluster []ClusterMember
	var sysStatus SystemStatus
	var sysdate, uptime, configdate, licensed, current, maxlast, cpu, swap, disk string
	var err error

	if sysdate, uptime, configdate, cluster, err = getStatus(); err != nil {
		c.Flash().Add("warning", "sslvpn status update failed")
	}
	if licensed, current, maxlast, err = getStats(); err != nil {
		c.Flash().Add("warning", "sslvpn statistics update failed")
	}
	if cpu, swap, disk, err = getHealth(); err != nil {
		c.Flash().Add("warning", "sslvpn health status update failed")
	}

	sysStatus = SystemStatus{sysdate, uptime, configdate, licensed, current, maxlast, cpu, swap, disk}

	c.Set("systemStatus", sysStatus)
	c.Set("cluster", cluster)

	return c.Render(http.StatusOK, r.HTML("html/dashboard.html"))
}

func ActiveHandler(c buffalo.Context) error {
	count := c.Param("count")
	fmt.Println(count)
	if count == "" {
		count = "200"
	}
	c.Set("count", count)
	activeUsers, err := getActiveUsers("all", count)
	if err != nil {
		c.Flash().Add("warning", "sslvpn Active-Users update failed")
		c.Redirect(301, "/")
	}

	c.Set("userRecords", activeUsers)

	return c.Render(http.StatusOK, r.HTML("html/mainpage.html"))
}

func StoreHandler(c buffalo.Context) error {
	realm := "store"
	user, days, err := getUsersdata(realm)
	if err != nil {
		c.Flash().Add("warning", "sslvpn Store-Users update failed")
		c.Redirect(301, "/")
	}
	c.Set("newUsers", user)
	c.Set("newDays", days)

	return c.Render(http.StatusOK, r.HTML("html/storeusers.html"))
}

func PartnerHandler(c buffalo.Context) error {
	realm := "partner"
	user, days, err := getUsersdata(realm)
	if err != nil {
		c.Flash().Add("warning", "sslvpn Partner-Users update failed")
		c.Redirect(301, "/")
	}
	c.Set("newUsers", user)
	c.Set("newDays", days)

	return c.Render(http.StatusOK, r.HTML("html/partnerusers.html"))
}

type User struct {
	Name  string `form:"name"`
	IP    string `form:"ip"`
	Mac   string `form:"mac"`
	Realm string `form:"realm"`
}

func CreateUserHandler(c buffalo.Context) error {

	user := &User{}
	if err := c.Bind(user); err != nil {
		return err
	}
	// fmt.Println(user)

	if user.Name != "" {
		err1, err2, err3 := userCreate(user.Realm, user.Name, user.IP, user.Mac)
		if err1 != nil {
			warning := "sslvpnuser " + err1.Error()
			c.Flash().Add("warning", warning)
		} else {
			inform := "sslvpnuser : " + user.Name + " is created"
			c.Flash().Add("success", inform)
		}
		if err2 != nil {
			if strings.Contains(err2.Error(), "empty") {
				goto Error3
			}
			warning := "static-ip user on profiler " + err2.Error()
			c.Flash().Add("warning", warning)
		} else {
			inform := "static-ip for user : " + user.Name + " is created"
			c.Flash().Add("success", inform)
		}
	Error3:
		if err3 != nil {
			if strings.Contains(err3.Error(), "empty") {
				goto Fin
			}
			warning := "mac-address on profiler " + err3.Error()
			c.Flash().Add("warning", warning)
		} else {
			inform := "mac-address : " + user.Mac + " is added"
			c.Flash().Add("success", inform)
		}
	Fin:
	}

	return c.Render(http.StatusOK, r.HTML("html/createuser.html"))
}

func ApproveHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		c.Flash().Add("warning", warning)
	} else if err := macApprove(user_mac); err != nil {
		c.Flash().Add("warning", err.Error())
	} else {
		inform := "mac-address : " + user_mac + " 승인처리가 완료 되었습니다!"
		c.Flash().Add("success", inform)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	// c.Redirect(302, "/")
	return nil
}
func UnapproveHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		c.Flash().Add("warning", warning)
	} else if err := macUnApprove(user_mac); err != nil {
		c.Flash().Add("warning", err.Error())
	} else {
		inform := "mac-address : " + user_mac + " 미승인처리 되었습니다!"
		c.Flash().Add("success", inform)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	// c.Redirect(302, "/")
	return nil
}
func PermitHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		c.Flash().Add("warning", warning)
	} else if err := macPermit(user_mac); err != nil {
		c.Flash().Add("warning", err.Error())
	} else {
		inform := "mac-address : " + user_mac + " USB 사용이 허용 되었습니다!"
		c.Flash().Add("success", inform)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	// c.Redirect(302, "/")
	return nil
}
func ProtectHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		c.Flash().Add("warning", warning)
	} else if err := macProtect(user_mac); err != nil {
		c.Flash().Add("warning", err.Error())
	} else {
		inform := "mac-address : " + user_mac + " USB 사용이 차단 되었습니다!"
		c.Flash().Add("success", inform)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	// c.Redirect(302, "/")
	return nil
}

func StoreUserHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "store"
	var singleDay string
	var singleUserRecord Records

	user, days, err := getUserdata(realm, user_id)
	if err != nil {
		fmt.Println(err)
		c.Redirect(301, "/")
	}
	activeUsers, err := getActiveUsers(user_id, "")
	if err != nil {
		fmt.Println(err)
		c.Redirect(301, "/")
	}

	if user.Username == "" {
		c.Redirect(301, "/")
	}
	singleDay = days[user.Username]

	if len(activeUsers) != 0 {
		singleUserRecord = activeUsers[0]
	} else {
		singleUserRecord = Records{"", "", "", "", "", "", ""}
	}

	staticIP := getStaticIP(realm, user_id)
	mac := readlog(user_id)
	userLog := userlog(user_id, mac)

	c.Set("userLog", userLog)
	c.Set("staticIP", staticIP)
	c.Set("singleUserRecord", singleUserRecord)
	c.Set("singleUser", user)
	c.Set("singleDay", singleDay)
	c.Set("realm", realm)

	return c.Render(http.StatusOK, r.HTML("html/user.plush.html"))
}

func PartnerUserHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "partner"
	var singleDay string
	var singleUserRecord Records

	user, days, err := getUserdata(realm, user_id)
	if err != nil {
		fmt.Println(err)
		c.Redirect(301, "/")
	}
	activeUsers, err := getActiveUsers(user_id, "")
	if err != nil {
		fmt.Println(err)
		c.Redirect(301, "/")
	}

	if user.Username == "" {
		c.Redirect(301, "/")
	}
	singleDay = days[user.Username]

	if len(activeUsers) != 0 {
		singleUserRecord = activeUsers[0]
	} else {
		singleUserRecord = Records{"", "", "", "", "", "", ""}
	}

	staticIP := getStaticIP(realm, user_id)
	mac := readlog(user_id)
	userLog := userlog(user_id, mac)

	c.Set("userLog", userLog)
	c.Set("staticIP", staticIP)
	c.Set("singleUserRecord", singleUserRecord)
	c.Set("singleUser", user)
	c.Set("singleDay", singleDay)
	c.Set("realm", realm)

	return c.Render(http.StatusOK, r.HTML("html/user.plush.html"))
}

func StoreDisableHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "store"
	updateStatus(realm, user_id, false)

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func StoreEnableHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "store"
	updateStatus(realm, user_id, true)

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func StoreResetHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "store"
	resetPW(realm, user_id)

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func PartnerDisableHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "partner"
	updateStatus(realm, user_id, false)

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func PartnerEnableHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "partner"
	updateStatus(realm, user_id, true)

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func PartnerResetHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := "partner"
	resetPW(realm, user_id)

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func PWDunlockHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	if err := deleteUser(realm, user_id); err == nil {
		if err1 := addUserPulse(realm, user_id); err1 != nil {
			fmt.Println(err1)
			warning := "sslvpnuser " + err1.Error()
			c.Flash().Add("warning", warning)
		} else {
			inform := "사용자 : " + user_id + " 의 비밀번호 잠김 해제 및 초기화가 완료 되었습니다!"
			c.Flash().Add("success", inform)
		}
	} else {
		fmt.Println(err)
		warning := "sslvpnuser " + err.Error()
		c.Flash().Add("warning", warning)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func DeleteSessionHandler(c buffalo.Context) error {
	session_id := c.Param("session_id")
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	respCode := deleteActiveUser(session_id)
	if respCode == 204 {
		inform := "사용자 " + user_id + " 로그아웃이 완료 되었습니다!"
		c.Flash().Add("success", inform)
	} else {
		warning := "사용자 " + user_id + " 로그아웃이 정상 처리되지 못했습니다. response code : " + strconv.Itoa(respCode)
		c.Flash().Add("warning", warning)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func TOTPresetHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	// fmt.Println(user_id)
	if err := totpReset(user_id); err != nil {
		c.Flash().Add("warning", err.Error())
	} else {
		inform := "TOTP user " + user_id + " under Authserver '01.PCS-TOTP' has been reset"
		c.Flash().Add("success", inform)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func TOTPunlockHandler(c buffalo.Context) error {
	user_id := c.Param("user_id")
	realm := c.Param("realm")
	// fmt.Println(user_id)
	if err := totpUnlock(user_id); err != nil {
		c.Flash().Add("warning", err.Error())
	} else {
		inform := "TOTP user " + user_id + " under Authserver '01.PCS-TOTP' has been unlocked"
		c.Flash().Add("success", inform)
	}

	c.Redirect(302, "/user/"+realm+"/"+user_id)
	return nil
}

func StaticIPHandler(c buffalo.Context) (err error) {
	realm := c.Param("realm")
	var users []PPSUserdata

	if users, err = getStaticUsers(realm); err != nil {
		return err
	}

	c.Set("staticusers", users)
	c.Set("realm", realm)

	return c.Render(http.StatusOK, r.HTML("html/staticip.html"))
}
