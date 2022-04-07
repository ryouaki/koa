package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/ryouaki/koa"
)

// Config struct
type Config struct {
	Name       string
	Store      Store
	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	SameSite   http.SameSite
	Raw        string
	Unparsed   []string // Raw text of unparsed attribute-value pairs
}

// KoaSession struct
type KoaSession struct {
	localAddr string
	id        uint64
	Config
}

// Store interface
type Store interface {
	Save(key string, value map[string]interface{}, second time.Duration) error
	Get(key string) (map[string]interface{}, error)
}

// MemStore struct
type MemStore struct {
	data map[string]interface{}
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
func (store *MemStore) Get(key string) (map[string]interface{}, error) {
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
func (store *MemStore) Save(key string, value map[string]interface{}, second time.Duration) error {
	if value != nil {
		store.data[key] = &MemInfo{
			value: value,
			time:  time.Now().Add(second),
		}
	}
	return nil
}

// Get func
func (store *RedisStore) Get(key string) (map[string]interface{}, error) {
	cmd := store.redisClient.Get(key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	b, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	value := make(map[string]interface{})

	json.Unmarshal(b, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Save func
func (store *RedisStore) Save(key string, value map[string]interface{}, second time.Duration) error {
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

// RedisStore struct
type RedisStore struct {
	redisClient redis.UniversalClient
}

var sess *KoaSession = nil

// Session func
func Session(conf *Config) func(error, *koa.Context, koa.Next) {
	addr := koa.GetIPAddr()
	id := koa.GetGoroutineID()
	name := "koa_sess_id"
	path := "/"

	if conf.Name != "" {
		name = conf.Name
	}

	if conf.Path != "" {
		path = conf.Path
	}

	sess = &KoaSession{
		localAddr: addr,
		id:        id,
		Config: Config{
			Store:      conf.Store,
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

	return func(err error, ctx *koa.Context, next koa.Next) {
		sessionID := ctx.GetCookie(sess.Name)
		sessID := fmt.Sprintf("%v%d%s", time.Now().UnixNano()/1e6, koa.GetGoroutineID(), koa.GetMD5ID([]byte(sess.localAddr)))
		if sessionID != nil {
			sessID = sessionID.Value
		}
		beforeSession, err := sess.Store.Get(sessID)

		if beforeSession != nil {
			ctx.SetData("session", beforeSession)
		} else {
			cookie := &http.Cookie{}
			koa.StructAssign(cookie, &sess.Config)
			cookie.Value = sessID
			ctx.SetCookie(cookie)
		}
		next(err)

		var afterSession map[string]interface{}
		sessionData := ctx.GetData("session")
		if sessionData != nil {
			afterSession = sessionData.(map[string]interface{})
		}
		sess.Store.Save(sessID, afterSession, time.Duration(sess.MaxAge)*time.Second)
	}
}
