package event_test

import (
	"testing"

	"github.com/freeverseio/laos-universal-node/internal/platform/event"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	if event.ERC721TransferEventSigHash != "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
		t.Errorf("ERC721TransferEventSigHash is not correct, got: %s, want: %s.", event.ERC721TransferEventSigHash, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	}
}
