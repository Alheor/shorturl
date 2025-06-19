package models

type contextKeyXAuthUser string

// CookiesName - имя ключа cookie.
const CookiesName = `authKey`

// ContextValueName - имя ключа cookie при передаче через контекст.
const ContextValueName contextKeyXAuthUser = `xAuthUser`

// User - структура авторизованного пользователя.
type User struct {
	ID string
}

// UserCookie - структура пользовательской cookie.
type UserCookie struct {
	User User
	Sign []byte
}

// EmptyUserIDErr - тип ошибки, обозначающий, что пользователь в cookie не авторизован.
type EmptyUserIDErr struct {
	Err error
}

// Error реализация интерфейса Error
func (e *EmptyUserIDErr) Error() string {
	return e.Err.Error()
}
