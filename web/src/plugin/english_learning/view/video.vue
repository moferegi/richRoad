<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="英语学习运营台（视频端）：支持视频分类、剧集、单集、字幕文件解析和用户资源授权管理。"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="视频分类" name="videoCategory">
          <div class="toolbar-row">
            <span class="toolbar-tip">用于维护视频大类（如影视英语、情景口语等）与存储目录键</span>
            <el-button type="primary" icon="plus" @click="openVideoCategoryDialog()">新增视频分类</el-button>
          </div>

          <el-table :data="videoCategoryTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="分类名称" min-width="240">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="storageKey" label="存储目录键" min-width="180" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.storageKey || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openVideoCategoryDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeVideoCategory(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="videoCategoryQuery.page"
              :page-size="videoCategoryQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="videoCategoryTotal"
              @current-change="handleVideoCategoryPageChange"
              @size-change="handleVideoCategorySizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="剧集管理" name="series">
          <div class="toolbar-row">
            <el-form :inline="true" :model="seriesQuery">
              <el-form-item label="所属分类">
                <el-select v-model="seriesQuery.categoryId" clearable placeholder="全部" style="width: 240px">
                  <el-option
                    v-for="item in videoCategoryOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleSeriesSearch">查询</el-button>
                <el-button icon="refresh" @click="resetSeriesSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openSeriesDialog()">新增剧集</el-button>
          </div>

          <el-table :data="seriesTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属分类" min-width="180">
              <template #default="scope">
                {{ videoCategoryNameMap[scope.row.categoryId] || `#${scope.row.categoryId}` }}
              </template>
            </el-table-column>
            <el-table-column label="剧集名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="coverUrl" label="封面地址" min-width="240" show-overflow-tooltip />
            <el-table-column prop="price" label="价格" width="120" />
            <el-table-column label="会员限制" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.needVip ? 'warning' : 'success'">{{ scope.row.needVip ? '需要会员' : '普通可见' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="viewCount" label="浏览量" width="100" />
            <el-table-column prop="userCount" label="观看人数" width="110" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openSeriesDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeSeries(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="seriesQuery.page"
              :page-size="seriesQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="seriesTotal"
              @current-change="handleSeriesPageChange"
              @size-change="handleSeriesSizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="单集与字幕" name="episode">
          <div class="toolbar-row">
            <el-form :inline="true" :model="episodeQuery">
              <el-form-item label="所属剧集">
                <el-select v-model="episodeQuery.seriesId" clearable filterable placeholder="全部" style="width: 280px">
                  <el-option
                    v-for="item in seriesOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleEpisodeSearch">查询</el-button>
                <el-button icon="refresh" @click="resetEpisodeSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openEpisodeDialog()">新增单集</el-button>
          </div>

          <el-table :data="episodeTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属剧集" min-width="220">
              <template #default="scope">
                {{ seriesNameMap[scope.row.seriesId] || `#${scope.row.seriesId}` }}
              </template>
            </el-table-column>
            <el-table-column label="单集名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="videoUrl" label="视频地址" min-width="260" show-overflow-tooltip />
            <el-table-column prop="trialPercent" label="试看比例(%)" width="120" />
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openEpisodeDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeEpisode(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="episodeQuery.page"
              :page-size="episodeQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="episodeTotal"
              @current-change="handleEpisodePageChange"
              @size-change="handleEpisodeSizeChange"
            />
          </div>

          <div class="subtitle-box">
            <el-divider content-position="left">字幕文件解析入库</el-divider>
            <el-form label-width="140px">
              <el-form-item label="目标单集">
                <el-select v-model="subtitleForm.episodeId" filterable placeholder="请选择单集" style="width: 480px">
                  <el-option
                    v-for="item in episodeOptions"
                    :key="item.ID"
                    :label="`${seriesNameMap[item.seriesId] || '-'} / ${formatI18nText(item.name)}`"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="英文字幕文件(必填)">
                <div class="upload-inline upload-inline-wide">
                  <el-input v-model="subtitleForm.englishSubtitleUrl" placeholder="SRT/VTT 文件URL" />
                  <el-upload
                    :show-file-list="false"
                    :http-request="(options) => uploadByRequest(options, subtitleUploadFolder, (url) => { subtitleForm.englishSubtitleUrl = url }, '英文字幕')"
                  >
                    <el-button type="primary" plain>上传英文字幕</el-button>
                  </el-upload>
                </div>
                <div class="dialog-hint">当前字幕上传目录：{{ subtitleUploadFolder }}</div>
              </el-form-item>

              <el-form-item
                v-for="lang in subtitleLanguages"
                :key="lang.code"
                :label="`${lang.name || lang.code} 字幕文件`"
              >
                <div class="upload-inline upload-inline-wide">
                  <el-input
                    v-model="subtitleForm.translationSubtitleMap[lang.code]"
                    :placeholder="`${lang.code} 字幕URL（可选）`"
                  />
                  <el-upload
                    :show-file-list="false"
                    :http-request="(options) => uploadByRequest(options, subtitleUploadFolder, (url) => { subtitleForm.translationSubtitleMap[lang.code] = url }, `${lang.code}字幕`)"
                  >
                    <el-button type="primary" plain>上传{{ lang.code }}</el-button>
                  </el-upload>
                </div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" @click="handleParseSubtitleFiles">提交字幕解析</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="用户授权" name="entitlement">
          <div class="toolbar-row">
            <el-form :inline="true" :model="entitlementQuery">
              <el-form-item label="用户ID">
                <el-input-number v-model="entitlementQuery.userId" :min="1" controls-position="right" />
              </el-form-item>
              <el-form-item label="资源类型">
                <el-select v-model="entitlementQuery.resourceType" clearable placeholder="全部" style="width: 220px">
                  <el-option label="分类" value="english_category" />
                  <el-option label="剧集" value="video_series" />
                  <el-option label="单集" value="video_episode" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleEntitlementSearch">查询</el-button>
                <el-button icon="refresh" @click="resetEntitlementSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openEntitlementDialog">新增授权</el-button>
          </div>

          <el-table :data="entitlementTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="userId" label="用户ID" width="100" />
            <el-table-column label="资源类型" width="140">
              <template #default="scope">
                <el-tag>{{ resourceTypeLabel(scope.row.resourceType) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="resourceId" label="资源ID" width="110" />
            <el-table-column prop="source" label="授权来源" width="140" />
            <el-table-column prop="grantedBy" label="授权人" width="100" />
            <el-table-column label="过期时间" min-width="180">
              <template #default="scope">
                {{ scope.row.expireAt || '永久有效' }}
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="scope">
                <el-button type="danger" link icon="delete" @click="removeEntitlement(scope.row)">撤销</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="entitlementQuery.page"
              :page-size="entitlementQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="entitlementTotal"
              @current-change="handleEntitlementPageChange"
              @size-change="handleEntitlementSizeChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="videoCategoryDialogVisible" :title="videoCategoryDialogMode === 'create' ? '新增视频分类' : '编辑视频分类'" width="760px">
      <el-form :model="videoCategoryForm" label-width="120px">
        <el-form-item label="分类名称">
          <MultiLangEditor
            :model="videoCategoryForm.nameI18n"
            title="视频分类多语言名称"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="存储目录键">
          <el-input v-model="videoCategoryForm.storageKey" placeholder="例如: liblib（用于视频/字幕目录 english-learn/video/liblib）" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="videoCategoryForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="videoCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitVideoCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="seriesDialogVisible" :title="seriesDialogMode === 'create' ? '新增剧集' : '编辑剧集'" width="860px">
      <el-form :model="seriesForm" label-width="120px">
        <el-form-item label="所属分类">
          <el-select v-model="seriesForm.categoryId" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in videoCategoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="剧集名称">
          <MultiLangEditor
            :model="seriesForm.nameI18n"
            title="剧集名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="封面地址">
          <div class="upload-inline upload-inline-wide">
            <el-input v-model="seriesForm.coverUrl" placeholder="上传后自动写入 URL" />
            <el-upload
              :show-file-list="false"
              :http-request="(options) => uploadByRequest(options, 'english-learn/pic', (url) => { seriesForm.coverUrl = url }, '封面')"
            >
              <el-button type="primary" plain>上传封面</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="seriesForm.price" :min="0" :precision="2" :step="1" />
        </el-form-item>
        <el-form-item label="需要会员">
          <el-switch v-model="seriesForm.needVip" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="seriesDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSeries">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="episodeDialogVisible" :title="episodeDialogMode === 'create' ? '新增单集' : '编辑单集'" width="860px">
      <el-form :model="episodeForm" label-width="120px">
        <el-form-item label="所属剧集">
          <el-select v-model="episodeForm.seriesId" filterable placeholder="请选择剧集" style="width: 100%">
            <el-option
              v-for="item in seriesOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="单集名称">
          <MultiLangEditor
            :model="episodeForm.nameI18n"
            title="单集名称多语言"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="视频地址">
          <div class="upload-inline upload-inline-wide">
            <el-input v-model="episodeForm.videoUrl" placeholder="上传后自动写入 URL" />
            <el-upload
              :show-file-list="false"
              :http-request="(options) => uploadByRequest(options, episodeUploadFolder, (url) => { episodeForm.videoUrl = url }, '视频文件')"
            >
              <el-button type="primary" plain>上传视频</el-button>
            </el-upload>
          </div>
          <div class="dialog-hint">当前视频上传目录：{{ episodeUploadFolder }}</div>
        </el-form-item>
        <el-form-item label="试看比例(%)">
          <el-input-number v-model="episodeForm.trialPercent" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="episodeForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="episodeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEpisode">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="entitlementDialogVisible" title="新增资源授权" width="760px">
      <el-form :model="entitlementForm" label-width="120px">
        <el-form-item label="用户ID">
          <el-input-number v-model="entitlementForm.userId" :min="1" controls-position="right" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="entitlementForm.resourceType" style="width: 100%">
            <el-option label="分类" value="english_category" />
            <el-option label="剧集" value="video_series" />
            <el-option label="单集" value="video_episode" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源ID">
          <el-select
            v-model="entitlementForm.resourceId"
            filterable
            clearable
            placeholder="可从列表选择，或手动输入"
            style="width: 100%"
          >
            <el-option
              v-for="item in resourceOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-input-number v-model="entitlementForm.resourceId" :min="1" class="mt-2 w-full" controls-position="right" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="entitlementForm.expireAt"
            type="datetime"
            placeholder="不填表示永久"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
            clearable
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="entitlementForm.remark" type="textarea" :rows="3" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="entitlementDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEntitlement">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
  import { uploadFile } from '@/api/fileUploadAndDownload'
  import { getLanguageList } from '@/api/client/language'
  import {
    createVideoCategory,
    createVideoEpisode,
    createVideoSeries,
    deleteVideoCategory,
    deleteVideoEpisode,
    deleteVideoSeries,
    getEntitlementList,
    getVideoCategoryList,
    getVideoEpisodeList,
    getVideoSeriesList,
    grantEntitlement,
    parseSubtitleFiles,
    revokeEntitlement,
    updateVideoCategory,
    updateVideoEpisode,
    updateVideoSeries
  } from '../api/english'

  defineOptions({
    name: 'EnglishLearningVideo'
  })

  const activeTab = ref('videoCategory')

  const normalizeI18nObject = (raw) => {
    if (!raw) return { zh: '' }
    if (typeof raw === 'object') {
      return Object.keys(raw).length > 0 ? { ...raw } : { zh: '' }
    }
    const text = String(raw).trim()
    if (!text) return { zh: '' }
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
      if (!lang) continue
      const text = String(value ?? '').trim()
      if (text) cleaned[lang] = text
    }
    return JSON.stringify(Object.keys(cleaned).length > 0 ? cleaned : { zh: '' })
  }

  const formatI18nText = (raw) => {
    if (!raw) return ''
    try {
      const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
      if (typeof parsed === 'object' && parsed !== null) {
        return parsed.zh || parsed.en || parsed.mn || Object.values(parsed)[0] || ''
      }
      return String(raw)
    } catch (e) {
      return String(raw)
    }
  }

  const normalizeStorageKey = (raw) => {
    return String(raw || '')
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9_-]+/g, '-')
      .replace(/^-+|-+$/g, '')
  }

  const videoCategoryQuery = ref({ page: 1, pageSize: 10 })
  const videoCategoryTable = ref([])
  const videoCategoryTotal = ref(0)

  const seriesQuery = ref({ page: 1, pageSize: 10, categoryId: undefined })
  const seriesTable = ref([])
  const seriesTotal = ref(0)

  const episodeQuery = ref({ page: 1, pageSize: 10, seriesId: undefined })
  const episodeTable = ref([])
  const episodeTotal = ref(0)

  const entitlementQuery = ref({ page: 1, pageSize: 10, userId: undefined, resourceType: undefined })
  const entitlementTable = ref([])
  const entitlementTotal = ref(0)

  const videoCategoryOptions = ref([])
  const seriesOptions = ref([])
  const episodeOptions = ref([])

  const managedLanguages = ref([])

  const subtitleLanguages = computed(() => {
    return managedLanguages.value.filter((lang) => String(lang.code || '').toLowerCase() !== 'en')
  })

  const videoCategoryNameMap = computed(() => {
    const map = {}
    for (const item of videoCategoryOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const seriesNameMap = computed(() => {
    const map = {}
    for (const item of seriesOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const videoCategoryDialogVisible = ref(false)
  const videoCategoryDialogMode = ref('create')
  const videoCategoryForm = ref({
    ID: 0,
    nameI18n: { zh: '' },
    storageKey: '',
    sort: 0
  })

  const seriesDialogVisible = ref(false)
  const seriesDialogMode = ref('create')
  const seriesForm = ref({
    ID: 0,
    categoryId: undefined,
    nameI18n: { zh: '' },
    coverId: 0,
    coverUrl: '',
    price: 0,
    needVip: false
  })

  const episodeDialogVisible = ref(false)
  const episodeDialogMode = ref('create')
  const episodeForm = ref({
    ID: 0,
    seriesId: undefined,
    nameI18n: { zh: '' },
    videoUrl: '',
    trialPercent: 8,
    sort: 0
  })

  const subtitleForm = ref({
    episodeId: undefined,
    englishSubtitleUrl: '',
    translationSubtitleMap: {}
  })

  const entitlementDialogVisible = ref(false)
  const entitlementForm = ref({
    userId: undefined,
    resourceType: 'video_series',
    resourceId: undefined,
    expireAt: '',
    remark: ''
  })

  const resourceOptions = computed(() => {
    if (entitlementForm.value.resourceType === 'english_category') {
      return videoCategoryOptions.value.map((item) => ({
        value: item.ID,
        label: `${item.ID} - ${formatI18nText(item.name)}`
      }))
    }
    if (entitlementForm.value.resourceType === 'video_series') {
      return seriesOptions.value.map((item) => ({
        value: item.ID,
        label: `${item.ID} - ${formatI18nText(item.name)}`
      }))
    }
    return episodeOptions.value.map((item) => ({
      value: item.ID,
      label: `${item.ID} - ${seriesNameMap.value[item.seriesId] || '-'} / ${formatI18nText(item.name)}`
    }))
  })

  const getCategoryByID = (categoryID) => videoCategoryOptions.value.find((item) => Number(item.ID) === Number(categoryID))
  const getSeriesByID = (seriesID) => seriesOptions.value.find((item) => Number(item.ID) === Number(seriesID))
  const getEpisodeByID = (episodeID) => episodeOptions.value.find((item) => Number(item.ID) === Number(episodeID))

  const resolveVideoFolderBySeriesID = (seriesID) => {
    const series = getSeriesByID(seriesID)
    const category = getCategoryByID(series?.categoryId)
    const storageKey = normalizeStorageKey(category?.storageKey || '')
    return storageKey ? `english-learn/video/${storageKey}` : 'english-learn/video/general'
  }

  const episodeUploadFolder = computed(() => {
    if (episodeForm.value.seriesId) {
      return resolveVideoFolderBySeriesID(episodeForm.value.seriesId)
    }
    return 'english-learn/video/general'
  })

  const subtitleUploadFolder = computed(() => {
    const episode = getEpisodeByID(subtitleForm.value.episodeId)
    if (!episode) {
      return 'english-learn/video/general/subtitle'
    }
    return `${resolveVideoFolderBySeriesID(episode.seriesId)}/subtitle`
  })

  const resourceTypeLabel = (type) => {
    if (type === 'english_category') return '分类'
    if (type === 'video_series') return '剧集'
    if (type === 'video_episode') return '单集'
    return type || '-'
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
        { code: 'zh', name: '中文' },
        { code: 'mn', name: '蒙文' },
        { code: 'en', name: '英文' }
      ]
    }
  }

  const loadVideoCategoryList = async () => {
    const res = await getVideoCategoryList({ ...videoCategoryQuery.value })
    if (res.code !== 0) return
    videoCategoryTable.value = res.data?.list || []
    videoCategoryTotal.value = Number(res.data?.total || 0)
    videoCategoryQuery.value.page = Number(res.data?.page || videoCategoryQuery.value.page)
    videoCategoryQuery.value.pageSize = Number(res.data?.pageSize || videoCategoryQuery.value.pageSize)
  }

  const loadSeriesList = async () => {
    const res = await getVideoSeriesList({ ...seriesQuery.value })
    if (res.code !== 0) return
    seriesTable.value = res.data?.list || []
    seriesTotal.value = Number(res.data?.total || 0)
    seriesQuery.value.page = Number(res.data?.page || seriesQuery.value.page)
    seriesQuery.value.pageSize = Number(res.data?.pageSize || seriesQuery.value.pageSize)
  }

  const loadEpisodeList = async () => {
    const res = await getVideoEpisodeList({ ...episodeQuery.value })
    if (res.code !== 0) return
    episodeTable.value = res.data?.list || []
    episodeTotal.value = Number(res.data?.total || 0)
    episodeQuery.value.page = Number(res.data?.page || episodeQuery.value.page)
    episodeQuery.value.pageSize = Number(res.data?.pageSize || episodeQuery.value.pageSize)
  }

  const loadEntitlementList = async () => {
    const params = {
      page: entitlementQuery.value.page,
      pageSize: entitlementQuery.value.pageSize
    }
    if (entitlementQuery.value.userId) {
      params.userId = entitlementQuery.value.userId
    }
    if (entitlementQuery.value.resourceType) {
      params.resourceType = entitlementQuery.value.resourceType
    }

    const res = await getEntitlementList(params)
    if (res.code !== 0) return
    entitlementTable.value = res.data?.list || []
    entitlementTotal.value = Number(res.data?.total || 0)
    entitlementQuery.value.page = Number(res.data?.page || entitlementQuery.value.page)
    entitlementQuery.value.pageSize = Number(res.data?.pageSize || entitlementQuery.value.pageSize)
  }

  const loadVideoCategoryOptions = async () => {
    const res = await getVideoCategoryList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    videoCategoryOptions.value = res.data?.list || []
  }

  const loadSeriesOptions = async () => {
    const res = await getVideoSeriesList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    seriesOptions.value = res.data?.list || []
  }

  const loadEpisodeOptions = async () => {
    const res = await getVideoEpisodeList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    episodeOptions.value = res.data?.list || []
  }

  const handleVideoCategoryPageChange = (page) => {
    videoCategoryQuery.value.page = page
    loadVideoCategoryList()
  }

  const handleVideoCategorySizeChange = (size) => {
    videoCategoryQuery.value.pageSize = size
    videoCategoryQuery.value.page = 1
    loadVideoCategoryList()
  }

  const handleSeriesSearch = () => {
    seriesQuery.value.page = 1
    loadSeriesList()
  }

  const resetSeriesSearch = () => {
    seriesQuery.value = { page: 1, pageSize: 10, categoryId: undefined }
    loadSeriesList()
  }

  const handleSeriesPageChange = (page) => {
    seriesQuery.value.page = page
    loadSeriesList()
  }

  const handleSeriesSizeChange = (size) => {
    seriesQuery.value.pageSize = size
    seriesQuery.value.page = 1
    loadSeriesList()
  }

  const handleEpisodeSearch = () => {
    episodeQuery.value.page = 1
    loadEpisodeList()
  }

  const resetEpisodeSearch = () => {
    episodeQuery.value = { page: 1, pageSize: 10, seriesId: undefined }
    loadEpisodeList()
  }

  const handleEpisodePageChange = (page) => {
    episodeQuery.value.page = page
    loadEpisodeList()
  }

  const handleEpisodeSizeChange = (size) => {
    episodeQuery.value.pageSize = size
    episodeQuery.value.page = 1
    loadEpisodeList()
  }

  const handleEntitlementSearch = () => {
    entitlementQuery.value.page = 1
    loadEntitlementList()
  }

  const resetEntitlementSearch = () => {
    entitlementQuery.value = { page: 1, pageSize: 10, userId: undefined, resourceType: undefined }
    loadEntitlementList()
  }

  const handleEntitlementPageChange = (page) => {
    entitlementQuery.value.page = page
    loadEntitlementList()
  }

  const handleEntitlementSizeChange = (size) => {
    entitlementQuery.value.pageSize = size
    entitlementQuery.value.page = 1
    loadEntitlementList()
  }

  const openVideoCategoryDialog = (row) => {
    videoCategoryDialogMode.value = row?.ID ? 'edit' : 'create'
    videoCategoryForm.value = {
      ID: row?.ID || 0,
      nameI18n: normalizeI18nObject(row?.name || ''),
      storageKey: row?.storageKey || '',
      sort: Number(row?.sort || 0)
    }
    videoCategoryDialogVisible.value = true
  }

  const submitVideoCategory = async () => {
    const payload = {
      ID: videoCategoryForm.value.ID,
      name: stringifyI18nObject(videoCategoryForm.value.nameI18n),
      storageKey: normalizeStorageKey(videoCategoryForm.value.storageKey),
      sort: Number(videoCategoryForm.value.sort || 0)
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('分类名称不能为空')
      return
    }

    const res = videoCategoryDialogMode.value === 'create'
      ? await createVideoCategory(payload)
      : await updateVideoCategory(payload)

    if (res.code !== 0) return
    ElMessage.success(videoCategoryDialogMode.value === 'create' ? '视频分类创建成功' : '视频分类更新成功')
    videoCategoryDialogVisible.value = false
    await Promise.all([loadVideoCategoryList(), loadVideoCategoryOptions()])
  }

  const removeVideoCategory = (row) => {
    ElMessageBox.confirm(`确认删除视频分类【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoCategory({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('视频分类删除成功')
      await Promise.all([
        loadVideoCategoryList(),
        loadVideoCategoryOptions(),
        loadSeriesList(),
        loadSeriesOptions(),
        loadEpisodeList(),
        loadEpisodeOptions()
      ])
    })
  }

  const openSeriesDialog = (row) => {
    seriesDialogMode.value = row?.ID ? 'edit' : 'create'
    seriesForm.value = {
      ID: row?.ID || 0,
      categoryId: row?.categoryId || seriesQuery.value.categoryId || undefined,
      nameI18n: normalizeI18nObject(row?.name || ''),
      coverId: Number(row?.coverId || 0),
      coverUrl: row?.coverUrl || '',
      price: Number(row?.price || 0),
      needVip: !!row?.needVip
    }
    seriesDialogVisible.value = true
  }

  const submitSeries = async () => {
    if (!seriesForm.value.categoryId) {
      ElMessage.warning('请选择所属分类')
      return
    }

    const payload = {
      ID: seriesForm.value.ID,
      categoryId: seriesForm.value.categoryId,
      name: stringifyI18nObject(seriesForm.value.nameI18n),
      coverId: Number(seriesForm.value.coverId || 0),
      coverUrl: String(seriesForm.value.coverUrl || '').trim(),
      price: Number(seriesForm.value.price || 0),
      needVip: !!seriesForm.value.needVip
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('剧集名称不能为空')
      return
    }

    const res = seriesDialogMode.value === 'create'
      ? await createVideoSeries(payload)
      : await updateVideoSeries(payload)

    if (res.code !== 0) return
    ElMessage.success(seriesDialogMode.value === 'create' ? '剧集创建成功' : '剧集更新成功')
    seriesDialogVisible.value = false
    await Promise.all([loadSeriesList(), loadSeriesOptions(), loadEpisodeOptions()])
  }

  const removeSeries = (row) => {
    ElMessageBox.confirm(`确认删除剧集【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoSeries({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('剧集删除成功')
      await Promise.all([loadSeriesList(), loadSeriesOptions(), loadEpisodeList(), loadEpisodeOptions()])
    })
  }

  const openEpisodeDialog = (row) => {
    episodeDialogMode.value = row?.ID ? 'edit' : 'create'
    episodeForm.value = {
      ID: row?.ID || 0,
      seriesId: row?.seriesId || episodeQuery.value.seriesId || undefined,
      nameI18n: normalizeI18nObject(row?.name || ''),
      videoUrl: row?.videoUrl || '',
      trialPercent: Number(row?.trialPercent || 8),
      sort: Number(row?.sort || 0)
    }
    episodeDialogVisible.value = true
  }

  const submitEpisode = async () => {
    if (!episodeForm.value.seriesId) {
      ElMessage.warning('请选择所属剧集')
      return
    }
    if (!String(episodeForm.value.videoUrl || '').trim()) {
      ElMessage.warning('视频地址不能为空')
      return
    }

    const payload = {
      ID: episodeForm.value.ID,
      seriesId: episodeForm.value.seriesId,
      name: stringifyI18nObject(episodeForm.value.nameI18n),
      videoUrl: String(episodeForm.value.videoUrl || '').trim(),
      trialPercent: Number(episodeForm.value.trialPercent || 8),
      sort: Number(episodeForm.value.sort || 0)
    }

    if (!formatI18nText(payload.name)) {
      ElMessage.warning('单集名称不能为空')
      return
    }

    const res = episodeDialogMode.value === 'create'
      ? await createVideoEpisode(payload)
      : await updateVideoEpisode(payload)

    if (res.code !== 0) return
    ElMessage.success(episodeDialogMode.value === 'create' ? '单集创建成功' : '单集更新成功')
    episodeDialogVisible.value = false
    await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
  }

  const removeEpisode = (row) => {
    ElMessageBox.confirm(`确认删除单集【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoEpisode({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('单集删除成功')
      await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
    })
  }

  const handleParseSubtitleFiles = async () => {
    if (!subtitleForm.value.episodeId) {
      ElMessage.warning('请选择目标单集')
      return
    }

    const englishSubtitleURL = String(subtitleForm.value.englishSubtitleUrl || '').trim()
    if (!englishSubtitleURL) {
      ElMessage.warning('英文字幕文件不能为空')
      return
    }

    const translationSubtitle = Object.entries(subtitleForm.value.translationSubtitleMap || {})
      .map(([language, subtitleUrl]) => ({
        language,
        subtitleUrl: String(subtitleUrl || '').trim()
      }))
      .filter((item) => item.language && item.subtitleUrl)

    const res = await parseSubtitleFiles({
      episodeId: subtitleForm.value.episodeId,
      englishSubtitleUrl: englishSubtitleURL,
      translationSubtitle
    })

    if (res.code !== 0) return
    ElMessage.success('字幕文件解析并入库成功')
  }

  const openEntitlementDialog = () => {
    entitlementForm.value = {
      userId: undefined,
      resourceType: 'video_series',
      resourceId: undefined,
      expireAt: '',
      remark: ''
    }
    entitlementDialogVisible.value = true
  }

  const submitEntitlement = async () => {
    if (!entitlementForm.value.userId) {
      ElMessage.warning('用户ID不能为空')
      return
    }
    if (!entitlementForm.value.resourceType) {
      ElMessage.warning('资源类型不能为空')
      return
    }
    if (!entitlementForm.value.resourceId) {
      ElMessage.warning('资源ID不能为空')
      return
    }

    const payload = {
      userId: Number(entitlementForm.value.userId),
      resourceType: entitlementForm.value.resourceType,
      resourceId: Number(entitlementForm.value.resourceId),
      remark: entitlementForm.value.remark || ''
    }
    if (entitlementForm.value.expireAt) {
      payload.expireAt = entitlementForm.value.expireAt
    }

    const res = await grantEntitlement(payload)
    if (res.code !== 0) return
    ElMessage.success('授权成功')
    entitlementDialogVisible.value = false
    loadEntitlementList()
  }

  const removeEntitlement = (row) => {
    ElMessageBox.confirm(`确认撤销用户 ${row.userId} 的授权记录 #${row.ID} 吗？`, '撤销确认', {
      type: 'warning'
    }).then(async () => {
      const res = await revokeEntitlement({
        userId: row.userId,
        resourceType: row.resourceType,
        resourceId: row.resourceId
      })
      if (res.code !== 0) return
      ElMessage.success('授权已撤销')
      loadEntitlementList()
    })
  }

  onMounted(async () => {
    await Promise.all([
      loadVideoCategoryList(),
      loadSeriesList(),
      loadEpisodeList(),
      loadEntitlementList(),
      loadVideoCategoryOptions(),
      loadSeriesOptions(),
      loadEpisodeOptions(),
      loadManagedLanguages()
    ])
  })
</script>

<style scoped>
  .toolbar-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }

  .toolbar-tip {
    color: #4b5563;
    font-size: 13px;
  }

  .subtitle-box {
    margin-top: 12px;
    padding: 12px 0;
  }

  .upload-inline {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    gap: 10px;
    width: 100%;
  }

  .upload-inline-wide {
    max-width: 680px;
  }

  .dialog-hint {
    margin-top: 6px;
    color: #6b7280;
    font-size: 12px;
    line-height: 1.4;
  }

  .mt-2 {
    margin-top: 8px;
  }

  .w-full {
    width: 100%;
  }
</style>
