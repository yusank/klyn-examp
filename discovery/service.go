package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type ServiceInfo struct {
	IP string
}

type Service struct {
	Name    string
	Info    ServiceInfo
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
}

// NewService -  new client v3
func NewService(name string, info ServiceInfo, endpoints []string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	s := &Service{
		Name:   name,
		Info:   info,
		stop:   make(chan error),
		client: cli,
	}

	return s, nil
}

// Start  -
func (s *Service) Start() error {
	log.Println("start")
	ch, err := s.keepAlive()
	go s.putMessage()
	if err != nil {
		log.Panic(err)
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			}

			log.Printf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)
		}
	}
}

// Stop - stop service
func (s *Service) Stop(err error) {
	s.stop <- err
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &s.Info
	key := "services/" + s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {
	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}

func (s *Service) putMessage() {
	var count int
	t := time.Now().Unix()
	for {
		key := "services/"
		info := &s.Info
		value, _ := json.Marshal(info)

		_, err := s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(s.leaseid))
		if err != nil {
			log.Fatal(err)
		}
		count++
		nt := time.Now().Unix()
		if nt-t >= 2 {
			fmt.Printf("[%s] put count:%d \n", s.Name, count)
		}
	}
}
