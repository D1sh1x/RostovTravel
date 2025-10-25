package request

type UserRequest struct {
	Name      string            `json:"name"`
	Username  string            `json:"username"`
	Password  string            `json:"password"`
	Favorites []FavoriteRequest `json:"favorites"`
	Role      string            `json:"role"`
}

type FavoriteRequest struct {
	Type   string `bson:"type" json:"type"`
	ItemID string `bson:"item_id" json:"item_id"`
}
