<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="handleAdd({...geo,level:-1})">新增省会</el-button>
        <span class="geo-tree-hint">当前展示国家及其子地区；点击左侧展开箭头可逐级查看，支持继续新增子地区</span>
      </div>
      <el-table
        ref="mainTable"
        :data="treeData"
        style="width: 100%"
        row-key="code"
        border
        lazy
        :load="load"
        :indent="8"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="id" label="编号" width="140" />
        <el-table-column prop="name" label="城市名称" min-width="220">
          <template #default="scope">
            <div class="geo-name-cell" :class="`geo-name-cell--level-${scope.row.level}`">
              <span class="geo-name-mark" />
              <span class="geo-name-text">{{ displayGeoName(scope.row) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="城市等级" width="100">
          <template #default="scope">
            <el-tag :type="levelTagType(scope.row.level)">{{ levelMap[scope.row.level] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="城市编码" width="120" />
        <el-table-column prop="geocode" label="地理编码" width="100" />
        <el-table-column prop="latitude" label="经度" width="100" />
        <el-table-column prop="longitude" label="纬度" width="100" />
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column width="220" label="操作">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button v-if="scope.row.level < 5" type="primary" link @click="handleAdd(scope.row)">新增子地区</el-button>
            <el-button type="primary" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="geoDialog" title="城市管理" width="600px">
      <el-form v-model="geo" label-width="120px">
        <el-form-item label="城市名称">
          <el-input v-model="geo.name" />
        </el-form-item>
        <el-form-item label="名称(多语言)">
          <MultiLangEditor
            :model="nameI18nObj"
            :languages="enabledLangs"
            title="城市名称多语言"
          />
        </el-form-item>
        <el-form-item label="城市编码">
          <el-input v-model="geo.code" :disabled="geo.id !== 0" />
        </el-form-item>
        <el-form-item label="地理编码">
          <el-input v-model="geo.geocode" />
        </el-form-item>
        <el-form-item label="经度">
          <el-input v-model="geo.latitude" />
        </el-form-item>
        <el-form-item label="纬度">
          <el-input v-model="geo.longitude" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="geo.sort" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="clear">取消</el-button>
        <el-button type="primary" @click="enter">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getGeos, getGeo, editGeo, createGeo, deleteGeo } from '@/plugin/geo/api/geo.js'
import { getEnabledLanguages } from '@/api/client/language'
import { ElMessageBox, ElMessage } from 'element-plus'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'

// 多语言编辑
const nameI18nObj = reactive({})
const enabledLangs = ref([])

const loadEnabledLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0 && res.data) {
      enabledLangs.value = Array.isArray(res.data) ? res.data : (res.data.list || [])
    }
  } catch(e) {}
}

const parseNameI18n = (jsonStr) => {
  Object.keys(nameI18nObj).forEach(k => delete nameI18nObj[k])
  try {
    const parsed = JSON.parse(jsonStr || '{}')
    Object.assign(nameI18nObj, parsed)
  } catch { /* ignore */ }
}

const serializeNameI18n = () => {
  const obj = {}
  Object.entries(nameI18nObj).forEach(([rawCode, rawText]) => {
    const code = String(rawCode || '').trim()
    if (!code) return
    const text = String(rawText ?? '').trim()
    if (!text) return
    obj[code] = text
  })
  return JSON.stringify(obj)
}

const displayGeoName = (row) => {
  const raw = row?.nameI18n
  let parsed = null
  if (raw && typeof raw === 'string') {
    try {
      parsed = JSON.parse(raw)
    } catch {
      parsed = null
    }
  } else if (raw && typeof raw === 'object') {
    parsed = raw
  }

  if (parsed && typeof parsed === 'object') {
    const codes = []
    const pushCode = (code) => {
      if (!code || codes.includes(code)) {
        return
      }
      codes.push(code)
    }

    pushCode('zh')
    pushCode('zh-TW')
    pushCode('ja')
    pushCode('ko')
    pushCode('th')
    pushCode('ar')
    pushCode('hi')
    pushCode('mn')
    pushCode('vi')
    for (const lang of enabledLangs.value) {
      pushCode(lang?.code)
    }
    pushCode('en')

    for (const code of codes) {
      const text = String(parsed[code] ?? '').trim()
      if (text) {
        return text
      }
    }
    const firstText = Object.values(parsed).find(item => String(item ?? '').trim())
    if (firstText) {
      return String(firstText)
    }
  }

  return String(row?.name ?? '')
}

