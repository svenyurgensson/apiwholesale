package system

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"time"
	"github.com/kylelemons/go-gypsy/yaml"
	"gopkg.in/mgo.v2"

	"strconv"
	"gopkg.in/mgo.v2/bson"
)

const Author = "Yury Batenko"
const Version = "0.0.1 / 2015-05-29"
const ApiVersion = "v1"
const ConfigFile = "./config/settings.yml"

const (
	Test          string = "test"
	Development   string = "development"
	Production    string = "production"
)

var (
	Boot_time     string
	Hostname      string
	Env           string

	DB            string
	ConnURL       string
	ConnOptions   string
	dbSession     *mgo.Session
)


func DEBUG(x interface{}){
	fmt.Printf("--> type: %T\n", x)
	fmt.Printf("--- value: %v\n", x)
}

func set_working_environment(){
	env := flag.String("e", "development", "Environment to run {test, development, production}, default: development")

	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `-e [environment]
Start server with environment [environment]
Environment string should be one of: {test, development, production}
Default environment is 'development'

`)
		flag.PrintDefaults()
	}

	flag.Parse()

	switch *env {
	case "test": Env = Test
	case "dev", "development": Env = Development
	case "prod", "production": Env = Production
	default:
		fmt.Println("Error! unknown environment: ", *env)
		os.Exit(1)
	}

	fmt.Println("set environment to: ", Env)
}


func load_parse_config(){
	config, err := yaml.ReadFile(ConfigFile)

	if err != nil {
		fmt.Println("Cannot load config! " + ConfigFile)
		fmt.Println(err)
		os.Exit(1)
	}

	DB, err = config.Get(Env + ".db.name")
	if err != nil {
		fmt.Println("Database name not present! " + ConfigFile)
		fmt.Println(err)
		os.Exit(1)
	}

	node := config.Root.(yaml.Map).Key(Env)
	hosts, err := yaml.Child(node, "db.hosts")
	if err != nil {
		fmt.Println("Database hosts not presents! " + ConfigFile)
		panic(err)
	}


	list := hosts.(yaml.List)
	conn_str_list := []string{}
	for i := 0; i < list.Len(); i++ {
		conn_str_list = append(conn_str_list, strings.TrimSpace(yaml.Render(list.Item(i))))
	}

	ConnURL = strings.Join(conn_str_list, ",") + "/" + DB

	ConnOptions, err = config.Get(Env + ".db.options")
	if err == nil {
		ConnURL = ConnURL + "?" + ConnOptions
	}
}

func GetSession() *mgo.Session {
	if dbSession == nil {
		db_establish_connection()
	}
	return dbSession.Clone()
}

func db_establish_connection() {
	if dbSession != nil {
		dbSession.Close()
	}

	fmt.Printf("(re)connecting to database: %s\n", ConnURL)

	var err error
	dbSession, err = mgo.Dial(ConnURL)
	if err != nil {
		panic(err)
	}
}

func Init(){
	set_working_environment()

	load_parse_config()

	db_establish_connection()

}




// ~~~~~~~ Helper functions
// https://gist.github.com/bsphere/8369aca6dde3e7b4392c
// https://medium.com/coding-and-deploying-in-the-cloud/time-stamps-in-golang-abcaf581b72f

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))

	return nil
}

func (t Timestamp) GetBSON() (interface{}, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}

	return time.Time(t), nil
}

func (t *Timestamp) SetBSON(raw bson.Raw) error {
	var tm time.Time

	if err := raw.Unmarshal(&tm); err != nil {
		return err
	}

	*t = Timestamp(tm)

	return nil
}

func (t *Timestamp) String() string {
	return time.Time(*t).String()
}

// https://gowalker.org/github.com/kylelemons/go-gypsy/yaml
