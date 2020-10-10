package errs

import (
	"errors"
	"fmt"
	"runtime"
)

func Param(s interface{}, code ...int) error {
	c := 400
	if len(code) > 0 {
		c = code[0]
	}

	switch s := s.(type) {
	case error:
		if IsBuiltinErrs(s) {
			return s
		}
		return &ErrParam{Code: c, err: s}
	default:
		return &ErrParam{Code: c, err: errors.New(fmt.Sprintf("%v", s))}
	}
}

func Internal(s interface{}, code ...int) error {
	c := 500
	if len(code) > 0 {
		c = code[0]
	}

	_, file, line, _ := runtime.Caller(1) //1表示取上一个函数栈的信息

	switch s := s.(type) {
	case error:
		if IsBuiltinErrs(s) {
			return s
		}
		return &ErrInternal{Code: c, err: errors.New(fmt.Sprintf("file: %v:%v [ %v ]", file, line, s))}
	default:
		return &ErrInternal{Code: c, err: errors.New(fmt.Sprintf("file: %v:%v [ %v ]", file, line, s))}
	}
}

func Business(s interface{}, code ...int) error {
	c := 600
	if len(code) > 0 {
		c = code[0]
	}

	switch s := s.(type) {
	case error:
		if IsBuiltinErrs(s) {
			return s
		}
		return &ErrBusiness{Code: c, err: s}
	default:
		return &ErrBusiness{Code: c, err: errors.New(fmt.Sprintf("%v", s))}
	}

}

func Custom(s interface{}, code int, data ...interface{}) error {
	var tmpData interface{}
	if len(data) > 0 {
		tmpData = data[0]
	}

	switch s := s.(type) {
	case error:
		if IsBuiltinErrs(s) {
			return s
		}
		return &ErrCustom{Code: code, err: s, Data: tmpData}
	default:
		return &ErrCustom{Code: code, err: errors.New(fmt.Sprintf("%v", s)), Data: tmpData}
	}

}

type ErrCustom struct {
	Code int
	Data interface{}
	err  error
}

func (p *ErrCustom) Error() string {
	return p.err.Error()
}

func (p *ErrCustom) UnWrap() error {
	if err := errors.Unwrap(p.err); err != nil {
		return err
	}
	return p.err
}

type ErrParam struct {
	Code int
	err  error
}

func (p *ErrParam) Error() string {
	return p.err.Error()
}

func (p *ErrParam) UnWrap() error {
	if err := errors.Unwrap(p.err); err != nil {
		return err
	}
	return p.err
}

type ErrInternal struct {
	Code int
	err  error
}

func (p *ErrInternal) Error() string {
	return p.err.Error()
}

func (p *ErrInternal) UnWrap() error {
	if err := errors.Unwrap(p.err); err != nil {
		return err
	}
	return p.err
}

type ErrBusiness struct {
	Code int
	err  error
}

func (p *ErrBusiness) Error() string {
	return p.err.Error()
}

func (p *ErrBusiness) UnWrap() error {
	if err := errors.Unwrap(p.err); err != nil {
		return err
	}
	return p.err
}

var (
	errInternal = &ErrInternal{}
	errParam    = &ErrParam{}
	errBusiness = &ErrBusiness{}
	errCustom   = &ErrCustom{}
)

func IsErrParam(err error) bool {
	return errors.As(err, &errParam)
}

func IsErrInternal(err error) bool {
	return errors.As(err, &errInternal)
}

func IsErrBusiness(err error) bool {
	return errors.As(err, &errBusiness)
}

func IsErrCustom(err error) bool {
	return errors.As(err, &errCustom)
}

func IsBuiltinErrs(err error) bool {
	switch {
	case IsErrBusiness(err), IsErrInternal(err), IsErrParam(err), IsErrCustom(err):
		return true
	default:
		return false
	}
}
