package fx

import (
	"context"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/google/uuid"
	xcontext "github.com/org/2112-space-lab/org/testing/pkg/x-context"
)

// CacheStore defines what methods should be provided by the cache storage
type CacheStore interface {
	store.StoreInterface
}

// Cached acts as a wrapper for getting/setting a value from cache if available or from getter
type Cached[T any] struct {
	cacheStore CacheStore
	cacheKey   string
	cacheOpts  []store.Option
}

// NewCached standard constructor
func NewCached[T any](ctx context.Context, s CacheStore, cacheKey string, options ...store.Option) *Cached[T] {
	uid, err := uuid.NewRandom()
	uidStr := "_failedUID"
	if err != nil {
		_, ctxLog := xcontext.ContextBuilder(ctx, "Cached.NewCached", nil)
		ctxLog.WithError(err).Warnf("failed to created cached container for [%s]", cacheKey)
	} else {
		uidStr = "_" + uid.String()
	}
	cacheOpts := options
	return &Cached[T]{
		cacheStore: s,
		cacheKey:   cacheKey + uidStr,
		cacheOpts:  cacheOpts,
	}
}

// Set puts the value in cache. In case of error or no cache store, a warning is logged and error is silenced
func (c *Cached[T]) Set(ctx context.Context, val T) {
	ctx, ctxLog := xcontext.ContextBuilder(ctx, "Cached.Set", nil)
	if c.cacheStore == nil {
		ctxLog.Debugf("nil cache store for [%s] - value will not be cached", c.cacheKey)
		return
	}
	err := c.cacheStore.Set(ctx, c.cacheKey, val, c.cacheOpts...)
	if err != nil {
		ctxLog.WithError(err).Warnf("failed to set cached value for [%s]", c.cacheKey)
	}
}

// Get looks-up the value from cache. In case of error or no cache store, the getter is used as fallback, a warning is logged and error is silenced
func (c *Cached[T]) Get(ctx context.Context, getter func(context.Context) (T, error)) (T, error) {
	ctx, ctxLog := xcontext.ContextBuilder(ctx, "Cached.Get", nil)
	if c.cacheStore == nil {
		ctxLog.Debugf("nil cache store for [%s] - will not be fetched from cache", c.cacheKey)
		return getter(ctx)
	}
	valFromCache, err := c.cacheStore.Get(ctx, c.cacheKey)
	if err == nil {
		typedVal, ok := valFromCache.(T)
		if !ok {
			ctxLog.Warnf("cached value for [%s] invalid type - will be fetched from getter", c.cacheKey)
		}
		return typedVal, nil
	}
	ctxLog.WithError(err).Tracef("value for [%s] not found in cache - will be fetched from getter", c.cacheKey)
	valFromGetter, err := getter(ctx)
	c.Set(ctx, valFromGetter)
	return valFromGetter, err
}
