package main

import "os"

// Пакет main, но без функции main - не должно быть проверок
func HelperFunction() {
	os.Exit(0) // OK - не в функции main
}

func AnotherHelper() {
	if true {
		os.Exit(1) // OK - не в функции main
	}
}
