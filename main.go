package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"voidwalker/blobsdownloader"
	"voidwalker/chainquery"
	"voidwalker/configs"
	ml "voidwalker/util"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/extras/util"
	"github.com/lbryio/lbry.go/v2/stream"
	"github.com/sirupsen/logrus"
)

var publishAddress string
var channelID string
var cqApi *chainquery.CQApi
var downloadsDir string
var uploadsDir string
var blobsDir string
var viewLock ml.MultipleLock
var publishLock ml.MultipleLock

func main() {
	err := configs.Init("./config.json")
	if err != nil {
		panic(err)
	}
	publishAddress = configs.Configuration.PublishAddress
	channelID = configs.Configuration.ChannelID
	if publishAddress == "" || channelID == "" {
		panic("publish_address or channel_id undefined!")
	}
	initLbrynet()
	cqApi, err = chainquery.Init()
	if err != nil {
		panic(err)
	}
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	uploadsDir = usr.HomeDir + "/Uploads/"
	downloadsDir = usr.HomeDir + "/Downloads/"
	blobsDir = usr.HomeDir + "/.lbrynet/blobfiles/"
	viewLock = ml.NewMultipleLock()
	publishLock = ml.NewMultipleLock()

	cache = &sync.Map{}

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	r.POST("/api/claim/publish", publish)
	r.GET("/view/:id/:claimname", view)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Run(":5000")
}

var daemon *jsonrpc.Client

