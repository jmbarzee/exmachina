package repeatable

import (
	"fmt"
	"testing"
	"time"

	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestOption(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	maxOptions := 1000
	totalRunsPerOption := 100

	for options := 2; options < maxOptions; options++ {
		name := fmt.Sprintf("%v Options", options)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			buckets := make([]int, options)
			totalRuns := options * totalRunsPerOption
			for i := 0; i < totalRuns; i++ {
				option := Option(aTime.Add(time.Duration(i)*time.Second), options)
				buckets[option]++
			}

			expectedPercentage := float64(1.0) / float64(options)
			for j, bucket := range buckets {
				actualPercentage := float64(bucket) / float64(totalRuns)
				if !helper.FloatsEqual(expectedPercentage, actualPercentage, 0.001) {
					t.Fatalf("bucket %v failed:\n\tExpected: %v\n\tActual: %v", j, expectedPercentage, actualPercentage)
				}
			}
		})
	}
}

func TestChance(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	maxChance := 1.0
	quantaChance := 0.01
	totalRuns := 10000

	for chance := 0.0; chance <= maxChance; chance = chance + quantaChance {
		name := fmt.Sprintf("%5.2g Chance", chance)
		t.Run(name, func(t *testing.T) {
			// t.Parallel()
			choiceCount := 0
			for i := 0; i < totalRuns; i++ {
				choice := Chance(aTime.Add(time.Duration(i)*time.Second), chance)
				if choice {
					choiceCount++
				}
			}
			actualChance := float64(choiceCount) / float64(totalRuns)
			if !helper.FloatsEqual(chance, actualChance, 0.01) {
				t.Fatalf("Percentages were not close enough:\n\tExpected: %5.2g\n\tActual: %5.2g", chance, actualChance)
			}
		})
	}
}

func TestRandDuration(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	runsPerPossibleDurations := 100
	maxPossibleDuration := time.Minute
	testQuanta := time.Second * 5

	for max := time.Second; max < maxPossibleDuration; max += testQuanta {
		for min := time.Duration(0); min < max; min += testQuanta {
			name := fmt.Sprintf("Range %v-%v", min, max)
			t.Run(name, func(t *testing.T) {
				max := max
				min := min
				t.Parallel()
				bucketSize := time.Millisecond
				possibleDurations := (max - min) / bucketSize
				buckets := make([]int, possibleDurations)
				totalRuns := int(possibleDurations) * runsPerPossibleDurations
				for i := 0; i < totalRuns; i++ {
					duration := RandDuration(aTime.Add(time.Duration(i)*time.Second), min, max)
					bucket := (duration - min) / bucketSize
					buckets[bucket]++
				}

				expectedPercentage := float64(1.0 / possibleDurations)
				for j, bucket := range buckets {
					actualPercentage := float64(bucket) / float64(totalRuns)
					if !helper.FloatsEqual(expectedPercentage, actualPercentage, 0.002) {
						t.Fatalf("bucket %v failed:\n\tExpected: %v\n\tActual: %v", j, expectedPercentage, actualPercentage)
					}
				}
			})
		}
	}
}

// TODO: Fix this test -> make loops with use ints instead of floats
// func TestRandShift(t *testing.T) {
// 	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
// 	runsPerPossibleShifts := 10000
// 	maxPossibleShift := 2.0
// 	testQuanta := 0.1

// 	for max := 0.0001; max < maxPossibleShift; max += testQuanta {
// 		for min := 0.0; min < max; min += testQuanta {
// 			name := fmt.Sprintf("Range %5.2g-%5.2g", min, max)
// 			t.Run(name, func(t *testing.T) {
// 				max := max
// 				min := min
// 				// t.Parallel()
// 				bucketSize := testQuanta
// 				numberOfBuckets := int((max - min) / bucketSize)
// 				buckets := make([]int, numberOfBuckets)
// 				totalRuns := int(numberOfBuckets) * runsPerPossibleShifts
// 				// fmt.Printf("RandShift(%4.2g, %4.2g) x %.2v times: %.2v buckets with %4.2g size\n", min, max, totalRuns, numberOfBuckets, bucketSize)
// 				for i := 0; i < totalRuns; i++ {
// 					shift := RandShift(aTime.Add(time.Duration(i)*time.Second), min, max, testQuanta)
// 					bucket := int((shift - min) / bucketSize)
// 					// fmt.Printf("\t %4.2g -> bucket:%v\n", shift, bucket)
// 					buckets[bucket]++
// 				}

// 				expectedPercentage := 2.0 / float64(numberOfBuckets)
// 				for j, bucket := range buckets {
// 					actualPercentage := float64(bucket) / float64(totalRuns)
// 					if !helper.FloatsEqual(expectedPercentage, actualPercentage, 0.01) {
// 						t.Fatalf("bucket %v failed:\n\tExpected: %v\n\tActual: %v", j, expectedPercentage, actualPercentage)
// 					}
// 				}
// 			})
// 		}
// 	}
// }
