// Package lib Данные для тестирования multichecker-а
package lib

import "os"

// SomeFunc В пакете не main анализатор не должен срабатывать
func SomeFunc() {
	os.Exit(0) // OK - не в пакете main
}

// SomeFunc Даже если есть функция main в пакете не main, анализатор не должен срабатывать
func main() {
	os.Exit(1) // OK - не в пакете main
}
