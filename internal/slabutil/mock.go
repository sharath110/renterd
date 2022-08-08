package slabutil

import (
	"context"
	"errors"
	"io"

	"go.sia.tech/renterd/internal/consensus"
	rhpv2 "go.sia.tech/renterd/rhp/v2"
	"go.sia.tech/renterd/slab"
)

type MockHost struct {
	sectors map[consensus.Hash256][]byte
}

func (h *MockHost) UploadSector(ctx context.Context, sector *[rhpv2.SectorSize]byte) (consensus.Hash256, error) {
	root := rhpv2.SectorRoot(sector)
	h.sectors[root] = append([]byte(nil), sector[:]...)
	return root, nil
}

func (h *MockHost) DownloadSector(ctx context.Context, w io.Writer, root consensus.Hash256, offset, length uint32) error {
	sector, ok := h.sectors[root]
	if !ok {
		return errors.New("unknown root")
	}
	_, err := w.Write(sector[offset:][:length])
	return err
}

func NewMockHost() *MockHost {
	return &MockHost{
		sectors: make(map[consensus.Hash256][]byte),
	}
}

type MockHostSet struct {
	Hosts map[consensus.PublicKey]*MockHost
}

func (hs *MockHostSet) AddHost() consensus.PublicKey {
	hostKey := consensus.GeneratePrivateKey().PublicKey()
	hs.Hosts[hostKey] = NewMockHost()
	return hostKey
}

func (hs *MockHostSet) SlabUploader() slab.SlabUploader {
	hosts := make(map[consensus.PublicKey]slab.SectorUploader)
	for hostKey, sess := range hs.Hosts {
		hosts[hostKey] = sess
	}
	return &slab.SerialSlabUploader{
		Hosts: hosts,
	}
}

func (hs *MockHostSet) SlabDownloader() slab.SlabDownloader {
	hosts := make(map[consensus.PublicKey]slab.SectorDownloader)
	for hostKey, sess := range hs.Hosts {
		hosts[hostKey] = sess
	}
	return &slab.SerialSlabDownloader{
		Hosts: hosts,
	}
}

func NewMockHostSet() *MockHostSet {
	return &MockHostSet{
		Hosts: make(map[consensus.PublicKey]*MockHost),
	}
}
