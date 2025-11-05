import axios from 'axios'

const API_BASE = '/api'

export const api = {
  // Recordings
  getRecordings() {
    return axios.get(`${API_BASE}/recordings`)
  },

  getRecording(id) {
    return axios.get(`${API_BASE}/recordings/${id}`)
  },

  deleteRecording(id) {
    return axios.delete(`${API_BASE}/recordings/${id}`)
  },

  getAudioUrl(id) {
    return `${API_BASE}/recordings/${id}/audio`
  },

  // Health check
  healthCheck() {
    return axios.get(`${API_BASE}/health`)
  }
}
