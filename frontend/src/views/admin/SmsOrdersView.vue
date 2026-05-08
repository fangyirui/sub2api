<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="space-y-3">
          <!-- Top action row -->
          <div class="flex flex-wrap items-center gap-3">
            <Select
              v-model="filters.status"
              :options="filterStatusOptions"
              class="w-40"
              @change="handleStatusChange"
            />

            <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
              <button
                @click="loadOrders"
                :disabled="loading"
                class="btn btn-secondary"
                :title="t('common.refresh')"
              >
                <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
              </button>
              <div class="flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 py-1.5 dark:border-dark-700 dark:bg-dark-800">
                <label class="text-xs font-medium text-gray-500 dark:text-dark-400">
                  {{ t('admin.smsOrders.count') }}
                </label>
                <input
                  v-model.number="generateCount"
                  type="number"
                  min="1"
                  max="50"
                  class="w-16 border-none bg-transparent p-0 text-sm focus:outline-none focus:ring-0 dark:text-white"
                />
              </div>
              <button
                @click="handleGenerate"
                :disabled="generating || generateCount < 1 || generateCount > 50"
                class="btn btn-primary"
              >
                <Icon name="plus" size="md" class="mr-1" />
                {{ generating ? t('common.creating') : t('admin.smsOrders.generate') }}
              </button>
            </div>
          </div>

          <!-- Copy template editor -->
          <div class="rounded-xl border border-gray-200 bg-white p-3 dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-2 flex items-center justify-between">
              <label class="text-xs font-medium text-gray-500 dark:text-dark-400">
                {{ t('admin.smsOrders.template.label') }}
                <span class="ml-1 font-mono text-xs text-primary-600 dark:text-primary-400">{order_no}</span>
              </label>
              <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('admin.smsOrders.template.hint') }}</span>
            </div>
            <textarea
              v-model="copyTemplate"
              rows="2"
              class="input w-full font-mono text-sm"
              :placeholder="t('admin.smsOrders.template.placeholder')"
            ></textarea>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="orders"
          :loading="loading"
          row-key="order_no"
        >
          <template #cell-order_no="{ value }">
            <code class="font-mono text-sm text-gray-900 dark:text-gray-100">{{ value }}</code>
          </template>

          <template #cell-phone_number="{ value }">
            <span v-if="value" class="font-mono text-sm text-primary-600 dark:text-primary-400">
              {{ value }}
            </span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">—</span>
          </template>

          <template #cell-status="{ value }">
            <span
              :class="[
                'badge',
                getStatusClass(value)
              ]"
            >
              {{ t(`admin.smsOrders.status.${value}`, value) }}
            </span>
          </template>

          <template #cell-sms_content="{ value }">
            <span
              v-if="value"
              class="line-clamp-1 max-w-md text-sm text-gray-700 dark:text-dark-200"
              :title="value"
            >
              {{ value }}
            </span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">—</span>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">
              {{ formatDateTime(value) }}
            </span>
          </template>

          <template #cell-actions="{ row }">
            <button
              @click="handleCopy(row.order_no)"
              :class="[
                'flex items-center gap-1 rounded-lg px-2 py-1 text-xs transition-colors',
                copiedOrderNo === row.order_no
                  ? 'text-green-600 dark:text-green-400'
                  : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-700'
              ]"
              :title="t('admin.smsOrders.copyWithTemplate')"
            >
              <Icon v-if="copiedOrderNo !== row.order_no" name="copy" size="sm" />
              <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>{{ t('admin.smsOrders.copy') }}</span>
            </button>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import type { Column } from '@/components/common/types'
import { useAppStore } from '@/stores'
import { useClipboard } from '@/composables/useClipboard'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'
import { adminAPI } from '@/api/admin'
import type { SmsOrder } from '@/api/admin/smsOrders'

const { t } = useI18n()
const appStore = useAppStore()
const { copyToClipboard } = useClipboard()

const TEMPLATE_STORAGE_KEY = 'sms-order-copy-template'
const ORDER_PLACEHOLDER = '{order_no}'

const orders = ref<SmsOrder[]>([])
const loading = ref(false)
const generating = ref(false)
const generateCount = ref(1)
const copiedOrderNo = ref<string | null>(null)

const filters = reactive({ status: '' })
const pagination = reactive({
  page: 1,
  page_size: getPersistedPageSize(),
  total: 0
})

