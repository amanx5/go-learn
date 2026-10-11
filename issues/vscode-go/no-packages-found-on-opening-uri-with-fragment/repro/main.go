package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	filePath := currentSourceFile()
	uri := fileURI(filePath, 3)
	fmt.Printf(`
Open this link below:
-----------
%v
-----------

`, uri)
}

// currentSourceFile returns the absolute path of this source file.
func currentSourceFile() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

// fileURI converts a local file path and line number into a file:// URI with an #L<line> fragment.
func fileURI(filePath string, line int) string {
	normalized := strings.TrimPrefix(filepath.ToSlash(filePath), "/")
	return fmt.Sprintf("file:///%s#L%d", normalized, line)
}
