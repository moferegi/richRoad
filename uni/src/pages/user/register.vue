<template>
  <view class="nf-page">
    <view class="nf-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <text class="nf-back-icon">&#xe603;</text>
        </view>
        <text class="nf-navbar-title">{{ $t('navRegister') }}</text>
        <view class="nf-lang-btn" @tap="showLangPicker = true">
          <text class="nf-lang-label">{{ langLabel }}</text>
        </view>
      </view>
    </view>

    <view class="nf-container">
      <!-- 标题 -->
      <view class="nf-header">
        <text class="nf-title">{{ $t('registerBtn') }}</text>
        <view class="nf-title-line"></view>
      </view>

      <!-- 表单卡片 -->
      <view class="nf-card">
        <view class="nf-field">
          <text class="nf-label">{{ $t('account') }}</text>
          <input class="nf-input" :placeholder="$t('accountPlaceholder')" maxlength="12" v-model="form.username" />
        </view>
        <view class="nf-field">
          <text class="nf-label">{{ $t('password') }}</text>
          <input class="nf-input" type="password" maxlength="18" :placeholder="$t('passwordPlaceholder')" v-model="form.password" />
        </view>
        <view class="nf-field">
          <text class="nf-label">{{ $t('repeatPassword') }}</text>
          <input class="nf-input" type="password" maxlength="18" :placeholder="$t('passwordPlaceholder')" v-model="form.rePassword" />
        </view>
      </view>

      <!-- 按钮 -->
      <view class="nf-actions">
        <button class="nf-btn nf-btn-primary" @tap="registerFunc()">{{ $t('registerBtn') }}</button>
        <button class="nf-btn nf-btn-ghost" @tap="toLogin()">{{ $t('goToLogin') }}</button>
      </view>
    </view>

    <!-- 语言弹窗 -->
    <lang-switch v-model="showLangPicker" />
  </view>
</template>
<script setup>
	import {
		onLoad,
		onReady
	} from '@dcloudio/uni-app';

	import {
		reactive,
		ref,
		computed
	} from 'vue';

	import {
		register
	} from "@/api/base.js"

	import {useUserStore} from "@/pinia/modules/user.js"
	import { useLangStore } from '@/pinia/modules/lang.js'
	import langSwitch from '@/components/lang-switch/lang-switch.vue'

	const langStore = useLangStore()
	const $t = computed(() => langStore.$t)
	const langLabel = computed(() => {
	  const map = { zh: '中', en: 'EN', mn: 'MN' }
	  return map[langStore.locale] || '中'
	})
	const showLangPicker = ref(false)

	const goBack = () => {
	  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/index' }) })
	}

	import { myRouter } from "@/utils/permission";

		const userStore = useUserStore()
		const token = userStore.token || ''

		if(token){
			myRouter("/pages/tabBar/index",true)
		}


	const form = reactive({
		username: "",
		password: "",
		rePassword: "",
		inviteCode: ""
	})

	// 从URL参数获取邀请码
	onLoad((options) => {
		if (options && options.inviteCode) {
			form.inviteCode = options.inviteCode
		}
	})

	const toLogin = () => {
		uni.navigateTo({
			url: '/pages/user/login'
		})
	}

	const registerFunc = async () => {
    if (!form.username) {
      uni.showToast({
        title: $t.value('enterUsername'),
        icon: 'none'
      });
      return;
    }
    if (!form.password) {
      uni.showToast({
        title: $t.value('enterPassword'),
        icon: 'none'
      });
      return;
    }
    if (!form.rePassword) {
      uni.showToast({
        title: $t.value('enterRePassword'),
        icon: 'none'
      });
      return;
    }
		const res = await register(form)
		if(res.code === 0){
			uni.showToast({
				title: $t.value('registerSuccess'),
				icon: 'none'
			});
			toLogin()
		}
	}
</script>
<style lang="scss" scoped>
page { background-color: #000; }

.nf-page {
  min-height: 100vh;
  background: #000;
  position: relative;
  overflow: hidden;
}

.nf-bg {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 0;
  background:
    radial-gradient(ellipse at 85% 25%, rgba(229, 9, 20, 0.25) 0%, transparent 50%),
    radial-gradient(ellipse at 15% 65%, rgba(229, 9, 20, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 50% 90%, rgba(229, 9, 20, 0.1) 0%, transparent 40%),
    #000;
}

/* 自定义导航栏 */
.nf-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(20px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}

.nf-navbar-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;

  &:active {
    background: rgba(255, 255, 255, 0.1);
  }
}

.nf-back-icon {
  font-size: 36rpx;
  color: #fff;
  font-family: 'iconfont';
  &::before { content: '←'; font-family: inherit; }
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

.nf-lang-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  border: 2rpx solid rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;

  &:active {
    background: rgba(255, 255, 255, 0.15);
    transform: scale(0.92);
  }
}

.nf-lang-label {
  font-size: 22rpx;
  font-weight: 700;
  color: #fff;
}

.nf-container {
  position: relative;
  z-index: 1;
  padding: 0 60rpx;
  padding-top: calc(var(--status-bar-height, 44px) + 88rpx + 40rpx);
}

.nf-header {
  margin-bottom: 64rpx;
}

.nf-title {
  font-size: 52rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 4rpx;
}

.nf-title-line {
  width: 80rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #e50914, #ff6b6b);
  margin-top: 20rpx;
}

.nf-card {
  background: rgba(255, 255, 255, 0.05);
  border: 2rpx solid rgba(255, 255, 255, 0.08);
  border-radius: 28rpx;
  padding: 40rpx 36rpx;
  backdrop-filter: blur(20px);
  box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.4);
}

.nf-field {
  margin-bottom: 36rpx;

  &:last-child { margin-bottom: 0; }
}

.nf-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
  font-weight: 600;
  letter-spacing: 2rpx;
  text-transform: uppercase;
  margin-bottom: 12rpx;
  display: block;
}

.nf-input {
  width: 100%;
  height: 88rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 2rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 16rpx;
  padding: 0 24rpx;
  color: #fff;
  font-size: 30rpx;
  box-sizing: border-box;
  transition: border-color 0.3s;

  &:focus {
    border-color: rgba(229, 9, 20, 0.6);
  }
}

.nf-actions {
  margin-top: 56rpx;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.nf-btn {
  width: 100%;
  height: 96rpx;
  border-radius: 48rpx;
  font-size: 32rpx;
  font-weight: 700;
  letter-spacing: 4rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  margin: 0;
  transition: all 0.3s;

  &:active { transform: scale(0.97); }
  &::after { border: none; }
}

.nf-btn-primary {
  background: linear-gradient(135deg, #e50914 0%, #ff4d4d 100%);
  color: #fff;
  box-shadow: 0 8rpx 32rpx rgba(229, 9, 20, 0.4);

  &:active {
    box-shadow: 0 4rpx 16rpx rgba(229, 9, 20, 0.5);
  }
}

.nf-btn-ghost {
  background: transparent;
  color: rgba(255, 255, 255, 0.6);
  border: 2rpx solid rgba(255, 255, 255, 0.12);

  &:active {
    background: rgba(255, 255, 255, 0.05);
    color: #fff;
  }
}
</style>
