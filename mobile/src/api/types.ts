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

export interface HeyueItem {
  id: number
  userid: number
  username: string
  symbol: string
  side: number
  num: number
  is_num: number
  status: number
  sellprice: number
  oneprice: number
  repeatprice: number
  rangetype: number
  rangeprice: number
  rangepercent: number
  rangeclosingpct: number
  rangeclosing: number
  closingprice: number
  risk: number
  risktime: number
  newprice: number
  newtime: number
  topprice: number
  reductionratio: number
  addtime: number
  updatetime: number
}

 export interface LoginParams {
   username?: string
   email?: string
   password: string
 }
