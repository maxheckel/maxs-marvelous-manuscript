<template>
  <div class="recording-detail">
    <button @click="goBack" class="btn-back">← Back to Recordings</button>

    <div v-if="loading" class="loading">Loading recording...</div>

    <div v-else-if="error" class="error">
      {{ error }}
    </div>

    <div v-else-if="recording" class="detail-container">
      <div class="detail-header">
        <h2>{{ recording.filename }}</h2>
        <span :class="['status', recording.status]">{{ recording.status }}</span>
      </div>

      <div class="audio-player">
        <audio controls :src="audioUrl" style="width: 100%">
          Your browser does not support the audio element.
        </audio>
      </div>

      <div class="detail-info">
        <div class="info-section">
          <h3>Recording Information</h3>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">File ID:</span>
              <span>{{ recording.file_id }}</span>
            </div>
            <div class="info-item">
              <span class="label">Duration:</span>
              <span>{{ formatDuration(recording.duration_seconds) }}</span>
            </div>
            <div class="info-item">
              <span class="label">File Size:</span>
              <span>{{ formatSize(recording.file_size_bytes) }}</span>
            </div>
            <div class="info-item">
              <span class="label">Created:</span>
              <span>{{ formatDate(recording.created_at) }}</span>
            </div>
            <div class="info-item" v-if="recording.completed_at">
              <span class="label">Completed:</span>
              <span>{{ formatDate(recording.completed_at) }}</span>
            </div>
            <div class="info-item">
              <span class="label">Transcription Status:</span>
              <span>{{ recording.transcription_status }}</span>
            </div>
          </div>
        </div>

        <div class="info-section" v-if="recording.notes">
          <h3>Notes</h3>
          <p>{{ recording.notes }}</p>
        </div>

        <div class="info-section placeholder">
          <h3>Session Summary</h3>
          <p>AI-powered session summaries coming soon...</p>
          <p>This will include:</p>
          <ul>
            <li>Key events and story beats</li>
            <li>NPCs encountered</li>
            <li>Locations visited</li>
            <li>Items obtained</li>
            <li>Combat encounters</li>
            <li>Important decisions</li>
          </ul>
        </div>

        <div class="info-section placeholder">
          <h3>Transcription</h3>
          <p>Audio transcription with speaker diarization coming soon...</p>
        </div>
      </div>

      <div class="actions">
        <button @click="downloadRecording" class="btn-download">Download Audio</button>
        <button @click="deleteRecording" class="btn-delete">Delete Recording</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, computed, Ref, ComputedRef } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { api, Recording } from '../services/api'

export default {
  name: 'RecordingDetail',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const recording: Ref<Recording | null> = ref(null)
    const loading = ref(true)
    const error: Ref<string | null> = ref(null)

    const audioUrl: ComputedRef<string> = computed(() => {
      if (!recording.value) return ''
      return api.getAudioUrl(recording.value.id)
    })

    const loadRecording = async (): Promise<void> => {
      try {
        loading.value = true
        const id = parseInt(route.params.id as string)
        const response = await api.getRecording(id)
        recording.value = response.data
        error.value = null
      } catch (err: any) {
        error.value = 'Failed to load recording: ' + err.message
      } finally {
        loading.value = false
      }
    }

    const goBack = (): void => {
      router.push('/')
    }

    const downloadRecording = (): void => {
      window.open(audioUrl.value, '_blank')
    }

    const deleteRecording = async (): Promise<void> => {
      if (!confirm('Are you sure you want to delete this recording?')) {
        return
      }

      try {
        if (recording.value) {
          await api.deleteRecording(recording.value.id)
          router.push('/')
        }
      } catch (err: any) {
        alert('Failed to delete recording: ' + err.message)
      }
    }

    const formatDuration = (seconds: number): string => {
      if (!seconds) return '0:00'
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      const secs = seconds % 60

      if (hours > 0) {
        return `${hours}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
      }
      return `${minutes}:${String(secs).padStart(2, '0')}`
    }

    const formatSize = (bytes: number): string => {
      if (!bytes) return '0 B'
      const kb = bytes / 1024
      const mb = kb / 1024
      const gb = mb / 1024

      if (gb >= 1) return `${gb.toFixed(2)} GB`
      if (mb >= 1) return `${mb.toFixed(2)} MB`
      if (kb >= 1) return `${kb.toFixed(2)} KB`
      return `${bytes} B`
    }

    const formatDate = (dateString: string): string => {
      const date = new Date(dateString)
      return date.toLocaleString()
    }

    onMounted(loadRecording)

    return {
      recording,
      loading,
      error,
      audioUrl,
      goBack,
      downloadRecording,
      deleteRecording,
      formatDuration,
      formatSize,
      formatDate
    }
  }
}
</script>

<style scoped>
.recording-detail {
  width: 100%;
}

.btn-back {
  background: #95a5a6;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 2rem;
  transition: opacity 0.2s;
}

.btn-back:hover {
  opacity: 0.8;
}

.loading, .error {
  padding: 2rem;
  text-align: center;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.error {
  color: #e74c3c;
}

.detail-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #ecf0f1;
}

.detail-header h2 {
  font-size: 1.8rem;
  color: #2c3e50;
  margin: 0;
}

.status {
  padding: 0.5rem 1rem;
  border-radius: 12px;
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status.completed {
  background: #d4edda;
  color: #155724;
}

.status.recording {
  background: #fff3cd;
  color: #856404;
}

.status.failed {
  background: #f8d7da;
  color: #721c24;
}

.audio-player {
  margin: 2rem 0;
  padding: 1.5rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.detail-info {
  margin: 2rem 0;
}

.info-section {
  margin-bottom: 2rem;
}

.info-section h3 {
  font-size: 1.3rem;
  color: #2c3e50;
  margin-bottom: 1rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.label {
  font-weight: 600;
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 0.25rem;
}

.placeholder {
  padding: 1.5rem;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #3498db;
}

.placeholder p {
  margin: 0.5rem 0;
  color: #666;
}

.placeholder ul {
  margin: 1rem 0 0 1.5rem;
  color: #666;
}

.actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 2px solid #ecf0f1;
}

.actions button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.actions button:hover {
  opacity: 0.8;
}

.btn-download {
  background: #3498db;
  color: white;
  flex: 1;
}

.btn-delete {
  background: #e74c3c;
  color: white;
}
</style>
