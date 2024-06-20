package utils

import "fmt"

const (
	BlueColor   = "\033[38;2;137;180;250m"
	GreenColor  = "\033[38;5;157m"
	OrangeColor = "\033[38;2;250;179;135m"
	PurpleColor = "\033[38;5;183m"
	RedColor    = "\033[38;2;243;139;168m"
	ResetColor  = "\033[0m"
)

func PrintRepoEntry(path string, msg string) {
	fmt.Printf("%s=== %s ===%s\n%s\n", BlueColor, path, ResetColor, msg)
}

func PrintSection(msg string) {
	fmt.Printf("%s%s%s\n", GreenColor, msg, ResetColor)
}
func PrintSeparator() {
	fmt.Printf("%s--------------------%s\n", PurpleColor, ResetColor)
}

func PrintErrorMsg(msg string) {
	fmt.Printf("%s%s%s\n", RedColor, msg, ResetColor)
}
