package domain

import (
	"fmt"
	"strings"
)

type ScoreBand string

const (
	BandCold     ScoreBand = "cold"
	BandWarm     ScoreBand = "warm"
	BandHot      ScoreBand = "hot"
	BandCritical ScoreBand = "critical"
)

type Score struct {
	Value   int       `json:"value"`
	Band    ScoreBand `json:"band"`
	Reasons []string  `json:"reasons"`
}

func ScoreRecord(record Record) Score {
	value := 0
	reasons := make([]string, 0, 8)
	switch record.Priority {
	case PriorityUrgent:
		value += 50
		reasons = append(reasons, "urgent priority")
	case PriorityHigh:
		value += 35
		reasons = append(reasons, "high priority")
	case PriorityNormal:
		value += 20
		reasons = append(reasons, "normal priority")
	default:
		value += 5
		reasons = append(reasons, "low priority")
	}
	if record.Status == StatusApproved {
		value += 20
		reasons = append(reasons, "approved")
	}
	if record.Status == StatusReview {
		value += 10
		reasons = append(reasons, "in review")
	}
	if record.Owner != "" {
		value += 8
		reasons = append(reasons, "assigned owner")
	}
	if len(record.Tags) > 0 {
		value += 5
		reasons = append(reasons, "tagged")
	}
	if strings.Contains(strings.ToLower(record.Need), "budget") {
		value += 7
		reasons = append(reasons, "budget signal")
	}
	if value > 100 {
		value = 100
	}
	return Score{Value: value, Band: bandFor(value), Reasons: reasons}
}

func bandFor(value int) ScoreBand {
	if value >= 80 {
		return BandCritical
	}
	if value >= 55 {
		return BandHot
	}
	if value >= 25 {
		return BandWarm
	}
	return BandCold
}

func ExplainScore(score Score) string {
	if len(score.Reasons) == 0 {
		return fmt.Sprintf("%s (%d)", score.Band, score.Value)
	}
	return fmt.Sprintf("%s (%d): %s", score.Band, score.Value, strings.Join(score.Reasons, ", "))
}

func IsActionable(score Score) bool { return score.Band == BandHot || score.Band == BandCritical }
