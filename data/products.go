package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// CREATING PRODUCT BLUE-PRINT
type Product struct {
	ID            int     `json:"id"` // will be generated automatically
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
	supposedRegexp := regexp.MustCompile(`[a-z]-+[a-z]+-[a-z]+`)

	passedString := supposedRegexp.FindAllString(fl.Field().String(), -1)

	if len(passedString) != 1 {
		return false
	}
	return true
}

// type which will hold slice of address of type Product
type Products []*Product

// FUNCTIONS FOR ABSTRACTION AND IMPROVING CODE QUALITY
// encoder
func (p *Products) ToJSON(w io.Writer) error {
	// package encoding convert data to and from byte-level and textual representations
	// encoding/json will convert to json
	encoder := json.NewEncoder(w)

	return encoder.Encode(p)
}

// decoder, just a single product
func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

// --------------------------------------- GET logic -------------------------------------------
func GetProducts() Products {
	return productList
}

// --------------------------------------- POST logic -------------------------------------------
func AddProduct(p *Product) {
	// overwriting ID to the incoming product
	p.ID = getNextID()

	// will append incoming product to the productList with ID attached to it
	productList = append(productList, p)
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

// --------------------------------------- PUT logic ---------------------------------------------
func UpdateProduct(productID int, prod *Product) error {
	// verifing for the product in our DB
	index, err := findProduct(productID)

	// product not found
	if err != nil {
		return err
	}

	// assigning passed productID to the passed product(empty), only after varifying that product exists
	prod.ID = productID
	productList[index] = prod
	return nil
}

var ErrorProductNotFound = fmt.Errorf("Product not Found in our DB")

func findProduct(productID int) (int, error) {
	for index, product := range productList {
		if product.ID == productID {
			return index, nil
		}
	}
	return -1, ErrorProductNotFound
}
