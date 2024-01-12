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
		expectedError         bool
		expectedModifiedSlice []uint64
	}{
		{
			name:                  "FindLowerBlock",
			currentBlock:          5,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 7},
			expectedBlockNumber:   4,
			expectedError:         false,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 7},
		},
		{
			name:                  "FindLowerBlock",
			currentBlock:          7,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 7},
			expectedBlockNumber:   6,
			expectedError:         false,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 7},
		},
		{
			name:                  "FindLowerBlock",
			currentBlock:          9,
			storedBlockNumbers:    []uint64{3, 4, 5, 6, 8},
			expectedBlockNumber:   8,
			expectedError:         false,
			expectedModifiedSlice: []uint64{3, 4, 5, 6, 8},
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			gotBlockNumber, err := getNextLowerBlockNumber(tc.currentBlock, tc.storedBlockNumbers)

			if tc.expectedError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if gotBlockNumber != tc.expectedBlockNumber {
					t.Errorf("got %v, expected %v", gotBlockNumber, tc.expectedBlockNumber)
				}
				if !reflect.DeepEqual(tc.storedBlockNumbers, tc.expectedModifiedSlice) {
					t.Errorf("slice was modified to %v, expected %v", tc.storedBlockNumbers, tc.expectedModifiedSlice)
				}
			}
		})
	}
}
