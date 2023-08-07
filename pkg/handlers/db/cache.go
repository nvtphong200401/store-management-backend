package db

type RedisKey string

const (
	ProductKey     RedisKey = "product"
	StoreKey       RedisKey = "store"
	SaleKey        RedisKey = "sale"
	ProductListKey RedisKey = "products"
	StoreListKey   RedisKey = "stores"
	SaleListKey    RedisKey = "sales"
)
