package models

type contextKeyXAuthUser string

const CookiesName = `authKey`
const ContextValueName contextKeyXAuthUser = `xAuthUser`

type User struct {
	ID string
}

type UserCookie struct {
	User User
	Sign []byte
}

type EmptyUserIDErr struct {
	Err error
}

func (e *EmptyUserIDErr) Error() string {
	return e.Err.Error()
}
