package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/gotnospirit/messageformat"
)

// localeToPlurals maps between Fider locale and gotnospirit/messageformat culture
var localeToPlurals = map[string]string{
	"en": "en",
	"pt-BR": "pt",
}

type Params map[string]interface{}

// cache for locale parser and file content to prevent excessive disk IO
var cache = make(map[string]localeData)
var mu sync.RWMutex
type localeData struct {
	file   map[string]string
	parser *messageformat.Parser
}

// getLocaleData returns the file content and culture specific parser
func getLocaleData(locale string) localeData {
	if item, ok := cache[locale]; ok {
		return item
	}

	mu.Lock()
	defer mu.Unlock()

	if item, ok := cache[locale]; ok {
		return item
	}

	content, err := ioutil.ReadFile(env.Path("locale/" + locale + ".json"))
	if err != nil {
		panic(errors.Wrap(err, "failed to read locale file"))
	}

	var file map[string]string
	err = json.Unmarshal(content, &file)
	if err != nil {
		panic(errors.Wrap(err, "failed unmarshal to json"))
	}

	parser, err := messageformat.NewWithCulture(localeToPlurals[locale])
	if err != nil {
		panic(errors.Wrap(err, "failed create parser"))
	}

	data := localeData{file, parser}

	if env.IsProduction() {
		cache[locale] = data
	}

	return data
}

// getMessage returns the translated message for a given locale
// If given key is not found, it'll fallback to english
func getMessage(locale, key string) (string, *messageformat.Parser) {
	localeData := getLocaleData(locale)
	if str, ok := localeData.file[key]; ok {
		return str, localeData.parser
	}

	enData := getLocaleData("en")
	if str, ok := enData.file[key]; ok {
		return str, enData.parser
	}

	return fmt.Sprintf("⚠️ Missing Translation: %s", key), enData.parser
}

// T translates a given key to current context locale
// If Locale is not set in context, the environment locale is used
// Params is used to replace variables and pluralize
func T(ctx context.Context, key string, params ...Params) string {
	locale, ok := ctx.Value(app.LocaleCtxKey).(string)
	if !ok {
		locale = env.Config.Locale
	}

	msg, parser := getMessage(locale, key)
	if len(params) == 0 {
		return msg
	}

	parsedMsg, err := parser.Parse(msg)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to parse msg '%s'", msg)))
	}

	str, err := parsedMsg.FormatMap(params[0])
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to format msg '%s' with params '%v'", msg, params[0])))
	}

	return str
}
