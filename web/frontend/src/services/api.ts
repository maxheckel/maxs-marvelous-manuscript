import axios, { AxiosResponse } from 'axios'

const API_BASE = '/api'

export interface Recording {
  id: number
  file_id: string
  filename: string
  file_path: string
  duration_seconds: number
  file_size_bytes: number
  status: 'recording' | 'completed' | 'failed'
  created_at: string
  completed_at?: string
  transcription_status: 'pending' | 'processing' | 'completed' | 'failed'
  notes?: string
}

export const api = {
  // Recordings
  getRecordings(): Promise<AxiosResponse<Recording[]>> {
    return axios.get<Recording[]>(`${API_BASE}/recordings`)
  },

  getRecording(id: number): Promise<AxiosResponse<Recording>> {
    return axios.get<Recording>(`${API_BASE}/recordings/${id}`)
  },

  deleteRecording(id: number): Promise<AxiosResponse<{ message: string }>> {
    return axios.delete(`${API_BASE}/recordings/${id}`)
  },

  getAudioUrl(id: number): string {
    return `${API_BASE}/recordings/${id}/audio`
  },

  // Health check
  healthCheck(): Promise<AxiosResponse<{ status: string }>> {
    return axios.get(`${API_BASE}/health`)
  }
}
