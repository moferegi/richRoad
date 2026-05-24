import { spawnSync } from 'node:child_process'

const npmCmd = process.platform === 'win32' ? 'npm.cmd' : 'npm'

const baseline = {
  total: Number(process.env.UNI_AUDIT_BASELINE_TOTAL || 51),
  high: Number(process.env.UNI_AUDIT_BASELINE_HIGH || 34),
  moderate: Number(process.env.UNI_AUDIT_BASELINE_MODERATE || 8),
  low: Number(process.env.UNI_AUDIT_BASELINE_LOW || 9)
}

function run(name, args, options = {}) {
  const result = spawnSync(name, args, {
    shell: process.platform === 'win32',
    encoding: 'utf8',
    maxBuffer: 20 * 1024 * 1024,
    ...options
  })

  if (result.error) {
    throw result.error
  }

  return result
}

function normalizeJson(text) {
  return text.replace(/^\uFEFF/, '').trim()
}

function fail(message, code = 1) {
  console.error(`[security-gate] ${message}`)
  process.exit(code)
}

console.log('[security-gate] Step 1/2: build h5')
const buildResult = run(npmCmd, ['run', 'build:h5'], { stdio: 'inherit' })
if (buildResult.status !== 0) {
  fail(`build failed with exit code ${buildResult.status}`, 10)
}

console.log('[security-gate] Step 2/2: audit json')
const auditResult = run(npmCmd, ['audit', '--json'])
if (!auditResult.stdout) {
  fail('audit returned empty stdout, unable to parse vulnerability metadata', 11)
}

let audit
try {
  audit = JSON.parse(normalizeJson(auditResult.stdout))
} catch (error) {
  fail(`failed to parse audit json: ${error.message}`, 12)
}

const meta = audit?.metadata?.vulnerabilities
if (!meta) {
  fail('audit metadata.vulnerabilities is missing', 13)
}

const current = {
  total: Number(meta.total || 0),
  high: Number(meta.high || 0),
  moderate: Number(meta.moderate || 0),
  low: Number(meta.low || 0)
}

console.log(
  `[security-gate] baseline total=${baseline.total}, high=${baseline.high}, moderate=${baseline.moderate}, low=${baseline.low}`
)
console.log(
  `[security-gate] current  total=${current.total}, high=${current.high}, moderate=${current.moderate}, low=${current.low}`
)

const violations = []
for (const key of ['total', 'high', 'moderate', 'low']) {
  if (current[key] > baseline[key]) {
    violations.push(`${key}: ${current[key]} > ${baseline[key]}`)
  }
}

if (violations.length > 0) {
  fail(`gate failed, vulnerability count increased (${violations.join(', ')})`, 20)
}

console.log('[security-gate] gate passed')
