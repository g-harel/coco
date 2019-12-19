package collectors

// Collector is a generic stateful interface to query and
// format data.
type Collector interface {
	// The collect step is executed in parallel for all
	// collectors. Fatal errors should be handled internally
	// (ex. send to handler + format to empty string).
	Collect(ErrorHandler)

	// The formatting step is executed in series for all
	// collectors. It should return an empty string if no
	// data is collected.
	Format() string
}

// ErrorHandler lets the caller handle any amount of
// collector errors asynchronously.
type ErrorHandler func(error)
