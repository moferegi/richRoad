
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="促销背景图:" prop="promotionImage">
    <SelectImage
     v-model="formData.promotionImage"
     file-type="image"
    />
</el-form-item>
        <el-form-item label="促销标题:" prop="title">
    <el-input v-model="formData.title" :clearable="true" placeholder="请输入促销标题" />
</el-form-item>
        <el-form-item label="促销文案:" prop="description">
    <el-input v-model="formData.description" :clearable="true" placeholder="请输入促销文案" />
</el-form-item>
        <el-form-item label="是否启用:" prop="isAction">
    <el-switch v-model="formData.isAction" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
        <el-form-item label="促销类别:" prop="category">
    <el-input v-model="formData.category" :clearable="true" placeholder="请输入促销类别" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createPromotion,
  updatePromotion,
  findPromotion
} from '@/api/shop/promotion'

defineOptions({
    name: 'PromotionForm'
})

// 自动获取字典
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            promotionImage: "",
            title: '',
            description: '',
            isAction: false,
            category: '',
        })
// 验证规则
const rule = reactive({
               promotionImage : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               title : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               isAction : [{
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
      const res = await findPromotion({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createPromotion(formData.value)
               break
             case 'update':
               res = await updatePromotion(formData.value)
               break
             default:
               res = await createPromotion(formData.value)
               break
           }
           btnLoading.value = false
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
