const fs = require("fs")
const path = require("path")
const translate = require("google-translate-api-x")

const args = process.argv.slice(2)
const dryRun = args.includes("--dry")
const localeArg = args.find((arg) => arg.startsWith("--locale="))?.split("=")[1]

const localesDir = path.resolve(__dirname, "../locale")
const baseLang = "en"
const files = ["client.json", "server.json"]

// 🧠 Safely freeze all top-level {...} blocks, even with nested ICU braces
function freezePlaceholders(input) {
  const replacements = {}
  let output = ""
  let index = 0
  let i = 0

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
      const token = `[[${index++}]]`
      replacements[token] = raw
      output += token
      i = j
    } else {
      output += input[i]
      i++
    }
  }

  output = output.replace(/\n/g, "[[*]]")

  return { frozen: output, replacements }
}

function restorePlaceholders(str, replacements) {
  let restored = Object.entries(replacements).reduce((out, [token, original]) => out.replace(token, original), str)
  restored = restored.replace(/\[\[\*\]\]/g, "\n")
  return restored
}

async function translatePreservingPlaceholders(text, from, to) {
  const { frozen, replacements } = freezePlaceholders(text)
  const { text: raw } = await translate(frozen, { from, to, forceTo: true })
  return restorePlaceholders(raw, replacements)
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
        if (locVal && locVal !== enVal) continue

        try {
          const translated = await translatePreservingPlaceholders(enVal, "en", targetLang)

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
