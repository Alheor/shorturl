// Данные для тестирования multichecker-а
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting application")

	// Этот вызов должен быть обнаружен анализатором
	os.Exit(0) // want "прямой вызов os.Exit в функции main запрещен"

	// Этот тоже
	if true {
		os.Exit(1) // want "прямой вызов os.Exit в функции main запрещен"
	}
}

// someOtherFunc Анализатор не должен срабатывать, так как не в main
func someOtherFunc() {
	os.Exit(1)
}

// helperFunc Анализатор не должен срабатывать
func helperFunc() {
	os.Exit(2)
}
