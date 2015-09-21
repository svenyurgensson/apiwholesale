package system

import (
    "fmt"
    "os"
    "flag"
    "strings"

    "github.com/kylelemons/go-gypsy/yaml"
    "time"

    "log/syslog"

    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
)

const Author = "Yury Batenko"
const Version = "1.3.5 / 2015-09-22"
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

    Log          *syslog.Writer

    RequestsTotal  int = 0
    RequestsFailed int = 0

    DB            string
    ConnURL       string
    ConnOptions   string
    dbSession     *mgo.Session

)


func DEBUG(x interface{}) {
    fmt.Printf("--> type: %T\n", x)
    fmt.Printf("--- value: %v\n", x)
}

func set_working_environment() {

    time.Local = time.UTC

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


func load_parse_config() {
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
    dbSession, err = mgo.DialWithTimeout(ConnURL, 3*time.Second)
    if err != nil {
        panic(err)
    }

    mgo.SetStats(true)

}

func initLog() {
    var err error
    if Log, err = syslog.New(syslog.LOG_INFO|syslog.LOG_ERR|syslog.LOG_LOCAL0, "[apiws]"); err != nil {
        panic(err)
    }
}


func Init(){
    set_working_environment()
    load_parse_config()
    initLog()

    db_establish_connection()

    Log.Info("Server started")
}




// ~~~~~~~ Helper functions
// https://gist.github.com/bsphere/8369aca6dde3e7b4392c
// https://medium.com/coding-and-deploying-in-the-cloud/time-stamps-in-golang-abcaf581b72f


// https://gowalker.org/github.com/kylelemons/go-gypsy/yaml
