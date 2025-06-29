// Данные для тестирования multichecker-а
package main

import "os"

// HelperFunction Пакет main, но без функции main - анализатор не должен срабатывать
func HelperFunction() {
	os.Exit(0) // OK - не в функции main
}

// AnotherHelper Пакет main, но без функции main - анализатор не должен срабатывать (дополнительная вложенность)
func AnotherHelper() {
	if true {
		os.Exit(1) // OK - не в функции main
	}
}
