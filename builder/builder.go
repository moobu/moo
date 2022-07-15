package builder

type Builder interface {
	// Build turns the source into a bundle which's then
	// delivered to the runtime that runs the bundle in
	// one or several pods. If Local is unset, we should
	// ask the retriever for the source using Remote.
	Build(*Source, ...BuildOption) (*Bundle, error)
	// Clean cleans up a bundle and the dependencies created
	// during its build.
	Clean(*Bundle, ...CleanOption) error
	// String returns the builder's name
	String() string
}

type Source struct {
	// Name of the source
	Name string
	// Type specifies which builder to use
	Type string
	// URL address of the source
	URL string
}

type Bundle struct {
	// Dir is the path of the bundle
	Dir string
	// Reference on which the bundle was built
	Ref string
	// Entry command and arguments
	Entry []string
	// Source built
	Source *Source
}

var Default Builder

func Build(s *Source, opts ...BuildOption) (*Bundle, error) {
	return Default.Build(s, opts...)
}

func Clean(b *Bundle, opts ...CleanOption) error {
	return Default.Clean(b, opts...)
}
