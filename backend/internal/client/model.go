package client

import "time"

// Client represents a rewards program participant.
type Client struct {
	ID           string
	Name         string
	Email        string
	PointBalance int64
}

// Award represents points awarded to a client for an action.
type Award struct {
	ID            string
	ClientID      string
	Type          AwardType
	PointsAwarded int64
	CreatedAt     time.Time
}

// AwardType enumerates valid award categories.
type AwardType string

const (
	AwardMonthlyContribution  AwardType = "MONTHLY_CONTRIBUTION"
	AwardIncreaseContribution AwardType = "INCREASE_CONTRIBUTION"
	AwardLoyalty12Months      AwardType = "LOYALTY_12_MONTHS"
	AwardFinancialAssessment  AwardType = "FINANCIAL_ASSESSMENT"
)

// AwardPointsMap defines points for each award type.
var AwardPointsMap = map[AwardType]int64{
	AwardMonthlyContribution:  100,
	AwardIncreaseContribution: 200,
	AwardLoyalty12Months:      1000,
	AwardFinancialAssessment:  500,
}

// ValidAwardType checks if a string is a valid AwardType.
func ValidAwardType(s string) bool {
	_, ok := AwardPointsMap[AwardType(s)]
	return ok
}

// PointsForAward returns points for an award type, or 0 if invalid.
func PointsForAward(t AwardType) int64 {
	return AwardPointsMap[t]
}
