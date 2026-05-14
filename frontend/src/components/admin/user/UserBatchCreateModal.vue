<template>
  <BaseDialog
    :show="show"
    :title="resultMode ? t('admin.users.batchCreateResult') : t('admin.users.batchCreate')"
    width="normal"
    @close="handleClose"
  >
    <!-- Form Mode -->
    <form v-if="!resultMode" id="batch-create-form" @submit.prevent="submit" class="space-y-5">
      <div>
        <label class="input-label">{{ t('admin.users.batchCount') }}</label>
        <input v-model.number="form.count" type="number" min="1" max="100" required class="input" />
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="input-label">{{ t('admin.users.columns.balance') }}</label>
          <input v-model.number="form.balance" type="number" step="any" class="input" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.users.columns.concurrency') }}</label>
          <input v-model.number="form.concurrency" type="number" class="input" />
        </div>
      </div>
    </form>

    <!-- Result Mode -->
    <div v-else class="space-y-4">
      <p class="text-sm text-gray-600 dark:text-dark-400">
        {{ t('admin.users.batchCreateSuccess', { count: results.length }) }}
      </p>
      <textarea
        ref="resultTextarea"
        readonly
        :value="resultText"
        class="input w-full h-64 font-mono text-sm resize-none"
      />
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button @click="handleClose" type="button" class="btn btn-secondary">
          {{ resultMode ? t('common.close') : t('common.cancel') }}
        </button>
        <button v-if="resultMode" @click="copyResult" type="button" class="btn btn-primary">
          {{ copied ? t('common.copied') : t('common.copy') }}
        </button>
        <button v-else type="submit" form="batch-create-form" :disabled="loading" class="btn btn-primary">
          {{ loading ? t('admin.users.creating') : t('common.create') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import BaseDialog from '@/components/common/BaseDialog.vue'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ close: []; success: [] }>()

const { t } = useI18n()

const form = reactive({ count: 10, balance: 10, concurrency: 5 })
const loading = ref(false)
const resultMode = ref(false)
const results = ref<{ email: string; password: string }[]>([])
const copied = ref(false)
const resultTextarea = ref<HTMLTextAreaElement | null>(null)

const resultText = ref('')

async function submit() {
  loading.value = true
  try {
    const res = await adminAPI.users.batchCreate(form.count, form.balance, form.concurrency)
    results.value = res.users
    resultText.value = res.users.map(u => `${u.email},${u.password}`).join('\n')
    resultMode.value = true
    emit('success')
  } catch {
    // error handled by global interceptor
  } finally {
    loading.value = false
  }
}

function copyResult() {
  navigator.clipboard.writeText(resultText.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

function handleClose() {
  resultMode.value = false
  results.value = []
  resultText.value = ''
  copied.value = false
  form.count = 10
  form.balance = 10
  form.concurrency = 5
  emit('close')
}
</script>
