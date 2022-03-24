package edbx

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

// dbConfig Cluser Base Config
type dbConfig struct {
	name string
	dsn  string
	user string
	pwd  string
}

// Open 开启连接
func (c *dbConfig) Open() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{c.dsn},
		Username:  c.user,
		Password:  c.pwd,
		Transport: NewTransport(),
	}

	// Create a client
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	log.Printf("[app.edbx] es success, name: %s", c.name)
	return client
}

func InitConn(m []map[string]interface{}) {
	for _, item := range m {
		var dbc = dbConfig{}

		if name, ok := item["name"]; ok {
			if nameStr, okInterface := name.(string); okInterface {
				if nameStr == "" {
					dbc.name = defaultEdbName
				} else {
					dbc.name = nameStr
				}
			}
		} else {
			dbc.name = defaultEdbName
		}

		if dsn, ok := item["dsn"]; ok {
			if dsnStr, okInterface := dsn.(string); okInterface {
				dbc.dsn = dsnStr
			}
		}

		if user, ok := item["user"]; ok {
			if userStr, okInterface := user.(string); okInterface {
				dbc.user = userStr
			}
		}

		if pwd, ok := item["pwd"]; ok {
			if pwdStr, okInterface := pwd.(string); okInterface {
				dbc.pwd = pwdStr
			}
		}

		setEdb(dbc.name, dbc.Open())
	}
}

func GetClient(name ...string) *elasticsearch.Client {
	if len(name) > 0 {
		return getEdb(name[0])
	} else {
		return getEdb(defaultEdbName)
	}
}
