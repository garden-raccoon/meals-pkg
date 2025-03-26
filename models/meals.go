package models

import (
	"fmt"
	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
	"github.com/goccy/go-json"
	"github.com/gofrs/uuid"
)

type Meal struct {
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	MealUuid     uuid.UUID `json:"meal_uuid"`
	Description  string    `json:"description"`
	MealSettings any       `json:"meal_settings"`
}

type MealSettings struct{}

// Proto is
func (mdb Meal) Proto() (*proto.Meal, error) {

	meal := &proto.Meal{
		MealUuid:    mdb.MealUuid.Bytes(),
		Name:        mdb.Name,
		Description: mdb.Description,
		Price:       float32(mdb.Price),
	}
	b, err := json.Marshal(mdb.MealSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal meal settings %w", err)
	}
	meal.MealSettings = b
	return meal, nil
}

func MealFromProto(pb *proto.Meal) (*Meal, error) {

	meal := &Meal{
		Name:        pb.Name,
		Price:       float64(pb.Price),
		Description: pb.Description,
		MealUuid:    uuid.FromBytesOrNil(pb.MealUuid),
	}
	var mealAny any
	err := json.Unmarshal(pb.MealSettings, mealAny)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal meal settings %w", err)
	}
	meal.MealSettings = mealAny
	return meal, nil
}

func MealsToProto(meals []*Meal) (*proto.Meals, error) {
	pb := &proto.Meals{}

	for _, b := range meals {
		meal, err := b.Proto()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		pb.Meals = append(pb.Meals, meal)
	}
	return pb, nil
}

func MealsFromProto(pb *proto.Meals) ([]*Meal, error) {
	var meals []*Meal
	for _, b := range pb.Meals {
		meal, err := MealFromProto(b)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		meals = append(meals, meal)
	}
	return meals, nil
}
