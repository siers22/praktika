const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

function getToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('access_token')
}

function getRefreshToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('refresh_token')
}

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = getRefreshToken()
  if (!refreshToken) return null

  try {
    const res = await fetch(`${BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
    if (!res.ok) return null
    const data = await res.json()
    const { access_token, refresh_token } = data.data
    localStorage.setItem('access_token', access_token)
    localStorage.setItem('refresh_token', refresh_token)
    return access_token
  } catch {
    return null
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  let token = getToken()

  const makeRequest = async (tkn: string | null) => {
    const headers: Record<string, string> = {
      ...(options.body && !(options.body instanceof FormData)
        ? { 'Content-Type': 'application/json' }
        : {}),
      ...(tkn ? { Authorization: `Bearer ${tkn}` } : {}),
      ...(options.headers as Record<string, string>),
    }
    return fetch(`${BASE_URL}${path}`, { ...options, headers })
  }

  let res = await makeRequest(token)

  if (res.status === 401) {
    const newToken = await refreshAccessToken()
    if (newToken) {
      res = await makeRequest(newToken)
    } else {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      window.location.href = '/login'
      throw new Error('Session expired')
    }
  }

  if (res.status === 204) return null as T

  const data = await res.json()
  if (!res.ok) {
    throw new Error(data.error?.message || 'Ошибка запроса')
  }
  return data
}

// Auth
export const auth = {
  login: (username: string, password: string) =>
    request<{ data: { access_token: string; refresh_token: string } }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),
  logout: (refreshToken: string) =>
    request('/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),
  changePassword: (oldPassword: string, newPassword: string) =>
    request('/auth/password', {
      method: 'PUT',
      body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
    }),
}

// Users
export const users = {
  list: (page = 1, perPage = 20) =>
    request<any>(`/users?page=${page}&per_page=${perPage}`),
  create: (data: any) =>
    request<any>('/users', { method: 'POST', body: JSON.stringify(data) }),
  getById: (id: string) => request<any>(`/users/${id}`),
  update: (id: string, data: any) =>
    request<any>(`/users/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  updateStatus: (id: string, isActive: boolean) =>
    request<any>(`/users/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ is_active: isActive }),
    }),
}

// Equipment
export const equipment = {
  list: (params: Record<string, string | number> = {}) => {
    const q = new URLSearchParams(
      Object.entries(params)
        .filter(([, v]) => v !== '' && v !== undefined)
        .map(([k, v]) => [k, String(v)])
    ).toString()
    return request<any>(`/equipment${q ? '?' + q : ''}`)
  },
  create: (data: any) =>
    request<any>('/equipment', { method: 'POST', body: JSON.stringify(data) }),
  getById: (id: string) => request<any>(`/equipment/${id}`),
  update: (id: string, data: any) =>
    request<any>(`/equipment/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  archive: (id: string) => request<any>(`/equipment/${id}`, { method: 'DELETE' }),
  uploadPhoto: (id: string, file: File) => {
    const form = new FormData()
    form.append('photo', file)
    return request<any>(`/equipment/${id}/photos`, { method: 'POST', body: form })
  },
  deletePhoto: (id: string, photoId: string) =>
    request<any>(`/equipment/${id}/photos/${photoId}`, { method: 'DELETE' }),
  exportCSV: () =>
    fetch(`${BASE_URL}/equipment/export/csv`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    }),
  movements: (id: string, page = 1) =>
    request<any>(`/equipment/${id}/movements?page=${page}`),
}

// Categories
export const categories = {
  list: () => request<any>('/categories'),
  create: (data: any) =>
    request<any>('/categories', { method: 'POST', body: JSON.stringify(data) }),
  update: (id: string, data: any) =>
    request<any>(`/categories/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (id: string) => request<any>(`/categories/${id}`, { method: 'DELETE' }),
}

// Departments
export const departments = {
  list: () => request<any>('/departments'),
  create: (data: any) =>
    request<any>('/departments', { method: 'POST', body: JSON.stringify(data) }),
  update: (id: string, data: any) =>
    request<any>(`/departments/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (id: string) => request<any>(`/departments/${id}`, { method: 'DELETE' }),
}

// Inventory
export const inventory = {
  listSessions: (page = 1) => request<any>(`/inventories?page=${page}`),
  createSession: (departmentId: string) =>
    request<any>('/inventories', {
      method: 'POST',
      body: JSON.stringify({ department_id: departmentId }),
    }),
  getSession: (id: string) => request<any>(`/inventories/${id}`),
  checkItem: (sessionId: string, data: any) =>
    request<any>(`/inventories/${sessionId}/items`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  updateItem: (sessionId: string, itemId: string, data: any) =>
    request<any>(`/inventories/${sessionId}/items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  complete: (sessionId: string) =>
    request<any>(`/inventories/${sessionId}/complete`, { method: 'POST' }),
  exportCSV: (sessionId: string) =>
    fetch(`${BASE_URL}/inventories/${sessionId}/export/csv`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    }),
}

// Movements
export const movements = {
  list: (params: Record<string, string | number> = {}) => {
    const q = new URLSearchParams(
      Object.entries(params)
        .filter(([, v]) => v !== '' && v !== undefined)
        .map(([k, v]) => [k, String(v)])
    ).toString()
    return request<any>(`/movements${q ? '?' + q : ''}`)
  },
  create: (data: any) =>
    request<any>('/movements', { method: 'POST', body: JSON.stringify(data) }),
}

// Reports
export const reports = {
  summary: () => request<any>('/reports/summary'),
  byDepartment: () => request<any>('/reports/by-department'),
  dashboard: () => request<any>('/reports/dashboard'),
}

// Audit
export const audit = {
  list: (params: Record<string, string | number> = {}) => {
    const q = new URLSearchParams(
      Object.entries(params)
        .filter(([, v]) => v !== '' && v !== undefined)
        .map(([k, v]) => [k, String(v)])
    ).toString()
    return request<any>(`/audit-logs${q ? '?' + q : ''}`)
  },
}
