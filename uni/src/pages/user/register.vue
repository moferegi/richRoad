<template>
	<view class="login">
		<view class="login-view">
			<view class="t-login">
				<form class="cl">
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
					<button @tap="registerFunc()">注 册</button>
					<view class="reg" @tap="toLogin()">前往登录</view>
				</form>
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
	}

	.txt {
		font-size: 32upx;
		font-weight: bold;
		color: #333333;
	}

	.reg {
		font-size: 28upx;
		color: #fff;
		height: 90upx;
		line-height: 90upx;
		border-radius: 50upx;
		font-weight: bold;
		background: #f5f6fa;
		color: #000000;
		text-align: center;
		margin-top: 30upx;
	}

	.login-view {
		padding-top: 40upx;
		width: 100%;
		position: relative;
		background-color: #ffffff;
		border-radius: 8% 8% 0% 0;
	}

	.t-login {
		width: 600upx;
		margin: 0 auto;
		font-size: 28upx;
		padding-top: 80upx;
	}

	.t-login button {
		font-size: 28upx;
		background: #2796f2;
		color: #fff;
		height: 90upx;
		line-height: 90upx;
		border-radius: 50upx;
		font-weight: bold;
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
		margin: 150upx 0 0 0;
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
</style>
