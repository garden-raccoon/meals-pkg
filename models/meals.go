package models

import (
	proto "github.com/garden-raccoon/meals-pkg/protocols/meals-pkg"
	"github.com/gofrs/uuid"
)

type Meal struct {
	MealUuid              uuid.UUID     `json:"meal_uuid"`
	LocationUuid          uuid.UUID     `json:"location_uuid"`
	Name                  string        `json:"name"`
	Description           string        `json:"description"`
	Category              *Category     `json:"category"`
	Price                 float64       `json:"price"`
	Tags                  []string      `json:"tags"`
	Weight                float64       `json:"weight"`
	MainIngredients       []*Ingredient `json:"main_ingredients"`
	AdditionalIngredients []*Ingredient `json:"additional_ingredients"`
}
type Category struct {
	Name          string      `json:"name"`
	SubCategories []*Category `json:"sub_categories"`
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
func (mdb Meal) Proto() *proto.Meal {

	meal := &proto.Meal{
		MealUuid:     mdb.MealUuid.Bytes(),
		LocationUuid: mdb.LocationUuid.Bytes(),
		Name:         mdb.Name,
		Description:  mdb.Description,
		Tags:         mdb.Tags,
		Weight:       float32(mdb.Weight),
	}
	if mdb.Category != nil {
		meal.Category = &proto.Category{
			Name: mdb.Category.Name,
		}
		if mdb.Category.SubCategories != nil {
			var subCategories []*proto.Category
			for _, sub := range mdb.Category.SubCategories {
				subCategories = append(subCategories, &proto.Category{
					Name: sub.Name,
				})
			}
			meal.Category.SubCategories = subCategories
		}
	}
	if mdb.MainIngredients != nil {
		meal.MainIngredients = mdb.IngredientsToProto(mdb.MainIngredients)
	}
	if mdb.AdditionalIngredients != nil {
		meal.AdditionalIngredients = mdb.IngredientsToProto(mdb.AdditionalIngredients)
	}

	return meal
}

func MealFromProto(pb *proto.Meal) *Meal {

	meal := &Meal{
		MealUuid:     uuid.FromBytesOrNil(pb.MealUuid),
		LocationUuid: uuid.FromBytesOrNil(pb.LocationUuid),
		Name:         pb.Name,
		Description:  pb.Description,
		Tags:         pb.Tags,
		Weight:       float64(pb.Weight),
	}
	if pb.Category != nil {
		meal.Category = &Category{
			Name: pb.Category.Name,
		}
		if pb.Category.SubCategories != nil {
			var subCategories []*Category
			for _, sub := range pb.Category.SubCategories {
				subCategories = append(subCategories, &Category{
					Name: sub.Name,
				})
			}
			meal.Category.SubCategories = subCategories
		}
	}
	if pb.MainIngredients != nil {
		meal.MainIngredients = IngredientsFromProto(pb.MainIngredients)
	}
	if pb.AdditionalIngredients != nil {
		meal.AdditionalIngredients = IngredientsFromProto(pb.AdditionalIngredients)
	}
	return meal
}

func MealsToProto(meals []*Meal) *proto.Meals {
	pb := &proto.Meals{}

	for _, b := range meals {
		meal := b.Proto()
		pb.Meals = append(pb.Meals, meal)
	}
	return pb
}

func MealsFromProto(pb *proto.Meals) []*Meal {
	var meals []*Meal
	for _, b := range pb.Meals {
		meal := MealFromProto(b)
		meals = append(meals, meal)
	}
	return meals
}
