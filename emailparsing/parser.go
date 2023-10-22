package emailparsing

import (
	"parser60/format"
	"parser60/schema"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)
const deliveryDateTimeText string = "Delivery Date & Time"

func GetInvoiceFromHtml(htmlData string) *schema.Invoice {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlData))

	invoice := schema.Invoice{DeliveryFee: 3500}

	doc.Find("td[colspan='4']").Each(func(i int, s *goquery.Selection) {
		table := s.Find("table[width='100%'][border='0'][cellpadding='0'][cellspacing='0']")
		item := schema.LineItem{}
		table.Find("tr").Each(func(j int, tr *goquery.Selection) {
			if j == 0 { // This is the first row (discount row)
				tr.Find("td").Each(func(k int, td *goquery.Selection) {
					if k == 3 { // This is the fourth cell (discount cell)
						text := strings.TrimSpace(td.Text())
						discount, _ := strconv.ParseFloat(strings.TrimPrefix(text, "- R "), 64)
						item.Discount = toCents(discount) // Convert to string with 2 decimal places
					}
				})
			} else if j == 1 { // This is the second row
				tr.Find("td").Each(func(k int, td *goquery.Selection) {
					text := strings.TrimSpace(td.Text())
					if text != "" {
						switch k {
						case 0:
							item.Name = text
						case 1:
							qty, _ := strconv.Atoi(text)
							item.Quantity = uint32(qty)
						case 2:
							price, _ := strconv.ParseFloat(strings.TrimPrefix(text, "R "), 64)
							item.Price = toCents(price) // Convert to cents
						case 3:
							total, _ := strconv.ParseFloat(strings.TrimPrefix(text, "R "), 64)
							item.Total = toCents(total) // Convert to cents
						}
					}
				})
			}
		})
		if item.Name != "" || item.Discount != 0 {
			invoice.Items = append(invoice.Items, item)
		}
	})
	doc.Find("table[width='100%'][border='0'][cellpadding='0'][cellspacing='0']").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, tr *goquery.Selection) {
			tr.Find("td").Each(func(k int, td *goquery.Selection) {
				td.Find("span").Each(func(l int, span *goquery.Selection) {
					text := span.Text()
					if text == deliveryDateTimeText {
						deliveryDateTime := td.Next().Find("span").Text()
						if deliveryDateTime != "" {
							invoice.Date = format.ToDate(deliveryDateTime)
						}
					}
				})

			})
		})
	})

	filteredItems := make([]schema.LineItem, 0)
	for i := 0; i < len(invoice.Items); i++ {
		item := invoice.Items[i]
		if item.Name != "" {
			if i < len(invoice.Items)-1 && invoice.Items[i+1].Name == "" {
				item.Discount = invoice.Items[i+1].Discount
				invoice.XtraSavings += item.Discount
				i++ // Skip next item
			}
			filteredItems = append(filteredItems, item)
		}
	}
	invoice.Items = filteredItems
	return &invoice
}

func toCents(input float64) uint64 {
	epsilon := 0.0001
	return uint64((input + epsilon) * 100)
}
