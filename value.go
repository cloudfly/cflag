package confo

// Value represents a configuration item interface.
// If a data type implements this interface, Confo will use it to assign the content from file(env, or flag) to the fields.
type Value interface {
	Set(string) error
}
