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

        <el-form-item label="订单ID" prop="orderID">

             <el-input v-model.number="searchInfo.orderID" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="商品ID" prop="goodID">

             <el-input v-model.number="searchInfo.goodID" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="SKUID" prop="SKUID">

             <el-input v-model.number="searchInfo.SKUID" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="商家回复创建时间" prop="shopReplyAt">

            <template #label>
            <span>
              商家回复创建时间
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
            <el-date-picker v-model="searchInfo.startShopReplyAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endShopReplyAt ? time.getTime() > searchInfo.endShopReplyAt.getTime() : false"></el-date-picker>
            —
            <el-date-picker v-model="searchInfo.endShopReplyAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startShopReplyAt ? time.getTime() < searchInfo.startShopReplyAt.getTime() : false"></el-date-picker>

        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
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
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="日期" prop="createdAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="用户ID" prop="userID" width="120" />
        <el-table-column align="left" label="订单ID" prop="orderID" width="120" />
        <el-table-column align="left" label="商品ID" prop="goodID" width="120" />
        <el-table-column align="left" label="SKUID" prop="SKUID" width="120" />
        <el-table-column align="left" label="用户评分" prop="rating" width="120" />
        <el-table-column align="left" label="评论内容" prop="content" width="120" />
        <el-table-column align="left" label="商家回复" prop="shopReply" width="120" />
         <el-table-column align="left" label="商家回复创建时间" prop="shopReplyAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.shopReplyAt) }}</template>
         </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateCommentFunc(scope.row)">变更</el-button>
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
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
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
            <el-form-item label="用户ID:"  prop="userID" >
              <el-input v-model.number="formData.userID" :clearable="true" placeholder="请输入用户ID" />
            </el-form-item>
            <el-form-item label="订单ID:"  prop="orderID" >
              <el-input v-model.number="formData.orderID" :clearable="true" placeholder="请输入订单ID" />
            </el-form-item>
            <el-form-item label="商品ID:"  prop="goodID" >
              <el-input v-model.number="formData.goodID" :clearable="true" placeholder="请输入商品ID" />
            </el-form-item>
            <el-form-item label="SKUID:"  prop="SKUID" >
              <el-input v-model.number="formData.SKUID" :clearable="true" placeholder="请输入SKUID" />
            </el-form-item>
            <el-form-item label="用户评分:"  prop="rating" >
              <el-input v-model.number="formData.rating" :clearable="true" placeholder="请输入用户评分" />
            </el-form-item>
            <el-form-item label="评论内容:"  prop="content" >
              <el-input v-model="formData.content" :clearable="true"  placeholder="请输入评论内容" />
            </el-form-item>
            <el-form-item label="商家回复:"  prop="shopReply" >
              <el-input v-model="formData.shopReply" :clearable="true"  placeholder="请输入商家回复" />
            </el-form-item>
            <el-form-item label="商家回复创建时间:"  prop="shopReplyAt" >
              <el-date-picker v-model="formData.shopReplyAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
            </el-form-item>
          </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createComment,
  deleteComment,
  deleteCommentByIds,
  updateComment,
  findComment,
  getCommentList
} from '@/api/shop/comment'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'Comment'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        userID: undefined,
        orderID: undefined,
        goodID: undefined,
        SKUID: undefined,
        rating: undefined,
        content: '',
        shopReply: '',
        shopReplyAt: new Date(),
        })



// 验证规则
const rule = reactive({
               orderID : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               rating : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               content : [{
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
        shopReplyAt : [{ validator: (rule, value, callback) => {
        if (searchInfo.value.startShopReplyAt && !searchInfo.value.endShopReplyAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startShopReplyAt && searchInfo.value.endShopReplyAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startShopReplyAt && searchInfo.value.endShopReplyAt && (searchInfo.value.startShopReplyAt.getTime() === searchInfo.value.endShopReplyAt.getTime() || searchInfo.value.startShopReplyAt.getTime() > searchInfo.value.endShopReplyAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change' }],
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
  const table = await getCommentList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteCommentFunc(row)
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
      const res = await deleteCommentByIds({ IDs })
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
const updateCommentFunc = async(row) => {
    const res = await findComment({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteCommentFunc = async (row) => {
    const res = await deleteComment({ ID: row.ID })
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

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        userID: undefined,
        orderID: undefined,
        goodID: undefined,
        SKUID: undefined,
        rating: undefined,
        content: '',
        shopReply: '',
        shopReplyAt: new Date(),
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createComment(formData.value)
                  break
                case 'update':
                  res = await updateComment(formData.value)
                  break
                default:
                  res = await createComment(formData.value)
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
