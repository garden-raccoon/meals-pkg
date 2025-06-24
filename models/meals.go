package models

import (
	"fmt"
	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
	"github.com/goccy/go-json"
	"github.com/gofrs/uuid"
)

type Meal struct {
	MealUuid     uuid.UUID `json:"meal_uuid"`
	MealSettings any       `json:"meal_settings"`
}
type UpdateMealRequest struct {
	MealUuid      uuid.UUID
	SettingsKey   string
	SettingsValue any
}

func UpdateMealToProto(u *UpdateMealRequest) (*proto.UpdateMealReq, error) {
	fields := &proto.UpdateMealReq{MealUuid: u.MealUuid.Bytes(), SettingsKey: u.SettingsKey}
	if u.SettingsValue != nil {
		b, err := json.Marshal(u.SettingsValue)
		if err != nil {
			return nil, fmt.Errorf("marshalling meal settings failed: %v", err)
		}
		fields.SettingsValue = b

	}
	return fields, nil
}

// Proto is
func (mdb Meal) Proto() (*proto.Meal, error) {

	meal := &proto.Meal{
		MealUuid: mdb.MealUuid.Bytes(),
	}
	if mdb.MealSettings != nil {
		b, err := json.Marshal(mdb.MealSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal meal settings %w", err)
		}
		meal.MealSettings = b
	}
	return meal, nil
}

func MealFromProto(pb *proto.Meal) (*Meal, error) {

	meal := &Meal{
		MealUuid: uuid.FromBytesOrNil(pb.MealUuid),
	}
	if pb.MealSettings != nil {
		var mealAny any
		err := json.Unmarshal(pb.MealSettings, &mealAny)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal meal settings %w", err)
		}
		meal.MealSettings = mealAny
	}
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
