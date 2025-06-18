/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */

const fs = require("fs")
const path = require("path")
const { TranslationServiceClient } = require("@google-cloud/translate").v3

const client = new TranslationServiceClient()

const args = process.argv.slice(2)
const dryRun = args.includes("--dry")
const localeArg = args.find((arg) => arg.startsWith("--locale="))?.split("=")[1]

const localesDir = path.resolve(__dirname, "../locale")
const baseLang = "en"
const files = ["client.json", "server.json"]

function extractPlaceholders(input) {
  const placeholders = []
  let output = ""
  let i = 0
  let tokenIndex = 0

  while (i < input.length) {
    if (input[i] === "{") {
      let depth = 1
      let j = i + 1
      while (j < input.length && depth > 0) {
        if (input[j] === "{") depth++
        else if (input[j] === "}") depth--
        j++
      }

      const raw = input.slice(i, j)
      const token = `__PH_${tokenIndex++}__`
      placeholders.push({ token, value: raw })
      output += token
      i = j
    } else {
      output += input[i]
      i++
    }
  }

  return { cleaned: output, placeholders }
}

function restorePlaceholders(text, placeholders) {
  return placeholders.reduce((acc, { token, value }) => acc.replace(token, value), text)
}

async function translateText(text, targetLang) {
  const projectId = await client.getProjectId()
  const parent = `projects/${projectId}/locations/global`

  const { cleaned, placeholders } = extractPlaceholders(text)

  console.log(`Translating "${cleaned}" to ${targetLang}...`)

  const [response] = await client.translateText({
    parent,
    contents: [cleaned],
    mimeType: "text/plain",
    targetLanguageCode: targetLang,
  })

  const translated = response.translations[0].translatedText
  const result = restorePlaceholders(translated, placeholders)
  console.log(`... done, result was ${translated}, final result is ${result}`)
  return result
}

;(async () => {
  for (const file of files) {
    const enPath = path.join(localesDir, baseLang, file)
    const enJson = JSON.parse(fs.readFileSync(enPath, "utf8"))

    for (const locale of fs.readdirSync(localesDir)) {
      if (locale === baseLang || !fs.statSync(path.join(localesDir, locale)).isDirectory()) continue
      if (localeArg && locale !== localeArg) continue

      const localePath = path.join(localesDir, locale, file)
      if (!fs.existsSync(localePath)) continue

      const localeJson = JSON.parse(fs.readFileSync(localePath, "utf8"))
      const targetLang = locale

      let changed = false

      for (const [key, enVal] of Object.entries(enJson)) {
        const locVal = localeJson[key]
        if (locVal) continue // There is already a translation, skip it

        try {
          console.log(`New key: translating "${enVal} (${key})" to ${targetLang}...`)
          const translated = await translateText(enVal, targetLang)

          if (translated && translated !== enVal) {
            if (dryRun) {
              console.log(`[DRY RUN] ${locale}/${file} - ${key}: "${enVal.replace(/\n/g, "\\n")}" → "${translated.replace(/\n/g, "\\n")}"`)
            } else {
              localeJson[key] = translated
              changed = true
            }
          }
        } catch (err) {
          console.warn(`⚠️ Failed to translate '${key}' to ${locale}: ${err.message}`)
        }
      }

      if (changed && !dryRun) {
        fs.writeFileSync(localePath, JSON.stringify(localeJson, null, 2))
        console.log(`✅ Wrote updated: ${locale}/${file}`)
      } else if (dryRun) {
        console.log(`[DRY RUN] ${locale}/${file} - no changes`)
      }
    }
  }
})()
