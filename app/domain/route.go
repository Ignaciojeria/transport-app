package domain

type Route struct {
	Destination NodeInfo
	Vehicle     Vehicle
	Operator    Operator
	Orders      []Order
}
