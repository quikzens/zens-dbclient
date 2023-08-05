package entity

type OrderQueryParam struct {
	SortBy  string
	OrderBy string
}

type Condition struct {
	Field       string `json:"field"`
	Operator    string `json:"operator"`
	FirstValue  string `json:"first_value"`
	SecondValue string `json:"second_value"`
}
