package dto

type QueryRoleGame struct {
	RoomId string `query:"room_id" validate:"required"`
	UserId string `query:"user_id" validate:"required"`
}

type Game struct {
	GameId    string `json:"game_id"`
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	SlotIndex int    `json:"slot_index"`
	IsDead    bool   `json:"is_dead"`
	Img       string `json:"img"`
	IsJoin    bool   `json:"is_join"`
	Role      string `json:"role"`
}

type GameResponse struct {
	GameId     string `json:"game_id"`
	UserId     string `json:"user_id"`
	UserName   string `json:"user_name"`
	SlotIndex  int    `json:"slot_index"`
	IsDead     bool   `json:"is_dead"`
	Img        string `json:"img"`
	IsJoin     bool   `json:"is_join"`
	IsWerewolf bool   `json:"is_werewolf"`
}

type Vote struct {
	UserID string `json:"user_id"`
	Vote   *int    `json:"vote"`
	Role   string `json:"role"`
}

type SummaryVote struct {
	Index int `json:"index"`
	Count int `json:"count"`
}
