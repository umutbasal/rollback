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

    // 1) needs rollback on its own error
    if err := yourawesomeproject.DoSomething(); err != nil {
        rb.Add(func() {
            return yourawesomeproject.UndoSomething()
        })
        return err
    }

    // 2) needs rollback on any error after this point including its own
    rb.Add(func() {
        return yourawesomeproject.UndoSomethingElse()
    })
    if err := yourawesomeproject.DoSomethingElse(); err != nil {
        return err
    }

    // 3) does not need rollback on its own error
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
```

### Clone
This way you can pass rollback functions to upper rollbackers without executing them.

```go
	// Create a new Rollbacker
	rollbacker := New()
	rollbacker.Add(func() {
		fmt.Println("Rollbacker 1 function 1")
	})
	rollbacker.Add(func() {
		fmt.Println("Rollbacker 1 function 2")
	})

	// Create another Rollbacker
	rollbacker2 := New()
	rollbacker2.Add(func() {
		fmt.Println("Rollbacker 2 function 1")
	})

	// Clone the Rollbackers
	rbfn := rollbacker.Clone().Rollback
	rbfn2 := rollbacker2.Clone().Rollback

	// Create an upper Rollbacker
	upperRollbacker := New()
	upperRollbacker.Add(rbfn)
	upperRollbacker.Add(rbfn2)

	// Rollback the upper Rollbacker
	upperRollbacker.Rollback()

	// Output:
	// Rollbacker 2 function 1
	// Rollbacker 1 function 2
	// Rollbacker 1 function 1
```
