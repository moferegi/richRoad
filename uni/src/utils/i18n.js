/**
 * 多语言国际化工具
 * 支持：中文(zh)、英语(en)、蒙古国语(mn)
 */

const messages = {
  zh: {
    // tabBar/index.vue
    searchPlaceholder: '请输入您想搜索的商品',
    searchEmpty: '请输入搜索关键词',
    reachedBottom: '已经到底啦',
    loading: '加载中...',
    noGoods: '暂无商品',
    all: '全部',

    // my/index.vue
    defaultUser: '默认用户',
    nicknamePlaceholder: '请填写昵称',
    loginPrompt: '请登录后查看',
    coupon: '优惠券',
    points: '积分',
    browseHistory: '浏览历史',
    myCollection: '我的收藏',
    onlineService: '在线客服',
    logout: '退出登录',
    goodsInfoIncomplete: '商品信息不完整',
    logoutSuccess: '退出成功',
    developing: '正在开发中...',
    switchLang: '切换语言',

    // kefu page
    kefuTitle: '在线客服',
    kefuOnline: '在线',
    kefuOffline: '离线',
    kefuBusy: '忙碌',
    kefuContact: '联系客服',
    kefuEmpty: '暂无客服',

    // player page
    playerEpisodes: '选集',
    playerDesc: '简介',
    playerRating: '评分',
    playerViews: '播放',
    playerCollect: '收藏',
    playerCollected: '已收藏',
    playerShare: '分享',
    playerFree: '免费',
    playerNowPlaying: '正在播放',
    playerNoEpisodes: '暂无剧集',
    playerEp: '第{{n}}集',
    playerNoVideo: '暂无视频资源',
    playerCollapse: '收起',
    playerExpand: '展开',
    playerEpCount: '{{n}}集',
    playerSpeed: '倍速',
    playerFullscreen: '全屏',
    playerExitFs: '退出',
    playerLandscape: '横屏',
    playerPortrait: '竖屏',
    playerFastFwd: '2x 快进中...',
    playerSwitchEp: '正在播放第{{n}}集',
    playerLoadFail: '加载失败',
    playerTapRetry: '点击重试',
    playerLoadSlow: '加载较慢，请稍候或切换网络',

    // seckilling.vue
    hotSelling: '近期热销',

    // no-pagin-row-good-list.vue
    sold: '售出',
    unit: '件',
    discountSmall: '小降',
    discountNormal: '优惠',
    discountGood: '特惠',
    discountGreat: '好价',
    discountLow: '低价',
    discountSpecial: '特价',
    discountDefault: '折扣',

    // login & register
    navLogin: '用户登录',
    navRegister: '用户注册',
    account: '账号',
    password: '密码',
    captcha: '验证码',
    repeatPassword: '重复密码',
    accountPlaceholder: '请输入您的账号',
    passwordPlaceholder: '请输入您的密码',
    captchaPlaceholder: '请输入您的验证码',
    loginBtn: '登 录',
    registerBtn: '注 册',
    goToRegister: '注 册',
    goToLogin: '前往登录',
    enterUsername: '请输入您的用户名',
    enterPassword: '请输入您的密码',
    enterCaptcha: '请输入验证码',
    enterRePassword: '请再次输入您的密码',
    loginSuccess: '登录成功',
    registerSuccess: '注册成功',
    wechatLogin: '微信登录',
    alipayLogin: '支付宝登录',

    // tabBar
    tabHome: '首页',
    tabCollect: '我的收藏',
    tabMy: '我的',

    // category-page
    categoryDetail: '分类详情',
    noProductsInCategory: '该分类暂无商品',
    loadFailed: '加载失败，请重试',
    allProducts: '全部商品',
    noOrderData: '暂无订单数据',

    // collect
    noCollectData: '还没有收藏商品',
    selfOperated: '自营',
    qualityAssured: '放心购',
    freeShipping: 'Plus免邮',
    collected: '已收藏',
    uncollected: '已取消收藏',
    loginFirst: '请登录后进行操作',
    noMoreData: '没有更多数据了',
    priceDrop: '小降',
    discount: '优惠',
    specialOffer: '特惠',
    goodPrice: '好价',
    lowPrice: '低价',
    bargain: '特价',
    saleTag: '折扣',

    // 语言选择弹窗
    selectLanguage: '选择语言',
    langZh: '中文',
    langEn: 'English',
    langMn: 'Монгол',
    cancel: '取消',
  },
  en: {
    searchPlaceholder: 'Search for products',
    searchEmpty: 'Please enter a keyword',
    reachedBottom: 'No more items',
    loading: 'Loading...',
    noGoods: 'No products yet',
    all: 'All',

    defaultUser: 'User',
    nicknamePlaceholder: 'Enter nickname',
    loginPrompt: 'Please login to view',
    coupon: 'Coupons',
    points: 'Points',
    browseHistory: 'History',
    myCollection: 'My Collection',
    onlineService: 'Customer Service',
    logout: 'Log Out',
    goodsInfoIncomplete: 'Product info incomplete',
    logoutSuccess: 'Logged out',
    developing: 'Coming soon...',
    switchLang: 'Language',

    kefuTitle: 'Customer Service',
    kefuOnline: 'Online',
    kefuOffline: 'Offline',
    kefuBusy: 'Busy',
    kefuContact: 'Contact',
    kefuEmpty: 'No agents available',

    playerEpisodes: 'Episodes',
    playerDesc: 'Description',
    playerRating: 'Rating',
    playerViews: 'Views',
    playerCollect: 'Collect',
    playerCollected: 'Collected',
    playerShare: 'Share',
    playerFree: 'Free',
    playerNowPlaying: 'Now Playing',
    playerNoEpisodes: 'No episodes',
    playerEp: 'EP {{n}}',
    playerNoVideo: 'No video available',
    playerCollapse: 'Less',
    playerExpand: 'More',
    playerEpCount: '{{n}} eps',
    playerSpeed: 'Speed',
    playerFullscreen: 'Full',
    playerExitFs: 'Exit',
    playerLandscape: 'Land',
    playerPortrait: 'Port',
    playerFastFwd: '2x Fast...',
    playerSwitchEp: 'Now playing EP {{n}}',
    playerLoadFail: 'Load failed',
    playerTapRetry: 'Tap to retry',
    playerLoadSlow: 'Loading slowly, please wait or switch network',

    hotSelling: 'Hot Selling',

    sold: 'Sold',
    unit: '',
    discountSmall: 'Sale',
    discountNormal: 'Deal',
    discountGood: 'Hot Deal',
    discountGreat: 'Best Price',
    discountLow: 'Low Price',
    discountSpecial: 'Special',
    discountDefault: 'Discount',

    account: 'Account',
    navLogin: 'Sign In',
    navRegister: 'Sign Up',
    password: 'Password',
    captcha: 'Captcha',
    repeatPassword: 'Confirm',
    accountPlaceholder: 'Enter your account',
    passwordPlaceholder: 'Enter your password',
    captchaPlaceholder: 'Enter captcha',
    loginBtn: 'Login',
    registerBtn: 'Register',
    goToRegister: 'Register',
    goToLogin: 'Go to Login',
    enterUsername: 'Please enter username',
    enterPassword: 'Please enter password',
    enterCaptcha: 'Please enter captcha',
    enterRePassword: 'Please confirm password',
    loginSuccess: 'Login successful',
    registerSuccess: 'Registered',
    wechatLogin: 'WeChat Login',
    alipayLogin: 'Alipay Login',

    tabHome: 'Home',
    tabCollect: 'Collection',
    tabMy: 'Me',

    categoryDetail: 'Category',
    noProductsInCategory: 'No products in this category',
    loadFailed: 'Load failed, please retry',
    allProducts: 'All Products',
    noOrderData: 'No order data',

    // collect
    noCollectData: 'No favorites yet',
    selfOperated: 'Official',
    qualityAssured: 'Verified',
    freeShipping: 'Free Ship',
    collected: 'Collected',
    uncollected: 'Removed',
    loginFirst: 'Please login first',
    noMoreData: 'No more data',
    priceDrop: 'Drop',
    discount: 'Deal',
    specialOffer: 'Offer',
    goodPrice: 'Hot',
    lowPrice: 'Low',
    bargain: 'Sale',
    saleTag: 'Off',

    selectLanguage: 'Select Language',
    langZh: '中文',
    langEn: 'English',
    langMn: 'Монгол',
    cancel: 'Cancel',
  },
  mn: {
    searchPlaceholder: 'Бараа хайх',
    searchEmpty: 'Түлхүүр үг оруулна уу',
    reachedBottom: 'Төгсгөлд хүрлээ',
    loading: 'Ачааллаж байна...',
    noGoods: 'Бараа байхгүй',
    all: 'Бүгд',

    defaultUser: 'Хэрэглэгч',
    nicknamePlaceholder: 'Нэрээ оруулна уу',
    loginPrompt: 'Нэвтэрнэ үү',
    coupon: 'Купон',
    points: 'Оноо',
    browseHistory: 'Түүх',
    myCollection: 'Миний цуглуулга',
    onlineService: 'Онлайн зөвлөгөө',
    logout: 'Гарах',
    goodsInfoIncomplete: 'Барааны мэдээлэл дутуу',
    logoutSuccess: 'Амжилттай гарлаа',
    developing: 'Тун удахгүй...',
    switchLang: 'Хэл солих',

    kefuTitle: 'Онлайн зөвлөгөө',
    kefuOnline: 'Онлайн',
    kefuOffline: 'Оффлайн',
    kefuBusy: 'Завгүй',
    kefuContact: 'Холбоо барих',
    kefuEmpty: 'Зөвлөгөө байхгүй',

    playerEpisodes: 'Анги',
    playerDesc: 'Танилцуулга',
    playerRating: 'Үнэлгээ',
    playerViews: 'Үзсэн',
    playerCollect: 'Хадгалах',
    playerCollected: 'Хадгалсан',
    playerShare: 'Хуваалцах',
    playerFree: 'Үнэгүй',
    playerNowPlaying: 'Тоглож байна',
    playerNoEpisodes: 'Анги байхгүй',
    playerEp: '{{n}}-р анги',
    playerNoVideo: 'Видео байхгүй',
    playerCollapse: 'Хураах',
    playerExpand: 'Дэлгэх',
    playerEpCount: '{{n}} анги',
    playerSpeed: 'Хурд',
    playerFullscreen: 'Бүтэн',
    playerExitFs: 'Гарах',
    playerLandscape: 'Хөндлөн',
    playerPortrait: 'Босоо',
    playerFastFwd: '2x Хурдан...',
    playerSwitchEp: '{{n}}-р анги тоглож байна',
    playerLoadFail: 'Ачаалж чадсангүй',
    playerTapRetry: 'Дахин оролдох',
    playerLoadSlow: 'Удаан ачаалж байна, түр хүлээнэ үү',

    hotSelling: 'Шилдэг борлуулалт',

    sold: 'Зарагдсан',
    unit: 'ш',
    discountSmall: 'Хямдрал',
    discountNormal: 'Урамшуулал',
    discountGood: 'Онцгой',
    discountGreat: 'Шилдэг үнэ',
    discountLow: 'Хямд',
    discountSpecial: 'Тусгай',
    discountDefault: 'Хөнгөлөлт',

    account: 'Бүртгэл',
    navLogin: 'Нэвтрэх',
    navRegister: 'Бүртгүүлэх',
    password: 'Нууц үг',
    captcha: 'Баталгаажуулах',
    repeatPassword: 'Давтах',
    accountPlaceholder: 'Бүртгэлээ оруулна уу',
    passwordPlaceholder: 'Нууц үгээ оруулна уу',
    captchaPlaceholder: 'Кодоо оруулна уу',
    loginBtn: 'Нэвтрэх',
    registerBtn: 'Бүртгүүлэх',
    goToRegister: 'Бүртгүүлэх',
    goToLogin: 'Нэвтрэх хуудас',
    enterUsername: 'Нэрээ оруулна уу',
    enterPassword: 'Нууц үгээ оруулна уу',
    enterCaptcha: 'Кодоо оруулна уу',
    enterRePassword: 'Нууц үгээ давтана уу',
    loginSuccess: 'Амжилттай нэвтэрлээ',
    registerSuccess: 'Амжилттай бүртгэгдлээ',
    wechatLogin: 'WeChat нэвтрэх',
    alipayLogin: 'Alipay нэвтрэх',

    tabHome: 'Нүүр',
    tabCollect: 'Дуртай',
    tabMy: 'Миний',

    categoryDetail: 'Ангилал',
    noProductsInCategory: 'Энэ ангилалд бараа байхгүй',
    loadFailed: 'Ачаалалт амжилтгүй',
    allProducts: 'Бүх бараа',
    noOrderData: 'Захиалга байхгүй',

    // collect
    noCollectData: 'Дуртай бараа байхгүй',
    selfOperated: 'Албан',
    qualityAssured: 'Баталгаат',
    freeShipping: 'Үнэгүй',
    collected: 'Хадгалсан',
    uncollected: 'Хассан',
    loginFirst: 'Нэвтэрнэ үү',
    noMoreData: 'Өгөгдөл байхгүй',
    priceDrop: 'Бууралт',
    discount: 'Хөнгөлөлт',
    specialOffer: 'Тусгай',
    goodPrice: 'Хямд',
    lowPrice: 'Хөнгө',
    bargain: 'Хямдрал',
    saleTag: 'Хөнгөлөлт',

    selectLanguage: 'Хэл сонгох',
    langZh: '中文',
    langEn: 'English',
    langMn: 'Монгол',
    cancel: 'Цуцлах',
  }
}

