# Rollback Package

The `rollback` package provides a simple way to manage and execute rollback operations in Go.

## Usage

```go
package main

import (
    "fmt"
    "github.com/umutbasal/rollback"
	yourawesomeproject "github.com/your/awesome/project/lib"
)

func main() {
	// Create a new rollbacker
    rb := rollback.New()
	// Execute rollback operations on defer
 	defer rb.Rollback()

	if err := yourawesomeproject.DoSomething(); err != nil {
		// needs rollback on its own error
		rb.Add(func() {
			return yourawesomeproject.UndoSomething()
		})
		return err
	}

	// needs rollback on any error after this point including its own
	rb.Add(func() {
		return yourawesomeproject.UndoSomethingElse()
	})
	if err := yourawesomeproject.DoSomethingElse(); err != nil {
		return err
	}

	// does not need rollback on its own error
	if err := yourawesomeproject.DoSomethingMore(); err != nil {
		return err
	}
	// but needs rollback if operations after this point fails
	rb.Add(func() {
		return yourawesomeproject.UndoSomethingMore()
	})

	// No error, finalize and reset rollbacker
	rb.Finalize()
}
