package disgord

//
//func TestCache2_users_gatewayEvent(t *testing.T) {
//	c := &usersCache{}
//	algs := []interfaces.CacheAlger{
//		lfu.NewCacheList(0),
//		lru.NewCacheList(0),
//		// tlfu
//	}
//
//	files, err := testdata.GetDataForDir("user", "user")
//	if err != nil {
//		panic(err)
//	}
//
//	for _, alg := range algs {
//		c.internal = alg
//
//		for _, file := range files {
//			data, err := ioutil.ReadFile(file)
//			if err != nil {
//				t.Fatal(err)
//			}
//			_, err = c.handleGatewayEvent("", data)
//			if err != nil {
//				t.Error(err)
//			}
//		}
//
//		userIDs := c.ListIDs()
//		if len(userIDs) != 1 {
//			t.Errorf("expected only one user. Got %d", userIDs)
//		}
//
//		usr := c.Get(userIDs[0])
//		if usr.Avatar != nil {
//			t.Error("expected avatar to be nil")
//		}
//	}
//}
