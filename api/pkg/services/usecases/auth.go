package usecases

type AuthUsecase interface {
	CreateUser(login, password string) (string, error)
	LoginUser(login, password string) (string, error)
	Authorization(token string) (bool, error)
}
