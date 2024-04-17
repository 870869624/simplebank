package api

import (
	"simplebank/util"

	"github.com/go-playground/validator/v10"
)

// 检验是否支持该货币，避免模型中的反射写的太多以及重复
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if str, ok := fl.Field().Interface().(string); ok {
		//确认能从反射转换为字符串，需要检查是否支持该货币
		return util.IsSupportCurrency(str)
	}
	return false
}
