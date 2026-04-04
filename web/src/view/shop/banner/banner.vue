<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAt">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
      </el-form-item>
      
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />
        
        <el-table-column align="left" label="日期" width="180" sortable>
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
        <el-table-column align="left" label="排序" prop="sort" width="80" />
        <el-table-column align="left" label="轮播标题" prop="title" width="120" />
          <el-table-column label="图片" width="200">
              <template #default="scope">
                <el-image style="width: 100px; height: 100px" :src="scope.row.externalPath || getUrl(scope.row.src)" fit="cover"/>
              </template>
          </el-table-column>
        <el-table-column align="left" label="跳转链接" prop="href" width="120" />
        <el-table-column align="left" label="商品ID" prop="goodID" width="80" />
        <el-table-column align="left" label="启用" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'" size="small">{{ scope.row.isEnabled ? '开' : '关' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="外部图片路径" prop="externalPath" width="160" show-overflow-tooltip />
        <el-table-column align="left" label="遮罩" width="80">
          <template #default="scope">
            <el-tag v-if="scope.row.maskEnabled" type="success" size="small">开</el-tag>
            <el-tag v-else type="info" size="small">关</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
                查看详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateBannerFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-drawer size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'添加':'修改'}}</span>
                <div>
                  <el-button type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="轮播标题:"  prop="title" >
                  <el-input v-model="formData.title" :clearable="true"  placeholder="请输入轮播标题" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="排序:" prop="sort">
                  <el-input-number v-model="formData.sort" :min="0" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="跳转链接:" prop="href">
                  <el-input v-model="formData.href" :clearable="true" placeholder="请输入跳转链接" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="8">
                <el-form-item label="关联商品ID:" prop="goodID">
                  <el-input-number v-model="formData.goodID" :min="0" placeholder="商品ID(点击跳商品详情)" style="width:100%;" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="启用:" prop="isEnabled">
                  <el-switch v-model="formData.isEnabled" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="图片(上传):"  prop="src" >
                <SelectImage
                 v-model="formData.src"
                 file-type="image"
                />
            </el-form-item>
            <el-form-item label="外部图片路径(优先于上传图片):" prop="externalPath">
              <el-input v-model="formData.externalPath" :clearable="true" placeholder="https://example.com/banner.jpg" />
            </el-form-item>

            <el-divider content-position="left">遮罩设置</el-divider>
            <el-row :gutter="12">
              <el-col :span="4">
                <el-form-item label="遮罩开关:" prop="maskEnabled">
                  <el-switch v-model="formData.maskEnabled" />
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="遮罩高度(px):" prop="maskHeight">
                  <el-input-number v-model="formData.maskHeight" :min="0" :max="500" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="背景色:" prop="maskBgColor">
                  <el-input v-model="formData.maskBgColor" placeholder="rgba(0,0,0,0.5)" />
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="文字颜色:" prop="maskTextColor">
                  <el-color-picker v-model="formData.maskTextColor" show-alpha />
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="文字大小:" prop="maskTextSize">
                  <el-input-number v-model="formData.maskTextSize" :min="8" :max="72" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="遮罩文字(JSON多语言):" prop="maskText">
              <div v-if="enabledLangs.length" class="w-full">
                <div v-for="lang in enabledLangs" :key="lang.code" class="flex items-center mb-2">
                  <span class="w-16 text-right mr-2">{{ lang.code }}:</span>
                  <el-input v-model="maskTextI18n[lang.code]" :placeholder="`${lang.name} 遮罩文字`" />
                </div>
              </div>
              <el-input v-else v-model="formData.maskText" placeholder="JSON格式多语言文字" />
            </el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer size="800" v-model="detailShow" :before-close="closeDetailShow" destroy-on-close>
          <template #header>
             <div class="flex justify-between items-center">
               <span class="text-lg">查看详情</span>
             </div>
         </template>
        <el-descriptions :column="1" border>
                <el-descriptions-item label="轮播标题">
                        {{ formData.title }}
                </el-descriptions-item>
                <el-descriptions-item label="图片">
                        <el-image style="width: 50px; height: 50px" :preview-src-list="ReturnArrImg(formData.externalPath || formData.src)" :src="formData.externalPath || getUrl(formData.src)" fit="cover" />
                </el-descriptions-item>
                <el-descriptions-item label="跳转链接">
                        {{ formData.href }}
                </el-descriptions-item>
                <el-descriptions-item label="外部图片路径">
                        {{ formData.externalPath || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="排序">
                        {{ formData.sort }}
                </el-descriptions-item>
                <el-descriptions-item label="遮罩开关">
                        {{ formData.maskEnabled ? '开' : '关' }}
                </el-descriptions-item>
                <el-descriptions-item v-if="formData.maskEnabled" label="遮罩高度">
                        {{ formData.maskHeight }}px
                </el-descriptions-item>
                <el-descriptions-item v-if="formData.maskEnabled" label="遮罩背景色">
                        <span :style="{ backgroundColor: formData.maskBgColor, padding: '2px 12px', borderRadius: '4px' }">{{ formData.maskBgColor }}</span>
                </el-descriptions-item>
                <el-descriptions-item v-if="formData.maskEnabled" label="遮罩文字">
                        {{ formData.maskText }}
                </el-descriptions-item>
                <el-descriptions-item v-if="formData.maskEnabled" label="遮罩文字颜色">
                        <span :style="{ color: formData.maskTextColor }">{{ formData.maskTextColor }}</span>
                </el-descriptions-item>
                <el-descriptions-item v-if="formData.maskEnabled" label="遮罩文字大小">
                        {{ formData.maskTextSize }}
                </el-descriptions-item>
        </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createBanner,
  deleteBanner,
  deleteBannerByIds,
  updateBanner,
  findBanner,
  getBannerList
} from '@/api/shop/banner'
import { getUrl } from '@/utils/image'
import { getEnabledLanguages } from '@/api/client/language'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({
    name: 'Banner'
})

// === i18n 多语言支持 ===
const enabledLangs = ref([])
const maskTextI18n = ref({})

const loadLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0) {
      enabledLangs.value = res.data || []
    }
  } catch (e) { /* ignore */ }
}

