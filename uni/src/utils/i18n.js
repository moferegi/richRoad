/**
 * 多语言国际化工具
 * 支持：中文(zh)、英语(en)、蒙古国语(mn)
 */

const messages = {
  zh: {
    // tabBar/index.vue
    searchPlaceholder: '搜索剧集、演员、导演...',
    searchEmpty: '请输入搜索关键词',
    searchBtn: '搜索',
    searchResultCount: '找到 {{n}}+ 个结果',
    searchNoResult: '未找到相关剧集',
    searchNoResultTip: '换个关键词试试',
    searchInitTitle: '搜索您喜欢的内容',
    searchInitTip: '输入剧集名、演员或导演名称',
    searchLoadAll: '已加载全部',
    reachedBottom: '已经到底啦',
    loading: '加载中...',
    noGoods: '暂无剧集',
    all: '全部',

    // my/index.vue
    defaultUser: '默认用户',
    nicknamePlaceholder: '请填写昵称',
    loginPrompt: '请登录后查看',
    coupon: '优惠券',
    points: '积分',
    browseHistory: '浏览历史',
    playHistory: '播放历史',
    playerContinueFrom: '从第{{ep}}集 {{time}} 继续',
    myCollection: '我喜欢的',
    onlineService: '在线客服',
    logout: '退出登录',
    goodsInfoIncomplete: '剧集信息不完整',
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
    playerCollect: '喜欢',
    playerCollected: '已喜欢',
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
    tabCollect: '我喜欢的',
    tabMy: '我的',

    // category-page
    categoryDetail: '分类详情',
    noProductsInCategory: '该分类暂无剧集',
    loadFailed: '加载失败，请重试',
    allProducts: '全部剧集',
    noOrderData: '暂无订单数据',

    // collect
    noCollectData: '还没有喜欢的剧集',
    selfOperated: '自营',
    qualityAssured: '放心购',
    freeShipping: 'Plus免邮',
    collected: '已喜欢',
    uncollected: '已取消喜欢',
    loginFirst: '请登录后进行操作',
    goLogin: '去登录',
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
    searchPlaceholder: 'Search series, actors, directors...',
    searchEmpty: 'Please enter a keyword',
    searchBtn: 'Search',
    searchResultCount: '{{n}}+ results found',
    searchNoResult: 'No results found',
    searchNoResultTip: 'Try a different keyword',
    searchInitTitle: 'Search what you love',
    searchInitTip: 'Enter series name, actor or director',
    searchLoadAll: 'All loaded',
    reachedBottom: 'No more items',
    loading: 'Loading...',
    noGoods: 'No series yet',
    all: 'All',

    defaultUser: 'User',
    nicknamePlaceholder: 'Enter nickname',
    loginPrompt: 'Please login to view',
    coupon: 'Coupons',
    points: 'Points',
    browseHistory: 'History',
    playHistory: 'Watch History',
    playerContinueFrom: 'Continue EP{{ep}} at {{time}}',
    myCollection: 'My Likes',
    onlineService: 'Customer Service',
    logout: 'Log Out',
    goodsInfoIncomplete: 'Series info incomplete',
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
    playerCollect: 'Like',
    playerCollected: 'Liked',
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
    tabCollect: 'My Likes',
    tabMy: 'Me',

    categoryDetail: 'Category',
    noProductsInCategory: 'No series in this category',
    loadFailed: 'Load failed, please retry',
    allProducts: 'All Series',
    noOrderData: 'No order data',

    // collect
    noCollectData: 'No liked series yet',
    selfOperated: 'Official',
    qualityAssured: 'Verified',
    freeShipping: 'Free Ship',
    collected: 'Liked',
    uncollected: 'Unliked',
    loginFirst: 'Please login first',
    goLogin: 'Sign In',
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
    searchPlaceholder: 'Кино, жүжигчин хайх...',
    searchEmpty: 'Түлхүүр үг оруулна уу',
    searchBtn: 'Хайх',
    searchResultCount: '{{n}}+ үр дүн олдлоо',
    searchNoResult: 'Үр дүн олдсонгүй',
    searchNoResultTip: 'Өөр түлхүүр үг ашиглана уу',
    searchInitTitle: 'Дуртай зүйлсээ хайна уу',
    searchInitTip: 'Цуврал нэр, жүжигчин оруулна уу',
    searchLoadAll: 'Бүгдийг ачааллалаа',
    reachedBottom: 'Төгсгөлд хүрлээ',
    loading: 'Ачааллаж байна...',
    noGoods: 'Цуврал байхгүй',
    all: 'Бүгд',

    defaultUser: 'Хэрэглэгч',
    nicknamePlaceholder: 'Нэрээ оруулна уу',
    loginPrompt: 'Нэвтэрнэ үү',
    coupon: 'Купон',
    points: 'Оноо',
    browseHistory: 'Түүх',
    playHistory: 'Үзсэн түүх',
    playerContinueFrom: '{{ep}}-р анги {{time}}-с үргэлжлүүлнэ',
    myCollection: 'Миний дуртай',
    onlineService: 'Онлайн зөвлөгөө',
    logout: 'Гарах',
    goodsInfoIncomplete: 'Цувралын мэдээлэл дутуу',
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
    playerCollect: 'Дуртай',
    playerCollected: 'Дурласан',
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
    tabCollect: 'Миний дуртай',
    tabMy: 'Миний',

    categoryDetail: 'Ангилал',
    noProductsInCategory: 'Энэ ангилалд цуврал байхгүй',
    loadFailed: 'Ачаалалт амжилтгүй',
    allProducts: 'Бүх цуврал',
    noOrderData: 'Захиалга байхгүй',

    // collect
    noCollectData: 'Дуртай цуврал байхгүй',
    selfOperated: 'Албан',
    qualityAssured: 'Баталгаат',
    freeShipping: 'Үнэгүй',
    collected: 'Дурласан',
    uncollected: 'Дургүй болсон',
    loginFirst: 'Нэвтэрнэ үү',
    goLogin: 'Нэвтрэх',
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
