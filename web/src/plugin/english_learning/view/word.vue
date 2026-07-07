<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="英语学习运营台（单词端）：支持分类、章节、单词的完整增删改查与资源上传。"
        type="info"
        :closable="false"
        show-icon
      />
      <div class="lang-switch-row">
        <span>展示语言：</span>
        <el-select v-model="displayLang" style="width: 140px">
          <el-option
            v-for="lang in displayLanguageOptions"
            :key="lang.value"
            :label="lang.label"
            :value="lang.value"
          />
        </el-select>
      </div>
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="分类管理" name="category">
          <div class="toolbar-row">
            <el-form :inline="true" :model="categoryQuery">
              <el-form-item label="会员分类">
                <el-select v-model="categoryQuery.needVip" clearable placeholder="全部" style="width: 160px">
                  <el-option :value="true" label="仅会员" />
                  <el-option :value="false" label="非会员" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleCategorySearch">查询</el-button>
                <el-button icon="refresh" @click="resetCategorySearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openCategoryDialog()">新增分类</el-button>
          </div>

          <el-table :data="categoryTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="分类名" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column label="Logo" min-width="220" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.logo || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="price" label="价格" width="120" />
            <el-table-column label="会员限制" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.needVip ? 'warning' : 'success'">{{ scope.row.needVip ? '需要会员' : '普通可见' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" fixed="right" width="180">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openCategoryDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeCategory(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="categoryQuery.page"
              :page-size="categoryQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="categoryTotal"
              @current-change="handleCategoryPageChange"
              @size-change="handleCategorySizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="章节管理" name="chapter">
          <div class="toolbar-row">
            <el-form :inline="true" :model="chapterQuery">
              <el-form-item label="所属分类">
                <el-select v-model="chapterQuery.categoryId" clearable placeholder="全部" style="width: 240px">
                  <el-option
                    v-for="item in categoryOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleChapterSearch">查询</el-button>
                <el-button icon="refresh" @click="resetChapterSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openChapterDialog()">新增章节</el-button>
          </div>

          <el-table :data="chapterTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属分类" min-width="200">
              <template #default="scope">
                {{ categoryNameMap[scope.row.categoryId] || `#${scope.row.categoryId}` }}
              </template>
            </el-table-column>
            <el-table-column label="章节名" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" fixed="right" width="180">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openChapterDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeChapter(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="chapterQuery.page"
              :page-size="chapterQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="chapterTotal"
              @current-change="handleChapterPageChange"
              @size-change="handleChapterSizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="单词管理" name="word">
          <div class="toolbar-row">
            <el-form :inline="true" :model="wordQuery">
              <el-form-item label="分类过滤">
                <el-select v-model="wordQuery.categoryId" clearable placeholder="全部" style="width: 220px">
                  <el-option
                    v-for="item in categoryOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="章节过滤">
                <el-select v-model="wordQuery.chapterId" clearable placeholder="全部" style="width: 260px">
                  <el-option
                    v-for="item in filteredWordChapterOptions"
                    :key="item.ID"
                    :label="`${categoryNameMap[item.categoryId] || '-'} / ${formatI18nText(item.name)}`"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="单词关键词">
                <el-input v-model="wordQuery.keyword" clearable placeholder="支持模糊查询" style="width: 220px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleWordSearch">查询</el-button>
                <el-button icon="refresh" @click="resetWordSearch">重置</el-button>
              </el-form-item>
            </el-form>
              <div class="toolbar-actions">
                <el-button type="warning" icon="upload" @click="openSqlImportDialog">批量SQL导入</el-button>
                <el-button type="primary" icon="plus" @click="openWordDialog()">新增单词</el-button>
              </div>
          </div>

          <el-table :data="wordTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="word" label="单词" min-width="150" />
            <el-table-column prop="phoneticUs" label="美式音标" min-width="160" show-overflow-tooltip />
            <el-table-column prop="phoneticUk" label="英式音标" min-width="160" show-overflow-tooltip />
            <el-table-column label="释义" min-width="260" show-overflow-tooltip>
              <template #default="scope">
                {{ formatI18nText(scope.row.explanation) }}
              </template>
            </el-table-column>
            <el-table-column label="发音" width="180">
              <template #default="scope">
                <el-tag :type="scope.row.audioUs ? 'success' : 'warning'">美音{{ scope.row.audioUs ? '已上传' : '缺失' }}</el-tag>
                <el-tag :type="scope.row.audioUk ? 'success' : 'warning'" class="ml-2">英音{{ scope.row.audioUk ? '已上传' : '缺失' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="260">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openWordDialog(scope.row)">编辑</el-button>
                <el-button type="warning" link icon="refresh" @click="regenerateWord(scope.row)">重生发音</el-button>
                <el-button
                  type="danger"
                  link
                  icon="delete"
                  :disabled="deletingWordIds.has(getEntityID(scope.row))"
                  @click="removeWord(scope.row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="wordQuery.page"
              :page-size="wordQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="wordTotal"
              @current-change="handleWordPageChange"
              @size-change="handleWordSizeChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="categoryDialogVisible" :title="categoryDialogMode === 'create' ? '新增分类' : '编辑分类'" width="760px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="categoryForm" label-width="120px">
        <el-form-item label="分类名称">
          <MultiLangEditor
            :model="categoryForm.nameI18n"
            title="分类名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="Logo地址">
          <div class="upload-inline">
            <el-input v-model="categoryForm.logo" placeholder="上传后自动写入 URL" />
            <el-upload
              :show-file-list="false"
              :http-request="(options) => uploadByRequest(options, 'english-learn/pic', (url) => { categoryForm.logo = url }, 'Logo')"
            >
              <el-button type="primary" plain>上传Logo</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="categoryForm.price" :min="0" :precision="2" :step="1" />
        </el-form-item>
        <el-form-item label="需要会员">
          <el-switch v-model="categoryForm.needVip" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="categoryForm.sort" :min="0" :step="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="chapterDialogVisible" :title="chapterDialogMode === 'create' ? '新增章节' : '编辑章节'" width="760px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="chapterForm" label-width="120px">
        <el-form-item label="所属分类">
          <el-select v-model="chapterForm.categoryId" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in categoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="章节名称">
          <MultiLangEditor
            :model="chapterForm.nameI18n"
            title="章节名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="chapterForm.sort" :min="0" :step="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="chapterDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitChapter">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="wordDialogVisible" :title="wordDialogMode === 'create' ? '新增单词' : '编辑单词'" width="860px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="wordForm" label-width="130px">
        <el-form-item label="单词本体">
          <el-input v-model="wordForm.word" placeholder="例如: destiny" />
        </el-form-item>
        <el-form-item label="美式音标">
          <el-input v-model="wordForm.phoneticUs" placeholder="例如: /ˈdestəni/" />
        </el-form-item>
        <el-form-item label="英式音标">
          <el-input v-model="wordForm.phoneticUk" placeholder="例如: /ˈdestəni/" />
        </el-form-item>
        <el-form-item label="释义">
          <MultiLangEditor
            :model="wordForm.explanationI18n"
            title="单词释义多语言"
            input-type="textarea"
            :rows="4"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="美式发音URL">
          <div class="upload-inline">
            <el-input v-model="wordForm.audioUs" placeholder="可上传覆盖，留空则尝试自动生成" />
            <el-upload
              :show-file-list="false"
              :http-request="(options) => uploadByRequest(options, 'english-learn/audio', (url) => { wordForm.audioUs = url }, '美式音频')"
            >
              <el-button type="primary" plain>上传美式</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="英式发音URL">
          <div class="upload-inline">
            <el-input v-model="wordForm.audioUk" placeholder="可上传覆盖，留空则尝试自动生成" />
            <el-upload
              :show-file-list="false"
              :http-request="(options) => uploadByRequest(options, 'english-learn/audio', (url) => { wordForm.audioUk = url }, '英式音频')"
            >
              <el-button type="primary" plain>上传英式</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="造句列表">
          <div class="sentence-editor">
            <div class="sentence-toolbar">
              <el-button type="primary" plain size="small" @click="appendWordSentence">新增造句</el-button>
            </div>
            <div v-if="wordForm.sentences.length === 0" class="dialog-hint">暂无造句，建议至少添加1条，便于跟打页展示。</div>
            <div v-for="(sentence, sentenceIndex) in wordForm.sentences" :key="sentence.__key" class="sentence-card">
              <div class="sentence-card-header">
                <span>造句 {{ sentenceIndex + 1 }}</span>
                <el-button type="danger" link @click="removeWordSentence(sentenceIndex)">删除</el-button>
              </div>
              <el-form-item label="英文句子" label-width="100px">
                <el-input v-model="sentence.source" placeholder="例如: The bowl is on the table." />
              </el-form-item>
              <el-form-item label="翻译" label-width="100px">
                <MultiLangEditor
                  :model="sentence.translateI18n"
                  title="造句翻译多语言"
                  input-type="textarea"
                  :rows="3"
                  :use-tabs="true"
                />
              </el-form-item>
              <el-form-item label="美式音频" label-width="100px">
                <el-input v-model="sentence.audioUs" placeholder="留空则后端自动生成" />
              </el-form-item>
              <el-form-item label="英式音频" label-width="100px">
                <el-input v-model="sentence.audioUk" placeholder="留空则后端自动生成" />
              </el-form-item>
              <el-form-item label="排序" label-width="100px">
                <el-input-number v-model="sentence.sort" :min="0" :step="1" />
              </el-form-item>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="归属分类">
          <el-select v-model="wordForm.categoryIds" multiple filterable clearable placeholder="可选多个分类" style="width: 100%">
            <el-option
              v-for="item in categoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="绑定章节">
          <el-select v-model="wordForm.chapterIds" multiple filterable clearable placeholder="章节可不选" style="width: 100%">
            <el-option
              v-for="item in filteredWordBindChapterOptions"
              :key="item.ID"
              :label="`${categoryNameMap[item.categoryId] || '-'} / ${formatI18nText(item.name)}`"
              :value="item.ID"
            />
          </el-select>
          <div class="dialog-hint">单词至少需要绑定一个分类或章节。章节属于可选维度，不强制。</div>
        </el-form-item>
        <el-form-item label="更新绑定" v-if="wordDialogMode === 'edit'">
          <el-checkbox v-model="wordForm.syncBindingsOnUpdate">
            按上面的分类/章节重建绑定关系（不勾选则保持原关系）
          </el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="wordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="wordSubmitting" @click="submitWord">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="sqlImportDialogVisible" title="批量SQL导入单词" width="980px" :close-on-click-modal="false" :close-on-press-escape="false">
      <el-form :model="sqlImportForm" label-width="140px">
        <el-form-item label="目标分类">
          <el-select v-model="sqlImportForm.categoryId" placeholder="请选择分类" style="width: 100%" filterable>
            <el-option
              v-for="item in categoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="SQL文件">
          <el-upload
            :show-file-list="false"
            :auto-upload="false"
            accept=".sql,text/plain"
            :on-change="handleSqlFileChange"
          >
            <el-button type="primary" plain>选择SQL文件</el-button>
          </el-upload>
          <div class="dialog-hint">支持 INSERT INTO junior (word,translate) VALUES ('word','中文'); 格式。</div>
        </el-form-item>
        <el-form-item label="SQL内容">
          <el-input
            v-model="sqlImportForm.sqlContent"
            type="textarea"
            :rows="8"
            placeholder="可粘贴SQL文本，或通过上方选择文件自动填充"
          />
        </el-form-item>
        <el-form-item label="自动分章节">
          <el-switch v-model="sqlImportForm.autoGenerateChapters" />
        </el-form-item>
        <el-form-item v-if="!sqlImportForm.autoGenerateChapters" label="指定章节（可选）">
          <el-select v-model="sqlImportForm.chapterId" clearable placeholder="不选则仅绑定分类" style="width: 100%" filterable>
            <el-option
              v-for="item in sqlImportChapterOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-else label="每章单词数">
          <el-input-number v-model="sqlImportForm.wordsPerChapter" :min="1" :step="1" />
          <div class="dialog-hint">将按导入顺序自动创建并分配到“第n章”。</div>
        </el-form-item>
        <el-form-item label="自动补发音">
          <el-switch v-model="sqlImportForm.generateAudio" />
        </el-form-item>
      </el-form>

      <div class="import-progress-box" v-if="sqlImportLogs.length > 0">
        <div class="import-progress-head">
          <span>导入进度日志（逐行）</span>
          <el-tag type="info">{{ sqlImportLogs.length }} 条</el-tag>
        </div>
        <div class="import-progress-list">
          <div v-for="(item, index) in sqlImportLogs" :key="`${item.time}_${index}`" class="import-progress-line">
            <el-tag :type="item.level === 'success' ? 'success' : item.level === 'error' ? 'danger' : 'warning'" size="small">{{ item.level }}</el-tag>
            <span class="import-progress-time">{{ item.time }}</span>
            <span class="import-progress-text">{{ item.message }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button :disabled="sqlImportSubmitting" @click="sqlImportDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="sqlImportSubmitting" @click="runSqlImport">开始导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
  import { uploadFile } from '@/api/fileUploadAndDownload'
  import { getLanguageList } from '@/api/client/language'
  import {
    createCategory,
    createChapter,
    createEnglishWord,
    deleteCategory,
    deleteChapter,
    deleteEnglishWord,
    findEnglishWord,
    getCategoryList,
    getChapterList,
    getEnglishWordList,
    regenerateWordAudio,
    upsertSqlWord,
    updateCategory,
    updateChapter,
    updateEnglishWord
  } from '../api/english'

  defineOptions({
    name: 'EnglishLearningWord'
  })

  const activeTab = ref('category')
  const displayLang = ref('zh')
  const managedLanguages = ref([])

  const displayLanguageOptions = computed(() => {
    if (managedLanguages.value.length === 0) {
      return [
        { value: 'zh', label: '中文' },
        { value: 'en', label: 'English' },
        { value: 'mn', label: 'Монгол' }
      ]
    }
    return managedLanguages.value.map((lang) => ({
      value: lang.code,
      label: lang.nativeName || lang.name || lang.code
    }))
  })

  const categoryQuery = ref({
    page: 1,
    pageSize: 10,
    needVip: undefined
  })
  const categoryTable = ref([])
  const categoryTotal = ref(0)

  const chapterQuery = ref({
    page: 1,
    pageSize: 10,
    categoryId: undefined
  })
  const chapterTable = ref([])
  const chapterTotal = ref(0)

  const wordQuery = ref({
    page: 1,
    pageSize: 10,
    categoryId: undefined,
    chapterId: undefined,
    keyword: ''
  })
  const wordTable = ref([])
  const wordTotal = ref(0)
  const deletingWordIds = ref(new Set())

  const categoryOptions = ref([])
  const chapterOptions = ref([])

  const getEntityID = (item) => Number(item?.ID || item?.id || 0)

  const categoryNameMap = computed(() => {
    const map = {}
    for (const item of categoryOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const filteredWordChapterOptions = computed(() => {
    const categoryId = Number(wordQuery.value.categoryId || 0)
    if (!categoryId) {
      return chapterOptions.value
    }
    return chapterOptions.value.filter((item) => Number(item.categoryId || 0) === categoryId)
  })

  watch(
    () => wordQuery.value.categoryId,
    () => {
      const validChapterIDs = new Set(filteredWordChapterOptions.value.map((item) => getEntityID(item)))
      if (!validChapterIDs.has(Number(wordQuery.value.chapterId || 0))) {
        wordQuery.value.chapterId = undefined
      }
    }
  )

  const categoryDialogVisible = ref(false)
  const categoryDialogMode = ref('create')
  const categoryForm = ref({
    ID: 0,
    nameI18n: { zh: '' },
    logo: '',
    price: 0,
    needVip: false,
    sort: 0
  })

  const chapterDialogVisible = ref(false)
  const chapterDialogMode = ref('create')
  const chapterForm = ref({
    ID: 0,
    categoryId: undefined,
    nameI18n: { zh: '' },
    sort: 0
  })

  const wordDialogVisible = ref(false)
  const wordDialogMode = ref('create')
  const wordSubmitting = ref(false)
  const wordForm = ref({
    ID: 0,
    word: '',
    phoneticUs: '',
    phoneticUk: '',
    audioUs: '',
    audioUk: '',
    explanationI18n: { zh: '' },
    sentences: [],
    categoryIds: [],
    chapterIds: [],
    syncBindingsOnUpdate: false
  })

  const sqlImportDialogVisible = ref(false)
  const sqlImportSubmitting = ref(false)
  const sqlImportForm = ref({
    categoryId: undefined,
    chapterId: undefined,
    autoGenerateChapters: false,
    wordsPerChapter: 30,
    generateAudio: false,
    sqlContent: ''
  })
  const sqlImportLogs = ref([])

  let sentenceKeySeed = 1
  const createWordSentenceFormItem = () => ({
    __key: `sentence_${Date.now()}_${sentenceKeySeed++}`,
    source: '',
    translateI18n: { zh: '' },
    audioUs: '',
    audioUk: '',
    sort: 0
  })

  const appendWordSentence = () => {
    wordForm.value.sentences.push(createWordSentenceFormItem())
  }

  const removeWordSentence = (index) => {
    wordForm.value.sentences.splice(index, 1)
  }

  const nowTimeText = () => {
    const d = new Date()
    const p = (n) => String(n).padStart(2, '0')
    return `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
  }

  const appendSqlLog = (level, message) => {
    sqlImportLogs.value.push({
      level,
      message,
      time: nowTimeText()
    })
  }

  const parseSqlWordRows = (sqlContent) => {
    const rows = []
    const lines = String(sqlContent || '').split(/\r?\n/)
    const pattern = /insert\s+into\s+[^\(]*\(\s*`?word`?\s*,\s*`?translate`?\s*\)\s*values\s*\(\s*'((?:''|[^'])*)'\s*,\s*'((?:''|[^'])*)'\s*\)/i
    for (const rawLine of lines) {
      const line = String(rawLine || '').trim()
      if (!line) {
        continue
      }
      const match = line.match(pattern)
      if (!match) {
        continue
      }
      const word = String(match[1] || '').replace(/''/g, "'").trim()
      const translateZh = String(match[2] || '').replace(/''/g, "'").trim()
      if (!word) {
        continue
      }
      rows.push({ word, translateZh })
    }
    return rows
  }

  const openSqlImportDialog = () => {
    sqlImportForm.value = {
      categoryId: wordQuery.value.categoryId || undefined,
      chapterId: undefined,
      autoGenerateChapters: false,
      wordsPerChapter: 30,
      generateAudio: false,
      sqlContent: ''
    }
    sqlImportLogs.value = []
    sqlImportDialogVisible.value = true
  }

  const handleSqlFileChange = (file) => {
    const raw = file?.raw
    if (!raw) {
      return
    }
    const reader = new FileReader()
    reader.onload = () => {
      sqlImportForm.value.sqlContent = String(reader.result || '')
      appendSqlLog('warning', `已加载文件：${file.name || 'unknown.sql'}`)
    }
    reader.onerror = () => {
      ElMessage.error('读取SQL文件失败')
    }
    reader.readAsText(raw, 'utf-8')
  }

  const findChapterByZhName = (categoryId, chapterNameZh) => {
    return chapterOptions.value.find((item) => {
      return Number(item.categoryId || 0) === Number(categoryId || 0) && formatI18nText(item.name) === chapterNameZh
    })
  }

  const ensureAutoChapter = async (categoryId, chapterIndex, cache) => {
    const chapterNameZh = `第${chapterIndex}章`
    if (cache.has(chapterNameZh)) {
      return cache.get(chapterNameZh)
    }

    let chapter = findChapterByZhName(categoryId, chapterNameZh)
    if (!chapter) {
      const createRes = await createChapter({
        categoryId,
        name: JSON.stringify({ zh: chapterNameZh, en: '', mn: '' }),
        sort: chapterIndex
      })
      if (createRes.code !== 0) {
        throw new Error(`自动创建章节失败: ${chapterNameZh}`)
      }
      await loadChapterOptions()
      chapter = findChapterByZhName(categoryId, chapterNameZh)
      if (!chapter) {
        throw new Error(`章节创建后未找到: ${chapterNameZh}`)
      }
      appendSqlLog('warning', `已自动创建章节：${chapterNameZh}`)
    }

    const chapterId = getEntityID(chapter)
    cache.set(chapterNameZh, chapterId)
    return chapterId
  }

  const runSqlImport = async () => {
    if (!sqlImportForm.value.categoryId) {
      ElMessage.warning('请选择目标分类')
      return
    }
    if (!String(sqlImportForm.value.sqlContent || '').trim()) {
      ElMessage.warning('请提供SQL内容')
      return
    }
    if (sqlImportForm.value.autoGenerateChapters && Number(sqlImportForm.value.wordsPerChapter || 0) <= 0) {
      ElMessage.warning('每章单词数必须大于0')
      return
    }

    const rows = parseSqlWordRows(sqlImportForm.value.sqlContent)
    if (rows.length === 0) {
      ElMessage.warning('未解析到有效的SQL单词行')
      return
    }

    sqlImportSubmitting.value = true
    sqlImportLogs.value = []
    appendSqlLog('warning', `开始导入，共 ${rows.length} 行`)

    const chapterCache = new Map()
    let successCount = 0
    let failCount = 0

    for (let i = 0; i < rows.length; i++) {
      const row = rows[i]
      try {
        let chapterId = Number(sqlImportForm.value.chapterId || 0)
        if (sqlImportForm.value.autoGenerateChapters) {
          const chapterIndex = Math.floor(i / Number(sqlImportForm.value.wordsPerChapter || 1)) + 1
          chapterId = await ensureAutoChapter(sqlImportForm.value.categoryId, chapterIndex, chapterCache)
        }

        const res = await upsertSqlWord({
          word: row.word,
          translateZh: row.translateZh,
          categoryId: sqlImportForm.value.categoryId,
          chapterId,
          generateAudio: sqlImportForm.value.generateAudio
        })

        if (res.code !== 0) {
          failCount++
          appendSqlLog('error', `第${i + 1}行 ${row.word} 失败: ${res.msg || '接口返回异常'}`)
          continue
        }

        successCount++
        const data = res.data || {}
        const actionText = data.action === 'create' ? '新增' : '修改'
        const chapterText = chapterId ? `章节#${chapterId}` : '仅分类'
        const audioText = sqlImportForm.value.generateAudio
          ? `美音:${data.audioUsGenerated ? '补齐' : '保留'} 英音:${data.audioUkGenerated ? '补齐' : '保留'}`
          : '未启用语音补齐'
        appendSqlLog('success', `第${i + 1}行 ${row.word} ${actionText}成功，${chapterText}，${audioText}`)
      } catch (error) {
        failCount++
        appendSqlLog('error', `第${i + 1}行 ${row.word} 失败: ${error?.message || '未知错误'}`)
      }
    }

    appendSqlLog('warning', `导入完成：成功 ${successCount}，失败 ${failCount}`)
    sqlImportSubmitting.value = false
    await Promise.all([loadWordList(), loadChapterList(), loadChapterOptions()])
    if (failCount === 0) {
      ElMessage.success(`导入完成，共 ${successCount} 条`)
    } else {
      ElMessage.warning(`导入完成，成功 ${successCount} 条，失败 ${failCount} 条`)
    }
  }

  const filteredWordBindChapterOptions = computed(() => {
    const selectedCategoryIDs = Array.isArray(wordForm.value.categoryIds)
      ? wordForm.value.categoryIds.map((id) => Number(id || 0)).filter((id) => id > 0)
      : []
    if (selectedCategoryIDs.length === 0) {
      return chapterOptions.value
    }
    const allowed = new Set(selectedCategoryIDs)
    return chapterOptions.value.filter((item) => allowed.has(Number(item.categoryId || 0)))
  })

  const sqlImportChapterOptions = computed(() => {
    const categoryId = Number(sqlImportForm.value.categoryId || 0)
    if (!categoryId) {
      return []
    }
    return chapterOptions.value.filter((item) => Number(item.categoryId || 0) === categoryId)
  })

  watch(
    () => wordForm.value.categoryIds,
    () => {
      const validChapterIDs = new Set(filteredWordBindChapterOptions.value.map((item) => getEntityID(item)))
      wordForm.value.chapterIds = (wordForm.value.chapterIds || []).filter((id) => validChapterIDs.has(Number(id || 0)))
    },
    { deep: true }
  )

  const formatI18nText = (raw) => {
    if (!raw) return ''
    try {
      const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
      if (typeof parsed === 'object' && parsed !== null) {
        const lang = String(displayLang.value || 'zh')
        return parsed[lang] || parsed.zh || parsed.en || parsed.mn || Object.values(parsed)[0] || ''
      }
      return String(raw)
    } catch (e) {
      return String(raw)
    }
  }

  const normalizeI18nObject = (raw) => {
    if (!raw) {
      return { zh: '' }
    }

    if (typeof raw === 'object') {
      return Object.keys(raw).length > 0 ? { ...raw } : { zh: '' }
    }

    const text = String(raw).trim()
    if (!text) {
      return { zh: '' }
    }

    try {
      const parsed = JSON.parse(text)
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return Object.keys(parsed).length > 0 ? { ...parsed } : { zh: '' }
      }
      return { zh: text }
    } catch (e) {
      return { zh: text }
    }
  }

  const stringifyI18nObject = (i18nObject) => {
    const source = i18nObject && typeof i18nObject === 'object' ? i18nObject : {}
    const cleaned = {}
    for (const [key, value] of Object.entries(source)) {
      const lang = String(key || '').trim()
      if (!lang) {
        continue
      }
      const text = String(value ?? '').trim()
      if (text) {
        cleaned[lang] = text
      }
    }
    if (Object.keys(cleaned).length === 0) {
      return JSON.stringify({ zh: '' })
    }
    return JSON.stringify(cleaned)
  }

  const extractUploadedURL = (res) => {
    return String(
      res?.data?.file?.url ||
      res?.data?.url ||
      res?.file?.url ||
      res?.url ||
      ''
    ).trim()
  }

  const uploadSingleFile = async (file, folder) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('folder', folder)

    const res = await uploadFile(formData)
    if (res.code !== 0) {
      throw new Error(res.msg || '上传失败')
    }

    const url = extractUploadedURL(res)
    if (!url) {
      throw new Error('上传成功但未返回URL')
    }

    return url
  }

  const uploadByRequest = async (options, folder, assignFn, resourceLabel) => {
    try {
      const url = await uploadSingleFile(options.file, folder)
      assignFn(url)
      ElMessage.success(`${resourceLabel}上传成功`)
      if (typeof options.onSuccess === 'function') {
        options.onSuccess({ url })
      }
    } catch (error) {
      ElMessage.error(error?.message || `${resourceLabel}上传失败`)
      if (typeof options.onError === 'function') {
        options.onError(error)
      }
    }
  }

  const loadCategoryList = async () => {
    const res = await getCategoryList({ ...categoryQuery.value })
    if (res.code !== 0) return
    categoryTable.value = res.data?.list || []
    categoryTotal.value = Number(res.data?.total || 0)
    categoryQuery.value.page = Number(res.data?.page || categoryQuery.value.page)
    categoryQuery.value.pageSize = Number(res.data?.pageSize || categoryQuery.value.pageSize)
  }

  const loadChapterList = async () => {
    const res = await getChapterList({ ...chapterQuery.value })
    if (res.code !== 0) return
    chapterTable.value = res.data?.list || []
    chapterTotal.value = Number(res.data?.total || 0)
    chapterQuery.value.page = Number(res.data?.page || chapterQuery.value.page)
    chapterQuery.value.pageSize = Number(res.data?.pageSize || chapterQuery.value.pageSize)
  }

  const loadWordList = async () => {
    const params = {
      page: wordQuery.value.page,
      pageSize: wordQuery.value.pageSize
    }
    if (wordQuery.value.categoryId) {
      params.categoryId = wordQuery.value.categoryId
    }
    if (wordQuery.value.chapterId) {
      params.chapterId = wordQuery.value.chapterId
    }
    if (String(wordQuery.value.keyword || '').trim()) {
      params.keyword = String(wordQuery.value.keyword || '').trim()
    }

    const res = await getEnglishWordList(params)
    if (res.code !== 0) return
    wordTable.value = res.data?.list || []
    wordTotal.value = Number(res.data?.total || 0)
    wordQuery.value.page = Number(res.data?.page || wordQuery.value.page)
    wordQuery.value.pageSize = Number(res.data?.pageSize || wordQuery.value.pageSize)
  }

  const loadCategoryOptions = async () => {
    const res = await getCategoryList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    categoryOptions.value = res.data?.list || []
  }

  const loadChapterOptions = async () => {
    const res = await getChapterList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    chapterOptions.value = res.data?.list || []
  }

  const loadManagedLanguages = async () => {
    try {
      const res = await getLanguageList({ page: 1, pageSize: 500 })
      const list = Array.isArray(res?.data) ? res.data : (res?.data?.list || [])
      managedLanguages.value = list
        .map((lang) => ({
          code: String(lang?.code || '').trim(),
          name: String(lang?.name || '').trim(),
          nativeName: String(lang?.nativeName || '').trim()
        }))
        .filter((lang) => lang.code)
    } catch (error) {
      managedLanguages.value = [
        { code: 'zh', name: '中文', nativeName: '中文' },
        { code: 'en', name: 'English', nativeName: 'English' },
        { code: 'mn', name: 'Монгол', nativeName: 'Монгол' }
      ]
    }

    if (!managedLanguages.value.some((lang) => lang.code === displayLang.value)) {
      displayLang.value = managedLanguages.value[0]?.code || 'zh'
    }
  }

  const handleCategorySearch = () => {
    categoryQuery.value.page = 1
    loadCategoryList()
  }

  const resetCategorySearch = () => {
    categoryQuery.value = { page: 1, pageSize: 10, needVip: undefined }
    loadCategoryList()
  }

  const handleCategoryPageChange = (page) => {
    categoryQuery.value.page = page
    loadCategoryList()
  }

  const handleCategorySizeChange = (size) => {
    categoryQuery.value.pageSize = size
    categoryQuery.value.page = 1
    loadCategoryList()
  }

  const handleChapterSearch = () => {
    chapterQuery.value.page = 1
    loadChapterList()
  }

  const resetChapterSearch = () => {
    chapterQuery.value = { page: 1, pageSize: 10, categoryId: undefined }
    loadChapterList()
  }

  const handleChapterPageChange = (page) => {
    chapterQuery.value.page = page
    loadChapterList()
  }

  const handleChapterSizeChange = (size) => {
    chapterQuery.value.pageSize = size
    chapterQuery.value.page = 1
    loadChapterList()
  }

  const handleWordSearch = () => {
    wordQuery.value.page = 1
    loadWordList()
  }

  const resetWordSearch = () => {
    wordQuery.value = { page: 1, pageSize: 10, categoryId: undefined, chapterId: undefined, keyword: '' }
    loadWordList()
  }

  const handleWordPageChange = (page) => {
    wordQuery.value.page = page
    loadWordList()
  }

  const handleWordSizeChange = (size) => {
    wordQuery.value.pageSize = size
    wordQuery.value.page = 1
    loadWordList()
  }

  const openCategoryDialog = (row) => {
    categoryDialogMode.value = row?.ID ? 'edit' : 'create'
    categoryForm.value = {
      ID: row?.ID || 0,
      nameI18n: normalizeI18nObject(row?.name || ''),
      logo: row?.logo || '',
      price: Number(row?.price || 0),
      needVip: !!row?.needVip,
      sort: Number(row?.sort || 0)
    }
    categoryDialogVisible.value = true
  }

  const submitCategory = async () => {
    const payload = {
      ID: categoryForm.value.ID,
      name: stringifyI18nObject(categoryForm.value.nameI18n),
      logo: String(categoryForm.value.logo || '').trim(),
      price: Number(categoryForm.value.price || 0),
      needVip: !!categoryForm.value.needVip,
      sort: Number(categoryForm.value.sort || 0)
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('分类名称不能为空')
      return
    }

    const res = categoryDialogMode.value === 'create'
      ? await createCategory(payload)
      : await updateCategory(payload)

    if (res.code !== 0) return
    ElMessage.success(categoryDialogMode.value === 'create' ? '分类创建成功' : '分类更新成功')
    categoryDialogVisible.value = false
    await Promise.all([loadCategoryList(), loadCategoryOptions()])
  }

  const removeCategory = (row) => {
    ElMessageBox.confirm(`确认删除分类【${formatI18nText(row.name)}】吗？删除后会级联清理章节关系。`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteCategory({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('分类删除成功')
      await Promise.all([loadCategoryList(), loadCategoryOptions(), loadChapterList(), loadChapterOptions(), loadWordList()])
    })
  }

  const openChapterDialog = (row) => {
    chapterDialogMode.value = row?.ID ? 'edit' : 'create'
    chapterForm.value = {
      ID: row?.ID || 0,
      categoryId: row?.categoryId || chapterQuery.value.categoryId || undefined,
      nameI18n: normalizeI18nObject(row?.name || ''),
      sort: Number(row?.sort || 0)
    }
    chapterDialogVisible.value = true
  }

  const submitChapter = async () => {
    if (!chapterForm.value.categoryId) {
      ElMessage.warning('请选择所属分类')
      return
    }

    const payload = {
      ID: chapterForm.value.ID,
      categoryId: chapterForm.value.categoryId,
      name: stringifyI18nObject(chapterForm.value.nameI18n),
      sort: Number(chapterForm.value.sort || 0)
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('章节名称不能为空')
      return
    }

    const res = chapterDialogMode.value === 'create'
      ? await createChapter(payload)
      : await updateChapter(payload)

    if (res.code !== 0) return
    ElMessage.success(chapterDialogMode.value === 'create' ? '章节创建成功' : '章节更新成功')
    chapterDialogVisible.value = false
    await Promise.all([loadChapterList(), loadChapterOptions()])
  }

  const removeChapter = (row) => {
    ElMessageBox.confirm(`确认删除章节【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteChapter({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('章节删除成功')
      await Promise.all([loadChapterList(), loadChapterOptions(), loadWordList()])
    })
  }

  const openWordDialog = async (row) => {
    const wordId = getEntityID(row)
    wordDialogMode.value = wordId ? 'edit' : 'create'
    wordForm.value = {
      ID: wordId,
      word: row?.word || '',
      phoneticUs: row?.phoneticUs || '',
      phoneticUk: row?.phoneticUk || '',
      audioUs: row?.audioUs || '',
      audioUk: row?.audioUk || '',
      explanationI18n: normalizeI18nObject(row?.explanation || ''),
      sentences: [],
      categoryIds: [],
      chapterIds: [],
      syncBindingsOnUpdate: false
    }

    if (wordId) {
      const res = await findEnglishWord({ ID: wordId })
      if (res.code !== 0) return
      const data = res.data || {}
      wordForm.value.word = data.word || ''
      wordForm.value.phoneticUs = data.phoneticUs || ''
      wordForm.value.phoneticUk = data.phoneticUk || ''
      wordForm.value.audioUs = data.audioUs || ''
      wordForm.value.audioUk = data.audioUk || ''
      wordForm.value.explanationI18n = normalizeI18nObject(data.explanation || '')
      wordForm.value.sentences = Array.isArray(data.sentences)
        ? data.sentences.map((item, idx) => ({
          __key: `sentence_${Date.now()}_${idx}_${sentenceKeySeed++}`,
          source: item?.source || '',
          translateI18n: normalizeI18nObject(item?.translate || ''),
          audioUs: item?.audioUs || '',
          audioUk: item?.audioUk || '',
          sort: Number(item?.sort || idx + 1)
        }))
        : []
      wordForm.value.categoryIds = Array.isArray(data.categoryIds) ? data.categoryIds.map((id) => Number(id || 0)).filter((id) => id > 0) : []
      wordForm.value.chapterIds = Array.isArray(data.chapterIds) ? data.chapterIds.map((id) => Number(id || 0)).filter((id) => id > 0) : []
    }

    wordDialogVisible.value = true
  }

  const submitWord = async () => {
    if (!String(wordForm.value.word || '').trim()) {
      ElMessage.warning('单词不能为空')
      return
    }

    const totalBindings = Number(wordForm.value.categoryIds.length || 0) + Number(wordForm.value.chapterIds.length || 0)
    if (wordDialogMode.value === 'create' && totalBindings === 0) {
      ElMessage.warning('请至少绑定一个分类或章节')
      return
    }

    const payload = {
      ID: wordForm.value.ID,
      word: String(wordForm.value.word || '').trim(),
      phoneticUs: String(wordForm.value.phoneticUs || '').trim(),
      phoneticUk: String(wordForm.value.phoneticUk || '').trim(),
      audioUs: String(wordForm.value.audioUs || '').trim(),
      audioUk: String(wordForm.value.audioUk || '').trim(),
      explanation: stringifyI18nObject(wordForm.value.explanationI18n),
      sentences: (wordForm.value.sentences || [])
        .map((item, idx) => ({
          source: String(item?.source || '').trim(),
          translate: stringifyI18nObject(item?.translateI18n),
          audioUs: String(item?.audioUs || '').trim(),
          audioUk: String(item?.audioUk || '').trim(),
          sort: Number(item?.sort || idx + 1)
        }))
        .filter((item) => item.source)
    }

    if (wordDialogMode.value === 'create') {
      payload.categoryIds = [...wordForm.value.categoryIds]
      payload.chapterIds = [...wordForm.value.chapterIds]
    } else if (wordForm.value.syncBindingsOnUpdate) {
      payload.categoryIds = [...wordForm.value.categoryIds]
      payload.chapterIds = [...wordForm.value.chapterIds]
    }

    wordSubmitting.value = true
    const res = wordDialogMode.value === 'create'
      ? await createEnglishWord(payload)
      : await updateEnglishWord(payload)
    wordSubmitting.value = false

    if (res.code !== 0) return
    ElMessage.success(wordDialogMode.value === 'create' ? '单词创建成功' : '单词更新成功')
    wordDialogVisible.value = false
    loadWordList()
  }

  const removeWord = (row) => {
    const wordId = getEntityID(row)
    if (!wordId) {
      ElMessage.warning('单词ID无效，无法删除')
      return
    }
    if (deletingWordIds.value.has(wordId)) {
      return
    }

    ElMessageBox.confirm(`确认删除单词【${row.word}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      deletingWordIds.value = new Set([...deletingWordIds.value, wordId])
      try {
        const res = await deleteEnglishWord({ ID: wordId })
        if (res.code !== 0) {
          await loadWordList()
          return
        }

        wordTable.value = wordTable.value.filter((item) => getEntityID(item) !== wordId)
        wordTotal.value = Math.max(0, Number(wordTotal.value || 0) - 1)
        if (wordTable.value.length === 0 && Number(wordQuery.value.page || 1) > 1) {
          wordQuery.value.page = Number(wordQuery.value.page || 1) - 1
        }
        await loadWordList()
        ElMessage.success('单词删除成功')
      } finally {
        const next = new Set(deletingWordIds.value)
        next.delete(wordId)
        deletingWordIds.value = next
      }
    })
  }

  const regenerateWord = (row) => {
    const wordId = getEntityID(row)
    if (!wordId) {
      ElMessage.warning('单词ID无效，无法重生发音')
      return
    }
    ElMessageBox.confirm(`确认重生成单词【${row.word}】的美式和英式发音吗？`, '重生成确认', {
      type: 'warning'
    }).then(async () => {
      const res = await regenerateWordAudio({ ID: wordId })
      if (res.code !== 0) return
      ElMessage.success('发音重生成任务执行成功')
      loadWordList()
    })
  }

  onMounted(async () => {
    await Promise.all([
      loadCategoryList(),
      loadChapterList(),
      loadWordList(),
      loadCategoryOptions(),
      loadChapterOptions(),
      loadManagedLanguages()
    ])
  })
</script>

<style scoped>
  .lang-switch-row {
    margin-top: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .toolbar-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    gap: 12px;
    flex-wrap: wrap;
  }

  .ml-2 {
    margin-left: 8px;
  }

  .upload-inline {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    gap: 10px;
    width: 100%;
  }

  .dialog-hint {
    margin-top: 6px;
    color: #6b7280;
    font-size: 12px;
    line-height: 1.4;
  }

  .sentence-editor {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .sentence-toolbar {
    display: flex;
    justify-content: flex-end;
  }

  .sentence-card {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 12px;
    background: #fafafa;
  }

  .sentence-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 13px;
    color: #374151;
    margin-bottom: 6px;
  }

  .toolbar-actions {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .import-progress-box {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    padding: 12px;
    background: #fafafa;
  }

  .import-progress-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    color: #111827;
    font-size: 13px;
    font-weight: 600;
  }

  .import-progress-list {
    max-height: 260px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .import-progress-line {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: #374151;
  }

  .import-progress-time {
    color: #6b7280;
    min-width: 60px;
  }

  .import-progress-text {
    line-height: 1.4;
    word-break: break-all;
  }
</style>
