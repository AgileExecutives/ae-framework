const fs = require('fs')
const path = require('path')

// jiti allows requiring TypeScript/ESM files without compilation
let jiti
try {
  jiti = require('jiti')(process.cwd())
} catch (e) {
  console.error('Please install dev dependency "jiti" to run the i18n checker (npm install).')
  process.exit(2)
}

const localesDir = path.join(__dirname, '..', 'src', 'locales')

function flatten(obj, prefix = '') {
  const res = {}
  for (const key of Object.keys(obj)) {
    const val = obj[key]
    const pathKey = prefix ? `${prefix}.${key}` : key
    if (val && typeof val === 'object' && !Array.isArray(val)) {
      Object.assign(res, flatten(val, pathKey))
    } else {
      res[pathKey] = true
    }
  }
  return res
}

function loadLocale(filePath) {
  const mod = jiti(filePath)
  return mod && mod.default ? mod.default : mod
}

function toTs(obj, indent = 0) {
  const pad = '  '.repeat(indent)
  if (obj === null) return 'null'
  if (typeof obj === 'string') return `'${String(obj).replace(/'/g, "\\'")}'`
  if (typeof obj === 'number' || typeof obj === 'boolean') return String(obj)
  if (Array.isArray(obj)) {
    const items = obj.map(i => toTs(i, indent + 1)).join(', ')
    return `[${items}]`
  }
  const entries = Object.keys(obj).map(key => {
    const safeKey = /^[A-Za-z_$][A-Za-z0-9_$]*$/.test(key) ? key : `'${key.replace(/'/g, "\\'")}'`
    const val = toTs(obj[key], indent + 1)
    if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
      return `${pad}  ${safeKey}: ${val}`
    }
    return `${pad}  ${safeKey}: ${val}`
  })
  return `{
${entries.join(',\n')}
${pad}}`
}

function writeLocaleFile(filePath, obj) {
  const content = `export default ${toTs(obj)}\n`
  fs.writeFileSync(filePath, content, 'utf8')
}

function main() {
  if (!fs.existsSync(localesDir)) {
    console.error('Locales directory not found:', localesDir)
    process.exit(1)
  }

  const files = fs.readdirSync(localesDir).filter(f => f.endsWith('.ts') || f.endsWith('.js'))
  if (!files.includes('en.ts') && !files.includes('en.js')) {
    console.error('Base locale `en` not found in', localesDir)
    process.exit(1)
  }

  const baseFile = files.find(f => f.startsWith('en.'))
  const base = loadLocale(path.join(localesDir, baseFile))
  const baseKeys = flatten(base)

  let failed = false
  const shouldFix = process.argv.includes('--fix') || process.argv.includes('fix') || process.argv.includes('--apply')

  for (const f of files) {
    if (f === baseFile) continue
    const locale = path.basename(f, path.extname(f))
    const obj = loadLocale(path.join(localesDir, f))
    const keys = flatten(obj)
    const missing = []
    for (const k of Object.keys(baseKeys)) {
      if (!keys[k]) missing.push(k)
    }

    if (missing.length) {
      if (!shouldFix) failed = true
      console.error(`\nLocale ${locale} is missing ${missing.length} keys:`)
      for (const m of missing) console.error('  -', m)

      if (shouldFix) {
        console.log(`Auto-fixing locale ${locale} by adding ${missing.length} keys.`)
        // build merged object based on base but starting from existing locale object
        const merged = JSON.parse(JSON.stringify(base))

        function applyMissing(target, source, path = '') {
          for (const k of Object.keys(source)) {
            const curPath = path ? `${path}.${k}` : k
            if (typeof source[k] === 'object' && source[k] !== null && !Array.isArray(source[k])) {
              target[k] = target[k] || {}
              applyMissing(target[k], source[k], curPath)
            } else {
              // leaf
              // if locale already has value keep it, else copy from base
              const parts = curPath.split('.')
              let node = obj
              for (let i = 0; i < parts.length - 1; i++) {
                node = node && node[parts[i]]
                if (!node) break
              }
              const leafKey = parts[parts.length - 1]
              const hasValue = node && Object.prototype.hasOwnProperty.call(node, leafKey)
              // set placeholder from base
              // navigate merged and set value at path
              let mnode = merged
              for (let i = 0; i < parts.length - 1; i++) {
                mnode = mnode[parts[i]] = mnode[parts[i]] || {}
              }
              if (!hasValue) mnode[leafKey] = source[k]
            }
          }
        }

        applyMissing(merged, base)

        // replace placeholders with original locale values where present
        function overlay(original, target) {
          for (const k of Object.keys(original)) {
            if (typeof original[k] === 'object' && original[k] !== null && !Array.isArray(original[k])) {
              target[k] = target[k] || {}
              overlay(original[k], target[k])
            } else {
              target[k] = original[k]
            }
          }
        }

        overlay(obj, merged)

        // write back to file
        const filePath = path.join(localesDir, f)
        try {
          writeLocaleFile(filePath, merged)
          console.log(`Wrote updated locale file: ${filePath}`)
        } catch (e) {
          console.error('Failed to write locale file', filePath, e)
          // if writing failed count as failure
          failed = true
        }
      }
    } else {
      console.log(`Locale ${locale}: OK`)
    }
  }

  if (failed) {
    console.error('\nI18n key check failed: some locales are missing keys compared to en.')
    process.exit(2)
  }

  console.log('\nI18n key check passed for all locales.')
}

main()
