package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

var (
	Server    *server
	Mysql     *mySQL
	Snowflake *snowflake
	Service   *service
	Etcd      *etcd
	Kafka     *kafka
	Redis     *redis
	OSS       *oss

	runtime_viper = viper.New()
)

func Init(path string, service string) {
	runtime_viper.SetConfigType("yaml")
	etcdAddr := os.Getenv("ETCD_ADDR")

	if etcdAddr == "" {
		panic(errors.New("ETCD_ADDR environment variable is required"))
	}

	Etcd = &etcd{Addr: etcdAddr}

	// Load shared configuration from etcd once during startup.
	err := runtime_viper.AddRemoteProvider("etcd3", Etcd.Addr, "/config/config.yaml")

	if err != nil {
		panic(err)
	}

	if err := runtime_viper.ReadRemoteConfig(); err != nil {
		panic(err)
	}

	configMapping(service)
}

func configMapping(srv string) {
	c := new(config)
	if err := runtime_viper.Unmarshal(&c); err != nil {
		panic(err)
	}
	Snowflake = &c.Snowflake

	Server = &c.Server
	Server.Secret = []byte(runtime_viper.GetString("server.jwt-secret"))

	Mysql = &c.MySQL
	Kafka = &c.Kafka
	Redis = &c.Redis
	OSS = &c.OSS
	Service = GetService(srv)
}

func GetService(srvname string) *service {
	addrlist := runtime_viper.GetStringSlice("services." + srvname + ".addr")

	return &service{
		Name:     runtime_viper.GetString("services." + srvname + ".name"),
		AddrList: addrlist,
		LB:       runtime_viper.GetBool("services." + srvname + ".load-balance"),
	}
}
