// ABOUTME: Tests for the generic stack data structure used by the interpreter
// ABOUTME: Verifies stack overflow protection and basic stack operations

package interpreter

import (
	"testing"

	"basic-interpreter/types"
)

func TestStack_Push_Overflow(t *testing.T) {
	// Create a small stack with capacity 2 for testing
	stack := NewStack[int](2)

	// Push first item - should succeed
	err := stack.Push(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Push second item - should succeed
	err = stack.Push(2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Push third item - should fail with overflow
	err = stack.Push(3)
	if err != ErrStackOverflow {
		t.Errorf("Expected ErrStackOverflow, got %v", err)
	}

	// Verify size is still 2
	if stack.Size() != 2 {
		t.Errorf("Expected size 2, got %d", stack.Size())
	}
}

func TestStack_ForLoopContext_Overflow(t *testing.T) {
	// Create a stack for FOR loop contexts with capacity 1
	stack := NewStack[ForLoopContext](1)

	// Create a FOR loop context
	forLoop := ForLoopContext{
		Variable:          "I",
		EndValue:          types.NewNumberValue(10),
		StepValue:         types.NewNumberValue(1),
		AfterForLineIndex: 1,
		AfterForStmtIndex: 0,
	}

	// Push first context - should succeed
	err := stack.Push(forLoop)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Push second context - should fail with overflow
	err = stack.Push(forLoop)
	if err != ErrStackOverflow {
		t.Errorf("Expected ErrStackOverflow, got %v", err)
	}
}

func TestStack_CallContext_Overflow(t *testing.T) {
	// Create a stack for call contexts with capacity 1
	stack := NewStack[CallContext](1)

	// Create a call context
	callCtx := CallContext{ReturnLineIndex: 10}

	// Push first context - should succeed
	err := stack.Push(callCtx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Push second context - should fail with overflow
	err = stack.Push(callCtx)
	if err != ErrStackOverflow {
		t.Errorf("Expected ErrStackOverflow, got %v", err)
	}
}

func TestStack_BasicOperations(t *testing.T) {
	stack := NewStack[string](10)

	// Test push and size
	err := stack.Push("first")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stack.Size() != 1 {
		t.Errorf("Expected size 1, got %d", stack.Size())
	}

	// Test peek
	top := stack.Peek()
	if top == nil || *top != "first" {
		t.Errorf("Expected 'first', got %v", top)
	}
	if stack.Size() != 1 {
		t.Errorf("Expected size still 1 after peek, got %d", stack.Size())
	}

	// Test pop
	popped := stack.Pop()
	if popped == nil || *popped != "first" {
		t.Errorf("Expected 'first', got %v", popped)
	}
	if stack.Size() != 0 {
		t.Errorf("Expected size 0 after pop, got %d", stack.Size())
	}

	// Test empty operations
	if !stack.IsEmpty() {
		t.Error("Expected stack to be empty")
	}

	if stack.Peek() != nil {
		t.Error("Expected nil when peeking empty stack")
	}

	if stack.Pop() != nil {
		t.Error("Expected nil when popping empty stack")
	}
}
