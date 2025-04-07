package p

import (
	"fmt"
)

func test() {
	fmt := fmt.Sprintf("%s", "shadow") // want "Variable 'fmt' collides with imported package name"
	_ = fmt
}
