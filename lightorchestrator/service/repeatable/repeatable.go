package repeatable

import (
	"hash/fnv"
	"math"
	"time"

	"github.com/jmbarzee/dominion/system"
)

// Option is an idempotent call based on a time and a number of options
// If the result of a fast hash of the time is within a certain
// option-range of possible hash values, it will return that option
func Option(t time.Time, options int) int {
	sum := float32(hashSum(t))
	max := float32(math.MaxUint32) // maximum possible hash
	portion := max / float32(options)

	return int(sum / portion)
}

// Chance is an idempotent call based on a time and a chance in the range (1, 0)
// it will return true if the result of a fast hash of the time is above
// a certain percentage (chance) of possible hash values
func Chance(t time.Time, chance float32) bool {
	sum := float32(hashSum(t))
	max := float32(math.MaxUint32) // maximum possible hash

	return sum > max*chance
}

// RandDuration is an idempotent call based on a time and a range of durations
// it will return a random duration from within that range
func RandDuration(t time.Time, min, max time.Duration) time.Duration {
	sum := time.Duration(hashSum(t))
	diff := max - min
	return sum%(diff) + min
}

// RandShift is an idempotent call based on a time and a range of shifts
// it will return a random shift from within that range
// unit provides the granularity of the randomness
func RandShift(t time.Time, min, max, unit float32) float32 {
	diff := max - min
	units := int(diff / unit)
	option := float32(Option(t, units))
	return option*unit + min
}

func hashSum(t time.Time) uint32 {
	hash := fnv.New32()
	timeBinary, err := t.MarshalBinary()
	if err != nil {
		system.Panic(err)
	}

	_, err = hash.Write(timeBinary)
	if err != nil {
		system.Panic(err)
	}

	return hash.Sum32()
}
