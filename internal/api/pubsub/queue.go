package pubsub

type node struct {
	value []byte
	next  *node
}

type Queue struct {
	head *node
	tail *node
	size int
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(value []byte) {
	newNode := &node{value: value}
	if q.tail != nil {
		q.tail.next = newNode
	} else {
		q.head = newNode
	}
	q.tail = newNode
	q.size++
}

func (q *Queue) Dequeue() ([]byte, bool) {
	if q.head == nil {
		return nil, false
	}
	value := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return value, true
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}
