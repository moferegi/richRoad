import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia

export { useUserStore } from './modules/user'
export { useLangStore } from './modules/lang'
