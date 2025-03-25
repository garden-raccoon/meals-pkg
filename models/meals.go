package models

import (
	"github.com/goccy/go-json"
	"github.com/gofrs/uuid"
	"log"

	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
)

type Meal struct {
	Name         string        `json:"name"`
	Price        float64       `json:"price"`
	MealUuid     uuid.UUID     `json:"meal_uuid"`
	Description  string        `json:"description"`
	MealSettings *MealSettings `json:"meal_settings"`
	//MealAny		[]byte		   `json:"meal_any"`
	MealAny any `json:"meal_any"`
}

type MealSettings struct{}

// Proto is
func (mdb Meal) Proto() *proto.Meal {

	meal := &proto.Meal{
		MealUuid:    mdb.MealUuid.Bytes(),
		Name:        mdb.Name,
		Description: mdb.Description,
		Price:       float32(mdb.Price),
	}
	b, err := json.Marshal(mdb.MealAny)
	if err != nil {
		log.Fatal(err)
	}
	meal.MealAny = b
	return meal
}

func MealFromProto(pb *proto.Meal) *Meal {

	meal := &Meal{
		Name:        pb.Name,
		Price:       float64(pb.Price),
		Description: pb.Description,
		MealUuid:    uuid.FromBytesOrNil(pb.MealUuid),
	}
	var mealAny any
	err := json.Unmarshal(pb.MealAny, mealAny)
	if err != nil {
		log.Fatal(err)
	}
	meal.MealAny = mealAny
	return meal
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
