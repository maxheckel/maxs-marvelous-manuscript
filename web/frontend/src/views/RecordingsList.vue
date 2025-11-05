<template>
  <div class="recordings-list">
    <h2>D&D Session Recordings</h2>

    <div v-if="loading" class="loading">Loading recordings...</div>

    <div v-else-if="error" class="error">
      {{ error }}
    </div>

    <div v-else-if="recordings.length === 0" class="empty">
      No recordings yet. Start a recording using the recorder app!
    </div>

    <div v-else class="recordings-grid">
      <div
        v-for="recording in recordings"
        :key="recording.id"
        class="recording-card"
        @click="viewRecording(recording.id)"
      >
        <div class="recording-header">
          <h3>{{ recording.filename }}</h3>
          <span :class="['status', recording.status]">{{ recording.status }}</span>
        </div>

        <div class="recording-info">
          <div class="info-row">
            <span class="label">Duration:</span>
            <span>{{ formatDuration(recording.duration_seconds) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Size:</span>
            <span>{{ formatSize(recording.file_size_bytes) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Created:</span>
            <span>{{ formatDate(recording.created_at) }}</span>
          </div>
        </div>

        <div class="recording-actions" @click.stop>
          <button @click="playRecording(recording.id)" class="btn-play">Play</button>
          <button @click="deleteRecording(recording.id)" class="btn-delete">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, Ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, Recording } from '../services/api'

export default {
  name: 'RecordingsList',
  setup() {
    const router = useRouter()
    const recordings: Ref<Recording[]> = ref([])
    const loading = ref(true)
    const error: Ref<string | null> = ref(null)

    const loadRecordings = async () => {
      try {
        loading.value = true
        const response = await api.getRecordings()
        recordings.value = response.data || []
        error.value = null
      } catch (err) {
        error.value = 'Failed to load recordings: ' + err.message
      } finally {
        loading.value = false
      }
    }

    const viewRecording = (id: number): void => {
      router.push(`/recordings/${id}`)
    }

    const playRecording = (id: number): void => {
      const url = api.getAudioUrl(id)
      window.open(url, '_blank')
    }

    const deleteRecording = async (id: number): Promise<void> => {
      if (!confirm('Are you sure you want to delete this recording?')) {
        return
      }

      try {
        await api.deleteRecording(id)
        await loadRecordings()
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

    onMounted(loadRecordings)

    return {
      recordings,
      loading,
      error,
      viewRecording,
      playRecording,
      deleteRecording,
      formatDuration,
      formatSize,
      formatDate
    }
  }
}
</script>

<style scoped>
.recordings-list {
  width: 100%;
}

h2 {
  font-size: 2rem;
  margin-bottom: 2rem;
  color: #2c3e50;
}

.loading, .error, .empty {
  padding: 2rem;
  text-align: center;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.error {
  color: #e74c3c;
}

.recordings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.recording-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.recording-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.recording-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 1rem;
}

.recording-header h3 {
  font-size: 1.1rem;
  color: #2c3e50;
  margin: 0;
  word-break: break-word;
  flex: 1;
}

.status {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  white-space: nowrap;
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

.recording-info {
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-row:last-child {
  border-bottom: none;
}

.label {
  font-weight: 600;
  color: #666;
}

.recording-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: opacity 0.2s;
  flex: 1;
}

button:hover {
  opacity: 0.8;
}

.btn-play {
  background: #3498db;
  color: white;
}

.btn-delete {
  background: #e74c3c;
  color: white;
}
</style>
