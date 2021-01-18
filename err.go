package errs

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	TypeCustom   = ErrType(1)
	TypeInternal = ErrType(2)
	TypeParam    = ErrType(3)
	TypeBusiness = ErrType(4)
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

type Error struct {
	typ  ErrType
	code int
	data interface{}
	err  error
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Data() interface{} {
	return e.data
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Type() ErrType {
	return e.typ
}

func (e Error) UnWrap() error {
	if err := errors.Unwrap(e.err); err != nil {
		return err
	}
	return e.err
}

func (e Error) String() string {
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

func custom(err error, code int, data interface{}, typ ErrType) error {
	if IsBuiltinErrs(err) {
		return err
	}

	return Error{code: code, err: err, data: data, typ: typ}
}

func IsBuiltinErrs(err error) bool {
	_, ok := err.(Errors)
	return ok
}
