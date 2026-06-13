export interface LoginRequest {
  username?: string
  email?: string
  password: string
}

export interface LoginResponse {
  status: number
  token?: string
  id?: number
  username?: string
  email?: string
  error?: string
}

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

/* ── ShellGroup CRUD ── */

export interface ShellGroupItem {
  id: number
  name: string
  mmurl: string
  addtime: number
}

export interface ShellGroupListResponse {
  status: number
  page: number
  pagesize: number
  total: number
  data: ShellGroupItem[]
  error?: string
}

export async function getShellGroups(page: number = 1): Promise<ShellGroupListResponse> {
  const res = await fetch(`/api/shell_group?page=${page}`, {
    headers: authHeaders(),
  })
  return res.json()
}

export async function createShellGroup(data: {
  name: string
  mmurl?: string
}): Promise<{ status: number; shell_group?: ShellGroupItem; error?: string }> {
  const res = await fetch('/api/shell_group', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function updateShellGroup(
  id: number,
  data: { name?: string; mmurl?: string }
): Promise<{ status: number; shell_group?: ShellGroupItem; error?: string }> {
  const res = await fetch(`/api/shell_group/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function deleteShellGroup(
  id: number
): Promise<{ status: number; message?: string; error?: string }> {
  const res = await fetch(`/api/shell_group/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return res.json()
}

/* ── Admin CRUD ── */

export interface AdminItem {
  id: number
  username: string
  email: string
  status: number
  register_time: string
}

export interface AdminListResponse {
  status: number
  page: number
  pagesize: number
  total: number
  data: AdminItem[]
  error?: string
}

function authHeaders(): Record<string, string> {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export async function getAdmins(page: number = 1): Promise<AdminListResponse> {
  const res = await fetch(`/api/admin?page=${page}`, {
    headers: authHeaders(),
  })
  return res.json()
}

export async function createAdmin(data: {
  username: string
  email: string
  password: string
  status?: number
}): Promise<{ status: number; admin?: AdminItem; error?: string }> {
  const res = await fetch('/api/admin', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function updateAdmin(
  id: number,
  data: { username?: string; email?: string; password?: string; status?: number }
): Promise<{ status: number; admin?: AdminItem; error?: string }> {
  const res = await fetch(`/api/admin/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function deleteAdmin(
  id: number
): Promise<{ status: number; message?: string; error?: string }> {
  const res = await fetch(`/api/admin/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return res.json()
}

/* ── AuthRule CRUD ── */

export interface AuthRuleItem {
  id: number
  pid: number
  name: string
  title: string
  icon: string
  type: number
  status: number
  condition: string
}

export interface AuthRuleListResponse {
  status: number
  page: number
  pagesize: number
  total: number
  data: AuthRuleItem[]
  error?: string
}

export async function getAuthRules(page: number = 1): Promise<AuthRuleListResponse> {
  const res = await fetch(`/api/auth_rule?page=${page}`, {
    headers: authHeaders(),
  })
  return res.json()
}

export async function createAuthRule(data: {
  pid?: number
  name: string
  title: string
  icon?: string
  type?: number
  status?: number
  condition?: string
}): Promise<{ status: number; auth_rule?: AuthRuleItem; error?: string }> {
  const res = await fetch('/api/auth_rule', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function updateAuthRule(
  id: number,
  data: {
    pid?: number
    name?: string
    title?: string
    icon?: string
    type?: number
    status?: number
    condition?: string
  }
): Promise<{ status: number; auth_rule?: AuthRuleItem; error?: string }> {
  const res = await fetch(`/api/auth_rule/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function deleteAuthRule(
  id: number
): Promise<{ status: number; message?: string; error?: string }> {
  const res = await fetch(`/api/auth_rule/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return res.json()
}

/* ── Shell CRUD ── */

export interface ShellItem {
  id: number
  host: string
  scheme: string
  group_id: number
  group?: {
    id: number
    name: string
    mmurl: string
    addtime: number
  }
  status: number
  num: number
  maxurl: string
  minurl: string
  dir: number
  lock: number
  remark: string
  addtime: number
  uptime: number
}

export interface ShellListResponse {
  status: number
  page: number
  pagesize: number
  total: number
  data: ShellItem[]
  error?: string
}

export async function getShells(page: number = 1): Promise<ShellListResponse> {
  const res = await fetch(`/api/shell?page=${page}`, {
    headers: authHeaders(),
  })
  return res.json()
}

export async function createShell(data: {
  host: string
  scheme?: string
  group_id?: number
  status?: number
  num?: number
  maxurl?: string
  minurl?: string
  dir?: number
  lock?: number
  remark?: string
}): Promise<{ status: number; shell?: ShellItem; error?: string }> {
  const res = await fetch('/api/shell', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function updateShell(
  id: number,
  data: {
    host?: string
    scheme?: string
    group_id?: number
    status?: number
    num?: number
    maxurl?: string
    minurl?: string
    dir?: number
    lock?: number
    remark?: string
  }
): Promise<{ status: number; shell?: ShellItem; error?: string }> {
  const res = await fetch(`/api/shell/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return res.json()
}

export async function deleteShell(
  id: number
): Promise<{ status: number; message?: string; error?: string }> {
  const res = await fetch(`/api/shell/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  return res.json()
}

/* ── ShellMax ── */

export interface ShellMaxItem {
  id: number
  shell_id: number
  url: string
  addtime: number
  status: number
}

export async function getShellMax(
  shellId: number
): Promise<{ status: number; shell_id: number; data: ShellMaxItem[]; error?: string }> {
  const res = await fetch(`/api/shell_max?shell_id=${shellId}`, {
    headers: authHeaders(),
  })
  return res.json()
}
 
 /* 鈹€鈹€ Coin CRUD 鈹€鈹€ */
 
 export interface CoinItem {
   id: number
   name: string
   symbol: string
   close: number
   priceprecision: number
   open: number
   low: number
   high: number
   status: number
   addtime: number
   updatetime: number
 }
 
 export interface CoinListResponse {
   status: number
   page: number
   pagesize: number
   total: number
   data: CoinItem[]
   error?: string
 }
 
export async function getCoins(page: number = 1, symbol?: string, status?: number): Promise<CoinListResponse> {
  let url = `/api/coin?page=${page}`
  if (symbol) {
    url += `&symbol=${encodeURIComponent(symbol)}`
  }
  if (status !== undefined) {
    url += `&status=${status}`
  }
  const res = await fetch(url, {
    headers: authHeaders(),
  })
  return res.json()
 }
 
 export async function createCoin(data: {
   name: string
   symbol: string
   close?: number
   priceprecision?: number
   open?: number
   low?: number
   high?: number
   status?: number
 }): Promise<{ status: number; coin?: CoinItem; error?: string }> {
   const res = await fetch('/api/coin', {
     method: 'POST',
     headers: { ...authHeaders(), 'Content-Type': 'application/json' },
     body: JSON.stringify(data),
   })
   return res.json()
 }
 
 export async function updateCoin(
   id: number,
   data: {
     name?: string
     symbol?: string
     close?: number
     priceprecision?: number
     open?: number
     low?: number
     high?: number
     status?: number
   }
 ): Promise<{ status: number; coin?: CoinItem; error?: string }> {
   const res = await fetch(`/api/coin/${id}`, {
     method: 'PUT',
     headers: { ...authHeaders(), 'Content-Type': 'application/json' },
     body: JSON.stringify(data),
   })
   return res.json()
 }
 
 export async function deleteCoin(
   id: number
 ): Promise<{ status: number; message?: string; error?: string }> {
   const res = await fetch(`/api/coin/${id}`, {
     method: 'DELETE',
     headers: authHeaders(),
   })
   return res.json()
 }
