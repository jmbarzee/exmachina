package shared

import (
	"hash/fnv"
	"time"
)

// RepeatableOption is an idempotent call based on a time and a number of options
// If the result of a fast hash of the time is within a certain
// option-range of possible hash values, it will return that option
func RepeatableOption(t time.Time, options int) int {
	hash := fnv.New32()
	timeBinary, err := t.MarshalBinary()
	if err != nil {
		panic(err)
	}

	_, err = hash.Write(timeBinary)
	if err != nil {
		panic(err)
	}

	sum := float32(hash.Sum32())
	max := float32(^uint(0)) // maximum possible hash
	portion := max / float32(options)

	return int(sum / portion)
}

// RepeatableChance is an idempotent call based on a time and a chance in the range (1, 0)
// it will return true if the result of a fast hash of the time is above
// a certain percentage (chance) of possible hash values
func RepeatableChance(t time.Time, chance float32) bool {
	hash := fnv.New32()
	timeBinary, err := t.MarshalBinary()
	if err != nil {
		panic(err)
	}

	_, err = hash.Write(timeBinary)
	if err != nil {
		panic(err)
	}

	sum := float32(hash.Sum32())

	max := float32(^uint(0))

	return sum > max*chance

}
