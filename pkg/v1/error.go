package nursery

var _ error = basic("")

type basic string

func (b basic) Error() string {
	return string(b)
}
