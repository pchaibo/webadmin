 import type { ApiResponse, UserInfo, LoginParams } from './types'
 
 const BASE = ''
 
 function getToken(): string {
   return localStorage.getItem('token') ?? ''
 }
 
 export function setToken(token: string) {
   localStorage.setItem('token', token)
 }
 
 export function clearToken() {
   localStorage.removeItem('token')
 }
 
 function authHeaders(): Record<string, string> {
   const t = getToken()
   return t ? { Authorization: 'Bearer ' + t } : {}
 }
 
 async function request<T>(method: string, url: string, body?: unknown): Promise<T> {
   const res = await fetch(BASE + url, {
     method,
     headers: { 'Content-Type': 'application/json', ...authHeaders() },
     body: body !== undefined ? JSON.stringify(body) : undefined,
   })
   return res.json() as Promise<T>
 }
 
 export function get<T>(url: string) {
   return request<T>('GET', url)
 }
 
 export function post<T>(url: string, body?: unknown) {
   return request<T>('POST', url, body)
 }
 
 export type { ApiResponse, UserInfo, LoginParams }
 
export async function login(params: LoginParams) {
  const res = await post<{ status: number; token?: string; id?: number; username?: string; email?: string }>('/user/login', params)
  if (res.status === 1 && res.token) {
    setToken(res.token)
  }
  return res
}
 
 export async function register(params: { username: string; password: string; email: string; mobile: string }) {
   return post<{ status: number; error?: string }>('/user/register', params)
 }
 
 export async function fetchCoinList() {
   return get<ApiResponse>('/user/coin')
 }
 
 export async function fetchHeyueList(page?: number, symbol?: string, username?: string, status?: number) {
   let url = '/user/heyue'
   const params: string[] = []
   if (page) params.push('page=' + page)
   if (symbol) params.push('symbol=' + encodeURIComponent(symbol))
   if (username) params.push('username=' + encodeURIComponent(username))
   if (status !== undefined) params.push('status=' + status)
   if (params.length) url += '?' + params.join('&')
   return get<{ status: number; data?: unknown[]; total?: number; page?: number; pagesize?: number; error?: string }>(url)
 }

export async function createHeyue(data: Record<string, unknown>) {
   return post<ApiResponse>('/user/heyue', data)
 }
 
 export async function updateHeyue(id: number, data: Record<string, unknown>) {
   return fetch(`${BASE}/user/heyue/${id}`, {
     method: 'PUT',
     headers: { 'Content-Type': 'application/json', ...authHeaders() },
     body: JSON.stringify(data),
   }).then(r => r.json()) as Promise<ApiResponse>
 }
 
export async function deleteHeyue(id: number) {
  const res = await fetch(`${BASE}/user/heyue/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return res.json() as Promise<ApiResponse>
}
 
 export async function fetchHeyueById(id: number) {
   return get<{ status: number; heyue?: Record<string, unknown>; error?: string }>('/user/heyue/' + id)
 }

export async function fetchUserInfo() {
  return get<ApiResponse<UserInfo>>('/user/info')
}
 
 export async function updatePassword(password: string) {
   return fetch(`${BASE}/user/password`, {
     method: 'PUT',
     headers: { 'Content-Type': 'application/json', ...authHeaders() },
     body: JSON.stringify({ password }),
   }).then(r => r.json()) as Promise<{ status: number; message?: string; error?: string }>
 }
 
 export async function updateBinance(bnaccess: string, bnasecret: string, margin: number) {
  return fetch(`${BASE}/user/binance`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
     body: JSON.stringify({ bnaccess, bnasecret, margin }),
  }).then(r => r.json()) as Promise<{ status: number; message?: string; error?: string }>
 }
 
 export async function fetchHeyueOrderList(page?: number, symbol?: string, ordertype?: number, side?: number) {
  let url = '/user/heyuesorder'
   const params: string[] = []
   if (page) params.push('page=' + page)
   if (symbol) params.push('symbol=' + encodeURIComponent(symbol))
   if (ordertype !== undefined) params.push('ordertype=' + ordertype)
   if (side !== undefined) params.push('side=' + side)
   if (params.length) url += '?' + params.join('&')
  return get<{
     status: number
     data?: unknown[]
     total?: number
     page?: number
     pagesize?: number
     total_usdt?: number
     total_usdt_long?: number
     total_usdt_short?: number
     error?: string
   }>(url)
 }
