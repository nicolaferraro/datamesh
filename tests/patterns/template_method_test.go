package patterns

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


type Algorithm struct {
	transform func(data int) int
}

func (a *Algorithm) Solve(input ...int) int {
	acc := 0
	for _, i := range input {
		acc += a.transform(i)
	}
	return acc
}


func TestAdder(t *testing.T) {
	adder := Algorithm{
		transform: func(data int) int {
			return data
		},
	}

	assert.Equal(t, 2, adder.Solve(1, 1, 0))
}

func TestSquareAdder(t *testing.T) {
	squareAdder := Algorithm{
		transform: func(data int) int {
			return data * data
		},
	}

	assert.Equal(t, 14, squareAdder.Solve(1, 2, 3))
}