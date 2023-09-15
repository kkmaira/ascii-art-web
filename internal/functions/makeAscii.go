package functions

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	standardHash   = "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf"
	shadowHash     = "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73"
	thinkertoyHash = "64285e4960d199f4819323c4dc6319ba34f1f0dd9da14d07111345f5d76c3fa3"
)

func MakeAscii(textToChange, banerToChange string) (string, error) {
	ASCIIbase, err := GetAsciiBase(banerToChange)
	if err != nil {
		return "", errors.New("500")

	}

	InputBase, err := InputBase(textToChange)
	if err != nil {
		return "", errors.New("400")
	}

	output, err := PrintAscii(InputBase, ASCIIbase, textToChange)
	if err != nil {
		return "", errors.New("400")
	}

	return output, nil
}

func GetAsciiBase(fileName string) ([]string, error) {
	var ASCIIbase []string
	var Oneline string
	hashError := errors.New("ERROR: Template file has been changed")
	// чтение файла

	FileContent, err := os.ReadFile(fmt.Sprintf("internal/fonts/%v.txt", fileName))
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v.txt", fileName)
	}

	// чекинг хэша

	hashString := fmt.Sprintf("%x", sha256.Sum256(FileContent))

	switch fileName {
	case "standard":
		if hashString != standardHash {
			return nil, hashError
		}
	case "shadow":
		if hashString != shadowHash {
			return nil, hashError
		}
	case "thinkertoy":
		if hashString != thinkertoyHash {
			return nil, hashError
		}
	}

	// создание массива из базы ASCII

	for i := 0; i < len(FileContent); i++ {
		el := string(FileContent[i])
		if el == "\r" {
			continue
		} else if el != "\n" {
			Oneline += el
		} else {
			ASCIIbase = append(ASCIIbase, Oneline)
			Oneline = ""
		}
	}

	return ASCIIbase, nil
}

// arg == argument
func InputBase(arg string) ([]string, error) {
	// работа с инпутом (в случае, если в нем будет присутствовать \n или \r, создаем массив из стрингов)
	ArrayElement := ""
	InputBase := []string{}

	for i := 0; i < len(arg); i++ {
		if arg[i] == 10 || arg[i] == 13 {
			InputBase = append(InputBase, ArrayElement)
			ArrayElement = ""
			continue
		}
		if i+1 < len(arg) && arg[i] == '\\' && arg[i+1] == 'n' {
			j := i
			countOfSlashes := 0
			for j >= 0 && arg[j] == '\\' {
				countOfSlashes++
				j--
			}
			if countOfSlashes%2 != 0 {
				if ArrayElement == "" {
					InputBase = append(InputBase, ArrayElement)
					i++
				} else {
					InputBase = append(InputBase, ArrayElement)
					ArrayElement = ""
					i++
				}
			}
		} else {
			ArrayElement += string(arg[i])
			if i == len(arg)-1 {
				InputBase = append(InputBase, ArrayElement)
			}
		}
	}

	return InputBase, nil
}

func PrintAscii(InputBase, ASCIIbase []string, arg string) (string, error) {
	var result string
	for j := 0; j < len(InputBase); j++ {
		for _, el := range InputBase[j] {
			if el < 32 || el > 126 {
				return "", errors.New("ERROR: The argument contains an invalid rune type")
			}
		}
	}
	for j := 0; j < len(InputBase); j++ {

		if InputBase[j] == "" {
			result += "\n"
			continue
		}
		for i := 1; i <= 8; i++ {

			for _, el := range InputBase[j] {
				result += (ASCIIbase[((el-rune(32))*9)+rune(i)])
			}

			result += "\n"
		}
	}
	if strings.HasSuffix(arg, "\\n") && !hasOnlyNewLines(InputBase) {
		result += "\n"
	}
	return result, nil
}
func hasOnlyNewLines(Input []string) bool {
	for i := range Input {
		if len(Input[i]) > 0 {
			return false
		}
	}
	return true
}