const copyTemplate = ref(loadTemplate())

watch(copyTemplate, (val) => {
  try {
    localStorage.setItem(TEMPLATE_STORAGE_KEY, val)
  } catch {
    // ignore quota / private mode errors
  }
})

function loadTemplate(): string {
  try {
    const stored = localStorage.getItem(TEMPLATE_STORAGE_KEY)
    if (stored) return stored
  } catch {
    // ignore
  }
  return t('admin.smsOrders.template.default')
}

const filterStatusOptions = computed(() => [
  { value: '', label: t('admin.smsOrders.statusFilter.all') },
  { value: 'created', label: t('admin.smsOrders.status.created') },
  { value: 'pending', label: t('admin.smsOrders.status.pending') },
  { value: 'received', label: t('admin.smsOrders.status.received') },
  { value: 'expired', label: t('admin.smsOrders.status.expired') },
  { value: 'failed', label: t('admin.smsOrders.status.failed') }
])

const columns = computed<Column[]>(() => [
  { key: 'order_no', label: t('admin.smsOrders.columns.orderNo') },
  { key: 'phone_number', label: t('admin.smsOrders.columns.phone') },
  { key: 'status', label: t('admin.smsOrders.columns.status') },
  { key: 'sms_content', label: t('admin.smsOrders.columns.smsContent') },
  { key: 'created_at', label: t('admin.smsOrders.columns.createdAt') },
  { key: 'actions', label: t('admin.smsOrders.columns.actions') }
])

function getStatusClass(status: string): string {
  switch (status) {
    case 'received':
      return 'badge-success'
    case 'pending':
      return 'badge-warning'
    case 'created':
      return 'badge-info'
    case 'failed':
      return 'badge-danger'
    case 'expired':
      return 'badge-gray'
    default:
      return 'badge-gray'
  }
}

function formatDateTime(value: string): string {
  if (!value) return ''
  return value
}

let abortController: AbortController | null = null

async function loadOrders() {
  if (abortController) abortController.abort()
  const ctrl = new AbortController()
  abortController = ctrl
  loading.value = true

  try {
    const res = await adminAPI.smsOrders.list(
      pagination.page,
      pagination.page_size,
      filters.status || undefined,
      { signal: ctrl.signal }
    )
    if (ctrl.signal.aborted || abortController !== ctrl) return
    orders.value = res.items
    pagination.total = res.total
  } catch (err: any) {
    if (
      ctrl.signal.aborted ||
      abortController !== ctrl ||
      err?.name === 'AbortError' ||
      err?.code === 'ERR_CANCELED'
    ) {
      return
    }
    appStore.showError(t('admin.smsOrders.loadFailed'))
    console.error('Error loading sms orders:', err)
  } finally {
    if (abortController === ctrl) {
      loading.value = false
      abortController = null
    }
  }
}

function handleStatusChange() {
  pagination.page = 1
  loadOrders()
}

function handlePageChange(page: number) {
  pagination.page = page
  loadOrders()
}

function handlePageSizeChange(pageSize: number) {
  pagination.page_size = pageSize
  pagination.page = 1
  loadOrders()
}

async function handleGenerate() {
  if (generateCount.value < 1 || generateCount.value > 50) return
  generating.value = true
  try {
    const res = await adminAPI.smsOrders.batchGenerate(generateCount.value)
    appStore.showSuccess(
      t('admin.smsOrders.generated', { count: res.items.length })
    )
    pagination.page = 1
    if (filters.status && filters.status !== 'created') {
      filters.status = ''
    }
    await loadOrders()
  } catch (err: any) {
    appStore.showError(err?.message || t('admin.smsOrders.generateFailed'))
  } finally {
    generating.value = false
  }
}

async function handleCopy(orderNo: string) {
  const template = copyTemplate.value || ORDER_PLACEHOLDER
  const text = template.includes(ORDER_PLACEHOLDER)
    ? template.split(ORDER_PLACEHOLDER).join(orderNo)
    : `${template}\n${orderNo}`
  const ok = await copyToClipboard(text)
  if (ok) {
    copiedOrderNo.value = orderNo
    setTimeout(() => {
      if (copiedOrderNo.value === orderNo) copiedOrderNo.value = null
    }, 2000)
  }
}

onMounted(() => {
  loadOrders()
})
</script>
