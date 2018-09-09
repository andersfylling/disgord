package disgord

type Cacher interface{}

func NewCache() *Cache {
	return &Cache{}
}

type Cache struct{}
