package port

type Repository interface {
	AuthRepo
	RoomRepo
	UserRepo
}

type Service interface {
	AuthService
	RoomService
	UserService
	GameService
}
