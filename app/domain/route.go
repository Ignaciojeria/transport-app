package domain

type Route struct {
	Organization
	Destination NodeInfo
	Vehicle     Vehicle
	Operator    Operator
	Orders      []Order
}
