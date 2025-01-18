package business

// Business represents a business in the database
type Business struct {
	ID             string   `bson:"id"`
	Name           string   `bson:"name"`
	Phone          string   `bson:"phone"`
	Email          string   `bson:"email"`
	PictureURLs    []string `bson:"picture_urls"`
	Description    string   `bson:"description"`
	GSTIndexNumber string   `bson:"gst_index_number"`
	Location       Location `bson:"location"`
	Industry       Industry `bson:"industry"`
}

// Industry represents the industry of a business
type Industry string

const (
	IndustryRetailAndEcommerce           Industry = "RetailAndEcommerce"
	IndustryFoodAndBeverages             Industry = "FoodAndBeverages"
	IndustryHealthAndWellness            Industry = "HealthAndWellness"
	IndustryEducationAndTraining         Industry = "EducationAndTraining"
	IndustryProfessionalServices         Industry = "ProfessionalServices"
	IndustryManufacturing                Industry = "Manufacturing"
	IndustryHospitalityAndTourism        Industry = "HospitalityAndTourism"
	IndustryArtsAndEntertainment         Industry = "ArtsAndEntertainment"
	IndustryTechnologyAndSoftware        Industry = "TechnologyAndSoftware"
	IndustryConstructionAndRealEstate    Industry = "ConstructionAndRealEstate"
	IndustryTransportationAndLogistics   Industry = "TransportationAndLogistics"
	IndustryAgricultureAndFarming        Industry = "AgricultureAndFarming"
	IndustryFinanceAndInsurance          Industry = "FinanceAndInsurance"
	IndustryEnergyAndUtilities           Industry = "EnergyAndUtilities"
	IndustryNonProfitAndSocialEnterprise Industry = "NonProfitAndSocialEnterprise"
	IndustryMediaAndPublishing           Industry = "MediaAndPublishing"
	IndustryAutomotive                   Industry = "Automotive"
	IndustryFashionAndLifestyle          Industry = "FashionAndLifestyle"
	IndustrySportsAndRecreation          Industry = "SportsAndRecreation"
	IndustryOther                        Industry = "Other"
)

// Location represents the location of a job
type Location struct {
	Coordinates Coordinates `bson:"coordinates" json:"coordinates"`
	Address     string      `bson:"address" json:"address"`
	Street      string      `bson:"street" json:"street"`
	PostalCode  string      `bson:"postal_code" json:"postalCode"`
	City        string      `bson:"city" json:"city"`
	State       string      `bson:"state" json:"state"`
}

// Coordinates represents the latitude and longitude of a location
type Coordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitute float64 `bson:"longitude" json:"longitude"`
}
