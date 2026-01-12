// ABOUTME: Tests for the SM-2 spaced repetition algorithm
// ABOUTME: Ensures correct interval and ease factor calculations for adaptive learning

package srs

import (
	"testing"
	"time"
)

func TestNewCard(t *testing.T) {
	card := NewCard("test-skill")

	if card.SkillID != "test-skill" {
		t.Errorf("expected skill ID 'test-skill', got '%s'", card.SkillID)
	}
	if card.EaseFactor != DefaultEaseFactor {
		t.Errorf("expected default ease factor %.2f, got %.2f", DefaultEaseFactor, card.EaseFactor)
	}
	if card.Interval != 0 {
		t.Errorf("expected initial interval 0, got %d", card.Interval)
	}
	if card.Repetitions != 0 {
		t.Errorf("expected initial repetitions 0, got %d", card.Repetitions)
	}
}

func TestReview_PerfectScore(t *testing.T) {
	card := NewCard("test-skill")

	// First review with perfect score (5)
	card.Review(5)

	if card.Interval != 1 {
		t.Errorf("expected interval 1 after first review, got %d", card.Interval)
	}
	if card.Repetitions != 1 {
		t.Errorf("expected 1 repetition, got %d", card.Repetitions)
	}
	// Ease factor should increase with perfect score
	if card.EaseFactor <= DefaultEaseFactor {
		t.Errorf("expected ease factor to increase, got %.2f", card.EaseFactor)
	}

	// Second review with perfect score
	card.Review(5)

	if card.Interval != 6 {
		t.Errorf("expected interval 6 after second review, got %d", card.Interval)
	}
	if card.Repetitions != 2 {
		t.Errorf("expected 2 repetitions, got %d", card.Repetitions)
	}

	// Third review - interval should multiply by ease factor
	prevInterval := card.Interval
	card.Review(5)

	expectedInterval := int(float64(prevInterval) * card.EaseFactor)
	if card.Interval < expectedInterval-1 || card.Interval > expectedInterval+1 {
		t.Errorf("expected interval ~%d after third review, got %d", expectedInterval, card.Interval)
	}
}

func TestReview_FailingScore(t *testing.T) {
	card := NewCard("test-skill")

	// Build up some progress
	card.Review(5)
	card.Review(5)
	card.Review(5)

	// Now fail (score < 3)
	card.Review(2)

	if card.Interval != 1 {
		t.Errorf("expected interval reset to 1 after fail, got %d", card.Interval)
	}
	if card.Repetitions != 0 {
		t.Errorf("expected repetitions reset to 0 after fail, got %d", card.Repetitions)
	}
}

func TestReview_EaseFactorBounds(t *testing.T) {
	card := NewCard("test-skill")

	// Repeatedly fail to drive down ease factor
	for i := 0; i < 10; i++ {
		card.Review(0) // Complete blackout
	}

	if card.EaseFactor < MinEaseFactor {
		t.Errorf("ease factor should not go below %.2f, got %.2f", MinEaseFactor, card.EaseFactor)
	}
}

func TestNextReviewDate(t *testing.T) {
	card := NewCard("test-skill")
	now := time.Now()

	card.Review(5)
	nextReview := card.NextReviewDate()

	expected := now.AddDate(0, 0, card.Interval)
	diff := nextReview.Sub(expected)

	// Allow 1 second tolerance
	if diff > time.Second || diff < -time.Second {
		t.Errorf("next review date off by %v", diff)
	}
}

func TestIsDue(t *testing.T) {
	card := NewCard("test-skill")

	// New card should be due immediately
	if !card.IsDue() {
		t.Error("new card should be due")
	}

	// Review it
	card.Review(5)

	// Should not be due right after review
	if card.IsDue() {
		t.Error("card should not be due right after review")
	}

	// Simulate time passing by manipulating LastReviewed
	card.LastReviewed = time.Now().AddDate(0, 0, -card.Interval-1)

	if !card.IsDue() {
		t.Error("card should be due after interval passed")
	}
}

func TestStrength(t *testing.T) {
	card := NewCard("test-skill")

	// New card has 0 strength
	if card.Strength() != 0 {
		t.Errorf("expected new card strength 0, got %.2f", card.Strength())
	}

	// Review improves strength
	card.Review(5)
	card.Review(5)
	card.Review(5)

	strength := card.Strength()
	if strength <= 0 || strength > 1.0 {
		t.Errorf("expected strength between 0 and 1, got %.2f", strength)
	}

	// Strength should decay over time
	card.LastReviewed = time.Now().AddDate(0, 0, -30)
	decayedStrength := card.Strength()

	if decayedStrength >= strength {
		t.Errorf("expected strength to decay over time, was %.2f now %.2f", strength, decayedStrength)
	}
}

func TestGradeQuality(t *testing.T) {
	tests := []struct {
		correct  bool
		timeMs   int64
		expected int
	}{
		{true, 500, 5},   // Fast correct = perfect
		{true, 2000, 4},  // Medium correct = good
		{true, 5000, 3},  // Slow correct = pass
		{false, 1000, 2}, // Incorrect with hesitation
		{false, 500, 1},  // Quick wrong guess
	}

	for _, tc := range tests {
		grade := CalculateGrade(tc.correct, tc.timeMs)
		if grade != tc.expected {
			t.Errorf("CalculateGrade(correct=%v, time=%d) = %d, want %d",
				tc.correct, tc.timeMs, grade, tc.expected)
		}
	}
}
