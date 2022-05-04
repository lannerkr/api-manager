package actions

import (
	"os/exec"
	"strings"
)

func readlog(user string) (mac string) {

	apilog := configuration.ApiLog
	readCmd := "tac " + apilog + " | grep PRO31457 | grep -m 1 " + user
	cmd, output := exec.Command("sh", "-c", readCmd), new(strings.Builder)
	cmd.Stdout = output
	cmd.Run()
	userProfile := strings.Split(output.String(), ",")

	if len(userProfile) >= 4 {
		device_attN1 := strings.Split(userProfile[3], ")")
		device_attN2 := strings.Split(device_attN1[0], "(")
		mac = device_attN2[1]
		mac = strings.Replace(mac, "-", ":", -1)
	}

	return mac
}

func userlog(user, mac string) (userLog string) {
	apilog := configuration.ApiLog
	readCmd := "tac " + apilog + " | grep -m 15 -E \"" + user + "|ADM31591.*" + mac + "\" | awk '{$1=\"\";$2=\"\";$3=\"\";$4=\"\";$5=\"\";print}'"
	cmd, output := exec.Command("sh", "-c", readCmd), new(strings.Builder)
	cmd.Stdout = output
	cmd.Run()

	userLog = output.String()
	//fmt.Println(userLog)

	// var newlog string
	// coloringLog := strings.Split(userLog, "\n")
	// for _, l := range coloringLog {
	// 	if strings.Contains(l, "Account Still LockedOut") {
	// 		slog := strings.Split(l, "Account Still LockedOut")
	// 		newlog = ""
	// 		newlog = newlog + slog[0] + "<strong style=\"color: red;\">Account Still LockedOut</strong>" + slog[1]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	} else if strings.Contains(l, "failed") {
	// 		slog := strings.Split(l, "failed")
	// 		newlog = ""
	// 		slen := len(slog) - 1
	// 		for i := 0; i < slen; i++ {
	// 			newlog = newlog + slog[i] + "<strong style=\"color: orange;\">failed</strong>"
	// 		}
	// 		newlog = newlog + slog[slen]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	} else if strings.Contains(l, "unlocked") {
	// 		slog := strings.Split(l, "unlocked")
	// 		newlog = ""
	// 		slen := len(slog) - 1
	// 		for i := 0; i < slen; i++ {
	// 			newlog = newlog + slog[i] + "<strong style=\"color: blue;\">unlocked</strong>"
	// 		}
	// 		newlog = newlog + slog[slen]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	} else if strings.Contains(l, "ocked") {
	// 		slog := strings.Split(l, "- ")
	// 		newlog = slog[0] + "- " + slog[1] + "- <strong style=\"color: red;\">" + slog[2] + "</strong> -" + slog[3]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	} else if strings.Contains(l, "Modified") {
	// 		slog := strings.Split(l, "Modified")
	// 		newlog = ""
	// 		slen := len(slog) - 1
	// 		for i := 0; i < slen; i++ {
	// 			newlog = newlog + slog[i] + "<strong style=\"color: blue;\">Modified</strong>"
	// 		}
	// 		newlog = newlog + slog[slen]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	} else if strings.Contains(l, "Created") {
	// 		slog := strings.Split(l, "Created")
	// 		newlog = ""
	// 		slen := len(slog) - 1
	// 		for i := 0; i < slen; i++ {
	// 			newlog = newlog + slog[i] + "<strong style=\"color: blue;\">Created</strong>"
	// 		}
	// 		newlog = newlog + slog[slen]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	} else if strings.Contains(l, "address cannot be allocated") {
	// 		slog := strings.Split(l, "address cannot be allocated")
	// 		newlog = ""
	// 		newlog = newlog + slog[0] + "<strong style=\"color: orange;\">address cannot be allocated</strong>" + slog[1]
	// 		userLog = strings.Replace(userLog, l, newlog, -1)
	// 	}
	// }

	return userLog
}
