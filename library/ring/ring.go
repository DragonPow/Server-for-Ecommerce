package ring

import (
	"errors"
	"sync"
)

var (
	ErrTooManyDataToWrite = errors.New("too many data to write")
	ErrIsFull             = errors.New("ringbuffer is full")
	ErrIsEmpty            = errors.New("ringbuffer is empty")
	ErrAccuqireLock       = errors.New("no lock to accquire")
)

// RingBuffer is a circular buffer that implement io.ReaderWriter interface.
type RingBuffer[T any] struct {
	buf         []T
	size        int
	r           int // next position to read
	w           int // next position to write
	isFull      bool
	SignalFull  chan struct{}
	SignalWrite chan struct{}
	SignalRead  chan struct{}
	mu          sync.Mutex
}

// New returns a new RingBuffer whose buffer has the given size.
func New[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buf:         make([]T, size),
		size:        size,
		SignalFull:  make(chan struct{}),
		SignalWrite: make(chan struct{}),
		SignalRead:  make(chan struct{}),
	}
}

func (r *RingBuffer[T]) setIsFull(isFull bool) {
	r.isFull = isFull
	if isFull {
		r.SignalFull <- struct{}{}
	}
}

func (r *RingBuffer[T]) setReaderPos(v int) {
	r.r = v
	//r.SignalRead <- struct{}{}
}

func (r *RingBuffer[T]) setWriterPos(v int) {
	r.w = v
	r.SignalWrite <- struct{}{}
}

// Read reads up to len(p) into p. It returns the number of read (0 <= n <= len(p)) and any error encountered. Even if Read returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, Read conventionally returns what is available instead of waiting for more.
// When Read encounters an error or end-of-file condition after successfully reading n > 0 bytes, it returns the number of read. It may return the (non-nil) error from the same call or return the error (and n == 0) from a subsequent call.
// Callers should always process the n > 0 returned before considering the error err. Doing so correctly handles I/O errors that happen after reading some and also both of the allowed EOF behaviors.
func (r *RingBuffer[T]) Read(p []T) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	r.mu.Lock()
	n, err = r.read(p)
	r.mu.Unlock()
	return n, err
}

// TryRead read up to len(p) into p like Read but it is not blocking.
// If it has not succeeded to accquire the lock, it return 0 as n and ErrAccuqireLock.
func (r *RingBuffer[T]) TryRead(p []T) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	ok := r.mu.TryLock()
	if !ok {
		return 0, ErrAccuqireLock
	}

	n, err = r.read(p)
	r.mu.Unlock()
	return n, err
}

func (r *RingBuffer[T]) read(p []T) (n int, err error) {
	if r.w == r.r && !r.isFull {
		return 0, ErrIsEmpty
	}

	if r.w > r.r {
		n = r.w - r.r
		if n > len(p) {
			n = len(p)
		}
		copy(p, r.buf[r.r:r.r+n])
		newR := (r.r + n) % r.size
		r.setReaderPos(newR)
		return
	}

	n = r.size - r.r + r.w
	if n > len(p) {
		n = len(p)
	}

	if r.r+n <= r.size {
		copy(p, r.buf[r.r:r.r+n])
	} else {
		c1 := r.size - r.r
		copy(p, r.buf[r.r:r.size])
		c2 := n - c1
		copy(p[c1:], r.buf[0:c2])
	}
	newR := (r.r + n) % r.size
	r.setReaderPos(newR)

	r.setIsFull(false)

	return n, err
}

// Write writes len(p) from p to the underlying buf.
// It returns the number of written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early.
// Write returns a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
func (r *RingBuffer[T]) Write(p []T) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.mu.Lock()
	n, err = r.write(p)
	r.mu.Unlock()

	return n, err
}

// TryWrite writes len(p) from p to the underlying buf like Write, but it is not blocking.
// If it has not succeeded to accquire the lock, it return 0 as n and ErrAccuqireLock.
func (r *RingBuffer[T]) TryWrite(p []T) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	ok := r.mu.TryLock()
	if !ok {
		return 0, ErrAccuqireLock
	}

	n, err = r.write(p)
	r.mu.Unlock()

	return n, err
}

