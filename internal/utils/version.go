package utils

import (
	"fmt"
	"io"
)

// Version инкапсулирует в себе информацию о версии, дате и рефе сборки
type Version struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

// Print выводит информацию о версии сборки
func (v Version) Print(writer io.Writer) {
	const template = "Build version: %s\n" +
		"Build date: %s\n" +
		"Build commit: %s\n"

	if _, err := fmt.Fprintf(
		writer,
		template,
		val(v.BuildVersion),
		val(v.BuildDate),
		val(v.BuildCommit),
	); err != nil {
		panic(err)
	}
}

func val(input string) string {
	if input == "" {
		return "N/A"
	}

	return input
}
