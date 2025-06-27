<template>
	<view class="login">
  <view class="gradient-background"></view>
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
    if (!form.rePassword) {
      uni.showToast({
        title: '请再次输入您的密码',
        icon: 'none'
      });
      return;
    }
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
  position: relative;
}

.gradient-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
  background: radial-gradient(circle at 20% 90%, rgba(255, 76, 125, 0.8) 0%, rgba(255, 76, 125, 0.4) 10%, transparent 60%),
  radial-gradient(circle at 80% 40%, rgba(255, 76, 125, 0.6) 0%, rgba(255, 76, 125, 0.3) 30%, transparent 60%),
  radial-gradient(circle at 40% 40%, rgba(255, 76, 125, 0.7) 0%, rgba(255, 76, 125, 0.2) 35%, transparent 70%),
  radial-gradient(circle at 90% 40%, rgba(255, 76, 125, 0.5) 0%, rgba(255, 76, 125, 0.1) 40%, transparent 80%),
  linear-gradient(135deg, rgba(255, 76, 125, 0.1) 0%, rgba(255, 76, 125, 0.05) 50%, rgba(255, 76, 125, 0.1) 100%);

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

/* 注册按钮样式 */
.register-btn {
  flex: 1;
  font-size: 28upx;
  background: linear-gradient(135deg, #ff4c7d 0%, #ff6b9d 50%, #ff8fb3 100%);
  color: #fff;
  height: 90upx;
  line-height: 90upx;
  border-radius: 50upx;
  font-weight: bold;
  margin: 0;
  border: none;
  box-shadow: 0 8upx 20upx rgba(255, 76, 125, 0.3);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.register-btn:active {
  transform: translateY(2upx);
  box-shadow: 0 4upx 12upx rgba(255, 76, 125, 0.4);
}

/* 登录按钮样式 */
.login-btn {
  flex: 1;
  font-size: 28upx;
  height: 90upx;
  line-height: 90upx;
  border-radius: 50upx;
  font-weight: bold;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(255, 255, 255, 0.7) 100%);
  color: #ff4c7d;
  text-align: center;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2upx solid rgba(255, 76, 125, 0.3);
  box-shadow: 0 6upx 16upx rgba(255, 76, 125, 0.15);
  transition: all 0.3s ease;
}

.login-btn:active {
  transform: translateY(2upx);
  background: linear-gradient(135deg, rgba(255, 76, 125, 0.1) 0%, rgba(255, 76, 125, 0.05) 100%);
  border-color: rgba(255, 76, 125, 0.5);
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
