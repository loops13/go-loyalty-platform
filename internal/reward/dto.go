package reward

// DTOs for reward domain.

// RewardResp is the response containing reward details.
type RewardResp struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PointsCost int64  `json:"pointsCost"`
}

// RedeemReq is the request to redeem a reward.
type RedeemReq struct {
	RewardID string `json:"rewardId"`
}

// RedeemResp is the response containing redemption details.
type RedeemResp struct {
	Reward  RewardResp `json:"reward"`
	Balance int64      `json:"balance"`
}
