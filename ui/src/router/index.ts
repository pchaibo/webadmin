import { createRouter, createWebHistory } from 'vue-router'

const isAuthed = () => Boolean(localStorage.getItem('admin_token'))

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
    path: '/',
    component: () => import('@/views/HomeView.vue'),
    meta: { requiresAuth: true },

    children: [
      {
        path: '',
        redirect: '/index',
      },
      {
        path: 'index',
        name: 'index',
        component: () => import('@/views/Index.vue'),
      },
      {
        path: 'admin',
        name: 'admin',
        component: () => import('@/views/admin.vue'),
      },
      {
        path: 'authrule',
        name: 'auth-rule',
        component: () => import('@/views/authrule.vue'),
      },
      {
        path: 'shellgroup',
        name: 'shell-group',
        component: () => import('@/views/ShellGroup.vue'),
      },
      {
        path: 'shell',
        name: 'shell',
        component: () => import('@/views/Shell.vue'),
      },
      {
      path: 'shell/add',
      name: 'shell-add',
      component: () => import('@/views/ShellAdd.vue'),
    },
     {
       path: 'coin',
       name: 'coin',
     component: () => import('@/views/Coin.vue'),
   },
    {
      path: 'heyue',
      name: 'heyue',
      component: () => import('@/views/Heyue.vue'),
    },
    //Heyueorder
    {
      path: 'heyueorder',
      name: 'heyueorder',
      component: () => import('@/views/Heyueorder.vue'),
    },
    {
      path: 'user',
       name: 'user',
       component: () => import('@/views/user.vue'),
     },
   
  ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach((to) => {
  if (to.meta.public) {
    if (to.name === 'login' && isAuthed()) {
      return { name: 'home' }
    }
    return true
  }

  if (to.meta.requiresAuth && !isAuthed()) {
    return { name: 'login' }
  }

  return true
})

export default router
