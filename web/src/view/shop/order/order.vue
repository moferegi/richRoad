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

        <el-form-item label="购买者ID" prop="userID">

             <el-input v-model.number="searchInfo.userID" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        >

        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />
        <el-table-column align="left" label="日期" width="180" sortable>
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="购买者ID" prop="userID" width="120" />
        <el-table-column align="left" label="订单价格（分）" prop="totalPrice" width="200" />
        <el-table-column align="left" label="订单状态" prop="status" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.status,orderStatusOptions) }}
            </template>
        </el-table-column>
          <el-table-column align="left" label="收货人" prop="name" width="120" />
          <el-table-column align="left" label="收货电话" prop="phone" width="120" />
          <el-table-column align="left" label="收货地址" prop="addr" width="120">
            <template #default="{row}">
               {{ row.province }}/{{ row.city }}/{{ row.area }}/{{ row.street }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="支付方式" prop="payMethod" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.payMethod === 'qrcode'" type="success" size="small">扫码支付</el-tag>
              <el-tag v-else-if="scope.row.payMethod === 'contact'" type="warning" size="small">客服收款</el-tag>
              <el-tag v-else type="info" size="small">未选择</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="快递单号" prop="express" width="160" />

          <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
                查看详情
            </el-button>
              <el-button v-if="scope.row.status === '0'" type="success" link class="table-button" @click="confirmPaymentFunc(scope.row)">确认收款</el-button>
              <el-button v-if="scope.row.status === '6'" type="danger" link class="table-button" @click="refundOrderFunc(scope.row)">同意退款</el-button>
              <el-button v-if="scope.row.status === '1'" type="primary" link icon="van" class="table-button" @click="sendOut(scope.row)">发货</el-button>
              <el-button v-if="scope.row.status === '2'" type="primary" link icon="van" class="table-button" @click="checkRoutersFunc(scope.row)">查看物流</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateOrderFunc(scope.row)">变更</el-button>
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
    <!-- 编辑订单抽屉 -->
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
            <el-form-item label="购买者ID:"  prop="userID" >
              <el-input v-model.number="formData.userID" :clearable="true" placeholder="请输入购买者ID" />
            </el-form-item>
            <el-form-item label="订单价格（分）:"  prop="totalPrice" >
              <el-input v-model.number="formData.totalPrice" :clearable="true" placeholder="请输入订单价格（分）" />
            </el-form-item>
            <el-form-item label="订单状态:"  prop="status" >
              <el-select v-model="formData.status" placeholder="请选择订单状态" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in orderStatusOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-form>
    </el-drawer>

    <!-- 订单详情抽屉 -->
    <el-drawer size="1000" v-model="detailDialogVisible" :show-close="false" :before-close="closeDetailDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">订单详情</span>
          <el-button @click="closeDetailDialog">关 闭</el-button>
        </div>
      </template>

      <div v-if="orderDetail" class="order-detail-container">
        <!-- 订单基本信息 -->
        <el-card class="mb-4" shadow="never">
          <template #header>
            <span class="text-base font-medium">订单信息</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">订单ID:</span>
                <span class="value">{{ orderDetail.ID }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">购买者ID:</span>
                <span class="value">{{ orderDetail.userID }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">订单状态:</span>
                <el-tag :type="getStatusType(orderDetail.status)">{{ filterDict(orderDetail.status, orderStatusOptions) }}</el-tag>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">创建时间:</span>
                <span class="value">{{ formatDate(orderDetail.CreatedAt) }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">支付方式:</span>
                <el-tag v-if="orderDetail.payMethod === 'qrcode'" type="success" size="small">扫码支付</el-tag>
                <el-tag v-else-if="orderDetail.payMethod === 'contact'" type="warning" size="small">客服收款</el-tag>
                <el-tag v-else type="info" size="small">未选择</el-tag>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">付款时间:</span>
                <span class="value">{{ orderDetail.paidAt ? formatDate(orderDetail.paidAt) : '未付款' }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">优惠券编号:</span>
                <span class="value">{{ orderDetail.couponNum || '无' }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">快递单号:</span>
                <span class="value">{{ orderDetail.express || '未发货' }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">关闭时间:</span>
                <span class="value">{{ orderDetail.closeTime ? formatDate(orderDetail.closeTime) : '-' }}</span>
              </div>
            </el-col>
          </el-row>
        </el-card>

        <!-- 价格信息 -->
        <el-card class="mb-4" shadow="never">
          <template #header>
            <span class="text-base font-medium">价格信息</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">原价:</span>
                <span class="value price">¥{{ (orderDetail.originPrice / 100).toFixed(2) }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">优惠金额:</span>
                <span class="value discount">-¥{{ (orderDetail.discount / 100).toFixed(2) }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">实际金额:</span>
                <span class="value total-price">¥{{ (orderDetail.totalPrice / 100).toFixed(2) }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">使用积分:</span>
                <span class="value">{{ orderDetail.usePoints ? `是 (${orderDetail.pointsUsed || 0} 积分)` : '否' }}</span>
              </div>
            </el-col>
          </el-row>
        </el-card>

        <!-- 收货信息 -->
        <el-card class="mb-4" shadow="never">
          <template #header>
            <span class="text-base font-medium">收货信息</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">收件人:</span>
                <span class="value">{{ orderDetail.name }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="detail-item">
                <span class="label">联系电话:</span>
                <span class="value">{{ orderDetail.phone }}</span>
              </div>
            </el-col>
            <el-col :span="24">
              <div class="detail-item">
                <span class="label">收货地址:</span>
                <span class="value">{{ orderDetail.province }}/{{ orderDetail.city }}/{{ orderDetail.area }}/{{ orderDetail.street }}</span>
              </div>
            </el-col>
          </el-row>
        </el-card>

        <!-- 退款信息 -->
        <el-card v-if="orderDetail?.refundReason || orderDetail?.status === '6' || orderDetail?.status === '5'" class="mb-4" shadow="never">
          <template #header>
            <span class="text-base font-medium">退款信息</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">退款状态:</span>
                <el-tag :type="getStatusType(orderDetail.status)">
                  {{ filterDict(orderDetail.status, orderStatusOptions) }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="16">
              <div class="detail-item">
                <span class="label">退款原因:</span>
                <span class="value">{{ orderDetail.refundReason || '无' }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">申请时间:</span>
                <span class="value">{{ orderDetail.refundAppliedAt ? formatDate(orderDetail.refundAppliedAt) : '无' }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">处理时间:</span>
                <span class="value">{{ orderDetail.refundHandledAt ? formatDate(orderDetail.refundHandledAt) : '未处理' }}</span>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="detail-item">
                <span class="label">处理备注:</span>
                <span class="value">{{ orderDetail.refundRemark || '无' }}</span>
              </div>
            </el-col>
            <el-col :span="24" v-if="getRefundImages(orderDetail.refundImages).length">
              <div class="detail-item">
                <span class="label">退款图片:</span>
                <div class="refund-images">
                  <el-image
                    v-for="(img, index) in getRefundImages(orderDetail.refundImages)"
                    :key="index"
                    :src="img"
                    style="width: 60px; height: 60px; margin-right: 8px"
                    fit="cover"
                    :preview-src-list="getRefundImages(orderDetail.refundImages)"
                  />
                </div>
              </div>
            </el-col>
          </el-row>
          <div v-if="orderDetail.status === '6'" class="refund-action">
            <el-button type="danger" @click="refundOrderFunc(orderDetail)">同意退款</el-button>
          </div>
        </el-card>

        <!-- 商品详情 -->
        <el-card shadow="never">
          <template #header>
            <span class="text-base font-medium">商品详情</span>
          </template>
          <el-table :data="orderDetail.detail" style="width: 100%" border>
             <el-table-column label="商品名称" width="200">
               <template #default="scope">
                 {{ scope.row.good?.title || '未知商品' }}
               </template>
             </el-table-column>
             <el-table-column label="SKU规格" width="200">
               <template #default="scope">
                 <div v-if="scope.row.sku">
                   <div>{{ scope.row.sku.name }}</div>
                   <div class="text-xs text-gray-500" v-if="scope.row.sku.specs && scope.row.sku.specs.length > 0">
                     <span v-for="(spec, index) in scope.row.sku.specs" :key="index">
                       {{ spec.label }}: {{ spec.value }}
                       <span v-if="index < scope.row.sku.specs.length - 1">; </span>
                     </span>
                   </div>
                 </div>
                 <span v-else>无规格信息</span>
               </template>
             </el-table-column>
             <el-table-column prop="quantity" label="购买数量" width="100" align="center" />
             <el-table-column prop="price" label="单价" width="120" align="center">
               <template #default="scope">
                 ¥{{ (scope.row.price / 100).toFixed(2) }}
               </template>
             </el-table-column>
             <el-table-column label="小计" width="120" align="center">
               <template #default="scope">
                 ¥{{ ((scope.row.price * scope.row.quantity) / 100).toFixed(2) }}
               </template>
             </el-table-column>
             <el-table-column label="商品图片" width="100" align="center">
               <template #default="scope">
                 <el-image 
                   v-if="scope.row.sku?.picture" 
                   :src="scope.row.sku.picture" 
                   style="width: 50px; height: 50px" 
                   fit="cover"
                   :preview-src-list="[scope.row.sku.picture]"
                 />
                 <el-image 
                   v-else-if="scope.row.good?.imageUrl" 
                   :src="scope.row.good.imageUrl" 
                   style="width: 50px; height: 50px" 
                   fit="cover"
                   :preview-src-list="[scope.row.good.imageUrl]"
                 />
                 <span v-else>无图片</span>
               </template>
             </el-table-column>
           </el-table>
        </el-card>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createOrder,
  deleteOrder,
  updateOrder,
  findOrder,
  getOrderList, checkRouters, refundOrder, confirmPayment
} from '@/api/shop/order'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'Order'
})

// 自动化生成的字典（可能为空）以及字段
const orderStatusOptions = ref([])
const formData = ref({
        userID: 0,
        totalPrice: 0,
        status: '',
        })


// 验证规则
const rule = reactive({
               userID : [{
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
  const table = await getOrderList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    orderStatusOptions.value = await getDictFunc('orderStatus')
}

// 获取需要的字典 可能为空 按需保留
setOptions()




// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteOrderFunc(row)
        })
    }


// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateOrderFunc = async(row) => {
    const res = await findOrder({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.reorder
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteOrderFunc = async (row) => {
    const res = await deleteOrder({ ID: row.ID })
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

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        userID: 0,
        totalPrice: 0,
        status: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createOrder(formData.value)
                  break
                case 'update':
                  res = await updateOrder(formData.value)
                  break
                default:
                  res = await createOrder(formData.value)
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

const sendOut = (row) => {
  ElMessageBox.prompt('请填写你的顺丰快递单号', '发货', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^SF[0-9]{13}$/,
    inputErrorMessage: '错误的顺丰单号',
  }).then(async(conf) => {
      await updateOrder({ ID: row.ID, status: '2',express:conf.value })
        ElMessage({
            type: 'success',
            message: '发货成功'
        })
    getTableData()
    })
}

const checkRoutersFunc = async (row) =>{
  const res = await checkRouters({ express: row.express })
  console.log(res)
}

const getRefundImages = (images) => {
  if (!images) return []
  let parsed = images
  if (typeof images === 'string') {
    try {
      parsed = JSON.parse(images)
    } catch (e) {
      parsed = [images]
    }
  }
  return ReturnArrImg(parsed)
}

const refundOrderFunc = (row) => {
  ElMessageBox.prompt('请输入退款备注（可选）', '退款处理', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPlaceholder: '退款备注'
  }).then(async({ value }) => {
    const res = await refundOrder({ orderID: row.ID, remark: value })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '退款成功'
      })
      if (detailDialogVisible.value && orderDetail.value?.ID === row.ID) {
        await getDetails(row)
      }
      getTableData()
    }
  })
}

const confirmPaymentFunc = (row) => {
  ElMessageBox.confirm(`确认订单 #${row.ID} 已收款？`, '确认收款', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await confirmPayment({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '确认收款成功' })
      if (detailDialogVisible.value && orderDetail.value?.ID === row.ID) {
        await getDetails(row)
      }
      getTableData()
    }
  })
}

// 订单详情相关
const detailDialogVisible = ref(false)
const orderDetail = ref(null)

// 获取订单详情
const getDetails = async (row) => {
  try {
    const res = await findOrder({ ID: row.ID })
    if (res.code === 0) {
      orderDetail.value = res.data.reorder
      detailDialogVisible.value = true
    } else {
      ElMessage.error('获取订单详情失败')
    }
  } catch (error) {
    console.error('获取订单详情错误:', error)
    ElMessage.error('获取订单详情失败')
  }
}

// 关闭订单详情弹窗
const closeDetailDialog = () => {
  detailDialogVisible.value = false
  orderDetail.value = null
}

// 获取订单状态对应的标签类型
const getStatusType = (status) => {
  const statusMap = {
    '0': 'info',     // 待支付
    '1': 'warning',  // 待发货
    '2': 'primary',  // 已发货
    '3': 'success',  // 已完成
    '4': 'danger',   // 已取消
    '5': 'danger',   // 已退款
    '6': 'warning'   // 退款中
  }
  return statusMap[status] || 'info'
}

</script>

<style scoped>
.order-detail-container {
  padding: 0;
}

.detail-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  min-height: 32px;
}

.detail-item .label {
  font-weight: 500;
  color: #606266;
  min-width: 80px;
  margin-right: 12px;
}

.detail-item .value {
  color: #303133;
  flex: 1;
}

.detail-item .value.price {
  color: #409EFF;
  font-weight: 500;
}

.detail-item .value.discount {
  color: #F56C6C;
  font-weight: 500;
}

.detail-item .value.total-price {
  color: #E6A23C;
  font-weight: 600;
  font-size: 16px;
}

.mb-4 {
  margin-bottom: 16px;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #EBEEF5;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-table) {
  font-size: 14px;
}

:deep(.el-table th) {
  background-color: #F5F7FA;
  color: #606266;
  font-weight: 500;
}

:deep(.el-image) {
  border-radius: 4px;
  border: 1px solid #DCDFE6;
}

.refund-images {
  display: flex;
  flex-wrap: wrap;
}

.refund-action {
  margin-top: 12px;
  text-align: right;
}
</style>
