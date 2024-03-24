package filereader

import (
	"bufio"
	"os"
	"strings"
)

func Reader(filename string) string {

	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	var slice []string
	for scanner.Scan() {
		str := scanner.Text()
		slice = append(slice, str)
	}
	return strings.Join(slice, " ")
}
