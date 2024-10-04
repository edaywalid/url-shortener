package app

import (
	"github.com/edaywalid/url-shortner/internal/config"
	"github.com/edaywalid/url-shortner/internal/handlers"
	"github.com/edaywalid/url-shortner/internal/models"
	"github.com/edaywalid/url-shortner/internal/routes"
	"github.com/edaywalid/url-shortner/internal/services"
	"github.com/edaywalid/url-shortner/utils/redis"
	"github.com/edaywalid/url-shortner/utils/zk"
	"github.com/rs/zerolog/log"
)

type App struct {
	Config  *config.Config
	Service *Service
	Handler *Handler
	Model   *Model
	Zk      *zk.Zookeeper
	Redis   *redis.Redis
	Route   *Route
}

type (
	Service struct {
		urlService *services.Service
	}
	Handler struct {
		urlHandler *handlers.Handler
	}
	Model struct {
		rangeModel *models.Range
	}
	Route struct {
		urlRoute *routes.Routes
	}
)

func NewApp() (*App, error) {
	log.Info().Msg("Loading the configuration")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return nil, err
	}

	zookeeper, err := zk.NewZookeeper(cfg.ZkAddr)
	if err != nil {
		return nil, err
	}

	redis, err := redis.NewRedis(cfg.RedisAddr)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: cfg,
		Zk:     zookeeper,
		Redis:  redis,
	}
	return app, nil
}

func (a *App) Close() {
	err := a.Service.urlService.Close()
	if err != nil {
		log.Error().Err(err).Msg("Couldnt close service")
	}
	a.Zk.Close()
	a.Redis.Close()
}

func (a *App) initService() error {
	a.Service = &Service{
		urlService: services.NewService(a.Zk, a.Config, a.Redis, a.Model.rangeModel),
	}

	if err := a.Service.urlService.RegisterService(); err != nil {
		return err
	}

	if err := a.Service.urlService.LoadRange(); err != nil {
		if _, ok := err.(*services.RangeNotFound); !ok {
			return err
		}
		err = a.Service.urlService.InitRange()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initHandler() {
	a.Handler = &Handler{
		urlHandler: handlers.NewHandler(a.Service.urlService, a.Config),
	}
}

func (a *App) initRoute() {
	a.Route = &Route{
		urlRoute: routes.NewRoutes(a.Handler.urlHandler),
	}
	a.Route.urlRoute.RegisterRoutes()
}

func (a *App) initModel() {
	a.Model = &Model{
		rangeModel: &models.Range{},
	}

}

func (a *App) Init() error {
	a.initModel()
	if err := a.initService(); err != nil {
		return err
	}
	a.initHandler()
	a.initRoute()
	return nil
}
