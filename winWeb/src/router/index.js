import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: 'login', layout: 'blank' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: { title: 'register', layout: 'blank' }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/home/index.vue'),
        meta: { title: 'home', requiresAuth: false }
      },
      {
        path: 'learning/diary',
        name: 'Diary',
        component: () => import('@/views/learning/diary/index.vue'),
        meta: { title: 'diary', requiresAuth: false }
      },
      {
        path: 'learning/diary-detail/:id',
        name: 'DiaryDetail',
        component: () => import('@/views/learning/diary-detail/index.vue'),
        meta: { title: 'diaryDetail', requiresAuth: false }
      },
      {
        path: 'learning/typing',
        name: 'Typing',
        component: () => import('@/views/learning/typing/index.vue'),
        meta: { title: 'typing', requiresAuth: false }
      },
      {
        path: 'learning/profile',
        name: 'Profile',
        component: () => import('@/views/learning/profile/index.vue'),
        meta: { title: 'profile', requiresAuth: true }
      },
      {
        path: 'learning/player/:episodeId',
        name: 'Player',
        component: () => import('@/views/learning/player/index.vue'),
        meta: { title: 'player', requiresAuth: false }
      },
      {
        path: 'learning/video-detail/:seriesId',
        name: 'VideoDetail',
        component: () => import('@/views/learning/video-detail/index.vue'),
        meta: { title: 'videoDetail', requiresAuth: false }
      },
      {
        path: 'learning/collections',
        name: 'Collections',
        component: () => import('@/views/learning/collections/index.vue'),
        meta: { title: 'collections', requiresAuth: true }
      },
      {
        path: 'learning/watch-history',
        name: 'WatchHistory',
        component: () => import('@/views/learning/watch-history/index.vue'),
        meta: { title: 'watchHistory', requiresAuth: true }
      },
      {
        path: 'learning/checkin-record',
        name: 'CheckinRecord',
        component: () => import('@/views/learning/checkin-record/index.vue'),
        meta: { title: 'checkinRecord', requiresAuth: true }
      },
      {
        path: 'learning/error-log',
        name: 'ErrorLog',
        component: () => import('@/views/learning/error-log/index.vue'),
        meta: { title: 'errorLog', requiresAuth: true }
      },
      {
        path: 'learning/point-history',
        name: 'PointHistory',
        component: () => import('@/views/learning/point-history/index.vue'),
        meta: { title: 'pointHistory', requiresAuth: true }
      },
      {
        path: 'learning/free-time-history',
        name: 'FreeTimeHistory',
        component: () => import('@/views/learning/free-time-history/index.vue'),
        meta: { title: 'freeTimeHistory', requiresAuth: true }
      },
      {
        path: 'learning/tag-filter',
        name: 'TagFilter',
        component: () => import('@/views/learning/tag-filter/index.vue'),
        meta: { title: 'tagFilter', requiresAuth: false }
      },
      {
        path: 'kefu',
        name: 'Kefu',
        component: () => import('@/views/kefu/index.vue'),
        meta: { title: 'kefu', requiresAuth: false }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/home'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('x-token')
  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
