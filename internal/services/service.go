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
