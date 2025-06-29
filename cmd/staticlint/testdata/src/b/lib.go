package lib

import "os"

// В пакете не main анализатор не должен срабатывать
func SomeFunc() {
	os.Exit(0) // OK - не в пакете main
}

func main() {
	// Даже если есть функция main в пакете не main
	os.Exit(1) // OK - не в пакете main
}
