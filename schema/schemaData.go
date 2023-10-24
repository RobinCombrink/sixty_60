package schema
import (
	"time"
)
type Invoice struct {
	Items []LineItem
	DeliveryFee uint64
	XtraSavings uint64
	Date time.Time
}
type LineItem struct {
	Name string
	Quantity uint32
	Price uint64
	Total uint64
	Discount uint64
}