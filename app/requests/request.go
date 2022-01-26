package requests

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/StubbornYouth/goblog/pkg/model"
	"github.com/thedevsaddam/govalidator"
)

// 自定义认证规则
// 假如 main 引入了 pkg1 最终依赖于 pkg3，pkg3 中的 init() 方法会优先被执行；
// 同一个包里，单文件的情况，init() 优先于其他方法执行，包括 main()；同一个包里的常量和变量声明会优先于 init() 方法执行；
// 同一个文件里允许多个 init() 存在，会按照自上而下的顺序执行；
// 同一个包，多个文件里存在 init() 的情况，执行顺序是按文件名的字母排序执行。
func init() {
	govalidator.AddCustomRule("not_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbField := rng[1]
		val := value.(string)

		var count int64
		model.DB.Table(tableName).Where(dbField+" = ?", val).Count(&count)
		if count != 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", val)
		}

		return nil
	})

	govalidator.AddCustomRule("is_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "is_exists:"), ",")

		tableName := rng[0]
		dbField := rng[1]
		val := value.(string)

		var count int64
		model.DB.Table(tableName).Where(dbField+" = ?", val).Count(&count)
		if count == 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 不存在", val)
		}

		return nil
	})

	// 中文版 字符串长度限制
	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:")) //handle other error
		if valLength > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})

	// min_cn:2
	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:")) //handle other error
		if valLength < l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}
		return nil
	})
}
