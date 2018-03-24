package core

import "time"

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	BLUE   = "\033[34m"
	YELLOW = "\033[33m"
	RESET  = "\033[0m"
)

func printElement(c string, s string) string {
	return c + s + RESET
}

func P_ERROR(s string) string {
	return printElement(RED, "[error]: "+s)
}

func P_SUCCESS(s string) string {
	return printElement(GREEN, "[success]: "+s)
}

func p_WARNING(s string) string {
	return printElement(YELLOW, "[warning]: "+s)
}

func P_INFO(s string) string {
	return printElement(BLUE, "[information]: "+s)
}

func P_TIME(c string, s string) string {
	now := string(time.Now().Format("2006-01-02 15:04:05"))
	s = "[" + now + "] " + s
	return printElement(c, s)
}

func P_SIMPLE(c string, s string) string {
	return printElement(c, s)
}
