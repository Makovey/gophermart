package model

import "github.com/Makovey/gophermart/internal/types"

type RewardType string

const (
	Percent RewardType = "%"
	Fixed   RewardType = "pt"
)

type Goods struct {
	Match      string             `json:"match"`
	Reward     types.FloatDecimal `json:"reward"`
	RewardType RewardType         `json:"reward_type"`
}
