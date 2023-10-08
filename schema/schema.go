package schema

type Invoice struct {
	Items []LineItem
	DeliveryFee uint16
	XtraSavings uint16
}
type LineItem struct {
	Name string
	Quantity uint32
	Price uint64
	Total uint64
	Discount uint64
}