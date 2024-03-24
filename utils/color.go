package utils

func Bold(text string) string {
	return "\u001B[1m" + text + "\033[0m"
}
func Purple(text string) string {
	return "\033[0;35m" + text + "\033[0m"
}
func Red(s string) string {
	return "\u001b[31m" + s + "\033[0m"
}

func Cyan(s string) string {
	return "\u001b[36m" + s + "\033[0m"
}
func Brown(s string) string {
	return "\033[0;33m" + s + "\033[0m"
}
func Blue(s string) string {
	return "\u001b[34m" + s + "\033[0m"
}
func Magenta(s string) string {
	return "\u001b[35m" + s + "\033[0m"
}
func Green(s string) string {
	return "\u001b[32m" + s + "\033[0m"
}

func Yellow(s string) string {
	return "\u001b[33m" + s + "\033[0m"
}

func BoldRed(s string) string {
	return "\033[1m\u001b[31m" + s + "\033[0m"
}
