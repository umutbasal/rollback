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
		// Add rollback operation on err
		rb.Add(func() {
			return yourawesomeproject.UndoSomething()
		})
		return err
	}

	if err := yourawesomeproject.DoSomethingElse(); err != nil {
		// Add rollback operation on err
		rb.Add(func() {
			return yourawesomeproject.UndoSomethingElse()
		})
		return err
	}

	// No error, finalize and reset rollbacker
	rb.Finalize()
}
