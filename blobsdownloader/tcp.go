package blobsdownloader

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"time"

	"voidwalker/configs"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/stream"
	"github.com/lbryio/reflector.go/peer"
	"github.com/lbryio/reflector.go/store"
)

func DownloadBlob(hash string, save bool, blobsDir string) (*stream.Blob, error) {
	bStore := GetBlobStore()
	blob, err := bStore.Get(hash)
	if err != nil {
		err = errors.Prefix(hash, err)
		return nil, errors.Err(err)
	}
	if save {
		err = os.MkdirAll(blobsDir, os.ModePerm)
		if err != nil {
			return nil, errors.Err(err)
		}
		err = ioutil.WriteFile(blobsDir+hash, blob, 0644)
		if err != nil {
			return nil, errors.Err(err)
		}
	}
	return &blob, nil
}

// GetBlobStore returns default pre-configured blob store.
func GetBlobStore() store.BlobStore {
	return peer.NewStore(peer.StoreOpts{
		Address: configs.Configuration.ReflectorServer,
		Timeout: 30 * time.Second,
	})
}

// downloads a stream and returns the speed in bytes per second
func DownloadStream(blob *stream.SDBlob, blobsDir string) error {
	hashes := GetStreamHashes(blob)
	for _, hash := range hashes {
		_, err := os.Stat(blobsDir + hash)
		if os.IsNotExist(err) {
			_, err := DownloadBlob(hash, true, blobsDir)
			if err != nil {
				return err
			}
		} else if err != nil {
			return errors.Err(err)
		}
	}
	return nil
}

func GetStreamHashes(blob *stream.SDBlob) []string {
	blobs := make([]string, 0, len(blob.BlobInfos))
	for _, b := range blob.BlobInfos {
		hash := hex.EncodeToString(b.BlobHash)
		if hash == "" {
			continue
		}
		blobs = append(blobs, hex.EncodeToString(b.BlobHash))
	}
	return blobs
}
