import { onMounted, shallowRef } from 'vue'

import { getCpuCheck, getRamCheck } from '@/api'
import type { HealthCheckRespData } from '@/types/api'
import { formatErrorMessage } from '@/utils/errors'

export function useMonitor() {
  const cpu = shallowRef<HealthCheckRespData | null>(null)
  const ram = shallowRef<HealthCheckRespData | null>(null)
  const loading = shallowRef(false)
  const error = shallowRef('')
  const lastCheckedAt = shallowRef('')

  async function refresh() {
    loading.value = true
    error.value = ''

    try {
      const [cpuResp, ramResp] = await Promise.all([getCpuCheck(), getRamCheck()])
      cpu.value = cpuResp
      ram.value = ramResp
      lastCheckedAt.value = new Date().toISOString()
    } catch (err) {
      error.value = formatErrorMessage(err)
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    void refresh()
  })

  return {
    cpu,
    ram,
    loading,
    error,
    lastCheckedAt,
    refresh,
  }
}
