package mysql

import (
	"testing"

	"github.com/innermond/sky/auth"
	"github.com/innermond/sky/config"
	"github.com/innermond/sky/sky"
)

func TestTokenService_good_bad_apikeys(t *testing.T) {
	db := config.DB()
	sign := config.PrivateKey()
	if err := config.Err(); err != nil {
		t.Fatal(err)
	}
	c := auth.NewTokenCreator(sign)
	s := NewSession(db)

	ts := NewTokenService(s, c)
	t.Log(ts)

	keys := map[string]bool{
		"b4d0c82da6495cb3de8f7d14182ebd27c7423ad7a0be9f89812e5f56dc84421d": true,
		"b42a93814bfac815a286fd90bdf019f26081ec4a7e833a30b8e3c7ada6f3b5eb": true,
		"294a488de53e14ac02111db1729ba0184f2077eec18532ed0518b2847354b37a": true,
		"fake": false,
	}
	for ks, ok := range keys {
		t.Run(ks, func(t *testing.T) {
			k := sky.ApiKey(ks)

			tok, err := ts.Create(k)

			if err != nil && ok {
				t.Error(err)
			}
			if err == nil && !ok {
				t.Error(err)
			}
			t.Log(ok, tok)
		})
	}
}
