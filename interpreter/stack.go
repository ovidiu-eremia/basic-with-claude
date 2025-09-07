// ABOUTME: Generic stack data structure for managing interpreter state stacks
// ABOUTME: Provides type-safe stack operations with bounds checking for FOR loops and GOSUB calls

package interpreter

// Stack provides a generic stack data structure with bounds checking
type Stack[T any] struct {
	items   []T
	maxSize int
}

// NewStack creates a new empty stack with a maximum size limit
func NewStack[T any](maxSize int) *Stack[T] {
	return &Stack[T]{
		items:   make([]T, 0),
		maxSize: maxSize,
	}
}

// Push adds an item to the top of the stack
// Returns an error if the stack would exceed its maximum size
func (s *Stack[T]) Push(item T) error {
	if len(s.items) >= s.maxSize {
		return ErrStackOverflow
	}
	s.items = append(s.items, item)
	return nil
}

// Pop removes and returns the top item from the stack
// Returns nil if the stack is empty
func (s *Stack[T]) Pop() *T {
	if len(s.items) == 0 {
		return nil
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return &top
}

// Peek returns the top item without removing it
// Returns nil if the stack is empty
func (s *Stack[T]) Peek() *T {
	if len(s.items) == 0 {
		return nil
	}
	return &s.items[len(s.items)-1]
}

// IsEmpty returns true if the stack has no items
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// FindByPredicate searches the stack from top to bottom for an item matching the predicate
// Returns a pointer to the first matching item, or nil if none found
func (s *Stack[T]) FindByPredicate(predicate func(T) bool) *T {
	for i := len(s.items) - 1; i >= 0; i-- {
		if predicate(s.items[i]) {
			return &s.items[i]
		}
	}
	return nil
}
