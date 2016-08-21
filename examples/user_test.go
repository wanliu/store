package examples

import (
	"github.com/wanliu/store"
	"testing"
)

// var store = NewStore(DefaultUserStore())

func TestUser(t *testing.T) {
	if _, err := store.Open("rbac.db"); err != nil {
		t.Fatalf("open database faield %s", err)
	}

	UserStore := store.NewStore(&User{}, nil)

	t.Logf("UserStore %+v", UserStore)

	usr, err := UserStore.Create(map[string]interface{}{
		"Login":  "hysios",
		"Email":  "hyysios@gmail.com",
		"Phone":  "07343-4380996",
		"Mobile": "118774781025",
		"Title":  "管理员",
		"Avatar": "/asdf",
	})

	if err != nil {
		t.Fatalf("create User failed %s", err)
	}
	t.Logf("UserStore created User object , %+v", usr)

	if usr, err = UserStore.Get(1); err != nil {
		t.Fatalf("read User failed %s", err)
	}

	t.Logf("UserStore read User object , %+v", usr)

	if err = UserStore.Put(1, map[string]interface{}{
		"Login": "hyysios",
	}); err != nil {
		t.Fatalf("put User failed %s", err)
	}

	if usr, err = UserStore.Get(1); err != nil {
		t.Fatalf("read User failed %s", err)
	}

	t.Logf("UserStore put User object , %+v", usr)
}

func init() {
	testPrepare()
}
