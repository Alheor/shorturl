package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting application")

	// Этот вызов должен быть обнаружен анализатором
	os.Exit(0) // want "прямой вызов os.Exit в функции main запрещен"

	// Еще один вызов для проверки
	if true {
		os.Exit(1) // want "прямой вызов os.Exit в функции main запрещен"
	}
}

func someOtherFunc() {
	// Этот вызов не должен быть обнаружен, так как не в main
	os.Exit(1)
}

func helperFunc() {
	// Этот тоже не должен быть обнаружен
	os.Exit(2)
}
