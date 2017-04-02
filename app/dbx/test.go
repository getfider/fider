package dbx

import (
	"io/ioutil"

	"github.com/WeCanHearYou/wechy/app/toolbox/env"
)

func setup(db Database) {
}

func teardown(db Database) {
	if !env.IsTest() {
		return
	}

	bytes, err := ioutil.ReadFile(env.Path("/app/dbx/teardown.sql"))
	if err != nil {
		panic(err)
	}

	db.Execute(string(bytes))
}
