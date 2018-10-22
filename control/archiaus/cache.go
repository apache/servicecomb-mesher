package archiaus

import "github.com/patrickmn/go-cache"

//save configs
var (
	//key is service name
	EgressConfigCache = cache.New(0, 0)
)
