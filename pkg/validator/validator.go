// Package validator 主要是将 gin 的默认表单验证模块替换为 validator.v9
package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"imsdk/pkg/errno"
	"reflect"
	"sync"
)

// Validator 验证器
type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

var (
	v                 binding.StructValidator = &Validator{}
	validatorMessages map[string]map[string]string
	//_                 = app.Config().Bind("validator-messages", "", &validatorMessages)
)

// ValidateStruct 验证结构体
func (v *Validator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

// Engine 获取验证器
func (v *Validator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

// lazyInit 延迟初始化
func (v *Validator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		// 获取form tag
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := field.Tag.Get("form")
			if name != "" {
				return name
			}
			name = field.Tag.Get("json")
			if name != "" {
				return name
			}
			return field.Name
		})
		//_ = v.validate.RegisterValidation("customRegular", func(fl validator.FieldLevel) bool {
		//	return false
		//})
	})
}

// ValidErrors 验证之后的错误信息
type ValidErrors struct {
	ErrorsInfo map[string]string
	triggered  bool
}

func (validErrors *ValidErrors) add(key, value string) {
	validErrors.ErrorsInfo[key] = value
	validErrors.triggered = true
}

// IsValid 是否验证成功
func (validErrors *ValidErrors) IsValid() bool {
	return !validErrors.triggered
}

func newValidErrors() *ValidErrors {
	return &ValidErrors{
		triggered:  false,
		ErrorsInfo: make(map[string]string),
	}
}

// Bind 自定义错误信息, 如果没有匹配需要在 configs/validator-messages.toml 中添加对应处理数据
func Bind(c *gin.Context, param interface{}) error {
	err := c.ShouldBind(param)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, value := range errs {
				return errno.Add("param-incorrect:"+value.Field(), errno.ParamsFormatErr)
			}
		}
	}
	return nil
}

func VerifyStruct(obj interface{}) error {
	err := v.ValidateStruct(obj)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, value := range errs {
				return errno.Add("param-incorrect:"+value.Field(), errno.ParamsFormatErr)
			}
		}
	}
	return nil
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
