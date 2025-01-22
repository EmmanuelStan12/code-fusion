package client

import (
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"log"
	"sync"
	"time"
)

type Debouncer struct {
	timer    *time.Timer
	lock     sync.Mutex
	delay    time.Duration
	callback func(code string)
}

func NewDebouncer(delay time.Duration, callback func(code string)) *Debouncer {
	return &Debouncer{
		delay:    delay,
		callback: callback,
	}
}

func (d *Debouncer) Debounce(code string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.delay, func() {
		go d.callback(code)
	})
}

type CodeSessionQueue struct {
	operations []dto.CodeSessionOperation
	signal     chan struct{}
	lock       sync.RWMutex
	codeLock   sync.Mutex
	current    int
	sessionId  model.SessionId
	code       string
	debounce   *Debouncer
}

func NewCodeSessionQueue(sessionId model.SessionId, code string) *CodeSessionQueue {
	return &CodeSessionQueue{
		signal:     make(chan struct{}),
		sessionId:  sessionId,
		code:       code,
		operations: make([]dto.CodeSessionOperation, 0),
	}
}

func (queue *CodeSessionQueue) RegisterDebounce(delay time.Duration, callback func(code string)) {
	queue.debounce = NewDebouncer(delay, callback)
}

func (queue *CodeSessionQueue) Enqueue(operation dto.CodeSessionOperation) {
	queue.lock.Lock()
	queue.operations = append(queue.operations, operation)
	queue.lock.Unlock()

	select {
	case queue.signal <- struct{}{}:
	default:
	}
}

func (queue *CodeSessionQueue) EnqueueAll(operations ...dto.CodeSessionOperation) {
	queue.lock.Lock()
	queue.operations = append(queue.operations, operations...)
	queue.lock.Unlock()

	select {
	case queue.signal <- struct{}{}:
	default:
	}
}

func (queue *CodeSessionQueue) ProcessQueueOperations() {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	queue.codeLock.Lock()
	defer queue.codeLock.Unlock()

	if len(queue.operations) == 0 || queue.current >= len(queue.operations) {
		return
	}

	for i := queue.current; i < len(queue.operations); i++ {
		queue.UpdateCode(queue.operations[i], i)
	}
	queue.current = len(queue.operations)

	if queue.debounce != nil {
		queue.debounce.Debounce(queue.code)

		log.Println("Clearing all operations")
		queue.operations = nil
		queue.current = 0

	}
}

func (queue *CodeSessionQueue) CloseQueue() {
	close(queue.signal)
}

func (queue *CodeSessionQueue) Start() {
	go func() {
		for {
			select {
			case _, ok := <-queue.signal:
				if !ok {
					return
				}
				queue.ProcessQueueOperations()
			}
		}
	}()
}
func (queue *CodeSessionQueue) UpdateCode(operation dto.CodeSessionOperation, currentPosition int) {
	newPosition := operation.Position

	switch operation.Type {
	case "insert":
		if newPosition >= 0 && newPosition <= len(queue.code) {
			queue.code = queue.code[:newPosition] + operation.Text + queue.code[newPosition:]
		} else {
			log.Printf("Invalid insert position %d\n", operation.Position)
		}
	case "delete":
		if newPosition >= 0 && newPosition <= len(queue.code) {
			var suffix, prefix string
			if newPosition != len(queue.code)-1 {
				suffix = queue.code[newPosition+1:]
			}
			index := newPosition - operation.Length + 1
			if index > 0 && index < len(queue.code) {
				prefix = queue.code[:newPosition-operation.Length+1] + suffix
			}
			queue.code = prefix + suffix
		} else {
			log.Printf("Invalid delete range: position %d, length %d\n", operation.Position, operation.Length)
		}
	default:
		log.Printf("Cannot process action %+v\n", operation)
	}
}
