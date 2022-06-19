package models

type Data struct {
	Id           int32  `json:"id" validate:"nonnil" bson:"_id"`
	FirstName    string `json:"first_name" validate:"nonnil" bson:"first_name"`
	LastName     string `json:"last_name" validate:"nonnil" bson:"last_name"`
	Email        string `json:"email" validate:"nonnil" bson:"email"`
	Gender       string `json:"gender" validate:"nonnil" bson:"gender"`
	Address      string `json:"address" validate:"nonnil" bson:"address"`
	Manufacturer string `json:"car_manufactur" validate:"nonnil" bson:"car_manufactur"`
	Model        string `json:"car_model" validate:"nonnil" bson:"car_model"`
	Year         int32  `json:"car_model_year" validate:"nonnil" bson:"car_model_year"`
}
