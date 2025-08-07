package telegram

import (
	"sync"
	"time"
)

type DurationLike = time.Duration
type TimeLike = time.Time

type botState struct {
	mu sync.RWMutex

	subscribers       map[int64]bool
	intervals         map[int64]time.Duration
	lastSent          map[int64]time.Time
	preferredCurrency map[int64]string
}

func newBotState() botState {
	return botState{
		subscribers:       make(map[int64]bool),
		intervals:         make(map[int64]time.Duration),
		lastSent:          make(map[int64]time.Time),
		preferredCurrency: make(map[int64]string),
	}
}

func (s *botState) ensureUser(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.preferredCurrency[chatID]; !ok {
		s.preferredCurrency[chatID] = "usd"
	}
	if _, ok := s.intervals[chatID]; !ok {
		s.intervals[chatID] = 5 * time.Minute
	}
}

func (s *botState) isSubscribed(chatID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.subscribers[chatID]
}

func (s *botState) setSubscribed(chatID int64, on bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if on {
		s.subscribers[chatID] = true
	} else {
		delete(s.subscribers, chatID)
	}
}

func (s *botState) snapshotSubscribers() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]int64, 0, len(s.subscribers))
	for id := range s.subscribers {
		out = append(out, id)
	}
	return out
}

func (s *botState) setInterval(chatID int64, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.intervals[chatID] = d
	delete(s.lastSent, chatID)
}

func (s *botState) getInterval(chatID int64) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.intervals[chatID]; ok {
		return v
	}
	return 5 * time.Minute
}

func (s *botState) getLastSent(chatID int64) time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSent[chatID]
}

func (s *botState) updateLastSent(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastSent[chatID] = time.Now()
}

func (s *botState) setCurrency(chatID int64, c string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.preferredCurrency[chatID] = c
}

func (s *botState) getCurrency(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if c, ok := s.preferredCurrency[chatID]; ok {
		return c
	}
	return "usd"
}
