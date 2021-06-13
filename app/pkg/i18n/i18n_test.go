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

	key := "error.pagenotfound.title"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("Page not found")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("Page not found")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("Página não encontrada")
}

func TestEnglish_ValidMessage_WithParams(t *testing.T) {
	RegisterT(t)

	key := "validation.settings.unknown"

	translated := i18n.T(emptyContext, key, i18n.Params{"name": "Notification"})
	Expect(translated).Equals("Unknown settings named 'Notification'")

	translated = i18n.T(enContext, key, i18n.Params{"name": "Notification"})
	Expect(translated).Equals("Unknown settings named 'Notification'")

	translated = i18n.T(ptBRContext, key, i18n.Params{"name": "Notification"})
	Expect(translated).Equals("Configuração 'Notification' é desconhecida")
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