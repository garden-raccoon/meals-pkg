package models

import (
	"github.com/gofrs/uuid"

	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
)

type Meal struct {
	Name         string        `json:"name"`
	Price        float64       `json:"price"`
	MealUuid     uuid.UUID     `json:"meal_uuid"`
	Description  string        `json:"description"`
	MealSettings *MealSettings `json:"meal_settings"`
}
type MealSettings struct{}

// Proto is
func (mdb Meal) Proto() *proto.Meal {

	order := &proto.Meal{
		MealUuid:    mdb.MealUuid.Bytes(),
		Name:        mdb.Name,
		Description: mdb.Description,
		Price:       float32(mdb.Price),
	}
	return order
}

func MealFromProto(pb *proto.Meal) *Meal {

	order := &Meal{
		Name:        pb.Name,
		Price:       float64(pb.Price),
		Description: pb.Description,
		MealUuid:    uuid.FromBytesOrNil(pb.MealUuid),
	}
	return order
}

func MealsToProto(meals []*Meal) *proto.Meals {
	pb := &proto.Meals{}

	for _, b := range meals {
		pb.Meals = append(pb.Meals, b.Proto())
	}
	return pb
}

func MealsFromProto(pb *proto.Meals) (meals []*Meal) {
	for _, b := range pb.Meals {
		meals = append(meals, MealFromProto(b))
	}
	return
}
