import { defineStore } from 'pinia'
import { login, phoneLogin, getUserInfo } from '@/api/base'
import storage from '@/utils/storage'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: storage.get('x-token') || '',
    userInfo: storage.get('userInfo') || {}
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo.nickname || state.userInfo.username || ''
  },
  actions: {
    setToken(token) {
      this.token = token
      storage.set('x-token', token)
    },
    setUserInfo(info) {
      this.userInfo = info
      storage.set('userInfo', info)
    },
    clearUserInfo() {
      this.token = ''
      this.userInfo = {}
      storage.remove('x-token')
      storage.remove('userInfo')
    },
    // 用户名登录
    async loginIn(params) {
      const res = await login(params)
      if (res.code === 0) {
        this.setToken(res.data.token)
        await this.getInfo()
      }
      return res
    },
    // 手机号登录
    async phoneLoginIn(params) {
      const res = await phoneLogin(params)
      if (res.code === 0) {
        this.setToken(res.data.token)
        await this.getInfo()
      }
      return res
    },
    // 获取用户信息
    async getInfo() {
      try {
        const res = await getUserInfo()
        if (res.code === 0) {
          this.setUserInfo(res.data)
        }
        return res
      } catch (e) {
        return e
      }
    },
    // 退出登录
    loginOut() {
      this.clearUserInfo()
    }
  }
})
