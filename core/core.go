package core

import (
	"bufio"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/graniet/physics-hardware/core/command"
)

func CheckRoot() {
	if os.Getenv("SUDO_USER") == "" {

		fmt.Println(P_ERROR("Please run software as root."))
		os.Exit(1)
	}
}

func checkPasswd() {

	input := bufio.NewScanner(os.Stdin)
	user_input := ""

	for passwordAsk == false {
		fmt.Print("(", P_SIMPLE(YELLOW, "physics-iot")+") please enter password: ")
		input.Scan()
		user_input = input.Text()

		h2 := sha1.New()
		h2.Write([]byte(user_input))
		sha1_check := hex.EncodeToString(h2.Sum(nil))

		if sha1_check == UserPassword {
			passwordAsk = true
		} else {
			fmt.Print(P_SIMPLE(RED, "password denied\n"))
		}
	}

}

func checkInternet() {

	time.Sleep(1000 * time.Millisecond)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", "https://facebook.com", nil)
	if err != nil {

		fmt.Println(P_TIME(RED, "We can't connect to wifi."))
		checkInternet()
	}

	_, errDo := http.DefaultClient.Do(req)
	if errDo != nil {

		fmt.Println(P_TIME(RED, "We can't connect to wifi."))
		checkInternet()
	} else {

		fmt.Println(P_TIME(GREEN, "Internet connection found"))
	}
}

func GenerateCap() {

	// Get physics directory
	_, programFile, _, _ := runtime.Caller(0)
	programDirectory = path.Dir(programFile)
	if command.CoreDirectory == "" {
		command.CoreDirectory = programDirectory
	}

	// Check if .cap exist
	_, err := os.Stat(programDirectory + "/../template/" + apiCAPTPL)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {

		completPath := programDirectory + "/../template/api_rest.txt"
		fileInformation, errors := ioutil.ReadFile(completPath)

		if errors != nil {

			P_ERROR("Error has been occurred.")
		} else {

			lines := strings.Split(string(fileInformation), "\n")
			configuration := 0
			for key, line := range lines {

				if strings.Contains(line, "{{USERNAME}}") {

					lines[key] = strings.Replace(line, "{{USERNAME}}", ApiUSERNAME, 1)
					configuration += 1
				} else if strings.Contains(line, "{{PASSWORD}}") {

					lines[key] = strings.Replace(line, "{{PASSWORD}}", ApiPASSWORD, 1)
					configuration += 1
				}
			}

			if configuration >= 2 {
				file, err := os.OpenFile(programDirectory+"/../template/start.cap", os.O_RDWR|os.O_CREATE, 0777)
				errChmod := os.Chmod(programDirectory+"/../template/start.cap", 0777)

				if errChmod != nil {
					fmt.Println(p_WARNING("file can't be writable."))
				}
				if err != nil {
					fmt.Println(err)
				}

				file.WriteString(strings.Join(lines, "\n"))
				file.Close()
				if userDebug {
					fmt.Println(P_SUCCESS("Configuration as been successfully ended"))
				}
			}
		}

	}
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Running() {

	CheckRoot()

	//CheckFlag

	if !checkFlags() {
		os.Exit(1)
	}

	// Check password
	if UserPasswordAsking {
		checkPasswd()
	}

	// Check device internet connection
	checkInternet()

	// Get a bettercap process id
	out, err := exec.Command("pgrep", "-f", "bettercap").Output()
	if err != nil {
		fmt.Println(err)
	}

	output := string(out)
	pids := strings.Split(output, "\n")

	// Check a count of process id
	if len(pids) < 2 {

		// if PID do not exists we start bettercap
		if userDebug {
			fmt.Println(p_WARNING("Bettercap Swiss army knife not started please wait..."))
		}
		GenerateCap()

		_, programFile, _, _ := runtime.Caller(0)
		programDirectory := path.Dir(programFile)
		cmd := exec.Command("bettercap", "-caplet", programDirectory+"/../template/start.cap", "-silent")

		// Configure output to /dev/null
		null, _ := os.Open(os.DevNull)
		defer null.Close()
		cmd.Stdout = null
		//cmd.Stdout = os.Stdout

		err := cmd.Start()
		if err != nil {

			// If bettercap not found
			fmt.Println(P_ERROR("Bettercap not found please 'go get github.com/bettercap/bettercap'"))
			os.Exit(1)
		}

		// Restart running function without this statement
		Running()
	} else {

		if pcapName == "" {

			pcapName = randSeq(6) + ".pcap"
		}
		generateTpl()
		localIP, err := getLocal()

		if err == false {

			fmt.Println(p_WARNING("We can't find local address"))
		}
		if userDebug {
			fmt.Println(P_INFO("Checking devices on network"))
		}
		time.Sleep(1000000 * time.Microsecond)
		command.SendEvent("Physics hardware started", "")
		getEvent(localIP, "50")

	}

}
