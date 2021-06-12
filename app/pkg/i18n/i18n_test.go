package i18n_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/i18n"

	. "github.com/getfider/fider/app/pkg/assert"
)

var emptyContext = context.Background()
var enContext = context.WithValue(emptyContext, app.LocaleCtxKey, "en")
var ptBRContext = context.WithValue(emptyContext, app.LocaleCtxKey, "pt-BR")

func TestValidMessage(t *testing.T) {
	RegisterT(t)

	key := "Search"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("Search")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("Search")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("Pesquisa")
}

func TestEnglish_ValidMessage_WithParams(t *testing.T) {
	RegisterT(t)

	key := "No similar posts matched '{title}'."

	translated := i18n.T(emptyContext, key, i18n.Params{"title": "Hello World"})
	Expect(translated).Equals("No similar posts matched 'Hello World'.")

	translated = i18n.T(enContext, key, i18n.Params{"title": "Hello World"})
	Expect(translated).Equals("No similar posts matched 'Hello World'.")

	translated = i18n.T(ptBRContext, key, i18n.Params{"title": "Olá Mundo"})
	Expect(translated).Equals("Nenhuma postagem semelhante à 'Olá Mundo'.")
}

func TestEnglish_ValidMessage_Plurals(t *testing.T) {
	RegisterT(t)

	key := "{count, plural, one {Send # invite} other {Send # invites}}"

	translated := i18n.T(emptyContext, key, i18n.Params{"count": 1})
	Expect(translated).Equals("Send 1 invite")
	translated = i18n.T(emptyContext, key, i18n.Params{"count": 2})
	Expect(translated).Equals("Send 2 invites")

	translated = i18n.T(enContext, key, i18n.Params{"count": 1})
	Expect(translated).Equals("Send 1 invite")
	translated = i18n.T(enContext, key, i18n.Params{"count": 2})
	Expect(translated).Equals("Send 2 invites")

	translated = i18n.T(ptBRContext, key, i18n.Params{"count": 1})
	Expect(translated).Equals("Enviar 1 convite")
	translated = i18n.T(ptBRContext, key, i18n.Params{"count": 2})
	Expect(translated).Equals("Enviar 2 convites")
}

func TestEnglish_InvalidMessage(t *testing.T) {
	RegisterT(t)

	key := "This message does not exist"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")
}