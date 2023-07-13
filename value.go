package confo

// Value represents a configuration item interface.
// If a data type implements this interface, Confo will use it to assign the content from file(env, or flag) to the fields.
type Value interface {
	String() string
	Set(string) error
}

type emptyValue struct{}

func (a emptyValue) String() string {
	return "empty value"
}

func (a *emptyValue) Set(string) error {
	return nil
}
