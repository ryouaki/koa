package session

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/ryouaki/koa"
)

// Store interface
type Store interface {
	Save(key string, value map[string]interface{}, second time.Duration) error
	Get(key string) (map[string]interface{}, error)
}

// KoaSession struct
type KoaSession struct {
	Store      Store
	CookieConf http.Cookie
}

// MemStore struct
type MemStore struct {
	data map[string]interface{}
}

// RedisStore struct
type RedisStore struct {
	redisClient redis.UniversalClient
}

// MemInfo struct
type MemInfo struct {
	value map[string]interface{}
	time  time.Time
}

// NewMemStore func
func NewMemStore() *MemStore {
	return &MemStore{
		data: make(map[string]interface{}),
	}
}

// Get func
func (store MemStore) Get(key string) (map[string]interface{}, error) {
	if v, ok := store.data[key]; ok {
		data := v.(*MemInfo)
		if data.time.Before(time.Now()) {
			return nil, nil
		}
		return data.value, nil
	}
	return nil, nil
}

// Save func
func (store MemStore) Save(key string, value map[string]interface{}, second time.Duration) error {
	if value != nil {
		store.data[key] = &MemInfo{
			value: value,
			time:  time.Now().Add(second),
		}
	}
	return nil
}

// Get func
func (store RedisStore) Get(key string) (map[string]interface{}, error) {
	cmd := store.redisClient.Get(key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	b, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	value := make(map[string]interface{})

	err = json.Unmarshal(b, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Save func
func (store RedisStore) Save(key string, value map[string]interface{}, second time.Duration) error {
	if value == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return store.redisClient.Set(key, data, second).Err()
}

// NewRedisStore func
func NewRedisStore(rds redis.UniversalClient) *RedisStore {
	return &RedisStore{
		redisClient: rds,
	}
}

var localAddrIp = ""

func init() {
	localAddrIp = koa.GetLocalAddrIp()
}

// Session func
func Session(conf http.Cookie, store Store) func(*koa.Context, koa.Next) {
	name := "koa_sessid"
	path := "/"

	if conf.Name != "" {
		name = conf.Name
	}

	if conf.Path != "" {
		path = conf.Path
	}

	sess := &KoaSession{
		Store: store,
		CookieConf: http.Cookie{
			Name:       name,
			Path:       path,
			Domain:     conf.Domain,
			Expires:    conf.Expires,
			RawExpires: conf.RawExpires,
			MaxAge:     conf.MaxAge,
			Secure:     conf.Secure,
			HttpOnly:   conf.HttpOnly,
			SameSite:   conf.SameSite,
			Raw:        conf.Raw,
			Unparsed:   conf.Unparsed,
		},
	}

	return func(ctx *koa.Context, next koa.Next) {
		cookie := ctx.GetCookie(sess.CookieConf.Name)
		beforeSession := make(map[string]interface{}, 0)
		var sessID string
		if cookie == nil {
			sessID = fmt.Sprintf("koa_sess-%04d-%s", koa.GetGoroutineID(), getMd5([]byte(localAddrIp+string(rune(time.Now().Unix())))))

			cookie = &http.Cookie{
				Name:       sess.CookieConf.Name,
				Path:       sess.CookieConf.Path,
				Domain:     sess.CookieConf.Domain,
				Expires:    sess.CookieConf.Expires,
				RawExpires: sess.CookieConf.RawExpires,
				MaxAge:     sess.CookieConf.MaxAge,
				Secure:     sess.CookieConf.Secure,
				HttpOnly:   sess.CookieConf.HttpOnly,
				SameSite:   sess.CookieConf.SameSite,
				Raw:        sess.CookieConf.Raw,
				Unparsed:   sess.CookieConf.Unparsed,
				Value:      sessID,
			}
			ctx.SetCookie(cookie)
		} else {
			sessID = cookie.Value
			session, err := sess.Store.Get(sessID)
			if err == nil && session != nil {
				beforeSession = session
			}
		}

		ctx.SetData("session", beforeSession)

		next()

		afterSession := ctx.GetData("session")
		if afterSession == nil {
			afterSession = make(map[string]interface{}, 0)
		}
		sess.Store.Save(sessID, afterSession.(map[string]interface{}), time.Duration(sess.CookieConf.MaxAge)*time.Second)
	}
}

// GetMD5ID func
func getMd5(b []byte) string {
	res := md5.Sum(b)
	return hex.EncodeToString(res[:])
}
