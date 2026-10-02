package streams

import "strconv"

type entry struct {
	id     id
	values []kvp
}

type kvp struct {
	key   string
	value string
}

type id struct {
	ms  string
	seq int
}

func (i id) String() string {
	return i.ms + "-" + strconv.Itoa(i.seq)
}
