package walmart

// An ItemDetails is returned when calling GetItemDetails, contains details about a single
// item such as name, price, and availability status.
type ItemDetails struct {
	ID                 string  // the remote ID of the item, given by the API
	Slug               string  // the item's slug
	Name               string  // name of the item
	Category           string  // the category path for the item, separated by '/'
	Price              float32 // the current price of the item in USD
	AvailabilityStatus string  // the availability of the item, eg. "IN_STOCK"
	InStock            bool    // whether or not the item is in stock (AvailabilityStatus but as a boolean)
}