func (r *RingBuffer[T]) write(p []T) (n int, err error) {
	if r.isFull {
		return 0, ErrIsFull
	}

	var avail int
	if r.w >= r.r {
		avail = r.size - r.w + r.r
	} else {
		avail = r.r - r.w
	}

	if len(p) > avail {
		err = ErrTooManyDataToWrite
		p = p[:avail]
	}
	n = len(p)

	if r.w >= r.r {
		c1 := r.size - r.w
		if c1 >= n {
			copy(r.buf[r.w:], p)
			r.w += n
		} else {
			copy(r.buf[r.w:], p[:c1])
			c2 := n - c1
			copy(r.buf[0:], p[c1:])
			r.setWriterPos(c2)
		}
	} else {
		copy(r.buf[r.w:], p)
		r.w += n
	}

	if r.w == r.size {
		r.setWriterPos(0)
	}
	if r.w == r.r {
		r.setIsFull(true)
	}

	return n, err
}

// WriteOne writes one byte into buffer, and returns ErrIsFull if buffer is full.
func (r *RingBuffer[T]) WriteOne(c T) error {
	r.mu.Lock()
	err := r.writeOne(c)
	r.mu.Unlock()
	return err
}

// TryWriteOne writes one into buffer without blocking.
// If it has not succeeded to accquire the lock, it return ErrAccuqireLock.
func (r *RingBuffer[T]) TryWriteOne(c T) error {
	ok := r.mu.TryLock()
	if !ok {
		return ErrAccuqireLock
	}

	err := r.writeOne(c)
	r.mu.Unlock()
	return err
}

func (r *RingBuffer[T]) writeOne(c T) (err error) {
	if r.w == r.r && r.isFull {
		return ErrIsFull
	}
	r.buf[r.w] = c
	r.setWriterPos(r.w + 1)

	if r.w == r.size {
		r.w = 0
	}
	if r.w == r.r {
		r.setIsFull(true)
	}

	return nil
}

// Length return the length of available reader.
func (r *RingBuffer[T]) Length() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.w == r.r {
		if r.isFull {
			return r.size
		}
		return 0
	}

	if r.w > r.r {
		return r.w - r.r
	}

	return r.size - r.r + r.w
}

// Capacity returns the size of the underlying buffer.
func (r *RingBuffer[T]) Capacity() int {
	return r.size
}

// Free returns the length of available position to write.
func (r *RingBuffer[T]) Free() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.w == r.r {
		if r.isFull {
			return 0
		}
		return r.size
	}

	if r.w < r.r {
		return r.r - r.w
	}

	return r.size - r.w + r.r
}

// List returns all available reader. It does not move the read pointer and only copy the available data.
func (r *RingBuffer[T]) List() []T {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.w == r.r {
		if r.isFull {
			buf := make([]T, r.size)
			copy(buf, r.buf[r.r:])
			copy(buf[r.size-r.r:], r.buf[:r.w])
			return buf
		}
		return nil
	}

	if r.w > r.r {
		buf := make([]T, r.w-r.r)
		copy(buf, r.buf[r.r:r.w])
		return buf
	}

	n := r.size - r.r + r.w
	buf := make([]T, n)

	if r.r+n < r.size {
		copy(buf, r.buf[r.r:r.r+n])
	} else {
		c1 := r.size - r.r
		copy(buf, r.buf[r.r:r.size])
		c2 := n - c1
		copy(buf[c1:], r.buf[0:c2])
	}

	return buf
}

// IsFull returns this ringbuffer is full.
func (r *RingBuffer[T]) IsFull() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.isFull
}

// IsEmpty returns this ringbuffer is empty.
func (r *RingBuffer[T]) IsEmpty() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return !r.isFull && r.w == r.r
}

// Reset the read pointer and writer pointer to zero.
func (r *RingBuffer[T]) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.setReaderPos(0)
	r.setWriterPos(0)
	r.setIsFull(false)
}
