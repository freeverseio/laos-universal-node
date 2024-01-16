package universal

import (
	"reflect"
	"testing"
)

func TestGetNextLowerBlockNumber(t *testing.T) {
	testCases := []struct {
		name                  string
		currentBlock          uint64
		storedBlockNumbers    []uint64
		expectedBlockNumber   uint64
		expectedModifiedSlice []uint64
		expectedFound         bool
	}{
		{
			name:                  "FindLowerBlock",
			currentBlock:          5,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 7},
			expectedBlockNumber:   4,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 7},
			expectedFound:         true,
		},
		{
			name:                  "FindLowerBlock",
			currentBlock:          7,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 7},
			expectedBlockNumber:   6,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 7},
			expectedFound:         true,
		},
		{
			name:                  "FindLowerBlock",
			currentBlock:          9,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 8},
			expectedBlockNumber:   8,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 8},
			expectedFound:         true,
		},
		{
			name:                  "BlockNotFound",
			currentBlock:          3,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 8},
			expectedBlockNumber:   0,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 8},
			expectedFound:         false,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			gotBlockNumber, found := getNextLowerBlockNumber(tc.currentBlock, tc.storedBlockNumbers)

			if tc.expectedFound != found {
				t.Errorf("got %v, expected %v", found, tc.expectedFound)
			}
			if gotBlockNumber != tc.expectedBlockNumber {
				t.Errorf("got %v, expected %v", gotBlockNumber, tc.expectedBlockNumber)
			}
			if !reflect.DeepEqual(tc.storedBlockNumbers, tc.expectedModifiedSlice) {
				t.Errorf("slice was modified to %v, expected %v", tc.storedBlockNumbers, tc.expectedModifiedSlice)
			}
		})
	}
}
