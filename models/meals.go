package models

import (
	"fmt"
	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
	"github.com/gofrs/uuid"
)

type Meal struct {
	MealUuid              uuid.UUID     `json:"meal_uuid"`
	LocationUuid          uuid.UUID     `json:"location_uuid"`
	Name                  string        `json:"name"`
	Description           string        `json:"description"`
	Category              string        `json:"category"`
	Price                 float64       `json:"price"`
	Tags                  []string      `json:"tags"`
	MainIngredients       []*Ingredient `json:"main_ingredients"`
	AdditionalIngredients []*Ingredient `json:"additional_ingredients"`
}

type Ingredient struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewIngredient(name string, price float64) *Ingredient {
	return &Ingredient{
		Name:  name,
		Price: price,
	}
}

func (mdb Meal) IngredientsToProto(ingrs []*Ingredient) []*proto.Ingredient {
	var ingredients []*proto.Ingredient
	for i := range ingrs {
		ingredients = append(ingredients, &proto.Ingredient{
			Name:  ingrs[i].Name,
			Price: float32(ingrs[i].Price),
		})
	}
	return ingredients
}
func IngredientsFromProto(ingrs []*proto.Ingredient) []*Ingredient {
	var ingredients []*Ingredient
	for i := range ingrs {
		ingredients = append(ingredients, &Ingredient{
			Name:  ingrs[i].Name,
			Price: float64(ingrs[i].Price),
		})
	}
	return ingredients
}

// Proto is
func (mdb Meal) Proto() (*proto.Meal, error) {

	meal := &proto.Meal{
		MealUuid:     mdb.MealUuid.Bytes(),
		LocationUuid: mdb.LocationUuid.Bytes(),
		Name:         mdb.Name,
		Description:  mdb.Description,
		Category:     mdb.Category,
		Price:        float32(mdb.Price),
		Tags:         mdb.Tags,
	}
	if mdb.MainIngredients != nil {
		meal.MainIngredients = mdb.IngredientsToProto(mdb.MainIngredients)
	}
	if mdb.AdditionalIngredients != nil {
		meal.AdditionalIngredients = mdb.IngredientsToProto(mdb.AdditionalIngredients)
	}

	return meal, nil
}

func MealFromProto(pb *proto.Meal) *Meal {

	meal := &Meal{
		MealUuid:     uuid.FromBytesOrNil(pb.MealUuid),
		LocationUuid: uuid.FromBytesOrNil(pb.LocationUuid),
		Name:         pb.Name,
		Description:  pb.Description,
		Category:     pb.Category,
		Price:        float64(pb.Price),
		Tags:         pb.Tags,
	}
	if pb.MainIngredients != nil {
		meal.MainIngredients = IngredientsFromProto(pb.MainIngredients)
	}
	if pb.AdditionalIngredients != nil {
		meal.AdditionalIngredients = IngredientsFromProto(pb.AdditionalIngredients)
	}
	return meal
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

func MealsFromProto(pb *proto.Meals) []*Meal {
	var meals []*Meal
	for _, b := range pb.Meals {
		meal := MealFromProto(b)
		meals = append(meals, meal)
	}
	return meals
}