func view(c *gin.Context) {
	id := c.Param("id")
	claimNameWithExt := c.Param("claimname")
	claimName := strings.TrimSuffix(claimNameWithExt, filepath.Ext(claimNameWithExt))
	if strings.HasSuffix(claimName, ".") {
		logrus.Errorf("claim %s#%s has an extra dot in the end of the name!", claimName, id)
		claimName = strings.TrimSuffix(claimName, ".")
	}
	viewLock.Lock(claimName + id)
	defer viewLock.Unlock(claimName + id)
	channelName := ""
	channelShortID := ""
	var claim *chainquery.Claim
	var err error
	if strings.Contains(id, "@") {
		parts := strings.Split(id, ":")
		channelName = parts[0]
		if len(parts) > 1 {
			channelShortID = parts[1]
		}

		claim, err = cqApi.ResolveClaimByChannel(claimName, channelShortID, channelName)
	} else {
		claim, err = cqApi.ResolveClaim(claimName, id)
	}
	if err != nil {
		if errors.Is(err, chainquery.ClaimNotFoundErr) {
			_ = c.AbortWithError(http.StatusNotFound, err)
		} else {
			logrus.Errorln(errors.FullTrace(err))
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	if !strings.Contains(claim.ContentType, "image/") {
		c.Redirect(301, fmt.Sprintf("https://cdn.lbryplayer.xyz/content/claims/%s/%s/stream", claimName, id))
		return
	}

	inUploads, err := isFileInDir(uploadsDir, claimNameWithExt)
	if err != nil {
		logrus.Errorln(errors.FullTrace(err))
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	inDownloads := false
	if !inUploads {
		inDownloads, err = isFileInDir(downloadsDir, claimNameWithExt)
		if err != nil {
			logrus.Errorln(errors.FullTrace(err))
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	mustDownload := !inUploads && !inDownloads
	if mustDownload {
		err = downloadStream(claim.SdHash, claimNameWithExt)
		if err != nil {
			logrus.Errorln(errors.FullTrace(err))
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	var reader *os.File
	if mustDownload || inDownloads {
		reader, err = os.Open(downloadsDir + claimNameWithExt)
		if err != nil {
			logrus.Errorln(errors.FullTrace(err))
			_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
			return
		}
	} else {
		reader, err = os.Open(uploadsDir + claimNameWithExt)
		if err != nil {
			logrus.Errorln(errors.FullTrace(err))
			_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
			return
		}
	}
	defer reader.Close()
	f, err := reader.Stat()
	if err != nil {
		logrus.Errorln(errors.FullTrace(err))
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
		return
	}
	c.DataFromReader(http.StatusOK, f.Size(), claim.ContentType, reader, nil)
	return

}

func isFileInDir(directory, fileName string) (bool, error) {
	_, err := os.Stat(directory + fileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, errors.Err(err)
	}
	return true, nil
}

func buildStream(sdBlob *stream.SDBlob, fileName string) error {
	tmpName := downloadsDir + fileName + ".tmp"
	finalName := downloadsDir + fileName
	f, err := os.Create(tmpName)
	if err != nil {
		return errors.Err(err)
	}
	w := bufio.NewWriter(f)
	for _, info := range sdBlob.BlobInfos {
		if info.Length == 0 {
			continue
		}
		hash := hex.EncodeToString(info.BlobHash)
		blobToDecrypt, err := ioutil.ReadFile(blobsDir + hash)
		if err != nil {
			return errors.Err(err)
		}
		decryptedBlob, err := stream.DecryptBlob(blobToDecrypt, sdBlob.Key, info.IV)
		if err != nil {
			return errors.Err(err)
		}
		_, err = w.Write(decryptedBlob)
		if err != nil {
			return errors.Err(err)
		}
		err = w.Flush()
		if err != nil {
			return errors.Err(err)
		}
	}
	err = os.Rename(tmpName, finalName)
	if err != nil {
		return errors.Err(err)
	}
	return nil
}

func downloadStream(sdHash string, fileName string) error {
	sdBlob, err := blobsdownloader.DownloadBlob(sdHash, false, blobsDir)
	if err != nil {
		return err
	}
	sdb := &stream.SDBlob{}
	err = sdb.FromBlob(*sdBlob)

	if err != nil {
		return err
	}
	err = blobsdownloader.DownloadStream(sdb, blobsDir)
	if err != nil {
		return err
	}
	return buildStream(sdb, fileName)
}

//var cache map[string]PublishResponse
var cache *sync.Map //map[string]PublishResponse

func publish(c *gin.Context) {
	startTime := time.Now()
	claimName := c.PostForm("name")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
		return
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)

	read, err := io.Copy(buf, file)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
		return
	}

	if read < 1 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err("the uploaded file was empty"))
		return
	}
	mimeType := mimetype.Detect(buf.Bytes())
	mimeString := mimetype.Detect(buf.Bytes()).String()
	if strings.Contains(mimeString, "image/") {
		checkSum := fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
		publishLock.Lock(checkSum)
		defer publishLock.Unlock(checkSum)
		filePath := uploadsDir + checkSum[:16] + mimeType.Extension()
		err = c.SaveUploadedFile(fileHeader, filePath)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
			return
		}
		resolveRsp, err := daemon.Resolve(checkSum[:16])
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
			return
		}
		cachedUncasted, ok := cache.Load(checkSum[:16])
		if ok {
			cached, _ := cachedUncasted.(PublishResponse)
			logrus.Infof("returning cached resource: %s", cached.Data.ServeURL)
			c.JSON(http.StatusOK, cached)
			return
		}
		for _, claim := range *resolveRsp {
			if claim.SigningChannel != nil && claim.SigningChannel.ClaimID == channelID {
				baseUrl := "https://spee.ch/" + claim.ClaimID[0:1] + "/" + checkSum[:16]
				extendedUrl := baseUrl + mimeType.Extension()
				response := PublishResponse{
					Success: true,
					Message: "publish completed successfully",
					Data: &PublishData{
						Name:     checkSum[:16],
						ClaimID:  claim.ClaimID,
						URL:      baseUrl,
						ShowURL:  baseUrl,
						ServeURL: extendedUrl,
						PushTo:   claim.ClaimID[0:1] + "/" + checkSum[:16],
						ClaimData: ClaimData{
							Name:          checkSum[:16],
							ClaimID:       claim.ClaimID,
							Title:         checkSum[:16],
							Description:   "",
							Address:       claim.Address,
							Outpoint:      claim.Txid + ":" + fmt.Sprintf("%d", claim.Nout),
							Height:        claim.Height,
							ContentType:   claim.Type,
							Amount:        claim.Amount,
							CertificateID: nil,
							ChannelName:   nil,
						},
					},
				}
				logrus.Infof("returning resolved resource: %s", extendedUrl)
				if strings.Contains(extendedUrl, "..") {
					logrus.Errorf("something is wrong with this crap: %+v", claim)
				}
				c.JSON(http.StatusOK, response)
				cache.Store(checkSum[:16], response)
				return
			}
		}
		tx, err := daemon.StreamCreate(checkSum[:16], filePath, 0.001, jsonrpc.StreamCreateOptions{
			ClaimCreateOptions: jsonrpc.ClaimCreateOptions{
				Title:        util.PtrToString(claimName),
				ClaimAddress: util.PtrToString(publishAddress),
			},
			Author:    util.PtrToString("voidwalker thumbnails"),
			ChannelID: util.PtrToString(channelID),
		})
		if err != nil {
			logrus.Errorf("failed publishing thumbnail: %s", errors.FullTrace(err))
			_ = c.AbortWithError(http.StatusInternalServerError, errors.Err(err))
			return
		}
		baseURL := "https://spee.ch/" + tx.Outputs[0].ClaimID[0:1] + "/" + checkSum[:16]
		extendedUrl := baseURL + mimeType.Extension()
		logrus.Infof("published thumbnail: %s in %s", extendedUrl, time.Since(startTime).String())
		response := PublishResponse{
			Success: true,
			Message: "publish completed successfully",
			Data: &PublishData{
				Name:     checkSum[:16],
				ClaimID:  tx.Outputs[0].ClaimID,
				URL:      baseURL,
				ShowURL:  baseURL,
				ServeURL: extendedUrl,
				PushTo:   tx.Outputs[0].ClaimID[0:1] + "/" + checkSum[:16],
				ClaimData: ClaimData{
					Name:          checkSum[:16],
					ClaimID:       tx.Outputs[0].ClaimID,
					Title:         checkSum[:16],
					Description:   "",
					Address:       tx.Outputs[0].Address,
					Outpoint:      tx.Outputs[0].Txid + ":" + fmt.Sprintf("%d", tx.Outputs[0].Nout),
					Height:        tx.Outputs[0].Height,
					ContentType:   tx.Outputs[0].Type,
					Amount:        tx.Outputs[0].Amount,
					CertificateID: nil,
					ChannelName:   nil,
				},
			},
		}
		if strings.Contains(extendedUrl, "..") {
			logrus.Errorf("something is wrong with this crap: %+v", tx.Outputs[0])
		}
		cache.Store(checkSum[:16], response)
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusBadRequest, PublishResponse{
		Success: false,
		Message: fmt.Sprintf("the provided content (%s) is not supported anymore. Only images allowed. use https://lbry.tv instead", mimeString),
		Data:    nil,
	})
}

type PublishResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *PublishData `json:"data"`
}

type PublishData struct {
	Name      string    `json:"name"`
	ClaimID   string    `json:"claimId"`
	URL       string    `json:"url"`
	ShowURL   string    `json:"showUrl"`
	ServeURL  string    `json:"serveUrl"`
	PushTo    string    `json:"pushTo"`
	ClaimData ClaimData `json:"claimData"`
}
type ClaimData struct {
	Name          string  `json:"name"`
	ClaimID       string  `json:"claimId"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Address       string  `json:"address"`
	Outpoint      string  `json:"outpoint"`
	Height        int     `json:"height"`
	ContentType   string  `json:"contentType"`
	Amount        string  `json:"amount"`
	CertificateID *string `json:"certificateId"`
	ChannelName   *string `json:"channelName"`
}

func initLbrynet() {
	daemon = jsonrpc.NewClient("")
	daemon.SetRPCTimeout(1 * time.Minute)
}
