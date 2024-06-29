package zk

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/knoxgao67/temporal-dispatcher-demo/config"
)

var client *zk.Conn

func Init() {
	var err error
	client, _, err = zk.Connect([]string{config.GetConfig().ZookeeperHostPort}, time.Minute, zk.WithLogInfo(true))
	if err != nil {
		log.Fatalln("unable to connect to zookeeper", err)
	}
}

func Lock(path string) func() {
	lock := zk.NewLock(client, path, zk.WorldACL(zk.PermAll))
	err := lock.Lock()
	if err != nil {
		log.Fatalln("lock failed", err)
	}
	return func() {
		_ = lock.Unlock()
	}
}

func Close() {
	client.Close()
}
