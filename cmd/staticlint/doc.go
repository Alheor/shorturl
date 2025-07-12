// Package staticlint это расширенный статический анализатор, содержащий multichecker
// для статического анализа кода, а так же локальный анализатора noexit.
//
// # Описание
//
// Staticlint - объединяющий множество анализаторов для комплексной
// проверки качества кода. Он включает стандартные анализаторы Go, анализаторы из
// staticcheck, а также дополнительные публичные и пользовательские анализаторы.
//
// Multichecker включает в себя:
//
// # Стандартные анализаторы golang.org/x/tools/go/analysis/passes:
//   - asmdecl - mismatches between assembly files and Go declarations.
//   - assign - detects useless assignments.
//   - atomic - checks for common mistakes using the sync/atomic package.
//   - bools - common mistakes involving boolean operators.
//   - buildtag - checks build tags.
//   - cgocall - detects some violations of the cgo pointer passing rules.
//   - composite - checks for unkeyed composite literals.
//   - copylock - checks for locks erroneously passed by value.
//   - errorsas - checks that the second argument to errors.As is a pointer to a type implementing error.
//   - httpresponse - checks for mistakes using HTTP responses.
//   - loopclosure - checks for references to enclosing loop variables from within nested functions.
//   - lostcancel - checks for failure to call a context cancellation function.
//   - nilfunc - checks for useless comparisons against nil.
//   - printf - consistency of Printf format strings and arguments.
//   - shift - checks for shifts that exceed the width of an integer.
//   - stdmethods - checks for misspellings in the signatures of methods similar to well-known interfaces.
//   - structtag - struct field tags are well formed.
//   - tests - checks for common mistaken usages of tests and examples.
//   - unmarshal - checks for passing non-pointer or non-interface types to unmarshal and decode functions.
//   - unreachable - checks for unreachable code.
//   - unsafeptr - checks for invalid conversions of uintptr to unsafe.Pointer.
//   - unusedresult - checks for unused results of calls to certain pure functions.
//
// # Анализаторы staticcheck.io:
//   - Все анализаторы класса SA (staticcheck)
//   - S1000 - Use plain channel send or receive instead of single-case select
//   - S1001 - Replace for loop with call to copy
//   - ST1000 - Incorrect or missing package comment
//   - ST1001 - Dot imports are discouraged
//
// # Публичные анализаторы:
//   - bodyclose (github.com/timakin/bodyclose/passes/bodyclose) - checks whether res.Body is correctly closed.
//   - durationcheck (github.com/charithe/durationcheck) - cases where two time.Duration values are being multiplied in possibly erroneous ways
//   - sqlrows (github.com/gostaticanalysis/sqlrows) - uncover bugs by reporting a diagnostic for mistakes of sql.Rows usage.
//
// # Пользовательский анализатор:
//   - noexit - запрещает прямой вызов os.Exit в функции main пакета main
//
// # Использование
//
// Запуск из корня проекта:
//
//	go run ./cmd/staticlint/... ./...
//
// Запуск глобально:
//
//	go install github.com/Alheor/shorturl/cmd/staticlint@latest
//	go run github.com/Alheor/shorturl/cmd/staticlint ./...
//
// # Конфигурация
//
// Staticlint не требует дополнительной конфигурации и готов к использованию
// сразу после установки. Все анализаторы включены по умолчанию.
package main
