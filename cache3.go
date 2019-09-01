package disgord

import (
	jp "github.com/buger/jsonparser"
)

// gatewayCacher allows cache repositories to handle event content.
//
// handleGatewayEvent should only be called for situations where a object
// will be created, or updated. Not deleted. For deleting an object use
// the Delete method of a BasicCacheRepo implemented cache repository.
type gatewayCacher interface {
	handleGatewayEvent(evt string, data []byte, flags Flag) (updated interface{}, err error)
}
type restCacher interface {
	handleRESTResponse(obj interface{}) error
}

type BasicCacheRepo interface {
	Size() uint
	Cap() uint
	ListIDs() []Snowflake
	// Get returns nil when no match was found
	Get(id Snowflake) interface{}
	Del(id Snowflake)
}

type cache3 struct {
}

//////////////////////////////////////////////////////
//
// JSON HELPERS
// All helpers must start with a lowercase "json".
//
//////////////////////////////////////////////////////

// jsonGetSnowflake
func jsonGetSnowflake(data []byte, keys ...string) (id Snowflake, err error) {
	var bytes []byte
	bytes, _, _, err = jp.Get(data, keys...)
	if err != nil {
		return 0, err
	}

	if err = id.UnmarshalJSON(bytes); err != nil {
		return 0, err
	}

	return id, nil
}

// jsonNumberOfKeys returns the number of json keys at depth 1.
func jsonNumberOfKeys(data []byte) (nrOfKeys uint) {
	jp.EachKey(data, func(i int, bytes []byte, valueType jp.ValueType, e error) {
		nrOfKeys++
	})
	return
}

func jsonArrayLen(data []byte, keys ...string) (len int) {
	_, _ = jp.ArrayEach(data, func(b []byte, _ jp.ValueType, _ int, _ error) {
		len++
	}, keys...)

	return
}
