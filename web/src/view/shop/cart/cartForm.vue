<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="用户ID:" prop="userID">
          <el-input v-model.number="formData.userID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="商品ID:" prop="goodID">
          <el-input v-model.number="formData.goodID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="skuID:" prop="skuID">
          <el-input v-model.number="formData.skuID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="数量:" prop="quantity">
          <el-input v-model.number="formData.quantity" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createCart,
  updateCart,
  findCart
} from '@/api/shop/cart'

defineOptions({
    name: 'CartForm'
})

// 自动获取字典
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            userID: 0,
            goodID: 0,
            skuID: 0,
            quantity: 0,
        })
// 验证规则
const rule = reactive({
               userID : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               goodID : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               skuID : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               quantity : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findCart({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.recart
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createCart(formData.value)
               break
             case 'update':
               res = await updateCart(formData.value)
               break
             default:
               res = await createCart(formData.value)
               break
           }
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
