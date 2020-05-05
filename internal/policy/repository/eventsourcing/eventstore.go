package eventsourcing

import (
	"github.com/caos/zitadel/internal/cache/config"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

type PolicyEventstore struct {
	es_int.Eventstore
	policyCache *PolicyCache
}

type PolicyConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartPolicy(conf PolicyConfig) (*PolicyEventstore, error) {
	policyCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	return &PolicyEventstore{
		Eventstore:  conf.Eventstore,
		policyCache: policyCache,
	}, nil
}
