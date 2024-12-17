package model

type RewardType string

const (
	Percent RewardType = "%"
	Fixed   RewardType = "pt"
)

type Goods struct {
	Match      string     `json:"match"`
	Reward     float64    `json:"reward"`
	RewardType RewardType `json:"reward_type"`
}
