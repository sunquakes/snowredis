package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	// Epoch Timestamp offset (2022-01-01 00:00:00 UTC)
	Epoch int64 = 1640995200000

	// Number of timestamp bits
	timestampBits = 41
	// Number of datacenter ID bits
	datacenterBits = 5
	// Number of worker ID bits
	workerBits = 5
	// Number of sequence bits
	sequenceBits = 12

	// Maximum datacenter ID and worker ID
	maxDatacenterID = -1 ^ (-1 << datacenterBits) // 31
	maxWorkerID     = -1 ^ (-1 << workerBits)     // 31

	// Maximum sequence number
	maxSequence = -1 ^ (-1 << sequenceBits) // 4095

	// Bit shifts for each part
	workerShift     = sequenceBits
	datacenterShift = sequenceBits + workerBits
	timestampShift  = sequenceBits + workerBits + datacenterBits
)

var (
	ErrInvalidNodeID     = errors.New("invalid node ID")
	ErrInvalidDatacenter = errors.New("invalid datacenter ID")
	ErrInvalidWorker     = errors.New("invalid worker ID")
	ErrOverFlow          = errors.New("sequence number exceeds maximum value")
)

// Node Snowflake algorithm node structure
type Node struct {
	sync.Mutex
	datacenterID  int64 // Datacenter ID for snowflake ID generation
	workerID      int64 // Worker ID for snowflake ID generation
	sequence      int64 // Sequence number for snowflake ID generation
	lastTimestamp int64 // Timestamp of last generated ID
}

// NewNode Creates a new snowflake algorithm node
func NewNode(datacenterID int64, workerID int64) (*Node, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, ErrInvalidDatacenter
	}

	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorker
	}

	return &Node{
		datacenterID:  datacenterID,
		workerID:      workerID,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

// currentTimeMillis Gets current timestamp in milliseconds
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
