package sqlcertcache_test

import (
	"context"
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"

	"golang.org/x/crypto/acme/autocert"

	"github.com/goenning/sqlcertcache"
	_ "github.com/lib/pq"
)

func expectNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Errorf("should not be nil")
	}
}

func expectNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Errorf("should be nil, but was %v", v)
	}
}

func expectEquals(t *testing.T, v interface{}, expected interface{}) {
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("should be %v, but was %v", expected, v)
	}
}

func getConnection() *sql.DB {
	conn, err := sql.Open("postgres", "postgres://pgcache_test:pgcache_test_pw@localhost:6543/pgcache_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	conn.Exec("drop table if exists autocert_cache")
	conn.Exec("drop table if exists cert_store")
	return conn
}

func TestNew(t *testing.T) {
	conn := getConnection()
	cache, err := sqlcertcache.New(conn, "autocert_cache")
	expectNotNil(t, cache)
	expectNil(t, err)

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM autocert_cache").Scan(&count)
	expectNil(t, err)
	expectEquals(t, count, 0)
}

func TestGet_UnkownKey(t *testing.T) {
	conn := getConnection()
	cache, _ := sqlcertcache.New(conn, "autocert_cache")
	data, err := cache.Get(context.Background(), "my-key")
	expectEquals(t, err, autocert.ErrCacheMiss)
	expectEquals(t, len(data), 0)
}

func TestGet_AfterPut(t *testing.T) {
	conn := getConnection()
	cache, _ := sqlcertcache.New(conn, "autocert_cache")

	actual, _ := ioutil.ReadFile("./LICENSE")
	err := cache.Put(context.Background(), "my-key", actual)
	expectNil(t, err)

	data, err := cache.Get(context.Background(), "my-key")
	expectNil(t, err)
	expectEquals(t, data, actual)

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM autocert_cache").Scan(&count)
	expectNil(t, err)
	expectEquals(t, count, 1)
}

func TestGet_AfterDelete(t *testing.T) {
	conn := getConnection()
	cache, _ := sqlcertcache.New(conn, "autocert_cache")

	actual := []byte{1, 2, 3, 4}
	err := cache.Put(context.Background(), "my-key", actual)
	expectNil(t, err)

	err = cache.Delete(context.Background(), "my-key")
	expectNil(t, err)

	data, err := cache.Get(context.Background(), "my-key")
	expectEquals(t, err, autocert.ErrCacheMiss)
	expectEquals(t, len(data), 0)
}

func TestDelete_UnkownKey(t *testing.T) {
	conn := getConnection()
	cache, _ := sqlcertcache.New(conn, "autocert_cache")

	var err error

	err = cache.Delete(context.Background(), "my-key1")
	expectNil(t, err)
	err = cache.Delete(context.Background(), "other-key")
	expectNil(t, err)
	err = cache.Delete(context.Background(), "hello-world")
	expectNil(t, err)
}

func TestPut_Overwrite(t *testing.T) {
	conn := getConnection()
	cache, _ := sqlcertcache.New(conn, "autocert_cache")

	data1 := []byte{1, 2, 3, 4}
	err := cache.Put(context.Background(), "thekey", data1)
	expectNil(t, err)
	data, err := cache.Get(context.Background(), "thekey")
	expectEquals(t, data, data1)

	data2 := []byte{5, 6, 7, 8}
	err = cache.Put(context.Background(), "thekey", data2)
	expectNil(t, err)
	data, err = cache.Get(context.Background(), "thekey")
	expectEquals(t, data, data2)
}

func TestDifferentTableName(t *testing.T) {
	conn := getConnection()
	cache, _ := sqlcertcache.New(conn, "cert_store")

	actual := []byte{1, 2, 3, 4}
	err := cache.Put(context.Background(), "thekey.hi", actual)
	expectNil(t, err)

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM cert_store").Scan(&count)
	expectNil(t, err)
	expectEquals(t, count, 1)

	err = conn.QueryRow("SELECT COUNT(*) FROM autocert_cache").Scan(&count)
	expectNotNil(t, err)
}

func TestGet_CancelledContext(t *testing.T) {
	conn := getConnection()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cache, _ := sqlcertcache.New(conn, "autocert_cache")
	data, err := cache.Get(ctx, "my-key")
	expectEquals(t, err, context.Canceled)
	expectEquals(t, len(data), 0)
}