/**
 * 获取翻译文本
 * @param {string} key - 翻译键
 * @param {string} lang - 语言代码
 * @returns {string}
 */
export function t(key, lang) {
  const locale = lang || uni.getStorageSync('app-lang') || 'zh'
  return (messages[locale] && messages[locale][key]) || messages.zh[key] || key
}

/**
 * 解析多语言字段值
 * 支持两种格式：
 * 1. 普通字符串 → 原样返回
 * 2. JSON对象 {"zh":"中文","en":"English","mn":"Монгол"} → 按当前语言返回
 * @param {string|object} value - 字段值（可能是字符串或JSON对象）
 * @param {string} [lang] - 语言代码，不传则自动获取当前语言
 * @returns {string}
 */
export function localText(value, lang) {
  if (!value) return ''
  const locale = lang || uni.getStorageSync('app-lang') || 'zh'
  // 已经是对象
  if (typeof value === 'object') {
    return value[locale] || value['zh'] || Object.values(value)[0] || ''
  }
  // 字符串，尝试解析 JSON
  if (typeof value === 'string') {
    if (value.charAt(0) === '{') {
      try {
        const obj = JSON.parse(value)
        return obj[locale] || obj['zh'] || Object.values(obj)[0] || value
      } catch (e) {
        return value
      }
    }
    return value
  }
  return String(value)
}

export default messages
