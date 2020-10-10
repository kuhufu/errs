package errs

import (
	"errors"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Printf("%v", errors.New("sdfsf"))
	fmt.Printf("%v", "sdfsdf")
	fmt.Printf("%v", 1)

}

func TestIsBuiltinErrs(t *testing.T) {
	e1 := Business("busss")
	e2 := Business(e1)

	fmt.Println(IsBuiltinErrs(e1))
	fmt.Println(IsBuiltinErrs(e2))
	fmt.Println(e1 == e2)
}
