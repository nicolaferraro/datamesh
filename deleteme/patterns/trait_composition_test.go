package patterns

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)


// A car can accelerate and steer

type Accelerable struct {
	pow int
}

func (a *Accelerable) Accelerate(c chan<- string) {
	c <- fmt.Sprintf("ACCELERATE-%v", a.pow)
}

func (a *Accelerable) InitAccelerable(pow int) {
	a.pow = pow
}

type Steerable struct {}

func (s *Steerable) InitSteerable() {
}

func (s *Steerable) Steer(c chan<- string) {
	c <- "STEER"
}

type Car struct {
	Accelerable
	Steerable
}

func NewCar(pow int) *Car {
	var c Car
	c.InitAccelerable(pow)
	c.InitSteerable()
	return &c
}

func TestCar(t *testing.T) {

	c := make(chan string, 2)

	car := NewCar(2)
	car.Accelerate(c)
	car.Steer(c)
	close(c)

	res := make([]string, 0, 2)
	for s := range c {
		res = append(res, s)
	}
	assert.Equal(t,[]string{"ACCELERATE-2", "STEER"}, res)
}