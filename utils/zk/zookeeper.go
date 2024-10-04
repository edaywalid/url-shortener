package zk

import (
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/rs/zerolog/log"
)

type Zookeeper struct {
	conn  *zk.Conn
	lock  *zk.Lock
	event <-chan zk.Event
}

func NewZookeeper(addr string) (*Zookeeper, error) {
	z := &Zookeeper{}
	err := z.connectToZookeeper(addr)
	if err != nil {
		return nil, err
	}
	return z, nil
}

func (z *Zookeeper) newLock(path string) *zk.Lock {
	return zk.NewLock(z.conn, path, zk.WorldACL(zk.PermAll))
}
func (z *Zookeeper) Lock() error {
	return z.lock.Lock()
}

func (z *Zookeeper) Unlock() error {
	return z.lock.Unlock()
}

func (z *Zookeeper) connectToZookeeper(addr string) error {
	servers := []string{addr}
	conn, event, err := zk.Connect(servers, time.Second*10)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Zookeeper")
		return err
	}
	fmt.Println(conn.State())
	z.conn = conn
	z.event = event
	z.lock = z.newLock("/url_shortener/range_lock")
	return nil
}

func (z *Zookeeper) Close() {
	z.conn.Close()
}

func (z *Zookeeper) Create(path string, data []byte) (string, error) {
	path, err := z.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	return path, err
}

func (z *Zookeeper) Exists(path string) (bool, error) {
	exists, _, err := z.conn.Exists(path)
	return exists, err
}

func (z *Zookeeper) Get(path string) ([]byte, error) {
	data, _, err := z.conn.Get(path)
	return data, err
}

func (z *Zookeeper) Set(path string, data []byte) error {
	_, err := z.conn.Set(path, data, -1)
	return err
}

func (z *Zookeeper) Watch() <-chan zk.Event {
	return z.event
}
