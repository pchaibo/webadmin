 export interface ApiResponse<T = unknown> {
   status: number
   data?: T
   [key: string]: unknown
 }
 
export interface UserInfo {
  id: number
  username: string
  email: string
  usdt?: number
   margin?: number
  bnaccess?: string
   bnasecret?: string
  token?: string
 }
 
 export interface LoginParams {
   username?: string
   email?: string
   password: string
 }
