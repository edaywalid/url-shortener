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

func (s *Service) InitRange() error {

	log.Info().Msg("Acquiring lock for range initialization")
	if err := s.zkConn.Lock(); err != nil {
		log.Error().Err(err).Msg("Failed to acquire lock")
		return err
	}
	defer func() {
		log.Info().Msg("Releasing lock for range initialization")
		if err := s.zkConn.Unlock(); err != nil {
			log.Error().Err(err).Msg("Failed to release lock")
		}
	}()

	path := s.nodePath + "/range"

	log.Info().Msg("Initializing the range")
	log.Info().Msg(path)

	if _, err := s.zkConn.Get(path); err == nil {
		log.Info().Msg("Range already exists")
		return nil
	}

	rangeExists, err := s.zkConn.Exists("/url_shortener/range/last")
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if the last range exists")
		return err
	}

	if !rangeExists {
		_, err := s.zkConn.Create("/url_shortener/range", []byte("range"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create the range node")
			if err == zkp.ErrNodeExists {
				log.Info().Msg("Range node already exists")
				goto last
			}
			return err
		}

		_, err = s.zkConn.Create("/url_shortener/range/last", []byte("0"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create the last range node")
			return err
		}
	}

last:
	lastRange, err := s.zkConn.Get("/url_shortener/range/last")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get the last range")
		return err
	}

	val, err := strconv.Atoi(string(lastRange))
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert the last range")
		return err
	}
	lastRangeData := uint64(val)

	s.Range = &models.Range{
		Start:   lastRangeData + 1,
		End:     lastRangeData + 1000000,
		Current: lastRangeData + 1,
	}

	data, err := json.Marshal(s.Range)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal the range data")
		return err
	}

	_, err = s.zkConn.Create(path, data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create the range node")
		return err
	}

	data = []byte(strconv.Itoa(int(s.Range.End)))

	err = s.zkConn.Set("/url_shortener/range/last", data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) LoadRange() error {
	path := s.nodePath + "/range"
	data, err := s.zkConn.Get(path)
	var rangeData models.Range
	if err != nil {
		log.Error().Err(err).Msg("Failed to get the range")
		return &RangeNotFound{path: path}
	}
	if err := json.Unmarshal(data, &rangeData); err != nil {
		return err
	}
	s.Range = &rangeData
	log.Info().Msg("Range loaded")
	log.Info().Msg(fmt.Sprintf("%+v", s.Range))
	return nil
}
