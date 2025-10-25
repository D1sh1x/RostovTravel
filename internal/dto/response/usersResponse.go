package response

import "time"

type UserResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	PasswordHash string             `json:"password_hash"`
	CreatedAt    time.Time          `json:"created_at"`
	Favorites    []FavoriteResponse `json:"favorites"`
	Role         string             `json:"role"`
}

type FavoriteResponse struct {
	Type   string `json:"type"`
	ItemID string `json:"item_id"`
}
