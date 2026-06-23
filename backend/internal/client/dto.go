package client

// DTOs for client domain.

// CreateReq is the request to create a new client.
type CreateReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ClientResp is the response containing client details.
type ClientResp struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PointBalance int64  `json:"pointBalance"`
}

// AwardReq is the request to award points.
type AwardReq struct {
	Type string `json:"type"`
}

// AwardResp is the response containing award details.
type AwardResp struct {
	ID            string `json:"id"`
	ClientID      string `json:"clientId"`
	Type          string `json:"type"`
	PointsAwarded int64  `json:"pointsAwarded"`
	CreatedAt     string `json:"createdAt"`
}