onMounted(() => {
  loadEnabledLangs()
})

const loadMap = new Map()
const mainTable = ref()// table的ref

const treeData = ref([

])

const baseGeo = {
  id: 0,
  name: '',
  nameI18n: '',
  level: 0,
  code: '0',
  geocode: '',
  latitude: '',
  longitude: '',
  sort: 0
}

const levelMap = {
  0: '国家',
  1: '一级行政区',
  2: '二级行政区',
  3: '三级行政区',
  4: '四级行政区',
  5: '五级行政区'
}

const levelTagType = (level) => {
  const typeMap = ['danger', 'warning', 'success', 'info', 'primary', '']
  return typeMap[level] || 'info'
}

const tempObj = {}

const geoDialog = ref(false)

const handleAdd = async(row) => {
  tempObj.row = row
  geo.value = {
    ...baseGeo, level: row.level, parentCode: row.code
  }
  parseNameI18n('{}')
  geoDialog.value = true
}

const enter = async() => {
  geo.value.nameI18n = serializeNameI18n()
  if (geo.value.id) {
    const res = await editGeo(geo.value)
    if (res.code === 0) {
      reload(tempObj.row.parentCode)
      ElMessage.success('编辑成功')
      clear()
    }
  } else {
    const res = await createGeo(geo.value)
    if (res.code === 0) {
      reload(tempObj.row.code)
      ElMessage.success('添加成功')
      clear()
    }
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该地区吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await deleteGeo(row)
    if (res.code === 0) {
      reload(row.parentCode)
      ElMessage.success('删除成功')
    }
  })
}

const clear = () => {
  geoDialog.value = false
  tempObj.row = {}
}

const geo = ref({ ...baseGeo })

const handleEdit = async(row) => {
  tempObj.row = row
  const res = await getGeo({ id: row.id, level: row.level })
  if (res.code === 0) {
    geo.value = res.data
    parseNameI18n(res.data.nameI18n)
    geoDialog.value = true
  }
}

const getTreeData = async(level, code) => {
  const res = await getGeos({ level, code })
  if ((level === '' || level === null || level === undefined) && Array.isArray(res.data)) {
    res.data.forEach((element) => {
      element.hasChildren = true
    })
  }
  return res.data
}

const init = async() => {
  treeData.value = await getTreeData('', '0')
}

init()

const reload = (code) => {
  code = code || 0
  if (loadMap.get(code)) {
    const { row, treeNode, resolve } = loadMap.get(code)
    mainTable.value.store.states.lazyTreeNodeMap[code] = []// 清空节点的数据
    load(row, treeNode, resolve)
  } else {
    // 根节点不在 loadMap 中，直接重新加载整棵树
    init()
  }
}

const load = async(row, treeNode, resolve) => {
  loadMap.set(row.code, { row, treeNode, resolve })
  resolve(await getTreeData(row.level, row.code) || [])
}

</script>

<style scoped>
.geo-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 10px;
  border-radius: 8px;
  background: linear-gradient(90deg, #1f2937 0%, #334155 100%);
  border: 1px solid rgba(255, 255, 255, 0.16);
}

.geo-name-mark {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  flex: 0 0 auto;
  background: #94a3b8;
  box-shadow: 0 0 0 1px rgba(15, 23, 42, 0.18);
}

.geo-name-text {
  font-size: 14px;
  line-height: 1.45;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 0.1px;
  text-rendering: optimizeLegibility;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
}

.geo-name-cell--level-0 .geo-name-mark { background: #ef4444; }
.geo-name-cell--level-0 .geo-name-text {
  font-size: 15px;
  font-weight: 700;
}

.geo-name-cell--level-1 .geo-name-mark { background: #f59e0b; }
.geo-name-cell--level-1 .geo-name-text { font-weight: 600; }

.geo-name-cell--level-2 .geo-name-mark { background: #10b981; }
.geo-name-cell--level-3 .geo-name-mark { background: #0ea5e9; }
.geo-name-cell--level-4 .geo-name-mark { background: #8b5cf6; }
.geo-name-cell--level-5 .geo-name-mark { background: #64748b; }

.geo-tree-hint {
  color: #0f172a;
  font-size: 13px;
  font-weight: 500;
  line-height: 1.5;
  align-self: center;
  margin-left: 12px;
}
.geo-name-cell--level-3,
.geo-name-cell--level-4,
.geo-name-cell--level-5 {
  margin-left: 2px;
}
</style>
