package util

// 添加更多支持货币
const (
	USD = "USD"
	RMB = "RMB"
	EUR = "EUR"
)

// 支持该货币返回为真，否则为假
func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, RMB, EUR:
		return true
	}
	return false
}
