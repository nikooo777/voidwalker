package compression

import (
	"io/ioutil"
	"mime"
	"path/filepath"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	giftowebp "github.com/sizeofint/gif-to-webp"
	"go.uber.org/atomic"
)

var inUse atomic.Bool
var AlreadyInUseErr = errors.Base("already busy compressing")
var UnsupportedErr = errors.Base("unsupported compression")

// returns a compressedPath, a mimeType or an error
func Compress(path, fileName, mimeType, storePath string) (string, string, error) {
	swapped := inUse.CAS(false, true)
	if !swapped {
		return "", "", AlreadyInUseErr
	}
	defer inUse.Store(false)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", errors.Err(err)
	}
	switch mimeType {
	case "image/gif":
		converter := giftowebp.NewConverter()
		converter.LoopCompatibility = false
		if len(file) > 500*1024 {
			converter.WebPConfig.SetTargetSize(500 * 1024)
		} else {
			converter.WebPConfig.SetTargetSize(len(file))
		}
		converter.WebPConfig.SetMethod(4)
		webpBin, err := converter.Convert(file)
		if err != nil {
			return "", "", errors.Err(err)
		}
		compressedPath := filepath.Join(storePath, fileName+".webp")
		err = ioutil.WriteFile(compressedPath, webpBin, 0600)
		if err != nil {
			return "", "", errors.Err(err)
		}
		//err = os.Remove(path)
		return fileName + ".webp", mime.TypeByExtension(".webp"), nil
	case mime.TypeByExtension("png"):
	case mime.TypeByExtension("jpeg"):
	case mime.TypeByExtension("jpg"):

	}
	return "", "", UnsupportedErr
}
