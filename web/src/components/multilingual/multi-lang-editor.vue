<template>
  <div class="multi-lang-editor">
    <div class="multi-lang-editor__toolbar">
      <div class="multi-lang-editor__title-wrap">
        <span class="multi-lang-editor__title">{{ title }}</span>
        <el-tag size="small" type="info">已填 {{ filledCount }}/{{ languageCount }}</el-tag>
      </div>
      <div class="multi-lang-editor__source-wrap">
        <span class="multi-lang-editor__source-label">源语言</span>
        <el-select v-model="activeSourceLang" size="small" class="multi-lang-editor__source-select">
          <el-option
            v-for="lang in validLanguages"
            :key="lang.code"
            :label="(lang.name || lang.nativeName || lang.code) + ' (' + lang.code + ')'"
            :value="lang.code"
          />
        </el-select>
      </div>
      <div class="multi-lang-editor__actions">
        <el-button size="small" @click="fillMissingSlots">一键补齐</el-button>
        <el-button size="small" @click="copySourceToAll">批量复制到全部</el-button>
        <el-button size="small" type="primary" :loading="translating && translatingMode === 'missing'" @click="translateMissingSlots">
          一键翻译空白
        </el-button>
        <el-button size="small" type="warning" :loading="translating && translatingMode === 'overwrite'" @click="translateAllSlots">
          覆盖翻译
        </el-button>
      </div>
    </div>

    <div v-if="validLanguages.length">
      <el-tabs v-if="useTabs" v-model="activeEditLang" :type="tabType" class="multi-lang-editor__tabs">
        <el-tab-pane
          v-for="lang in validLanguages"
          :key="lang.code"
          :label="lang.name || lang.nativeName || lang.code"
          :name="lang.code"
          :lazy="inputType === 'richtext'"
        >
          <div class="multi-lang-editor__pane">
            <div class="multi-lang-editor__item-head">
              <el-tag size="small">{{ lang.code }}</el-tag>
              <span class="multi-lang-editor__lang-name">{{ lang.name || lang.nativeName || lang.code }}</span>
            </div>

            <slot name="editor" :lang="lang" :model="model">
              <el-input
                v-if="inputType === 'input'"
                v-model="model[lang.code]"
                :placeholder="(lang.name || lang.code) + ' 文案'"
                clearable
              />

              <RichEdit
                v-else-if="inputType === 'richtext'"
                v-model="model[lang.code]"
                :upload-folder="richUploadFolder"
                class="multi-lang-editor__rich-edit"
              />

              <el-input
                v-else
                v-model="model[lang.code]"
                type="textarea"
                :rows="rows"
                :placeholder="(lang.name || lang.code) + ' 文案'"
              />

              <div v-if="inputType === 'richtext'" class="multi-lang-editor__rich-tip">
                支持富文本编辑，一键翻译后请复核格式与换行。
              </div>
            </slot>
          </div>
        </el-tab-pane>
      </el-tabs>

      <div v-else class="multi-lang-editor__grid">
        <div v-for="lang in validLanguages" :key="lang.code" class="multi-lang-editor__item">
          <div class="multi-lang-editor__item-head">
            <el-tag size="small">{{ lang.code }}</el-tag>
            <span class="multi-lang-editor__lang-name">{{ lang.name || lang.nativeName || lang.code }}</span>
          </div>

          <slot name="editor" :lang="lang" :model="model">
            <el-input
              v-if="inputType === 'input'"
              v-model="model[lang.code]"
              :placeholder="(lang.name || lang.code) + ' 文案'"
              clearable
            />

            <RichEdit
              v-else-if="inputType === 'richtext'"
              v-model="model[lang.code]"
              :upload-folder="richUploadFolder"
              class="multi-lang-editor__rich-edit"
            />

            <el-input
              v-else
              v-model="model[lang.code]"
              type="textarea"
              :rows="rows"
              :placeholder="(lang.name || lang.code) + ' 文案'"
            />

            <div v-if="inputType === 'richtext'" class="multi-lang-editor__rich-tip">
              支持富文本编辑，一键翻译后请复核格式与换行。
            </div>
          </slot>
        </div>
      </div>
    </div>

    <div v-else class="multi-lang-editor__empty">请先在语言管理中配置语言</div>
  </div>
</template>

