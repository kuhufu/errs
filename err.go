package errs

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	TypeCustom = ErrType(iota + 1)
	TypeInternal
	TypeParam
	TypeBusiness
)

type ErrType int

var strMap = map[ErrType]string{
	TypeCustom:   "Custom",
	TypeInternal: "Internal",
	TypeParam:    "Param",
	TypeBusiness: "Business",
}

func (t ErrType) String() string {
	return strMap[t]
}

type Errors interface {
	error
	Code() int
	Data() interface{}
	Type() ErrType
}

type err struct {
	typ  ErrType
	code int
	data interface{}
	err  error
}

func (e err) Code() int {
	return e.code
}

func (e err) Data() interface{} {
	return e.data
}

func (e err) Error() string {
	return e.err.Error()
}

func (e err) Type() ErrType {
	return e.typ
}

func (e err) UnWrap() error {
	if err := errors.Unwrap(e.err); err != nil {
		return err
	}
	return e.err
}

func (e err) String() string {
	return fmt.Sprintf("type: %s, code: %v, err: %v", e.typ, e.code, e.err.Error())
}

func Param(s interface{}, code ...int) error {
	c := 400
	if len(code) > 0 {
		c = code[0]
	}

	err, ok := s.(error)
	if !ok {
		err = errors.New(fmt.Sprintf("%v", s))
	}

	return custom(err, c, nil, TypeParam)
}

func Internal(s interface{}, code ...int) error {
	c := 500
	if len(code) > 0 {
		c = code[0]
	}

	_, file, line, _ := runtime.Caller(1) //1表示取上一个函数栈的信息
	err := errors.New(fmt.Sprintf("file: %v:%v [ %v ]", file, line, s))

	return custom(err, c, nil, TypeInternal)
}

func Business(s interface{}, code ...int) error {
	c := 600
	if len(code) > 0 {
		c = code[0]
	}

	err, ok := s.(error)
	if !ok {
		err = errors.New(fmt.Sprintf("%v", s))
	}

	return custom(err, c, nil, TypeBusiness)
}

func Custom(s interface{}, code int, data interface{}) error {
	err, ok := s.(error)
	if !ok {
		err = errors.New(fmt.Sprintf("%v", s))
	}

	return custom(err, code, data, TypeCustom)
}

func custom(e error, code int, data interface{}, typ ErrType) error {
	if IsBuiltinErrs(e) {
		return e
	}

	return err{code: code, err: e, data: data, typ: typ}
}

func IsBuiltinErrs(err error) bool {
	_, ok := err.(Errors)
	return ok
}
