package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	// co "api-gateway-dc/controllers"

	goredis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/ipsusila/opt"
	"github.com/joho/godotenv"
	"github.com/nitishm/go-rejson/v4"
)

// type DataWsUser struct {
// 	Token   string      `json:"token"`
// 	Data    interface{} `json:"data"`
// 	Request string      `json:request`
// 	Uuid    uuid.UUID   `json:"-"`
// }
type Error string

type DataAuth struct {
	Token string `json:"token"`
}

type Token struct {
	Id          uint
	Name        string
	Root        int
	Office      int
	Departement int
	Company     int
	jwt.StandardClaims
}

type RequestWs struct {
	Payload interface{} `json:"payload"`
	Request string      `json:"request"`
}

type ResponKafka struct {
	Payload interface{} `json:"payload"`
	Request string      `json:"request"`
	Token   string      `json:"token"`
	Output  interface{} `json:"output"`
}

type MessageWs struct {
	Data    *ResponKafka `json:"data"`
	Status  bool         `json:"status"`
	Message string       `json:"message"`
}

type UserContext struct {
	Id          uint
	Name        string
	Root        int
	Office      int
	Departement int
	Company     int
}

var TestMsg string

// var DataAuth2 []DataWsUser
var configTokenExpiration int
var configAppName string
var secretKey []byte
var configBroker string
var configTopicRespon string
var configTopicRequest string
var configGroup string
var rh *rejson.Handler
var cli *goredis.Client

const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"

var Bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// var DataRequest []interface{}
var DataRequest map[string]interface{}

func init() {
	// with loggger lumberjack
	ex, err := os.Executable()
	if err != nil {
		log.Println("Error Executable: ", err)
	}
	exPath := filepath.Dir(ex) + "/"
	confPath := exPath + "config.hjson"
	// re-open file
	file, err := os.Open(confPath)
	if err != nil {
		confPath = "config.hjson"
	}
	defer file.Close()

	// get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		log.Println("Failed to load secret key file -> ", readErr)
		return
	}

	secretKey = key

	//parse configurationf file
	cfgFile := flag.String("conf", confPath, "Configuration file")
	flag.Parse()
	log.Print("masuk config")

	//load options
	config, err := opt.FromFile(*cfgFile, opt.FormatAuto)
	if err != nil {
		log.Printf("Error while loading configuration file %v -> %v\n", *cfgFile, err)
		return
	}

	DataAuth2 := make([]DataAuth, 0)
	log.Printf("%v", DataAuth2)

	// DataRequest := make(map[string]interface{}, 0)
	// log.Printf("DataRequest : %v", DataRequest)

	//app_name config
	configAppName = config.Get("server").GetString("appName", "api-gateway-dc")
	log.Println("init() -> configAppName: ", configAppName)

	//token_expiration config
	wib := 7 * 60
	configTokenExpiration = config.Get("server").GetInt("tokenExpiration", 15) + wib
	log.Println("init() -> configTokenExpiration: ", configTokenExpiration)

	configBroker = config.Get("messagebroker").GetString("brokerServer", "")
	configTopicRespon = config.Get("messagebroker").GetString("topicRespon", "")
	configGroup = config.Get("messagebroker").GetString("brokerGroup", "")
	configTopicRequest = config.Get("messagebroker").GetString("topicRequest", "")
	log.Printf("brokerServer 2 : %v", configBroker)

	//REDIS =============================================================
	var addr = flag.String("Server", "localhost:6379", "Redis server address")
	rh = rejson.NewReJSONHandler()
	flag.Parse()

	cli = goredis.NewClient(&goredis.Options{Addr: *addr})

	// //testing
	// ctx := context.Background()
	// test := cli.Do(ctx, "FT.SEARCH", "user-json-idx", "@status:1").Val()
	// log.Printf("masuk 1111")
	// msge, err := json.Marshal(test)
	// if err != nil {
	// 	fmt.Printf("Error : %s\n", err)
	// }
	// log.Printf(string(msge))

	// //-----

	rh.SetGoRedisClient(cli)

	//redis ============================================
}

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

func GetTokenExpiration() int {
	return configTokenExpiration
}

func GetAppName() string {
	return configAppName
}

func GetJwtSecretKey() []byte {
	return secretKey
}

func GetBroker() string {
	return configBroker
}

func GetTopicRespon() string {
	return configTopicRespon
}

func GetTopicRequest() string {
	return configTopicRequest
}

func GetGroup() string {
	return configGroup
}

func GetReJson() *rejson.Handler {
	return rh
}

func GetRedis() *goredis.Client {
	return cli
}

func Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case int:
		x := int(reply)
		if int(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case float64:
		x := int(reply)
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case nil:
		return 0, err
	case Error:
		return 0, err
	}
	return 0, fmt.Errorf("redigo: unexpected type for Int, got type %T", reply)
}
