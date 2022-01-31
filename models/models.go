package models

type User struct {
    ID              string `json:"id"`
    Streak          int64  `json:"streak"`
    HighestStreak   int64  `json:"highestStreak"`
    TotalLearned    int64  `json:"totalLearned"`
		ReviewPoints    int64  `json:"reviewPoints"`
		LastLearned     int64  `json:"lastLearned"`
		UsesTraditional bool   `json:"usesTraditional"`
}