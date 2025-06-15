const fs = require("fs")
const path = require("path")
const translate = require("google-translate-api-x")

const args = process.argv.slice(2)
const dryRun = args.includes("--dry")
const localeArg = args.find((arg) => arg.startsWith("--locale="))?.split("=")[1]

const localesDir = path.resolve(__dirname, "../locale")
const baseLang = "en"
const files = ["client.json", "server.json"]

function extractPlaceholders(str) {
  return [...str.matchAll(/{[^}]+}/g)].map((m) => m[0])
}

async function translateWithPlaceholders(text, from, to) {
  const placeholders = extractPlaceholders(text)
  let safeInput = text

  placeholders.forEach((p, i) => {
    safeInput = safeInput.replace(p, `__VAR_${i}__`)
  })

  const { text: raw } = await translate(safeInput, { from, to })

  let restored = raw
  placeholders.forEach((p, i) => {
    restored = restored.replace(new RegExp(`__VAR_${i}__`, "i"), p)
  })

  return restored
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
      const targetLang = locale.split("-")[0]

      let changed = false

      for (const [key, enVal] of Object.entries(enJson)) {
        const locVal = localeJson[key]
        if (locVal && locVal !== enVal) continue

        try {
          const translated = await translateWithPlaceholders(enVal, "en", targetLang)

          if (translated && translated !== enVal) {
            if (dryRun) {
              console.log(`[DRY RUN] ${locale}/${file} - ${key}: "${enVal}" → "${translated}"`)
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
