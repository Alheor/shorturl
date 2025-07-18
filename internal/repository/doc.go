// Package repository - сервис репозиториев.
//
// # Описание
//
// В зависимости от конфигурации, репозиторий может работать либо с БД PostgreSQL, либо хранить данные в файле или памяти.
// Согласно загрузившейся конфигурации, создается необходимый экземпляр репозитория. Все экземпляры имплементируют интерфейс IRepository,
// что позволяет работать с любым из них независимо.
//
// При использовании базы банных, схема будет создана автоматически, если в указанной БД ее нет.
// Для начала работы достаточно пустой БД, все остальное сервис сделает сам.
//
// BUG(repository): Обратите внимание, что методы RemoveBatch и RemoveByOriginalURL поддерживаются только при работе с БД.
package repository
