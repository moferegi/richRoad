<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="英语学习运营台（单词端）：支持分类、章节、单词的完整增删改查。"
        type="info"
        :closable="false"
        show-icon
      />
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
            <el-table-column prop="logo" label="Logo" min-width="180" show-overflow-tooltip />
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
              <el-form-item label="章节过滤">
                <el-select v-model="wordQuery.chapterId" clearable placeholder="全部" style="width: 240px">
                  <el-option
                    v-for="item in chapterOptions"
                    :key="item.ID"
                    :label="`${categoryNameMap[item.categoryId] || '-'} / ${formatI18nText(item.name)}`"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleWordSearch">查询</el-button>
                <el-button icon="refresh" @click="resetWordSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openWordDialog()">新增单词</el-button>
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
                <el-tag :type="scope.row.audioUs ? 'success' : 'warning'">美音{{ scope.row.audioUs ? '已生成' : '缺失' }}</el-tag>
                <el-tag :type="scope.row.audioUk ? 'success' : 'warning'" class="ml-2">英音{{ scope.row.audioUk ? '已生成' : '缺失' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="260">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openWordDialog(scope.row)">编辑</el-button>
                <el-button type="warning" link icon="refresh" @click="regenerateWord(scope.row)">重生发音</el-button>
                <el-button type="danger" link icon="delete" @click="removeWord(scope.row)">删除</el-button>
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

    <el-dialog v-model="categoryDialogVisible" :title="categoryDialogMode === 'create' ? '新增分类' : '编辑分类'" width="640px">
      <el-form :model="categoryForm" label-width="120px">
        <el-form-item label="分类名称(JSON)">
          <el-input
            v-model="categoryForm.name"
            type="textarea"
            :rows="3"
            placeholder='例如: {"zh":"影视英语","en":"Movie English"}'
          />
        </el-form-item>
        <el-form-item label="Logo URL">
          <el-input v-model="categoryForm.logo" placeholder="可选" />
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

    <el-dialog v-model="chapterDialogVisible" :title="chapterDialogMode === 'create' ? '新增章节' : '编辑章节'" width="640px">
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
        <el-form-item label="章节名称(JSON)">
          <el-input
            v-model="chapterForm.name"
            type="textarea"
            :rows="3"
            placeholder='例如: {"zh":"第一章","en":"Chapter 1"}'
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

    <el-dialog v-model="wordDialogVisible" :title="wordDialogMode === 'create' ? '新增单词' : '编辑单词'" width="760px">
      <el-form :model="wordForm" label-width="120px">
        <el-form-item label="单词本体">
          <el-input v-model="wordForm.word" placeholder="例如: destiny" />
        </el-form-item>
        <el-form-item label="美式音标">
          <el-input v-model="wordForm.phoneticUs" placeholder="例如: /ˈdestəni/" />
        </el-form-item>
        <el-form-item label="英式音标">
          <el-input v-model="wordForm.phoneticUk" placeholder="例如: /ˈdestəni/" />
        </el-form-item>
        <el-form-item label="释义(JSON)">
          <el-input
            v-model="wordForm.explanation"
            type="textarea"
            :rows="4"
            placeholder='例如: {"zh":"命运","en":"fate"}'
          />
        </el-form-item>
        <el-form-item label="绑定章节" v-if="wordDialogMode === 'create'">
          <el-select v-model="wordForm.chapterIds" multiple filterable placeholder="创建时建议至少绑定一个章节" style="width: 100%">
            <el-option
              v-for="item in chapterOptions"
              :key="item.ID"
              :label="`${categoryNameMap[item.categoryId] || '-'} / ${formatI18nText(item.name)}`"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="重绑章节" v-else>
          <el-select v-model="wordForm.chapterIds" multiple filterable clearable placeholder="不选择表示保持原章节绑定不变" style="width: 100%">
            <el-option
              v-for="item in chapterOptions"
              :key="item.ID"
              :label="`${categoryNameMap[item.categoryId] || '-'} / ${formatI18nText(item.name)}`"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="wordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="wordSubmitting" @click="submitWord">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
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
    updateCategory,
    updateChapter,
    updateEnglishWord
  } from '../api/english'

  defineOptions({
    name: 'EnglishLearningWord'
  })

  const activeTab = ref('category')

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
    chapterId: undefined
  })
  const wordTable = ref([])
  const wordTotal = ref(0)

  const categoryOptions = ref([])
  const chapterOptions = ref([])

  const categoryNameMap = computed(() => {
    const map = {}
    for (const item of categoryOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const categoryDialogVisible = ref(false)
  const categoryDialogMode = ref('create')
  const categoryForm = ref({
    ID: 0,
    name: '{}',
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
    name: '{}',
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
    explanation: '{}',
    chapterIds: []
  })

  const formatI18nText = (raw) => {
    if (!raw) return ''
    try {
      const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
      if (typeof parsed === 'object' && parsed !== null) {
        return parsed.zh || parsed.en || Object.values(parsed)[0] || ''
      }
      return String(raw)
    } catch (e) {
      return String(raw)
    }
  }

  const ensureJsonText = (raw) => {
    const text = String(raw || '').trim()
    if (!text) {
      return '{}'
    }
    JSON.parse(text)
    return text
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
    const res = await getEnglishWordList({ ...wordQuery.value })
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
    wordQuery.value = { page: 1, pageSize: 10, chapterId: undefined }
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
      name: row?.name || '{}',
      logo: row?.logo || '',
      price: Number(row?.price || 0),
      needVip: !!row?.needVip,
      sort: Number(row?.sort || 0)
    }
    categoryDialogVisible.value = true
  }

  const submitCategory = async () => {
    if (!categoryForm.value.name) {
      ElMessage.warning('分类名称不能为空')
      return
    }
    try {
      categoryForm.value.name = ensureJsonText(categoryForm.value.name)
    } catch (e) {
      ElMessage.error('分类名称必须是合法 JSON')
      return
    }

    const payload = {
      ID: categoryForm.value.ID,
      name: categoryForm.value.name,
      logo: categoryForm.value.logo,
      price: Number(categoryForm.value.price || 0),
      needVip: !!categoryForm.value.needVip,
      sort: Number(categoryForm.value.sort || 0)
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
      await Promise.all([loadCategoryList(), loadCategoryOptions(), loadChapterList(), loadChapterOptions()])
    })
  }

  const openChapterDialog = (row) => {
    chapterDialogMode.value = row?.ID ? 'edit' : 'create'
    chapterForm.value = {
      ID: row?.ID || 0,
      categoryId: row?.categoryId || chapterQuery.value.categoryId || undefined,
      name: row?.name || '{}',
      sort: Number(row?.sort || 0)
    }
    chapterDialogVisible.value = true
  }

  const submitChapter = async () => {
    if (!chapterForm.value.categoryId) {
      ElMessage.warning('请选择所属分类')
      return
    }
    if (!chapterForm.value.name) {
      ElMessage.warning('章节名称不能为空')
      return
    }
    try {
      chapterForm.value.name = ensureJsonText(chapterForm.value.name)
    } catch (e) {
      ElMessage.error('章节名称必须是合法 JSON')
      return
    }

    const payload = {
      ID: chapterForm.value.ID,
      categoryId: chapterForm.value.categoryId,
      name: chapterForm.value.name,
      sort: Number(chapterForm.value.sort || 0)
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
      await Promise.all([loadChapterList(), loadChapterOptions()])
    })
  }

  const openWordDialog = async (row) => {
    wordDialogMode.value = row?.ID ? 'edit' : 'create'
    wordForm.value = {
      ID: row?.ID || 0,
      word: row?.word || '',
      phoneticUs: row?.phoneticUs || '',
      phoneticUk: row?.phoneticUk || '',
      explanation: row?.explanation || '{}',
      chapterIds: []
    }

    if (row?.ID) {
      const res = await findEnglishWord({ ID: row.ID })
      if (res.code !== 0) return
      const data = res.data || {}
      wordForm.value.word = data.word || ''
      wordForm.value.phoneticUs = data.phoneticUs || ''
      wordForm.value.phoneticUk = data.phoneticUk || ''
      wordForm.value.explanation = data.explanation || '{}'
    }

    wordDialogVisible.value = true
  }

  const submitWord = async () => {
    if (!wordForm.value.word) {
      ElMessage.warning('单词不能为空')
      return
    }
    if (wordDialogMode.value === 'create' && wordForm.value.chapterIds.length === 0) {
      ElMessage.warning('请至少绑定一个章节')
      return
    }
    try {
      wordForm.value.explanation = ensureJsonText(wordForm.value.explanation)
    } catch (e) {
      ElMessage.error('释义必须是合法 JSON')
      return
    }

    const payload = {
      ID: wordForm.value.ID,
      word: wordForm.value.word,
      phoneticUs: wordForm.value.phoneticUs,
      phoneticUk: wordForm.value.phoneticUk,
      explanation: wordForm.value.explanation
    }

    if (wordDialogMode.value === 'create' || wordForm.value.chapterIds.length > 0) {
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
    ElMessageBox.confirm(`确认删除单词【${row.word}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteEnglishWord({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('单词删除成功')
      loadWordList()
    })
  }

  const regenerateWord = (row) => {
    ElMessageBox.confirm(`确认重生成单词【${row.word}】的美式和英式发音吗？`, '重生成确认', {
      type: 'warning'
    }).then(async () => {
      const res = await regenerateWordAudio({ ID: row.ID })
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
      loadChapterOptions()
    ])
  })
</script>

<style scoped>
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
</style>
