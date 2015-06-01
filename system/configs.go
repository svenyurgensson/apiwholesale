package system

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"github.com/kylelemons/go-gypsy/yaml"
	"gopkg.in/mgo.v2"
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
	DBSession     * mgo.Session
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
	err:= 22
}

func db_establish_connection(){
	fmt.Printf("connecting to database: %s\n", ConnURL)

	session, err := mgo.Dial(ConnURL)
	if err != nil {
		panic(err)
	}
	DBSession = session
}

func Init(){
	set_working_environment()

	load_parse_config()

	db_establish_connection()

}


// https://gowalker.org/github.com/kylelemons/go-gypsy/yaml