<script setup>
import { computed, defineAsyncComponent, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getLanguageList, translateI18n } from '@/api/client/language'

const RichEdit = defineAsyncComponent(() => import('@/components/richtext/rich-edit.vue'))

const props = defineProps({
  model: {
    type: Object,
    required: true,
  },
  languages: {
    type: Array,
    default: () => [],
  },
  sourceLang: {
    type: String,
    default: 'zh',
  },
  title: {
    type: String,
    default: '多语言内容',
  },
  inputType: {
    type: String,
    default: 'input',
  },
  rows: {
    type: Number,
    default: 3,
  },
  richUploadFolder: {
    type: String,
    default: '',
  },
  useTabs: {
    type: Boolean,
    default: false,
  },
  tabType: {
    type: String,
    default: 'border-card',
  },
})

const translating = ref(false)
const translatingMode = ref('')
const activeSourceLang = ref('')
const activeEditLang = ref('')
const managedLanguages = ref([])

const DEFAULT_MULTILINGUAL_LANGS = [
  { code: 'zh', name: '中文' },
  { code: 'en', name: '英文' },
  { code: 'mn', name: '蒙文' },
  { code: 'zh-TW', name: '繁体' },
  { code: 'th', name: '泰语' },
  { code: 'hi', name: '印地语' },
  { code: 'id', name: '印尼语' },
  { code: 'vi', name: '越南语' },
  { code: 'ar', name: '阿拉伯语' },
  { code: 'ja', name: '日语' },
  { code: 'ko', name: '韩语' },
  { code: 'ms', name: '马来语' },
]

const normalizeLanguages = (list = []) => {
  return list
    .map((lang) => ({
      code: String(lang?.code || '').trim(),
      name: String(lang?.name || '').trim(),
      nativeName: String(lang?.nativeName || '').trim(),
    }))
    .filter((lang) => lang.code)
}

const mergeLanguageSources = (...sources) => {
  const map = new Map()
  sources.forEach((source) => {
    normalizeLanguages(source).forEach((lang) => {
      const existing = map.get(lang.code)
      if (!existing) {
        map.set(lang.code, lang)
        return
      }
      map.set(lang.code, {
        ...existing,
        name: existing.name || lang.name,
        nativeName: existing.nativeName || lang.nativeName,
      })
    })
  })
  return Array.from(map.values())
}

const validLanguages = computed(() => {
  return mergeLanguageSources(
    managedLanguages.value,
    props.languages,
    DEFAULT_MULTILINGUAL_LANGS
  )
})

const richTextLikeEmptyPatterns = [
  /^<p><br><\/p>$/i,
  /^<p>\s*<\/p>$/i,
  /^<div><br><\/div>$/i,
]

const hasMeaningfulContent = (value) => {
  const raw = String(value ?? '').trim()
  if (!raw) {
    return false
  }

  if (props.inputType !== 'richtext') {
    return !!raw
  }

  if (richTextLikeEmptyPatterns.some((pattern) => pattern.test(raw))) {
    return false
  }

  const plain = raw
    .replace(/<style[\s\S]*?<\/style>/gi, ' ')
    .replace(/<script[\s\S]*?<\/script>/gi, ' ')
    .replace(/<[^>]+>/g, ' ')
    .replace(/&nbsp;/gi, ' ')
    .replace(/\s+/g, ' ')
    .trim()

  return !!plain
}

const languageCount = computed(() => validLanguages.value.length)

const filledCount = computed(() => {
  if (!props.model || typeof props.model !== 'object') return 0
  return validLanguages.value.reduce((count, lang) => {
    return hasMeaningfulContent(props.model[lang.code]) ? count + 1 : count
  }, 0)
})

const fillMissingSlots = () => {
  if (!props.model || typeof props.model !== 'object') return
  let filled = 0
  validLanguages.value.forEach((lang) => {
    if (props.model[lang.code] === undefined || props.model[lang.code] === null) {
      props.model[lang.code] = ''
      filled++
    }
  })
  if (filled > 0) {
    ElMessage.success(`已补齐 ${filled} 个缺失语言槽位`)
  } else {
    ElMessage.info('暂无缺失语言槽位')
  }
}

const loadManagedLanguages = async () => {
  try {
    const res = await getLanguageList({ page: 1, pageSize: 500 })
    const list = Array.isArray(res?.data) ? res.data : (res?.data?.list || [])
    managedLanguages.value = normalizeLanguages(list)
  } catch (error) {
    managedLanguages.value = []
  }
}

