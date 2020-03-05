package micro

import (
	"strings"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	log "github.com/micro/go-micro/v2/logger"
	micro_zerolog "github.com/micro/go-plugins/logger/zerolog/v2"
)

// Service represents a microservice definition.
type Service struct {
	micro.Service
	Name string
	ID   string
}

// NewService initializes a new microservice instance with the given identifier
// (e.g. com.example.svc.greeter), and given a configuration struct containing
// default values, loads configuration values from the given sources into the
// struct.
func NewService(id string, conf interface{}, sources ...source.Source) (*Service, error) {
	parts := strings.Split(id, ".")
	name := parts[len(parts)-1]

	service := micro.NewService(
		micro.Name(name),
		micro.WrapClient(requestIDClientWrapper),
		micro.WrapHandler(requestIDHandlerWrapper),
	)
	service.Init()

	log.DefaultLogger = micro_zerolog.NewLogger(
		micro_zerolog.UseAsDefault(),
		log.WithFields(map[string]interface{}{
			"node":    service.Server().Options().Id,
			"service": id,
		}),
		micro_zerolog.WithTimeFormat(time.RFC3339Nano),
	)

	if conf != nil {
		if err := config.Load(sources...); err != nil {
			return nil, err
		}

		if err := config.Scan(conf); err != nil {
			return nil, err
		}
	}

	return &Service{service, name, id}, nil
}

