package order

// Order represents a field and the sorting order (ascending or descending).
type Order struct {
	Field string
	Order string
}

// Asc creates an ascending order on the given field.
func Asc(field string) Order {
	return Order{
		Field: field,
		Order: "ASC",
	}
}

// Desc creates a descending order on the given field.
func Desc(field string) Order {
	return Order{
		Field: field,
		Order: "DESC",
	}
}
