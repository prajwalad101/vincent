package util

import (
	"fmt"
	"strings"
	"time"
)

func GenerateFileName(originalFilename string) string {
	// split filename and extension
	fileNameSlice := strings.Split(originalFilename, ".")
	fileName := fileNameSlice[0]
	extension := fileNameSlice[1]

	currentTime := time.Now()
	parsedTime := currentTime.Format(time.RFC3339)

	fileName = fmt.Sprintf(
		"%s(%s)",
		fileName,
		parsedTime,
	)
	return fmt.Sprintf("%s.%s", fileName, extension)
}
