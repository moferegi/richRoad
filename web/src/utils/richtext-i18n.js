const richTextLikeEmptyPatterns = [
  /^\s*$/,
  /^<p><br><\/p>$/i,
  /^<p>\s*(<br\s*\/?>|&nbsp;|&#160;|\s)*<\/p>$/i,
  /^(<br\s*\/?>|&nbsp;|&#160;|\s)+$/i,
]

const structureTrackTags = new Set([
  'a',
  'img',
  'video',
  'audio',
  'iframe',
  'table',
  'thead',
  'tbody',
  'tfoot',
  'tr',
  'td',
  'th',
  'ul',
  'ol',
  'li',
  'h1',
  'h2',
  'h3',
  'h4',
  'h5',
  'h6',
  'blockquote',
  'pre',
  'code',
])

const normalizeRichTextValue = (value) => {
  return String(value || '')
    .replace(/<!--[\s\S]*?-->/g, '')
    .trim()
}

const hasMeaningfulRichTextContent = (value) => {
  const raw = normalizeRichTextValue(value)
  if (!raw) return false
  if (richTextLikeEmptyPatterns.some((pattern) => pattern.test(raw))) return false

  const plain = raw
    .replace(/<[^>]+>/g, '')
    .replace(/&nbsp;|&#160;/gi, ' ')
    .trim()

  if (plain.length > 0) return true
  return /<(img|video|audio|iframe|table|ul|ol|li|blockquote|code|pre|a|h[1-6])\b/i.test(raw)
}

const collectTrackedTags = (node, collector) => {
  if (!node?.children?.length) return
  Array.from(node.children).forEach((child) => {
    const tag = String(child.tagName || '').toLowerCase()
    if (structureTrackTags.has(tag)) {
      collector.push(tag)
    }
    collectTrackedTags(child, collector)
  })
}

const buildTagCountSignature = (tags) => {
  const counter = {}
  tags.forEach((tag) => {
    counter[tag] = (counter[tag] || 0) + 1
  })

  return Object.entries(counter)
    .sort(([left], [right]) => left.localeCompare(right))
    .map(([tag, count]) => `${tag}:${count}`)
    .join('|')
}

const buildRichTextStructureSignature = (value) => {
  const raw = normalizeRichTextValue(value)
  if (!hasMeaningfulRichTextContent(raw)) {
    return {
      empty: true,
      skip: false,
      tagsSequenceSignature: '',
      tagsCountSignature: '',
      anchorCount: 0,
      mediaCount: 0,
      listCount: 0,
      tableCount: 0,
    }
  }

  if (typeof DOMParser === 'undefined') {
    return {
      empty: false,
      skip: true,
      tagsSequenceSignature: '',
      tagsCountSignature: '',
      anchorCount: 0,
      mediaCount: 0,
      listCount: 0,
      tableCount: 0,
    }
  }

  const parser = new DOMParser()
  const doc = parser.parseFromString(`<article>${raw}</article>`, 'text/html')
  const root = doc.body?.firstElementChild
  if (!root) {
    return {
      empty: false,
      skip: false,
      tagsSequenceSignature: '',
      tagsCountSignature: '',
      anchorCount: 0,
      mediaCount: 0,
      listCount: 0,
      tableCount: 0,
    }
  }

  const tags = []
  collectTrackedTags(root, tags)

  return {
    empty: false,
    skip: false,
    tagsSequenceSignature: tags.join('>'),
    tagsCountSignature: buildTagCountSignature(tags),
    anchorCount: root.querySelectorAll('a').length,
    mediaCount: root.querySelectorAll('img,video,audio,iframe').length,
    listCount: root.querySelectorAll('ul,ol').length,
    tableCount: root.querySelectorAll('table').length,
  }
}

export const validateRichTextI18nStructure = (model, options = {}) => {
  const fieldLabel = String(options.fieldLabel || '富文本')
  const preferredBaseLang = String(options.preferredBaseLang || 'zh').trim()

  const entries = Object.entries(model || {}).map(([lang, value]) => ({
    lang: String(lang || '').trim(),
    value: String(value || ''),
  }))

  const meaningfulEntries = entries.filter((entry) => entry.lang && hasMeaningfulRichTextContent(entry.value))
  if (meaningfulEntries.length <= 1) {
    return { valid: true, message: '' }
  }

  const baseEntry = meaningfulEntries.find((entry) => entry.lang === preferredBaseLang) || meaningfulEntries[0]
  const baseSignature = buildRichTextStructureSignature(baseEntry.value)

  if (baseSignature.skip) {
    return { valid: true, message: '' }
  }

  const mismatchLangs = []
  meaningfulEntries.forEach((entry) => {
    if (entry.lang === baseEntry.lang) return

    const currentSignature = buildRichTextStructureSignature(entry.value)
    if (currentSignature.skip) return

    const structureMatched =
      currentSignature.tagsSequenceSignature === baseSignature.tagsSequenceSignature &&
      currentSignature.tagsCountSignature === baseSignature.tagsCountSignature &&
      currentSignature.anchorCount === baseSignature.anchorCount &&
      currentSignature.mediaCount === baseSignature.mediaCount &&
      currentSignature.listCount === baseSignature.listCount &&
      currentSignature.tableCount === baseSignature.tableCount

    if (!structureMatched) {
      mismatchLangs.push(entry.lang)
    }
  })

  if (!mismatchLangs.length) {
    return { valid: true, message: '' }
  }

  return {
    valid: false,
    message: `${fieldLabel}多语言结构校验失败：以下语种与基准语种(${baseEntry.lang})的关键HTML结构不一致：${mismatchLangs.join(', ')}。请检查翻译后是否误删或新增了链接、列表、表格或媒体标签。`,
  }
}
