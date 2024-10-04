package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/edaywalid/url-shortner/internal/config"
	"github.com/edaywalid/url-shortner/internal/models"
	"github.com/edaywalid/url-shortner/utils"
	"github.com/edaywalid/url-shortner/utils/redis"
	"github.com/edaywalid/url-shortner/utils/zk"
	"github.com/rs/zerolog/log"

)

type Service struct {
	zkConn     *zk.Zookeeper
	redis      *redis.Redis
	config     *config.Config
	nodePath   string
	Range      *models.Range
	rangeMutex *sync.Mutex
}

func NewService(
	zkConn *zk.Zookeeper,
	config *config.Config,
	redis *redis.Redis,
	rangeData *models.Range,
) *Service {
	return &Service{
		zkConn:     zkConn,
		config:     config,
		redis:      redis,
		rangeMutex: &sync.Mutex{},
		Range:      rangeData,
	}
}

type RangeNotFound struct {
	path string
}

func (r *RangeNotFound) Error() string {
	return fmt.Sprintf("range not found: %s", r.path)
}


func (s *Service) RegisterService() error {
	// check if "/url_shortener" Exists
	exists, err := s.zkConn.Exists("/url_shortener")
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if the root node exists")
		return err
	}
	if !exists {
		_, err := s.zkConn.Create("/url_shortener", []byte("root"))
		if err != nil {

			log.Error().Err(err).Msg("Failed to create the root node")
			if err == zkp.ErrNodeExists {
				log.Info().Msg("Root node already exists")
				goto servers
			}
			return err
		}
	}

servers:
	// check if "/url_shortener/servers" exists
	exists, err = s.zkConn.Exists("/url_shortener/servers")
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if the servers node exists")
		return err
	}
	if !exists {
		_, err := s.zkConn.Create("/url_shortener/servers", []byte("servers"))
		if err != nil {
			if err == zkp.ErrNodeExists {
				log.Info().Msg("Servers node already exists")
				goto path
			}
			log.Error().Err(err).Msg("Failed to create the servers node")
			return err
		}
	}
path:
	path := fmt.Sprintf("/url_shortener/servers/%s", s.config.ServerID)
	exists, err = s.zkConn.Exists(path)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if the server exists")
		return err
	}
	if !exists {
		s.nodePath, err = s.zkConn.Create(path, []byte(s.config.ServerID))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create the server node")
			return err
		}
	} else {
		s.nodePath = path
		log.Info().Msg("Server already exists")
		return nil
	}
	return nil
}
