// ABOUTME: SM-2 spaced repetition algorithm implementation
// ABOUTME: Core learning engine that schedules reviews based on forgetting curves

package srs

import (
	"math"
	"time"
)

const (
	// DefaultEaseFactor is the starting ease factor for new cards (2.5 is standard SM-2)
	DefaultEaseFactor = 2.5
	// MinEaseFactor prevents cards from becoming impossibly hard
	MinEaseFactor = 1.3
)

// Card represents a single learnable skill tracked by the SRS system
type Card struct {
	SkillID      string
	EaseFactor   float64
	Interval     int // days until next review
	Repetitions  int
	LastReviewed time.Time
}

// NewCard creates a new card for tracking a skill
func NewCard(skillID string) *Card {
	return &Card{
		SkillID:      skillID,
		EaseFactor:   DefaultEaseFactor,
		Interval:     0,
		Repetitions:  0,
		LastReviewed: time.Time{}, // zero value = never reviewed
	}
}

// Review processes a review session with the given quality grade (0-5)
// Grade meanings:
//
//	5 - perfect response
//	4 - correct after hesitation
//	3 - correct with difficulty
//	2 - incorrect, but easy recall
//	1 - incorrect, remembered
//	0 - complete blackout
func (c *Card) Review(grade int) {
	// Clamp grade to valid range
	if grade < 0 {
		grade = 0
	}
	if grade > 5 {
		grade = 5
	}

	// Update ease factor using SM-2 formula
	// EF' = EF + (0.1 - (5 - grade) * (0.08 + (5 - grade) * 0.02))
	c.EaseFactor = c.EaseFactor + (0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02))

	// Enforce minimum ease factor
	if c.EaseFactor < MinEaseFactor {
		c.EaseFactor = MinEaseFactor
	}

	// Handle failed reviews (grade < 3)
	if grade < 3 {
		c.Repetitions = 0
		c.Interval = 1
	} else {
		// Successful review - calculate next interval
		switch c.Repetitions {
		case 0:
			c.Interval = 1
		case 1:
			c.Interval = 6
		default:
			c.Interval = int(math.Round(float64(c.Interval) * c.EaseFactor))
		}
		c.Repetitions++
	}

	c.LastReviewed = time.Now()
}

// NextReviewDate returns when this card should next be reviewed
func (c *Card) NextReviewDate() time.Time {
	if c.LastReviewed.IsZero() {
		return time.Now() // New card is due immediately
	}
	return c.LastReviewed.AddDate(0, 0, c.Interval)
}

// IsDue returns true if the card is due for review
func (c *Card) IsDue() bool {
	if c.LastReviewed.IsZero() {
		return true // New cards are always due
	}
	return time.Now().After(c.NextReviewDate()) || time.Now().Equal(c.NextReviewDate())
}

// Strength returns a 0.0-1.0 value representing mastery of this skill
// Takes into account repetitions, ease factor, and time since last review
func (c *Card) Strength() float64 {
	if c.Repetitions == 0 && c.LastReviewed.IsZero() {
		return 0 // Never practiced
	}

	// Base strength from repetitions (asymptotic approach to 1.0)
	// More reps = stronger, but diminishing returns
	baseStrength := 1.0 - math.Pow(0.7, float64(c.Repetitions))

	// Decay based on time since last review
	if !c.LastReviewed.IsZero() {
		daysSince := time.Since(c.LastReviewed).Hours() / 24
		// Forgetting curve: strength decays exponentially
		// Decay rate is slower for higher ease factors (well-learned items)
		decayRate := 0.1 / c.EaseFactor
		decay := math.Exp(-decayRate * daysSince)
		baseStrength *= decay
	}

	// Clamp to valid range
	if baseStrength < 0 {
		baseStrength = 0
	}
	if baseStrength > 1.0 {
		baseStrength = 1.0
	}

	return baseStrength
}

// CalculateGrade determines the quality grade based on correctness and response time
// This translates user performance into SM-2 grade (0-5)
func CalculateGrade(correct bool, responseTimeMs int64) int {
	if correct {
		// Fast correct (< 1s) = perfect
		if responseTimeMs < 1000 {
			return 5
		}
		// Medium (1-3s) = good
		if responseTimeMs < 3000 {
			return 4
		}
		// Slow but correct = pass
		return 3
	}

	// Incorrect responses
	// Hesitated before wrong answer = some recall
	if responseTimeMs >= 1000 {
		return 2
	}
	// Quick wrong = guessing
	return 1
}
