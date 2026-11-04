package utils

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch             int64 = 1577808000000 // Custom epoch in Unix milliseconds; the 41-bit timestamp spans about 69 years.
	timestampBits     int64 = 41            // Bits allocated to the timestamp
	datacenteridBits  int64 = 5             // Bits allocated to the data center ID
	workeridBits      int64 = 5             // Bits allocated to the worker ID
	sequenceBits      int64 = 12            // Bits allocated to the sequence number
	timestampMax      int64 = 2199023255551 // Maximum timestamp offset: (1 << timestampBits) - 1
	datacenteridMax   int64 = 31            // Maximum data center ID: (1 << datacenteridBits) - 1
	workeridMax       int64 = 31            // Maximum worker ID: (1 << workeridBits) - 1
	sequenceMask      int64 = 4095          // Sequence mask: (1 << sequenceBits) - 1
	workeridShift     int64 = 12            // Worker ID bit offset
	datacenteridShift int64 = 17            // Data center ID bit offset
	timestampShift    int64 = 22            // Timestamp bit offset
)

type Snowflake struct {
	sync.Mutex
	timestamp    int64
	workerid     int64
	datacenterid int64
	sequence     int64
}

func NewSnowflake(datacenterid, workerid int64) (*Snowflake, error) {
	if datacenterid < 0 || datacenterid > datacenteridMax {
		return nil, fmt.Errorf("datacenterid must be between 0 and %d", datacenteridMax-1)
	}
	if workerid < 0 || workerid > workeridMax {
		return nil, fmt.Errorf("workerid must be between 0 and %d", workeridMax-1)
	}
	return &Snowflake{
		timestamp:    0,
		datacenterid: datacenterid,
		workerid:     workerid,
		sequence:     0,
	}, nil
}

// NextVal combines the timestamp, data center ID, worker ID, and sequence number.
func (s *Snowflake) NextVal() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000 // Convert to milliseconds.
	if s.timestamp == now {
		// Increment the sequence for IDs generated in the same millisecond.
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// Wait for the next millisecond when the 12-bit sequence wraps around.
			// The next millisecond starts with sequence zero.
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// Reset the sequence when the millisecond changes.
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		return 0
	}
	s.timestamp = now
	r := (t)<<timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence)
	s.Unlock()
	return r
}

// GetDeviceID extracts the data center and worker IDs.
func GetDeviceID(sid int64) (datacenterid, workerid int64) {
	datacenterid = (sid >> datacenteridShift) & datacenteridMax
	workerid = (sid >> workeridShift) & workeridMax
	return
}

// GetTimestamp extracts the timestamp offset in milliseconds.
func GetTimestamp(sid int64) (timestamp int64) {
	timestamp = (sid >> timestampShift) & timestampMax
	return
}

// GetGenTimestamp returns the ID creation time in Unix milliseconds.
func GetGenTimestamp(sid int64) (timestamp int64) {
	timestamp = GetTimestamp(sid) + epoch
	return
}

// GetGenTime formats the ID creation time with second precision.
func GetGenTime(sid int64) (t string) {
	// Convert the millisecond timestamp to seconds for time.Unix.
	t = time.Unix(GetGenTimestamp(sid)/1000, 0).Format("2006-01-02 15:04:05")
	return
}

// GetTimestampStatus returns the fraction of the timestamp range elapsed since the epoch.
func GetTimestampStatus() (state float64) {
	state = float64((time.Now().UnixNano()/1000000 - epoch)) / float64(timestampMax)
	return
}