onMounted(() => {
  loadManagedLanguages()
})

watch(
  () => [props.sourceLang, validLanguages.value.map((lang) => lang.code).join(',')],
  () => {
    const source = String(props.sourceLang || '').trim()
    const codes = validLanguages.value.map((lang) => lang.code)
    if (source && codes.includes(source)) {
      activeSourceLang.value = source
      return
    }
    if (!codes.includes(activeSourceLang.value)) {
      activeSourceLang.value = codes[0] || ''
    }
  },
  { immediate: true }
)

watch(
  () => validLanguages.value.map((lang) => lang.code).join(','),
  () => {
    const codes = validLanguages.value.map((lang) => lang.code)
    if (!codes.includes(activeEditLang.value)) {
      activeEditLang.value = codes[0] || ''
    }
  },
  { immediate: true }
)

const resolveSourceText = () => {
  const sourceCode = String(activeSourceLang.value || '').trim()
  if (!sourceCode) return null

  const sourceValue = String(props.model?.[sourceCode] ?? '')
  if (!hasMeaningfulContent(sourceValue)) return null

  return {
    text: sourceValue,
    source: sourceCode,
  }
}

const copySourceToAll = () => {
  if (!props.model || typeof props.model !== 'object') return
  const sourceInfo = resolveSourceText()
  if (!sourceInfo) {
    ElMessage.warning('请先在源语言填写内容')
    return
  }

  validLanguages.value
    .map((lang) => lang.code)
    .filter((code) => code !== sourceInfo.source)
    .forEach((code) => {
      props.model[code] = sourceInfo.text
    })

  ElMessage.success('已按源语言复制到全部')
}

const translateByMode = async (overwrite = false) => {
  if (!props.model || typeof props.model !== 'object') return
  if (!validLanguages.value.length) {
    ElMessage.warning('暂无可翻译语言')
    return
  }

  const sourceInfo = resolveSourceText()
  if (!sourceInfo) {
    ElMessage.warning('请先在源语言填写内容')
    return
  }

  const targets = validLanguages.value
    .map((lang) => lang.code)
    .filter((code) => code !== sourceInfo.source)
    .filter((code) => overwrite || !hasMeaningfulContent(props.model[code]))

  if (!targets.length) {
    ElMessage.info(overwrite ? '没有可覆盖翻译的语言' : '没有需要翻译的空白语言')
    return
  }

  translating.value = true
  translatingMode.value = overwrite ? 'overwrite' : 'missing'
  try {
    const res = await translateI18n({
      text: sourceInfo.text,
      source: sourceInfo.source,
      targets,
    })

    const translations = res?.data?.translations || {}
    targets.forEach((code) => {
      const text = String(translations[code] || '').trim()
      if (text) {
        props.model[code] = text
      }
    })

    ElMessage.success(overwrite ? '覆盖翻译完成' : '翻译完成')
  } catch (error) {
    ElMessage.error('翻译失败，请稍后重试')
  } finally {
    translating.value = false
    translatingMode.value = ''
  }
}

const translateMissingSlots = async () => {
  await translateByMode(false)
}

const translateAllSlots = async () => {
  await translateByMode(true)
}
</script>

<style scoped>
.multi-lang-editor {
  width: 100%;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  padding: 10px;
  background: var(--el-bg-color);
}

.multi-lang-editor__toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 10px;
}

.multi-lang-editor__title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.multi-lang-editor__title {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.multi-lang-editor__actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.multi-lang-editor__source-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
}

.multi-lang-editor__source-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.multi-lang-editor__source-select {
  width: 190px;
}

.multi-lang-editor__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 10px;
}

.multi-lang-editor__tabs {
  width: 100%;
}

.multi-lang-editor__tabs :deep(.el-tabs__content) {
  padding: 0;
}

.multi-lang-editor__pane {
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
  padding: 10px;
}

.multi-lang-editor__item {
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
  padding: 8px;
}

.multi-lang-editor__item-head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.multi-lang-editor__lang-name {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.multi-lang-editor__rich-edit {
  width: 100%;
}

.multi-lang-editor__rich-tip {
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

.multi-lang-editor__empty {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}
</style>
