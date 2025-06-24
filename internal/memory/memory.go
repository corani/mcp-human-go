package memory

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        string
	Question  string
	Context   string
	Answer    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemoryDB is an in-memory database for mapping (question, context) to answer.
type entry struct {
	Question Question
	Callback func(answer string)
}

type MemoryDB struct {
	mu      sync.RWMutex
	entries map[string]*entry // key: question ID
}

// No longer needed

// NewMemoryDB creates a new MemoryDB instance.
func NewMemoryDB() *MemoryDB {
	db := &MemoryDB{
		entries: make(map[string]*entry),
	}

	// Add some dummy questions for testing
	db.Add("What is the capital of France?", "Geography", nil)
	db.Add("Who wrote '1984'?", "Literature", nil)
	db.Add("What is 2 + 2?", "Math", nil)

	return db
}

// Add stores a new (question, context) with a callback, generates and returns the unique question ID.
func (db *MemoryDB) Add(question, context string, callback func(answer string)) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.NewString()
	now := time.Now()

	db.entries[id] = &entry{
		Question: Question{
			ID:        id,
			Question:  question,
			Context:   context,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Callback: callback,
	}

	return id, nil
}

// Get retrieves the entry by question ID.
func (db *MemoryDB) Get(id string) (Question, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	e, found := db.entries[id]
	if !found {
		return Question{}, false
	}

	return e.Question, true
}

// ListQuestions returns a slice of all questions in the database.
func (db *MemoryDB) ListQuestions() []Question {
	db.mu.RLock()
	defer db.mu.RUnlock()

	questions := make([]Question, 0, len(db.entries))
	for _, e := range db.entries {
		questions = append(questions, e.Question)
	}

	sort.Slice(questions, func(i, j int) bool {
		return questions[i].CreatedAt.Before(questions[j].CreatedAt)
	})

	return questions
}

// UpdateAnswer sets the answer for a question by ID, invokes the callback, and removes the callback.
func (db *MemoryDB) UpdateAnswer(id, answer string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	e, found := db.entries[id]
	if !found {
		return errors.New("question ID not found")
	}

	e.Question.Answer = answer
	e.Question.UpdatedAt = time.Now()

	if e.Callback != nil {
		cb := e.Callback
		e.Callback = nil
		// Unlock before invoking callback to avoid deadlocks if callback interacts with db
		db.mu.Unlock()
		cb(answer)
		db.mu.Lock()
	}

	return nil
}
