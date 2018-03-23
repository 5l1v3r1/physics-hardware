package core

import (
	"net"
	"net/http"
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"time"
	"os"
	"strings"
	"runtime"
	"path"
)

func getLocal() (string, bool){
	addrs, _ := net.InterfaceAddrs()

	for _, networdAddress := range addrs {

		addr, status := networdAddress.(*net.IPNet)

		if status && !addr.IP.IsLoopback(){
			if addr.IP.To4() != nil {
				return addr.IP.String(), true
			}
		}
	}
	return "", false
}

func generateTpl() bool{

	_, programFile, _, _ := runtime.Caller(0)
	programDirectory := path.Dir(programFile)
	fileName := programDirectory + "/../webserver/index.tpl"
	fileInformation, errors := ioutil.ReadFile(fileName)

	if errors != nil {

		fmt.Println(P_ERROR("Error with packet file generation"))
		os.Exit(1)
	} else {

		lines := strings.Split(string(fileInformation), "\n")
		configuration := 0
		for key, line := range lines {

			if strings.Contains(line, "{{PCAP_NAME}}") {

				lines[key] = strings.Replace(line, "{{PCAP_NAME}}",  pcapName, 2)
				configuration += 1
			}
		}

		file, err := os.OpenFile(programDirectory + "/../webserver/index.html", os.O_RDWR|os.O_CREATE, 0777)
		errChmod := os.Chmod(programDirectory + "/../webserver/index.html", 0777)

		if errChmod != nil {
			fmt.Println(p_WARNING("file can't be writable."))
		}
		if err != nil {
			fmt.Println(err)
		}

		file.WriteString(strings.Join(lines, "\n"))
		file.Close()
		if userDebug {
			fmt.Println(P_SUCCESS("PCAP name Configuration as been successfully ended"))
		}
	}
	return false
}

func loadModules(api_url string) {
	if serverInit == false && userWebService {

		sendCmd(api_url, "set http.server.path "+ programDirectory +"/../webserver/")
		sendCmd(api_url, "http.server on ")
		fmt.Println(P_TIME(GREEN, "Web server has been started."))
		//generateTpl()
		serverInit = true
	}
	if sniffingInit == false {

		_, err := os.OpenFile(programDirectory + "/../webserver/output/" + pcapName, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			fmt.Println(P_ERROR("Can't write PCAP file."))
			os.Exit(1)
		}
		errChmod := os.Chmod(programDirectory + "/../webserver/output/"+pcapName, 0777)

		if errChmod != nil {
			fmt.Println(P_ERROR("pcap file can't be writable."))
			os.Exit(1)
		}
		sendCmd(api_url, "set net.sniff.output "+ programDirectory +"/../webserver/output/"+pcapName)
		sendCmd(api_url, "set net.sniff.verbose false")
		sendCmd(api_url, "set net.sniff.local true")
		sendCmd(api_url, "net.sniff on")
		fmt.Println(P_TIME(GREEN, "Sniffing has been started."))
		sniffingInit = true
	}
}

func getEvent(api_url string, nb string) {
	apiUrlCpl = api_url
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", "https://"+ api_url +":8083/api/events?n=" + string(nb), nil)
	if err != nil {
		fmt.Println(P_ERROR("We can't connect to wifi (1)."))
		loadInteraction()
	}
	req.SetBasicAuth(ApiUSERNAME, ApiPASSWORD)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(P_ERROR("We can't connect to wifi (2)."))
		loadInteraction()
	} else {
		startedCorrect = true
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var devices bettercapDevices

		json.Unmarshal([]byte(body), &devices)
		time.Sleep(1000 * time.Millisecond)
		loadModules(api_url)
		loadInteraction()

	}

}

func sendGET(api_url string, command string) []byte{

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", "https://"+ api_url +":8083" + command, nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth(ApiUSERNAME, ApiPASSWORD)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		return []byte(body)

	}

	return []byte("Error")
}

func sendCmd(api_url string, command string) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	type Payload struct {
		Cmd string `json:"cmd"`
	}

	data := Payload{
		Cmd: command,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
	} else {
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", "https://" + api_url + ":8083/api/session", body)
		if err != nil {
			fmt.Printf("server not responding %s", err.Error())
		} else {
			req.SetBasicAuth(ApiUSERNAME, ApiPASSWORD)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("server not responding %s", err.Error())
			} else {
				defer resp.Body.Close()
			}
		}
	}
}

