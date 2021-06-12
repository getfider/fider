// import { t, activate } from "./i18n"

// ;[
//   { locale: "", expected: "Search" },
//   { locale: "en", expected: "Search" },
//   { locale: "pt-BR", expected: "Pesquisa" },
// ].forEach((testCase) => {
//   test(`${testCase.locale}: valid message`, async () => {
//     await activate(testCase.locale)
//     expect(t("Search")).toBe(testCase.expected)
//   })
// })
// ;[
//   { locale: "", expected: "No similar posts matched 'Hello World'." },
//   { locale: "en", expected: "No similar posts matched 'Hello World'." },
//   { locale: "pt-BR", expected: "Nenhuma postagem semelhante à 'Hello World'." },
// ].forEach((testCase) => {
//   test(`${testCase.locale}: valid message with params`, async () => {
//     await activate(testCase.locale)
//     const translated = t("No similar posts matched '{title}'.", { title: "Hello World" })
//     expect(translated).toBe(testCase.expected)
//   })
// })
// ;[
//   { locale: "", expected1: "Send 1 invite", expected2: "Send 2 invites" },
//   { locale: "en", expected1: "Send 1 invite", expected2: "Send 2 invites" },
//   { locale: "pt-BR", expected1: "Enviar 1 convite", expected2: "Enviar 2 convites" },
// ].forEach((testCase) => {
//   test(`${testCase.locale}: valid plural message`, async () => {
//     await activate(testCase.locale)
//     const translated1 = t("{count, plural, one {Send # invite} other {Send # invites}}", { count: 1 })
//     expect(translated1).toBe(testCase.expected1)

//     const translated2 = t("{count, plural, one {Send # invite} other {Send # invites}}", { count: 2 })
//     expect(translated2).toBe(testCase.expected2)
//   })
// })
// ;[{ locale: "" }, { locale: "en" }, { locale: "pt-BR" }].forEach((testCase) => {
//   test(`${testCase.locale}: invalid message`, async () => {
//     await activate(testCase.locale)
//     const translated = t("This message does not exist")
//     expect(translated).toBe("⚠️ Missing Translation: This message does not exist")
//   })
// })
