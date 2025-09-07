package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &EndStatement{Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
	assert.True(t, mock.endRequested)
}

func TestStopStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &StopStatement{Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
	assert.True(t, mock.stopRequested)
}

func TestRunStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &RunStatement{Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
}

func TestGotoStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &GotoStatement{TargetLine: 50, Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
	assert.True(t, mock.gotoRequested)
	assert.Equal(t, 50, mock.gotoTarget)
}
