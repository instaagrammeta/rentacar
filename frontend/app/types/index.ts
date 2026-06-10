// Domain types mirroring the JSON contract returned by the Go backend.

export interface User {
  id: number
  username: string
  email: string | null
  full_name: string
  role: string
  role_label: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Client {
  id: number
  client_code: string
  first_name: string
  last_name: string
  full_name: string
  phone: string
  email: string | null
  date_of_birth: string | null
  passport_number: string | null
  driver_license_number: string | null
  driver_license_issue_date: string | null
  driver_experience_years: number
  passport_scan: string | null
  passport_front: string | null
  passport_back: string | null
  driver_license_scan: string | null
  driver_photo: string | null
  qr_code_path: string | null
  is_vip: boolean
  notes: string | null
  status: string
  status_label: string
  created_at: string
  updated_at: string
}

export interface Car {
  id: number
  brand: string
  model: string
  year: number
  color: string | null
  vin: string | null
  plate_number: string
  mileage: number
  daily_price: number
  weekly_price: number
  monthly_price: number
  deposit_amount: number
  status: string
  status_label: string
  display_name: string
  photos: string[]
  created_at: string
  updated_at: string
}

export interface Reservation {
  id: number
  client_id: number
  car_id: number
  start_date: string
  end_date: string
  deposit: number
  notes: string | null
  status: string
  status_label: string
  client_name?: string
  car_name?: string
  created_at: string
  updated_at: string
}

export interface Rental {
  id: number
  contract_number: string
  client_id: number
  car_id: number
  reservation_id: number | null
  employee_id: number | null
  rental_start: string
  rental_end: string
  deposit: number
  daily_price: number
  total_price: number
  start_mileage: number | null
  pdf_path: string | null
  notes: string | null
  status: string
  status_label: string
  client_name?: string
  car_name?: string
  employee_name?: string
  has_return: boolean
  vehicle_return?: VehicleReturn
  created_at: string
  updated_at: string
}

export interface VehicleReturn {
  id: number
  rental_id: number
  return_date: string
  mileage: number
  fuel_level: number
  damages: string | null
  extra_days: number
  late_fee: number
  damage_cost: number
  penalties: number
  final_payment: number
}

export interface Payment {
  id: number
  receipt_number: string
  client_id: number
  rental_id: number | null
  cashier_id: number | null
  amount: number
  method: string
  method_label: string
  payment_type: string
  payment_type_label: string
  paid_at: string
  receipt_path: string | null
  notes: string | null
  client_name?: string
}

export interface BlacklistEntry {
  id: number
  client_id: number
  reason: string
  reason_label: string
  comment: string | null
  created_by: number | null
  client_name?: string
  created_at: string
}

export interface Accident {
  id: number
  car_id: number
  client_id: number | null
  rental_id: number | null
  accident_date: string
  description: string | null
  repair_cost: number
  photos: string[]
  car_name?: string
  client_name?: string
  created_at: string
}

export interface DashboardSummary {
  cars_available: number
  cars_rented: number
  cars_reserved: number
  cars_maintenance: number
  today_revenue: number
  monthly_revenue: number
  active_rentals: number
  upcoming_returns: number
}

export interface Paginated<T> {
  items: T[]
  total: number
  page: number
  per_page: number
}

export interface Option {
  value: string
  label: string
}
