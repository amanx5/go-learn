package main

// Every package import statement (normal and blank) follows the exact same initialization sequence defined by the Go runtime:
//
// For example: import "fmt"
// 	1. Dependency Order:
// 			All packages imported by fmt (and transitively down the tree) are initialized first.
// 	2. Variable Initialization:
// 			Package-level variables in fmt are evaluated and initialized.
// 	3. init() Execution:
// 			All init() functions inside fmt are executed automatically before main() begins.
//

import (
	// normal import
	"fmt"

	// blank import / side effect import
	// we are importing this package just for its side effects
	// "_" tells compiler to ignore this import and do not complain that this import is not used
	_ "io"
)

func main() {
	fmt.Print("no op")
}
