<template>
  <view class="clothes-page">
    <view class="header">
      <view class="search-row">
        <input class="search-input" v-model="searchName" :placeholder="$t('clothesSearchPlaceholder')" confirm-type="search" @confirm="onSearch" />
        <view class="search-btn" @tap="onSearch">{{ $t('searchAction') }}</view>
      </view>
      <view class="gender-tabs">
        <view class="gender-tab" :class="{ active: genderTab === 'female' }" @tap="switchGender('female')">{{ $t('clothesGenderFemale') }}</view>
        <view class="gender-tab" :class="{ active: genderTab === 'male' }" @tap="switchGender('male')">{{ $t('clothesGenderMale') }}</view>
      </view>
    </view>

    <scroll-view class="list-wrap" scroll-y>
      <view class="grid">
        <view class="card" v-for="item in displayList" :key="item.ID || item.id || item.name">
          <image class="card-image" :src="mainImage(item)" mode="aspectFill" />
          <view class="card-body">
            <text class="card-name">{{ goodName(item) || $t('unnamedGoods') }}</text>
            <view class="tag-list">
              <text class="tag" v-for="(tag, idx) in extractTags(item)" :key="idx">{{ tag }}</text>
            </view>
            <view class="btn-row">
              <view class="btn tryon" @tap="chooseTryon(item)">{{ $t('tryOnAction') }}</view>
              <view class="btn buy" @tap="goDetail(item)">{{ $t('buyAction') }}</view>
            </view>
          </view>
        </view>
      </view>

      <view class="load-more-wrap">
        <view class="load-more" v-if="hasMore" @tap="loadMore">{{ $t('loadMoreAction') }}</view>
        <text class="all-loaded" v-else>{{ $t('allLoadedText') }}</text>
      </view>
      <view style="height: 40rpx"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getGoodList } from '@/api/homePage.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { localText } from '@/utils/i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { setSelectedClothes } from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = langStore.$t
const searchName = ref('')
const genderTab = ref('female')
const page = ref(1)
const pageSize = 10
const hasMore = ref(true)
const loading = ref(false)
const goodsList = ref([])
const displayList = ref([])

const goodName = (item) => {
  return localText(item?.nameI18n || item?.name || item?.titleI18n || item?.title, langStore.locale)
}

const categoryText = (item) => {
  return localText(item?.categoryNameI18n || item?.categoryName, langStore.locale)
}

const genderHit = (item, gender) => {
  const text = [
    item.gender,
    item.sex,
    goodName(item),
    categoryText(item),
    Array.isArray(item.tags) ? item.tags.join(' ') : item.tags,
  ].filter(Boolean).join(' ').toLowerCase()

  const maleHit = /(男|man|male|men)/i.test(text)
  const femaleHit = /(女|woman|female|women)/i.test(text)

  if (gender === 'male') {
    if (femaleHit && !maleHit) return false
    return true
  }
  if (maleHit && !femaleHit) return false
  return true
}

const refreshDisplay = () => {
  displayList.value = goodsList.value.filter(item => genderHit(item, genderTab.value))
}

const mainImage = (item) => {
  if (item.externalImagePath) return getExternalUrl(item.externalImagePath)
  return getUrl(item.imageUrl || item.picture || item.image || '')
}

const extractTags = (item) => {
  const tags = []
  if (Array.isArray(item.tags)) tags.push(...item.tags)
  if (typeof item.tags === 'string' && item.tags) tags.push(item.tags)
  if (item.tag) tags.push(item.tag)
  const category = categoryText(item)
  if (category) tags.push(category)
  return tags.filter(Boolean).slice(0, 2)
}

