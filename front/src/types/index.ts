export type Role = 'admin' | 'inventory' | 'viewer'

export type EquipmentStatus = 'in_use' | 'in_storage' | 'in_repair' | 'written_off' | 'reserved'

export type InventoryStatus = 'in_progress' | 'completed'

export type ActualStatus = 'found' | 'not_found' | 'damaged'

export interface User {
  id: string
  username: string
  full_name: string
  email: string
  role: Role
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
}

export interface AuthUser {
  id: string
  role: Role
}

export interface Category {
  id: string
  name: string
  description: string | null
  created_at: string
}

export interface Department {
  id: string
  name: string
  location: string | null
  created_at: string
}

export interface EquipmentPhoto {
  id: string
  equipment_id: string
  file_path: string
  uploaded_at: string
}

export interface Equipment {
  id: string
  inventory_number: string
  name: string
  description: string | null
  category_id: string
  category_name: string | null
  serial_number: string | null
  model: string | null
  manufacturer: string | null
  purchase_date: string | null
  purchase_price: number | null
  warranty_expiry: string | null
  status: EquipmentStatus
  department_id: string
  department_name: string | null
  responsible_person_id: string | null
  responsible_person_name: string | null
  notes: string | null
  is_archived: boolean
  photos: EquipmentPhoto[]
  created_at: string
  updated_at: string
}

export interface Movement {
  id: string
  equipment_id: string
  equipment_name: string | null
  inventory_number: string | null
  from_department_id: string
  from_department_name: string | null
  to_department_id: string
  to_department_name: string | null
  moved_by: string
  moved_by_name: string | null
  moved_at: string
  reason: string | null
}

export interface InventorySession {
  id: string
  department_id: string
  department_name: string | null
  status: InventoryStatus
  created_by: string
  created_by_name: string | null
  started_at: string
  finished_at: string | null
}

export interface InventoryItem {
  id: string
  session_id: string
  equipment_id: string
  equipment_name: string | null
  inventory_number: string | null
  expected_status: string
  actual_status: ActualStatus
  comment: string | null
  checked_at: string
}

export interface AuditLog {
  id: string
  user_id: string
  username: string | null
  action: string
  entity_type: string
  entity_id: string | null
  details: Record<string, unknown> | null
  created_at: string
}

export interface Pagination {
  page: number
  per_page: number
  total: number
  total_pages: number
}

export interface ApiResponse<T> {
  data: T
  meta?: Pagination
}

export interface DashboardData {
  total_equipment: number
  by_status: Record<string, number>
  by_category: { category: string; count: number }[]
  recent_movements: Movement[]
  warranty_expiring_soon: Equipment[]
}

export const STATUS_LABELS: Record<EquipmentStatus, string> = {
  in_use: 'В эксплуатации',
  in_storage: 'На складе',
  in_repair: 'В ремонте',
  written_off: 'Списано',
  reserved: 'Зарезервировано',
}

export const STATUS_COLORS: Record<EquipmentStatus, string> = {
  in_use: 'success',
  in_storage: 'info',
  in_repair: 'warning',
  written_off: 'danger',
  reserved: 'primary',
}

export const ROLE_LABELS: Record<Role, string> = {
  admin: 'Администратор',
  inventory: 'Инвентаризатор',
  viewer: 'Наблюдатель',
}

export const ACTUAL_STATUS_LABELS: Record<ActualStatus, string> = {
  found: 'Найдено',
  not_found: 'Не найдено',
  damaged: 'Повреждено',
}
