<template>
	<view class="login">
		<view class="login-view">
			<view class="t-login">
				<form class="cl">
					<view class="form-box">
            <view class="t-a">
              <text class="txt">账号</text>
              <input placeholder="请输入您的账号" maxlength="12" v-model="form.username" />
            </view>
            <view class="t-a">
              <text class="txt">密码</text>
              <input type="password" name="code" maxlength="18" placeholder="请输入您的密码" v-model="form.password" />
            </view>
            <view class="t-a">
              <text class="txt">重复密码</text>
              <input type="password" name="code" maxlength="18" placeholder="请输入您的密码" v-model="form.rePassword" />
            </view>
          </view>
				</form>
        <view class="button-container">
          <button class="register-btn" @tap="registerFunc()">注 册</button>
          <button class="login-btn" @tap="toLogin()">前往登录</button>
        </view>
			</view>
		</view>
	</view>
</template>
<script setup>
	import {
		reactive,
		ref
	} from 'vue';

	import {
		register
	} from "@/api/base.js"

	import {useUserStore} from "@/pinia/modules/user.js"
	import { myRouter } from "@/utils/permission";

		const userStore = useUserStore()
		const token = userStore.token || ''

		if(token){
			myRouter("/pages/tabBar/index",true)
		}


	const form = reactive({
		username: "",
		password: "",
		rePassword: ""
	})

	const toLogin = () => {
		uni.navigateTo({
			url: '/pages/user/login'
		})
		// myRouter("/pages/user/login",true)
	}
	//去登陆
	const registerFunc = async () => {
		const res = await register(form)
		if(res.code === 0){
			uni.showToast({
				title: '注册成功',
				icon: 'none'
			});
			toLogin()
		}
	}
</script>
<style>
.login {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 80rpx);
  background: url("@/static/background.jpg") no-repeat center center;
  background-size: cover;
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
