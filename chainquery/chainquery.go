package chainquery

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"voidwalker/configs"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

type CQApi struct {
	dbConn *sql.DB
	cache  *sync.Map
}

var instance *CQApi
var ClaimNotFoundErr = errors.Base("claim not found")

func Init() (*CQApi, error) {
	if instance != nil {
		return instance, nil
	}
	db, err := connect()
	if err != nil {
		return nil, err
	}
	instance = &CQApi{
		cache:  &sync.Map{},
		dbConn: db,
	}
	return instance, nil
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", configs.Configuration.Chainquery.User, configs.Configuration.Chainquery.Password, configs.Configuration.Chainquery.Host, configs.Configuration.Chainquery.Database))
	return db, errors.Err(err)
}

type Claim struct {
	ClaimID     string `json:"claim_id"`
	SdHash      string `json:"sd_hash"`
	ContentType string `json:"content_type"`
}
type ChainqueryCache struct {
	lastFetched time.Time
	claim       *Claim
}

func (c *CQApi) ResolveClaim(claimName, shortID string) (*Claim, error) {
	cached, ok := c.cache.Load(claimName + shortID)
	if ok {
		cachedCasted, _ := cached.(ChainqueryCache)
		if cachedCasted.lastFetched.After(time.Now().Add(-5 * time.Minute)) {
			logrus.Infoln("loaded from cache")
			return cachedCasted.claim, nil
		}
	}
	query := "SELECT claim_id, sd_hash, content_type FROM claim where name = ? AND claim_id LIKE ?"
	rows, err := c.dbConn.Query(query, claimName, fmt.Sprintf("%s%%", shortID))
	if err != nil {
		return nil, errors.Err(err)
	}
	defer rows.Close()

	for rows.Next() {
		var claim Claim
		err = rows.Scan(&claim.ClaimID, &claim.SdHash, &claim.ContentType)
		if err != nil {
			return nil, errors.Err(err)
		}
		c.cache.Store(claimName+shortID, ChainqueryCache{
			lastFetched: time.Now(),
			claim:       &claim,
		})
		return &claim, nil
	}
	return nil, ClaimNotFoundErr
}

func (c *CQApi) ResolveClaimByChannel(claimName, channelID, channelName string) (*Claim, error) {
	cached, ok := c.cache.Load(claimName + channelID)
	if ok {
		cachedCasted, _ := cached.(ChainqueryCache)
		if cachedCasted.lastFetched.After(time.Now().Add(-5 * time.Minute)) {
			logrus.Infoln("loaded from cache")
			return cachedCasted.claim, nil
		}
	}
	query := "SELECT claim_id, sd_hash, content_type FROM claim where name = ? AND publisher_id = (select claim_id from claim where claim_id like ? and name = ? order by id desc limit 1)"
	rows, err := c.dbConn.Query(query, claimName, fmt.Sprintf("%s%%", channelID), channelName)
	if err != nil {
		return nil, errors.Err(err)
	}
	defer rows.Close()

	for rows.Next() {
		var claim Claim
		err = rows.Scan(&claim.ClaimID, &claim.SdHash, &claim.ContentType)
		if err != nil {
			return nil, errors.Err(err)
		}
		c.cache.Store(claimName+channelID, ChainqueryCache{
			lastFetched: time.Now(),
			claim:       &claim,
		})
		return &claim, nil
	}
	return nil, ClaimNotFoundErr
}
