package port

type Repository interface {
	AuthRepo
	RoomRepo
	UserRepo
	GameRepo
}

type Service interface {
	AuthService
	RoomService
	UserService
	GameService
}
