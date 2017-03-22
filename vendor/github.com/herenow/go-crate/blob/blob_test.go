package blob

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

var (
	driver *Driver
)

func init() {
	var err error
	driver, err = New("http://localhost:4200")
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func TestTableCreation(t *testing.T) {
	table, err := driver.GetTable("myblobs")
	if err != nil {
		table, err = driver.NewTable("myblobs", 1, 1)
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(table.Name)

	table, err = driver.NewTable("testTable", 1, 1)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	t.Log(table.Name)
	if err = table.Drop(); err != nil {
		t.Error(err)
	}
}

func TestUpload(t *testing.T) {
	data := strings.Repeat("asdfasdfa", rand.Intn(1000))
	r := strings.NewReader(data)
	digest := Sha1Digest(r)
	_, err := r.Seek(0, 0)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("data size:", r.Len())
	table, err := driver.GetTable("myblobs")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	t.Log(table.Name)
	_, err = table.Upload(digest, r)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
}

func TestUploadDelete(t *testing.T) {
	data := strings.Repeat("asdfasdfa", rand.Intn(1000))
	r := strings.NewReader(data)
	digest := Sha1Digest(r)
	_, err := r.Seek(0, 0)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("data size:", r.Len())
	table, err := driver.GetTable("myblobs")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	t.Log(table.Name)
	record, err := table.Upload(digest, r)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if err = table.Delete(record.Digest); err != nil {
		t.Error(err.Error())
	}
}

func TestUploadEx(t *testing.T) {
	table, err := driver.GetTable("myblobs")
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	data := strings.Repeat("asdfasdfa", rand.Intn(1000))
	r := strings.NewReader(data)
	record, err := table.UploadEx(r)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if err = table.Delete(record.Digest); err != nil {
		t.Error(err.Error())
	}
}
