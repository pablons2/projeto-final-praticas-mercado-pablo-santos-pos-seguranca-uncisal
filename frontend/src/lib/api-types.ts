/** API contract types — mirror the backend DTOs (docs/plano.md §6). */

export interface User {
  id: string
  name: string
  email: string
  created_at: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  name: string
  email: string
  password: string
}

export interface Incident {
  id: string
  titulo: string
  descricao: string
  categoria: string
  severidade: string
  status: string
  owner_id: string
  created_at: string
  updated_at: string
}

export interface IncidentCreatePayload {
  titulo: string
  descricao: string
  categoria: string
  severidade: string
  status: string
}

export type IncidentUpdatePayload = IncidentCreatePayload
