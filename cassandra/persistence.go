package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"

	"time"
)

var (
	Config *Configuration
)

type Configuration struct {
	URL       string        `json:"url" mapstructure:"url"`
	Keyspace  string        `json:"keyspace" mapstructure:"keyspace"`
	Username  string        `json:"username" mapstructure:"username"`
	Password  string        `json:"password" mapstructure:"password"`
	NumConns  int           `json:"numConns" mapstructure:"numConns"`
	KeepAlive time.Duration `json:"keepAliveDuration" mapstructure:"keepAliveDuration"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("cassandra.Configuration URL=%v Keyspace=%v Username=%v Password=%v NumConns=%d KeepAlive=%s",
		c.URL, c.Keyspace, c.Username, c.Password, c.NumConns, c.KeepAlive,
	)
}

func Setup(cfg *Configuration) error {
	cluster := gocql.NewCluster(cfg.URL)
	cluster.NumConns = cfg.NumConns
	cluster.SocketKeepalive = cfg.KeepAlive
	cluster.ProtoVersion = 4
	cluster.Keyspace = cfg.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	_, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("cassandra.CreateSessionErr err=%v", err.Error())
	}
	Config = cfg
	return nil
}
