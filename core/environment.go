package core

var devicesArray []string

type bettercapDevices []map[string]map[string]string

type DeviceListEvents map[int]map[string]string

var serverInit = false
var sniffingInit = false
var passwordAsk = false
var apiUrlCpl = ""
var programDirectory = ""
var startedCorrect = false

var pcapName = ""
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var (
	UserPassword       = "51abb9636078defbf888d8457a7c76f85c8f114c" // Password (sha1)
	UserPasswordAsking = false                                      // Ask password before run
	UserDebug          = true                                       // Enable debug verbose (fmt.println())
	userWebService     = false                                      // (pwd) : /webserver/
	AskFlag            = true
)

var (
	ApiUSERNAME = "PhysicsBotnet" // Bettercap API username
	ApiPASSWORD = "MyPassw0rds"   // Bettercap API password
	apiPORT     = 3030            // Bettercap API port
	apiCAPTPL   = "api_rest.txt"  // /template/api_rest.txt (.cap) starter template
)
