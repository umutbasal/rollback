package rollback

import (
	"fmt"
	"testing"
)

// ExampleRollbacker shows how to use Rollbacker
// It creates a new Rollbacker, adds two rollback functions, and calls Rollbacker.Rollback
// The rollback functions are executed in reverse order.
func ExampleRollbacker_Rollback() {
	// Create a new Rollbacker
	rollbacker := New()
	defer rollbacker.Rollback()

	fmt.Println("Out of rollback operations")

	// Add a rollback function
	rollbacker.Add(func() {
		fmt.Println("Rollback function 1")
	})

	// Add another rollback function
	rollbacker.Add(func() {
		fmt.Println("Rollback function 2")
	})

	// Output:
	// Out of rollback operations
	// Rollback function 2
	// Rollback function 1
}

// ExampleRollbacker_Finalize shows how to use Rollbacker.Finalize
// It creates a new Rollbacker, adds two rollback functions, and calls Rollbacker.Finalize
// The rollback functions are not executed when Rollbacker.Finalize is called.
func ExampleRollbacker_Finalize() {
	// Create a new Rollbacker
	rollbacker := New()
	defer rollbacker.Rollback()

	fmt.Println("Out of rollback operations")

	// Add a rollback function
	rollbacker.Add(func() {
		fmt.Println("Rollback function 1")
	})

	// Add another rollback function
	rollbacker.Add(func() {
		fmt.Println("Rollback function 2")
	})

	// Finalize the Rollbacker
	rollbacker.Finalize()

	// Output:
	// Out of rollback operations
}

// TestRollbacker tests the Rollbacker
func TestRollbacker(t *testing.T) {
	// Create a new Rollbacker
	rollbacker := New()

	// Add a rollback function
	rollbacker.Add(func() {
		t.Log("Rollback function 1")
	})

	// Add another rollback function
	rollbacker.Add(func() {
		t.Log("Rollback function 2")
	})

	// Test Finalize
	rollbacker.Finalize()

	if len(rollbacker.rollbackFns) != 0 {
		t.Error("Finalize() should clear rollback functions")
	}
}

// TestRollbackerRollback tests the Rollbacker.Rollback method
func TestRollbackerRollback(t *testing.T) {
	// Create a new Rollbacker
	rollbacker := New()

	want := 0
	field := struct {
		value int
	}{
		value: 0,
	}

	field.value = 1

	// Add a rollback function
	rollbacker.Add(func() {
		field.value = 0
	})

	rollbacker.Rollback()

	if field.value != want {
		t.Errorf("Rollbacker.Rollback() = %d, want %d", field.value, want)
	}
}

// TestRollbackerClone tests the Rollbacker.Clone method
func TestRollbackerClone(t *testing.T) {
	// Create a new Rollbacker
	rollbacker := New()
	defer rollbacker.Rollback()

	// Add a rollback function
	rollbacker.Add(func() {
		t.Log("Rollback function 1")
	})

	// Add another rollback function
	rollbacker.Add(func() {
		t.Log("Rollback function 2")
	})

	// Clone the Rollbacker
	clone := rollbacker.Clone()

	// Add a rollback function to the clone
	clone.Add(func() {
		t.Log("Rollback function 3")
	})

	// Test Finalize
	clone.Finalize()

	if len(clone.rollbackFns) != 0 {
		t.Error("Finalize() should clear rollback functions")
	}

	if len(rollbacker.rollbackFns) != 2 {
		t.Error("Rollbacker should not be affected by clone")
	}
}

// ExampleRollbackerClone_2
func ExampleRollbacker_Clone() {
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
}

func sub() (rollbackFn func()) {
	return func() {
		fmt.Println("rollback sub")
	}
}

func NativeDefers() {}

// ExampleNativeDefers
func ExampleNativeDefers() {
	var err error
	var err1, err2, err3 error

	err = err1
	defer func() {
		if err != nil {
			fmt.Printf("Error1: %s\n", err)
		}
	}()

	sfn := sub()
	defer func() {
		sfn()
	}()

	err = err2
	defer func() {
		if err != nil {
			fmt.Printf("Error2: %s\n", err)
		}
	}()

	err3 = fmt.Errorf("error foo")
	err = err3
	defer func() {
		if err != nil {
			fmt.Printf("Error3: %s\n", err)
		}
	}()

	// Output:
	// Error3: error foo
	// Error2: error foo
	// rollback sub
	// Error1: error foo
}
