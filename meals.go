package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/garden-raccoon/meals-pkg/models"
	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
)

// TODO need to set timeout via lib initialisation
// timeOut is  hardcoded GRPC requests timeout value
const timeOut = 60

// Debug on/off
var Debug = true

type MealsPkgAPI interface {
	CreateOrUpdateMeals(s []*models.Meal) error
	GetMeals() ([]*models.Meal, error)
	DeleteMeal(mealUuid uuid.UUID) error
	MealByName(name string) (*models.Meal, error)
	HealthCheck() error
	// Close GRPC Api connection
	Close() error
}

// Api is profile-service GRPC Api
// structure with client Connection
type Api struct {
	addr    string
	timeout time.Duration
	*grpc.ClientConn
	mu sync.Mutex
	proto.MealsServiceClient
	grpc_health_v1.HealthClient
}

// New create new Battles Api instance
func New(addr string) (MealsPkgAPI, error) {
	api := &Api{timeout: timeOut * time.Second}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create MealsApi:  %w", err)
	}

	api.MealsServiceClient = proto.NewMealsServiceClient(api.ClientConn)
	return api, nil
}

func (api *Api) DeleteMeal(mealUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	req := &proto.MealDeleteReq{
		MealUuid: mealUuid.Bytes(),
	}
	_, err := api.MealsServiceClient.DeleteMeal(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteMeal api request: %w", err)
	}
	return nil
}

func (api *Api) GetMeals() ([]*models.Meal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	var resp *proto.Meals
	empty := new(proto.MealsEmpty)
	resp, err := api.MealsServiceClient.GetMeals(ctx, empty)
	if err != nil {
		return nil, fmt.Errorf("GetMeals api request: %w", err)
	}

	meals := models.MealsFromProto(resp)

	return meals, nil
}

func (api *Api) CreateOrUpdateMeals(s []*models.Meal) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	meals := models.MealsToProto(s)

	_, err = api.MealsServiceClient.CreateOrUpdateMeals(ctx, meals)
	if err != nil {
		return fmt.Errorf("create meals api request: %w", err)
	}
	return nil
}

// initConn initialize connection to Grpc servers
func (api *Api) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	api.ClientConn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	return
}

func (api *Api) MealByName(name string) (*models.Meal, error) {
	getter := &proto.MealGetter{
		Getter: &proto.MealGetter_Name{
			Name: name,
		},
	}
	return api.getMeal(getter)
}

// getMeal is
func (api *Api) getMeal(getter *proto.MealGetter) (*models.Meal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	resp, err := api.MealsServiceClient.MealByName(ctx, getter)
	if err != nil {
		return nil, fmt.Errorf("MealAPI getMeal request failed: %w", err)
	}

	return models.MealFromProto(resp), nil
}

func (api *Api) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	api.mu.Lock()
	defer api.mu.Unlock()

	resp, err := api.HealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "mealsapi"})
	if err != nil {
		return fmt.Errorf("healthcheck error: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("node is %s", errors.New("service is unhealthy"))
	}
	//api.status = service.StatusHealthy
	return nil
}
