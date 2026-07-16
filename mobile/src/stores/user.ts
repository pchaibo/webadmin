 import { defineStore } from 'pinia'
 
 export const useUserStore = defineStore('user', {
   state: () => ({
     token: localStorage.getItem('token') ?? '',
     username: localStorage.getItem('username') ?? '',
     email: localStorage.getItem('email') ?? '',
   }),
   getters: {
     isLoggedIn: (state) => !!state.token,
   },
   actions: {
     login(token: string, username: string, email: string) {
       this.token = token
       this.username = username
       this.email = email
       localStorage.setItem('token', token)
       localStorage.setItem('username', username)
       localStorage.setItem('email', email)
     },
     logout() {
       this.token = ''
       this.username = ''
       this.email = ''
       localStorage.removeItem('token')
       localStorage.removeItem('username')
       localStorage.removeItem('email')
     },
   },
 })
