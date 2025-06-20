<template>
  <view class="login">
    <image class="background-image" src="@/static/background.jpg" mode="aspectFill"></image>
    <view class="login-view">
      <view class="t-login">
        <form class="cl">
          <view class="form-box">
            <view class="t-a">
              <text class="txt">账号</text>
              <input name="username" placeholder="请输入您的账号" maxlength="11" v-model="form.username" />
            </view>
            <view class="t-a">
              <text class="txt">密码</text>
              <input type="password" name="password" maxlength="18" placeholder="请输入您的密码" v-model="form.password" />
            </view>
            <view class="t-a">
              <text class="txt">验证码</text>
              <input name="captcha" maxlength="18" placeholder="请输入您的验证码" v-model="form.captcha" />
              <image class="t-image" @tap="getCaptchaFunc()" :src="captchaImg"></image>
            </view>
          </view>
        </form>

        <!-- 按钮容器置底 -->
        <view class="button-container">
          <button class="btn login-btn" @tap="login()">登 录</button>
          <button class=" register-btn" @tap="toRegister()">注 册</button>
        </view>

        <!--				<view class="t-f"><text>—————— 第三方账号登录 ——————</text></view>
                <view class="t-e cl">
                  <view class="t-g" @tap="wxLogin()">
                    <image src="@/static/wx.png"></image>
                  </view>
                  <view class="t-g" @tap="zfbLogin()">
                    <image src="@/static/qq.png"></image>
                  </view>
                </view>-->
      </view>
    </view>
  </view>
</template>
<script setup>
	import {getCaptcha} from "@/api/base.js"
	import {
		reactive,
		ref
	} from 'vue';

	import {useUserStore} from "@/pinia/modules/user.js"

	const userStore = useUserStore()
/* 	const token = userStore.token || ''

	if(token){
		uni.navigateTo({
			url: '/pages/tabBar/index'
		})
		// myRouter("/pages/tabBar/index",true)
	} */

	const form = reactive({
		username: '', //手机号码
		password: '' ,//密码
		captcha: '',
		captchaId: ''
	})

	const captchaImg = ref("")

	const getCaptchaFunc = async () =>{
		const res = await getCaptcha()
		if(res.code === 0){
			captchaImg.value = res.data.picPath
			form.captchaId = res.data.captchaId
		}
	}

	getCaptchaFunc()


	//当前登录按钮操作
	const login = async () => {
		if (!form.username) {
			uni.showToast({
				title: '请输入您的用户名',
				icon: 'none'
			});
			return;
		}
		if (!form.password) {
			uni.showToast({
				title: '请输入您的密码',
				icon: 'none'
			});
			return;
		}

		if (!form.captcha) {
			uni.showToast({
				title: '请输入验证码',
				icon: 'none'
			});
			return;
		}

	    const flag = await userStore.loginIn(form)
		if(flag){
			uni.showToast({
				title: '登录成功',
			})
			uni.navigateTo({
				url: '/pages/tabBar/index'
			})
			// myRouter("/pages/tabBar/tabBar",true)
			return
		}
		getCaptchaFunc()
	}
	//注册按钮点击
	const toRegister = () => {
		uni.navigateTo({
			url: '/pages/user/register'
		})
		// myRouter("/pages/user/register",true)
	}
	//等三方微信登录
	const wxLogin = () => {
		uni.showToast({
			title: '微信登录',
			icon: 'none'
		});
	}
	//第三方支付宝登录
	const zfbLogin = () => {
		uni.showToast({
			title: '支付宝登录',
			icon: 'none'
		});
	}
</script>
<style>
.login {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 80rpx);
  position: relative;
}

.background-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
}

.txt {
  font-size: 32upx;
  font-weight: bold;
  color: #333333;
}

.login-view {
  padding-top: 40upx;
  width: 100%;
  position: relative;
  border-radius: 8% 8% 0% 0;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.t-login {
  width: 600upx;
  margin: 0 auto;
  font-size: 28upx;
  padding-top: 80upx;
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 表单区域 */
.cl {
  flex: 1;
}
.form-box{
  background-color: rgba(255, 255, 255, 0.6);
  padding:40rpx;
  border-radius: 40rpx;
}
/* 按钮容器样式 */
.button-container {
  display: flex;
  justify-content: space-between;
  gap: 20upx;
  margin-top: auto;
  padding-bottom: 40upx;
}

/* 登录按钮样式 */
.login-btn {
  flex: 1;
  font-size: 28upx;
  background: #2796f2;
  color: #fff;
  height: 90upx;
  line-height: 90upx;
  border-radius: 50upx;
  font-weight: bold;
  margin: 0;
  border: none;
}

/* 注册按钮样式 */
.register-btn {
  flex: 1;
  font-size: 28upx;
  height: 90upx;
  line-height: 90upx;
  border-radius: 50upx;
  font-weight: bold;
  background: #f5f6fa;
  color: #000000;
  text-align: center;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.t-login input {
  height: 90upx;
  line-height: 90upx;
  margin-bottom: 50upx;
  border-bottom: 1px solid #e9e9e9;
  font-size: 28upx;
}

.t-login .t-a {
  position: relative;
}

.t-b {
  padding: 130upx 0 0 70upx;
}

.t-login .t-c {
  position: absolute;
  right: 22upx;
  top: 22upx;
  background: #5677fc;
  color: #fff;
  font-size: 24upx;
  border-radius: 50upx;
  height: 50upx;
  line-height: 50upx;
  padding: 0 25upx;
}

.t-login .t-d {
  text-align: center;
  color: #999;
  margin: 80upx 0;
}

.t-login .t-e {
  text-align: center;
  width: 250upx;
  margin: 80upx auto 0;
}

.t-login .t-g {
  float: left;
  width: 50%;
}

.t-login .t-e image {
  width: 50upx;
  height: 50upx;
}

.t-login .t-f {
  text-align: center;
  margin: 100upx 0 0 0;
  color: #666;
}

.t-login .t-f text {
  margin-left: 20upx;
  color: #aaaaaa;
  font-size: 27upx;
}

.t-login .uni-input-placeholder {
  color: #aeaeae;
}

.cl {
  zoom: 1;
}

.cl:after {
  clear: both;
  display: block;
  visibility: hidden;
  height: 0;
  content: '\20';
}

.t-a {
  position: relative;
}

.t-image {
  position: absolute;
  width: 200upx;
  height: 60upx;
  top: 60upx;
  right: 0upx;
  z-index: 99;
}
</style>
