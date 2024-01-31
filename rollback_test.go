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
