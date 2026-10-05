package main

import "fmt"

type Bike struct{}

func (b Bike) Start() {
	fmt.Println("Bike started")
}

type Car struct{}

func (c Car) Start() {
	fmt.Println("Car started")
}

type Vehicle interface {
	Start()
}

func main() {
	var v1 Vehicle = Bike{}
	v1.Start()

	var v2 Vehicle = Car{}
	v2.Start()
}
