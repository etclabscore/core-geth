package rawdb

import (
	"context"
)

type ExternalFreezerRemoteAPI interface {
	HasAncient(ctx context.Context, kind string, number uint64) (bool, error)
	Ancient(ctx context.Context, kind string, number uint64) ([]byte, error)
	Ancients(ctx context.Context) (uint64, error)
	AncientSize(ctx context.Context, kind string) (uint64, error)

	AppendAncient(ctx context.Context, number uint64, hash, header, body, receipt, td []byte)
	TruncateAncients(ctx context.Context, n uint64) error
	Sync(ctx context.Context) error
}

// FreezerRemoteAPI exposes a JSONRPC related API
type FreezerRemoteAPI struct {
	freezerRemote *freezerRemote
}

// NewFreezerRemoteAPI exposes an endpoint to create a remote service
func NewFreezerRemoteAPI(freezerStr string, namespace string) (*FreezerRemoteAPI, error) {
	freezer, err := newFreezerRemote(freezerStr, namespace)
	if err != nil {
		return nil, err
	}

	freezerAPI := FreezerRemoteAPI{
		freezerRemote: freezer,
	}
	return &freezerAPI, nil
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (freezerRemoteAPI *FreezerRemoteAPI) HasAncient(ctx context.Context, kind string, number uint64) (bool, error) {
	return freezerRemoteAPI.freezerRemote.HasAncient(kind, number)
}
