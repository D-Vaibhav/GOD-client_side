package data

import (
	"encoding/json"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// CREATING PRODUCT BLUE-PRINT
type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"`
	Price         float32 `json:"price" validate:"gt=0"`
	GST           float32 `json:"gst" validate:"gt=0"`
	LicenceNumber string  `json:"licenceNumber" validate:"required,licenceNumberTag"`
	CreatedOn     string  `json:"-"`
	UpdatedOn     string  `json:"-"`
	DeletedOn     string  `json:"-"`
}

// PRODUCT DATABASE
// sending address so can be updated in future
var productList = []*Product{
	&Product{
		ID:            1,
		Name:          "Coffee",
		Description:   "best cooffee description",
		Price:         18.75,
		GST:           1.25,
		LicenceNumber: "upe123",
		CreatedOn:     time.Now().UTC().String(),
		UpdatedOn:     time.Now().UTC().String(),
	},
	&Product{
		ID:            2,
		Name:          "Tea",
		Description:   "best tea description",
		Price:         8.75,
		GST:           1.25,
		LicenceNumber: "rjb123",
		CreatedOn:     time.Now().UTC().String(),
		UpdatedOn:     time.Now().UTC().String(),
	},
}

// VALIDATING PRODUCT
func (p *Product) Validate() error {
	// creating new validator object
	validate := validator.New()

	// tieing tag to func
	validate.RegisterValidation("licenceNumberTag", validateLicenceNumber) // second para func should return bool

	// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise
	return validate.Struct(p)
}

func validateLicenceNumber(fl validator.FieldLevel) bool {
	supposedRegexp := regexp.MustComplie(`[a-z]+-[a-z]+-[a-z]+`)

	passedString := supposedRegexp.FindAllStrings(fl.Field().String(), -1)

	if len(passedString) != 1 {
		return false
	}
	return true
}

// type which will hold slice of address of type Product
type Products []*Product

// FUNCTIONS FOR ABSTRACTION AND IMPROVING CODE QUALITY
func (p *Products) ToJSON(w io.Writer) error {
	// package encoding convert data to and from byte-level and textual representations
	// encoding/json will convert to json
	encoder := json.NewEncoder(w)

	return encoder.Encode(p)
}

func GetProducts() Products {
	return productList
}
