package rollback

// Rollback represents a function that can be used to undo a set of operations.
type Rollback func()

// Rollbacker represents a struct that holds a collection of rollback functions.
type Rollbacker struct {
	rollbackFns []Rollback
}

// New creates a new instance of Rollbacker.
func New() *Rollbacker {
	return &Rollbacker{}
}

// Clone rollbacker
// Clone creates a deep copy of the Rollbacker instance.
// It returns a new Rollbacker instance with the same rollback functions as the original.
func (r *Rollbacker) Clone() *Rollbacker {
	clone := New()
	clone.rollbackFns = append(clone.rollbackFns, r.rollbackFns...)
	return clone
}

// Add adds a rollback function to the Rollbacker.
// The rollback function will be executed when Rollbacker.Rollback is called.
func (r *Rollbacker) Add(fn Rollback) {
	r.rollbackFns = append(r.rollbackFns, fn)
}

// Finalize clears the rollback functions slice, indicating that the rollback process has been completed.
func (r *Rollbacker) Finalize() {
	r.rollbackFns = nil
}

// Rollback is a method of the Rollbacker struct that executes the rollback functions in reverse order.
// It iterates through the rollbackFns slice and calls each function in reverse order.
func (r *Rollbacker) Rollback() {
	for i := len(r.rollbackFns) - 1; i >= 0; i-- {
		r.rollbackFns[i]()
	}
}
