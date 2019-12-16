package collectors

type Collector interface {
	Collect(ErrorHandler)
	Format() string
}

type ErrorHandler func(error)