const parseMaskTextI18n = (jsonStr) => {
  const obj = {}
  if (!jsonStr) return obj
  try { return JSON.parse(jsonStr) } catch { return obj }
}

const serializeMaskTextI18n = (obj) => {
  const filtered = {}
  for (const [k, v] of Object.entries(obj)) {
    if (v) filtered[k] = v
  }
  return Object.keys(filtered).length ? JSON.stringify(filtered) : ''
}

onMounted(() => { loadLangs() })

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        title: '',
        src: "",
        href: '',
        externalPath: '',
        maskEnabled: false,
        maskHeight: 40,
        maskBgColor: 'rgba(0,0,0,0.5)',
        maskText: '',
        maskTextColor: '#FFFFFF',
        maskTextSize: 14,
        sort: 0,
        isEnabled: true,
        goodID: null,
        })


// 验证规则
const rule = reactive({
               title : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               src : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getBannerList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteBannerFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteBannerByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateBannerFunc = async(row) => {
    const res = await findBanner({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.rebanner
        maskTextI18n.value = parseMaskTextI18n(formData.value.maskText)
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteBannerFunc = async (row) => {
    const res = await deleteBanner({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)


// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findBanner({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.rebanner
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
          title: '',
          href: '',
          externalPath: '',
          maskEnabled: false,
          maskHeight: 40,
          maskBgColor: 'rgba(0,0,0,0.5)',
          maskText: '',
          maskTextColor: '#FFFFFF',
          maskTextSize: 14,
          sort: 0,
          }
}


// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    maskTextI18n.value = {}
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    maskTextI18n.value = {}
    formData.value = {
        title: '',
        href: '',
        externalPath: '',
        maskEnabled: false,
        maskHeight: 40,
        maskBgColor: 'rgba(0,0,0,0.5)',
        maskText: '',
        maskTextColor: '#FFFFFF',
        maskTextSize: 14,
        sort: 0,
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              // 序列化遮罩文字i18n
              if (enabledLangs.value.length) {
                formData.value.maskText = serializeMaskTextI18n(maskTextI18n.value)
              }
              let res
              switch (type.value) {
                case 'create':
                  res = await createBanner(formData.value)
                  break
                case 'update':
                  res = await updateBanner(formData.value)
                  break
                default:
                  res = await createBanner(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

</script>

<style>

</style>
