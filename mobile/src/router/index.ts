 import { createRouter, createWebHistory } from 'vue-router'
 
 const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL),
   routes: [
     {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginPage.vue'),
      meta: { public: true },
    },
     {
       path: '/register',
       name: 'register',
       component: () => import('@/views/RegisterPage.vue'),
       meta: { public: true },
     },
    {
      path: '/',
       component: () => import('@/views/MainPage.vue'),
       meta: { requiresAuth: true },
       children: [
         { path: '', redirect: '/coin' },
        { path: 'coin', name: 'coin', component: () => import('@/views/CoinPage.vue') },
        { path: 'heyue', name: 'heyue', component: () => import('@/views/HeyuePage.vue') },
         { path: 'heyueorder', name: 'heyueorder', component: () => import('@/views/HeyueOrderPage.vue') },
        { path: 'profile', name: 'profile', component: () => import('@/views/ProfilePage.vue') },
      ],
    },
     {
      path: '/heyue-add',
      name: 'heyue-add',
      component: () => import('@/views/HeyueAdd.vue'),
      meta: { requiresAuth: true },
    },
     {
       path: '/config',
       name: 'config',
       component: () => import('@/views/ConfigPage.vue'),
       meta: { requiresAuth: true },
     },
     {
       path: '/binan',
       name: 'binan',
       component: () => import('@/views/BinanPage.vue'),
       meta: { requiresAuth: true },
     },
    { path: '/:pathMatch(.*)*', redirect: '/coin' },
   ],
 })
 
 router.beforeEach((to) => {
   const token = localStorage.getItem('token')
   if (to.meta.public) {
     if (to.name === 'login' && token) {
       return { name: 'coin' }
     }
     return true
   }
   if (to.meta.requiresAuth && !token) {
     return { name: 'login' }
   }
   return true
 })
 
 export default router
