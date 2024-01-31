package main

import (
	"fmt"
	"math/rand"

	"github.com/umutbasal/rollback"
)

type Container struct {
	Name   string
	Image  string
	Status ContainerStatus
}

const (
	ContainerStatusRunning ContainerStatus = "Running"
	ContainerStatusPending ContainerStatus = "Pending"
)

type ContainerStatus string

func DeployContainer(container *Container) (rollback.Rollback, error) {
	rb := rollback.New()
	defer rb.Rollback()

	rb.Add(func() {
		// rollback container function here
		fmt.Printf("Rolling back container %s\n", container.Name)
	})

	// Deploy container
	var err error
	if container.Image != "nginx:latest" {
		err = fmt.Errorf("invalid image")
		fmt.Printf("Error: %s\n", err)
	}
	if err != nil {
		return nil, err
	}

	rbfn := rb.Clone().Rollback
	rb.Finalize()
	return rbfn, err
}

func IsRunning(container *Container) bool {
	r := rand.Intn(2)
	fmt.Printf("random running status: %d\n", r)
	return r == 1
}

func randomrunningstatus() {
	deployed := map[string]Container{}
	defer fmt.Printf("Deployed containers: %+v\n", deployed)

	rb := rollback.New()
	defer rb.Rollback()

	c := &Container{
		Name:   "nginx",
		Image:  "nginx:latest",
		Status: ContainerStatusPending,
	}

	// Handles its own rollback if return error, so does not need to be added to rollbacker on its error, or before execution
	rbfn, err := DeployContainer(c)
	if err != nil {
		return
	}
	// so we add it rollback function after execution
	// when IsRunning(c) is false, we call rollback function
	rb.Add(rbfn)

	if !IsRunning(c) {
		rb.Add(func() {
			fmt.Printf("delete container %s\n", c.Name)
			delete(deployed, c.Name)
		})

		return
	}

	rb.Finalize()
	deployed[c.Name] = *c
}

func invalidimage() {
	deployed := map[string]Container{}
	defer fmt.Printf("Deployed containers: %+v\n", deployed)

	rb := rollback.New()
	defer rb.Rollback()

	c := &Container{
		Name:   "nginx",
		Image:  "nginxx:latest",
		Status: ContainerStatusPending,
	}

	// Handles its own rollback if return error, so does not need to be added to rollbacker on its error, or before execution
	rbfn, err := DeployContainer(c)
	if err != nil {
		return
	}
	// so we add it rollback function after execution
	// when IsRunning(c) is false, we call rollback function
	rb.Add(rbfn)

	if !IsRunning(c) {
		rb.Add(func() {
			fmt.Printf("delete container %s\n", c.Name)
			delete(deployed, c.Name)
		})

		return
	}

	rb.Finalize()
	deployed[c.Name] = *c
}

func main() {
	fmt.Println("rollback example")
	randomrunningstatus()
	fmt.Println("\n\ninvalid image example")
	invalidimage()

}
