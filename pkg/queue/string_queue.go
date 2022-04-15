package queue

type StringQueue []string

func (q *StringQueue) Add(n string) {
	*q = append(*q, n)
}

func (q *StringQueue) Poll() string {
	p := (*q)[0]
	*q = (*q)[1:]

	return p
}

func (q StringQueue) Size() int {
	return len(q)
}

func (q StringQueue) IsEmpty() bool {
	return q.Size() == 0
}
