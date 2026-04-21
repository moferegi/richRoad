/**
 * Cloudflare Worker - 付费视频防盗链签名验证
 *
 * Worker 环境变量（Settings  Variables and Secrets）：
 *   SIGN_KEY  (Secret) = config.yaml 中 hotlink.sign-key 的值
 *   B2_ORIGIN (Text)   = https://s3.us-east-005.backblazeb2.com/Moffuu
 *
 * 签名算法：sign = HMAC-SHA256(SIGN_KEY, path + ":" + t_hex)
 */

const PROTECTED_EXTENSIONS = ['.mp4', '.mov', '.webm', '.m3u8', '.ts']
const PUBLIC_EXTENSIONS = ['.jpg', '.jpeg', '.png', '.webp', '.gif', '.svg']

export default {
  async fetch(request, env) {
    const url = new URL(request.url)
    const path = url.pathname
    const ext = path.substring(path.lastIndexOf('.')).toLowerCase()

    // 公开资源：从 B2 源站直接拉取（video.mnmovie.icu 无独立源站，不能 fetch(request) 回自身）
    if (PUBLIC_EXTENSIONS.includes(ext)) {
      const b2OriginPub = (env.B2_ORIGIN || '').replace(/\/$/, '')
      if (b2OriginPub) {
        const b2BucketPub = b2OriginPub.split('/').pop()
        const b2PrefixPub = b2BucketPub ? `/file/${b2BucketPub}` : ''
        const pubPath = (b2PrefixPub && path.startsWith(b2PrefixPub))
          ? path.slice(b2PrefixPub.length) : path
        const pubOriginUrl = b2OriginPub + pubPath
        const pubHeaders = new Headers()
        const keepPub = ['range', 'if-modified-since', 'if-none-match']
        keepPub.forEach(k => { if (request.headers.has(k)) pubHeaders.set(k, request.headers.get(k)) })
        const pubResp = await fetch(pubOriginUrl, { method: 'GET', headers: pubHeaders })
        const pubRespHeaders = new Headers(pubResp.headers)
        pubRespHeaders.set('Access-Control-Allow-Origin', '*')
        pubRespHeaders.set('Cache-Control', 'public, max-age=86400')
        return new Response(pubResp.body, {
          status: pubResp.status,
          statusText: pubResp.statusText,
          headers: pubRespHeaders,
        })
      }
      return fetch(request)
    }

    // 受保护资源需要验证签名
    if (PROTECTED_EXTENSIONS.includes(ext)) {
      const sign = url.searchParams.get('sign')
      const t = url.searchParams.get('t')

      if (!sign || !t) {
        return new Response('Forbidden: missing signature', {
          status: 403,
          headers: { 'X-Worker-Error': 'missing-signature' }
        })
      }

      const deadline = parseInt(t, 16)
      const now = Math.floor(Date.now() / 1000)
      if (isNaN(deadline) || now > deadline) {
        return new Response(`Forbidden: URL expired (deadline=${deadline}, now=${now})`, {
          status: 403,
          headers: { 'X-Worker-Error': 'url-expired' }
        })
      }

      const signKey = env.SIGN_KEY
      if (!signKey) {
        return new Response('Forbidden: SIGN_KEY not configured', {
          status: 403,
          headers: { 'X-Worker-Error': 'no-sign-key' }
        })
      }

      // B2 CDN 会在路径前加 /file/{bucket}/ 前缀，但 Go 签名时用的是不带前缀的原始路径
      // 例：Worker 收到 /file/Moffuu/video.mp4，Go 签名的是 /video.mp4
      const b2Origin = (env.B2_ORIGIN || '').replace(/\/$/, '')
      const b2BucketName = b2Origin ? b2Origin.split('/').pop() : ''
      const b2Prefix = b2BucketName ? `/file/${b2BucketName}` : ''
      const signPath = (b2Prefix && path.startsWith(b2Prefix))
        ? path.slice(b2Prefix.length)
        : path

      const valid = await verifyHmac(signKey, signPath + ':' + t, sign)
      if (!valid) {
        return new Response('Forbidden: invalid signature', {
          status: 403,
          headers: { 'X-Worker-Error': 'invalid-signature' }
        })
      }

      // 签名有效：从 B2 源站直接拉取
      // signPath 已去掉 /file/{bucket} 前缀，直接拼 B2_ORIGIN
      let originUrl
      if (b2Origin) {
        originUrl = b2Origin + signPath
      } else {
        const cleanUrl = new URL(request.url)
        cleanUrl.searchParams.delete('sign')
        cleanUrl.searchParams.delete('t')
        originUrl = cleanUrl.toString()
      }

      const proxyHeaders = new Headers()
      const keep = ['range', 'if-range', 'if-modified-since', 'if-none-match']
      keep.forEach(k => { if (request.headers.has(k)) proxyHeaders.set(k, request.headers.get(k)) })

      const resp = await fetch(originUrl, {
        method: request.method,
        headers: proxyHeaders,
      })

      const respHeaders = new Headers(resp.headers)
      respHeaders.set('Access-Control-Allow-Origin', '*')
      respHeaders.set('Cache-Control', 'private, no-store')

      return new Response(resp.body, {
        status: resp.status,
        statusText: resp.statusText,
        headers: respHeaders,
      })
    }

    return fetch(request)
  }
}

async function verifyHmac(key, msg, sign) {
  const expected = await computeHmac(key, msg)
  return expected === sign
}

async function computeHmac(key, msg) {
  const encoder = new TextEncoder()
  const cryptoKey = await crypto.subtle.importKey(
    'raw',
    encoder.encode(key),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign']
  )
  const buf = await crypto.subtle.sign('HMAC', cryptoKey, encoder.encode(msg))
  return Array.from(new Uint8Array(buf))
    .map(b => b.toString(16).padStart(2, '0'))
    .join('')
}