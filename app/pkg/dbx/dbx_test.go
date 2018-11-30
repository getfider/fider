package dbx_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
)

type user struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	IgnoreThis string
}

type userWithTenant struct {
	ID     int     `db:"id"`
	Name   string  `db:"name"`
	Tenant *tenant `db:"tenant"`
}

type userProvider struct {
	Provider string          `db:"provider"`
	User     *userWithTenant `db:"user"`
}

type tenant struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func TestMain(m *testing.M) {
	db := dbx.New()
	db.Seed()
	db.Close()

	code := m.Run()
	os.Exit(code)
}

func TestBind_SimpleStruct(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := user{}

	err := trx.Get(&u, "SELECT id, name FROM users LIMIT 1")
	Expect(err).IsNil()
	Expect(u.ID).Equals(1)
	Expect(u.Name).Equals("Jon Snow")
	Expect(u.IgnoreThis).Equals("")
}

func TestBind_DeepNestedStruct(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := userProvider{}

	err := trx.Get(&u, `SELECT up.provider,
														up.user_id,
														u.tenant_id AS user_tenant_id
										 FROM user_providers up
										 INNER JOIN users u
										 ON up.user_id = u.id
										 WHERE provider_uid = 'FB2222'`)
	Expect(err).IsNil()
	Expect(u.Provider).Equals("facebook")
	Expect(u.User.ID).Equals(4)
	Expect(u.User.Tenant.ID).Equals(2)
}

func TestBind_SimpleStruct_SingleField(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := user{}

	err := trx.Get(&u, "SELECT name FROM users LIMIT 1")
	Expect(err).IsNil()
	Expect(u.ID).Equals(0)
	Expect(u.Name).Equals("Jon Snow")
}

func TestBind_SimpleStruct_Multiple(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()
	u := []*user{}

	err := trx.Select(&u, "SELECT name FROM users WHERE tenant_id = 1 ORDER BY id")
	Expect(err).IsNil()

	Expect(u).HasLen(3)
	Expect(u[0].Name).Equals("Jon Snow")
	Expect(u[1].Name).Equals("Arya Stark")
	Expect(u[2].Name).Equals("Sansa Stark")
}

func TestBind_NestedStruct(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	u := userWithTenant{}

	err := trx.Get(&u, `
		SELECT u.id, u.name, t.id AS tenant_id, t.name AS tenant_name
		FROM users u
		INNER JOIN tenants t
		ON t.id = u.tenant_id
		WHERE u.id = 1
		LIMIT 1
	`)
	Expect(err).IsNil()
	Expect(u.ID).Equals(1)
	Expect(u.Name).Equals("Jon Snow")
	Expect(u.Tenant).IsNotNil()
	Expect(u.Tenant.ID).Equals(1)
	Expect(u.Tenant.Name).Equals("Demonstration")
}

func TestBind_NestedStruct_Multiple(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	u := []*userWithTenant{}

	err := trx.Select(&u, `
		SELECT u.id, u.name, t.id AS tenant_id, t.name AS tenant_name
		FROM users u
		INNER JOIN tenants t
		ON t.id = u.tenant_id
		WHERE u.tenant_id = 1
		ORDER BY u.id
	`)
	Expect(err).IsNil()
	Expect(u).HasLen(3)
	Expect(u[0].Name).Equals("Jon Snow")
	Expect(u[0].Tenant.Name).Equals("Demonstration")
	Expect(u[1].Name).Equals("Arya Stark")
	Expect(u[1].Tenant.Name).Equals("Demonstration")
	Expect(u[2].Name).Equals("Sansa Stark")
	Expect(u[2].Tenant.Name).Equals("Demonstration")
}

func TestExists_True(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	exists, err := trx.Exists("SELECT 1 FROM users WHERE id = 1")
	Expect(err).IsNil()
	Expect(exists).IsTrue()
}

func TestExists_False(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	exists, err := trx.Exists("SELECT 1 FROM users WHERE id = 0")
	Expect(err).IsNil()
	Expect(exists).IsFalse()
}

func TestCount(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	count, err := trx.Count("SELECT 1 FROM users WHERE id = 1")
	Expect(err).IsNil()
	Expect(count).Equals(1)
}

func TestCount_Empty(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	count, err := trx.Count("SELECT 1 FROM users WHERE id = 0")
	Expect(err).IsNil()
	Expect(count).Equals(0)
}

func TestScalar(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	var value int
	err := trx.Scalar(&value, "SELECT id FROM users WHERE id = 1")
	Expect(err).IsNil()
	Expect(value).Equals(1)
}

func TestArray(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	type postTags struct {
		ID   int     `db:"id"`
		Tags []int64 `db:"tags"`
	}

	result := postTags{}
	err := trx.Get(&result, "SELECT 1 as id, array[5,10] as tags")
	Expect(err).IsNil()
	Expect(result.ID).Equals(1)
	Expect(result.Tags).HasLen(2)
	Expect(result.Tags[0]).Equals(int64(5))
	Expect(result.Tags[1]).Equals(int64(10))
}

func TestArray_Empty(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	type postTags struct {
		ID   int    `db:"id"`
		Slug string `db:"slug"`
	}

	result := []*postTags{}
	err := trx.Select(&result, "SELECT id, name FROM Tags WHERE id = -1")
	Expect(err).IsNil()
	Expect(result).HasLen(0)
	bytes, err := json.Marshal(result)
	Expect(err).IsNil()
	Expect(string(bytes)).Equals("[]")
}

func TestByteArray(t *testing.T) {
	RegisterT(t)
	db := dbx.New()
	defer db.Close()

	trx, _ := db.Begin()
	defer trx.Rollback()

	type file struct {
		ContentType string `db:"content_type"`
		Size        int    `db:"size"`
		Content     []byte `db:"file"`
	}

	fileContent, err := ioutil.ReadFile(env.Path("/favicon.ico"))
	Expect(err).IsNil()

	_, err = trx.Execute(`
		INSERT INTO uploads (tenant_id, size, content_type, file)
		VALUES (1, $1, 'text/plain', $2)
	`, len(fileContent), fileContent)

	Expect(err).IsNil()

	theFile := file{}
	err = trx.Get(&theFile, "SELECT content_type, size, file FROM uploads WHERE id = 1")
	Expect(err).IsNil()
	Expect(theFile.ContentType).Equals("text/plain")
	Expect(theFile.Content).Equals(fileContent)
	Expect(theFile.Size).Equals(len(theFile.Content))
}
