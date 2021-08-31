package receiver

import (
	"time"

	"github.com/bfoody/Walmart-Scraper/domain"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api"
	"github.com/bfoody/Walmart-Scraper/services/client/internal/api/walmart"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

// MaxTries defines the maximum number of attempts to fetch info.
const MaxTries = 50

// A TaskService provides methods for executing tasks.
type TaskService struct {
	client *walmart.Client
	log    *zap.Logger
	rl     ratelimit.Limiter
}

// NewTaskService creates and returns a *TaskService with the provided Walmart client.
func NewTaskService(logger *zap.Logger) *TaskService {
	http := api.NewHTTPClient()
	client := walmart.NewClient(http)

	return &TaskService{
		client: client,
		log:    logger,
		rl:     ratelimit.New(10),
	}
}

// FetchProductInfo fetches the info for a single product from the API and returns it as a *ProductInfo.
func (s *TaskService) FetchProductInfo(productLocation *domain.ProductLocation) (*domain.ProductInfo, error) {
	// Wait until ratelimiter allows new operations.
	s.rl.Take()

	var id *walmart.ItemDetails
	var err error

	for i := 0; i < MaxTries; i++ {
		id, err = s.client.GetItemDetails(productLocation.Slug, productLocation.LocalID)
		if err != nil {
			s.log.Error(
				"error fetching product info, retrying in 1 second",
				zap.Int("attempt", i+1),
				zap.String("productLocationID", productLocation.ID),
				zap.Error(err),
			)
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	pi := itemDetailsToProductInfo(productLocation.ID, productLocation.ProductID, *id)

	return &pi, nil
}
