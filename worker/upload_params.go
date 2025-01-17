package worker

import (
	"go.sia.tech/renterd/api"
	"go.sia.tech/renterd/build"
	"go.sia.tech/renterd/object"
)

type uploadParameters struct {
	bucket string
	path   string

	multipart  bool
	uploadID   string
	partNumber int

	ec               object.EncryptionKey
	encryptionOffset uint64

	rs          api.RedundancySettings
	bh          uint64
	contractSet string
	packing     bool
	mimeType    string
}

func defaultParameters(bucket, path string) uploadParameters {
	return uploadParameters{
		bucket: bucket,
		path:   path,

		ec:               object.GenerateEncryptionKey(), // random key
		encryptionOffset: 0,                              // from the beginning

		rs: build.DefaultRedundancySettings,
	}
}

func multipartParameters(bucket, path, uploadID string, partNumber int) uploadParameters {
	return uploadParameters{
		bucket: bucket,
		path:   path,

		multipart:  true,
		uploadID:   uploadID,
		partNumber: partNumber,

		ec:               object.GenerateEncryptionKey(), // random key
		encryptionOffset: 0,                              // from the beginning

		rs: build.DefaultRedundancySettings,
	}
}

type UploadOption func(*uploadParameters)

func WithBlockHeight(bh uint64) UploadOption {
	return func(up *uploadParameters) {
		up.bh = bh
	}
}

func WithContractSet(contractSet string) UploadOption {
	return func(up *uploadParameters) {
		up.contractSet = contractSet
	}
}

func WithCustomKey(ec object.EncryptionKey) UploadOption {
	return func(up *uploadParameters) {
		up.ec = ec
	}
}

func WithCustomEncryptionOffset(offset uint64) UploadOption {
	return func(up *uploadParameters) {
		up.encryptionOffset = offset
	}
}

func WithMimeType(mimeType string) UploadOption {
	return func(up *uploadParameters) {
		up.mimeType = mimeType
	}
}

func WithPacking(packing bool) UploadOption {
	return func(up *uploadParameters) {
		up.packing = packing
	}
}

func WithRedundancySettings(rs api.RedundancySettings) UploadOption {
	return func(up *uploadParameters) {
		up.rs = rs
	}
}