const loadList = async (reset = false) => {
  if (loading.value) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      hasMore.value = true
      goodsList.value = []
    }

    const res = await getGoodList({
      page: page.value,
      pageSize,
      name: searchName.value,
      keyword: searchName.value,
      sort: 'id',
      order: 'desc',
    })

    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data.list) ? res.data.list : []
      if (reset) {
        goodsList.value = list
      } else {
        goodsList.value = goodsList.value.concat(list)
      }

      const total = Number(res.data.total || 0)
      if (total > 0) {
        hasMore.value = goodsList.value.length < total
      } else {
        hasMore.value = list.length >= pageSize
      }
      refreshDisplay()
    }
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  loadList(true)
}

const switchGender = (gender) => {
  genderTab.value = gender
  refreshDisplay()
}

const loadMore = async () => {
  if (!hasMore.value) return
  page.value += 1
  await loadList(false)
}

const chooseTryon = (item) => {
  const image = item.externalImagePath ? getExternalUrl(item.externalImagePath) : (item.imageUrl || item.picture || '')
  if (!image) {
    uni.showToast({ title: $t('goodsImageMissing'), icon: 'none' })
    return
  }
  const text = `${goodName(item) || ''}${categoryText(item) || ''}`
  const isLower = /(裤|裙|下装|pants|skirt|bottom)/i.test(text)
  setSelectedClothes(isLower ? { lowerImage: image } : { upperImage: image })
  uni.switchTab({ url: '/pages/tabBar/index' })
}

const goDetail = (item) => {
  const goodID = item.ID || item.id
  if (!goodID) {
    uni.showToast({ title: $t('goodsInfoInvalid'), icon: 'none' })
    return
  }
  uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${goodID}` })
}

onShow(() => {
  if (goodsList.value.length === 0) {
    loadList(true)
  }
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.clothes-page {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
  color: #0f172a;
  padding-top: calc(var(--status-bar-height, 0px) + 16rpx);
}

.header {
  padding: 0 20rpx 16rpx;
}

.search-row {
  display: flex;
  gap: 12rpx;
}

.search-input {
  flex: 1;
  height: 72rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.95);
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  padding: 0 24rpx;
  color: #0f172a;
  font-size: 24rpx;
}

.search-btn {
  width: 120rpx;
  height: 72rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  color: #fff;
  box-shadow: 0 10rpx 26rpx rgba(37, 99, 235, 0.28);
}

.gender-tabs {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
}

.gender-tab {
  flex: 1;
  height: 64rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(15, 23, 42, 0.58);
  background: rgba(255, 255, 255, 0.8);
}

.gender-tab.active {
  color: #0f172a;
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(219, 234, 254, 0.75);
}

.list-wrap {
  height: calc(100vh - var(--status-bar-height, 0px) - 170rpx);
  padding: 0 20rpx;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16rpx;
}

.card {
  background: rgba(255, 255, 255, 0.95);
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 14rpx;
  overflow: hidden;
  box-shadow: 0 8rpx 22rpx rgba(15, 23, 42, 0.05);
}

.card-image {
  width: 100%;
  height: 260rpx;
}

.card-body {
  padding: 12rpx;
}

.card-name {
  font-size: 24rpx;
  min-height: 68rpx;
  line-height: 34rpx;
}

.tag-list {
  min-height: 42rpx;
}

.tag {
  display: inline-block;
  margin-right: 8rpx;
  margin-bottom: 8rpx;
  padding: 4rpx 10rpx;
  border-radius: 999rpx;
  font-size: 18rpx;
  background: rgba(219, 234, 254, 0.8);
  color: rgba(30, 64, 175, 0.9);
}

.btn-row {
  display: flex;
  gap: 8rpx;
}

.btn {
  flex: 1;
  height: 54rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
}

.btn.tryon {
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
}

.btn.buy {
  background: rgba(15, 23, 42, 0.06);
  color: #0f172a;
}

.load-more-wrap {
  margin-top: 24rpx;
  text-align: center;
}

.load-more {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 280rpx;
  height: 66rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.92);
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  color: #0f172a;
  font-size: 24rpx;
}

.all-loaded {
  color: rgba(15, 23, 42, 0.45);
  font-size: 22rpx;
}
</style>
