<template>
  <div
    class="relative flex min-h-screen flex-col overflow-hidden bg-gradient-to-br from-gray-50 via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
  >
    <!-- Background Decorations -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute -right-40 -top-40 h-96 w-96 rounded-full bg-primary-400/20 blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-primary-500/15 blur-3xl"></div>
      <div class="absolute inset-0 bg-[linear-gradient(rgba(20,184,166,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(20,184,166,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-bold text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-3">
          <LocaleSwitcher />
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white">
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-12">
      <div class="mx-auto max-w-2xl">
        <!-- Title -->
        <div class="mb-10 text-center">
          <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 shadow-lg shadow-primary-500/30">
            <svg class="h-8 w-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" />
            </svg>
          </div>
          <h1 class="mb-3 text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">
            {{ t('smsQuery.title') }}
          </h1>
          <p class="text-gray-500 dark:text-dark-400">
            {{ t('smsQuery.subtitle') }}
          </p>
        </div>

        <!-- Search Box -->
        <div class="mb-8 overflow-hidden rounded-2xl border border-gray-200/60 bg-white/80 p-6 shadow-sm backdrop-blur-sm dark:border-dark-700/60 dark:bg-dark-800/80">
          <div class="flex gap-3">
            <div class="relative flex-1">
              <svg class="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400 dark:text-dark-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
              </svg>
              <input
                id="sms-order-no-input"
                v-model="orderNo"
                type="text"
                :placeholder="t('smsQuery.inputPlaceholder')"
                class="w-full rounded-xl border border-gray-200 bg-white py-3 pl-11 pr-4 text-sm text-gray-900 transition-all placeholder:text-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder:text-dark-500 dark:focus:border-primary-500"
                @keyup.enter="doQuery"
              />
            </div>
            <button
              id="sms-query-btn"
              @click="doQuery"
              :disabled="loading || !orderNo.trim()"
              class="inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-primary-500 to-primary-600 px-6 py-3 text-sm font-medium text-white shadow-md shadow-primary-500/30 transition-all hover:from-primary-600 hover:to-primary-700 hover:shadow-lg disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-none"
            >
              <svg v-if="loading" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              {{ t('smsQuery.queryBtn') }}
            </button>
          </div>
        </div>

        <!-- Result Card -->
        <div v-if="result" class="overflow-hidden rounded-2xl border border-gray-200/60 bg-white/80 shadow-sm backdrop-blur-sm dark:border-dark-700/60 dark:bg-dark-800/80">
          <!-- Status Bar -->
          <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <div class="flex items-center gap-2">
              <span
                class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold"
                :class="statusClass(result.status)"
              >
                {{ t(`smsQuery.status.${result.status}`) }}
              </span>
            </div>
            <div class="flex items-center gap-3">
              <button
                v-if="canRefresh(result.status)"
                @click="doQuery"
                :disabled="loading"
                class="inline-flex items-center gap-1.5 rounded-lg bg-primary-50 px-3 py-1 text-xs font-medium text-primary-700 transition-colors hover:bg-primary-100 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-primary-900/30 dark:text-primary-400 dark:hover:bg-primary-900/40"
              >
                <svg
                  class="h-3.5 w-3.5"
                  :class="loading ? 'animate-spin' : ''"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M4.5 4.5a8.25 8.25 0 0114.32 4.5m.68-4.5v4.5h-4.5m.32 8.25A8.25 8.25 0 014.68 13.5m-.68 4.5v-4.5h4.5"
                  />
                </svg>
                {{ t('smsQuery.refreshBtn') }}
              </button>
              <span class="text-xs text-gray-400 dark:text-dark-500">
                {{ formatTime(result.created_at) }}
              </span>
            </div>
          </div>

          <!-- Content -->
          <div class="space-y-4 p-6">
            <!-- Order No -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-dark-500">{{ t('smsQuery.orderNo') }}</label>
              <p class="font-mono text-sm text-gray-900 dark:text-white">{{ result.order_no }}</p>
            </div>

            <!-- Phone -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-dark-500">{{ t('smsQuery.phoneNumber') }}</label>
              <p
                v-if="result.phone_number"
                class="text-lg font-semibold text-primary-600 dark:text-primary-400"
              >
                {{ result.phone_number }}
              </p>
              <p v-else class="text-sm text-gray-500 dark:text-dark-400">{{ t('smsQuery.waitingPhone') }}</p>
            </div>

            <!-- SMS Content -->
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-dark-500">{{ t('smsQuery.smsContent') }}</label>
              <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-700/50">
                <p
                  v-if="result.sms_content"
                  class="whitespace-pre-wrap text-sm leading-relaxed text-gray-700 dark:text-dark-200"
                >
                  {{ result.sms_content }}
                </p>
                <p v-else-if="result.status === 'pending'" class="text-sm text-gray-500 dark:text-dark-400">
                  {{ t('smsQuery.waitingCode') }}
                </p>
                <p v-else class="text-sm text-gray-500 dark:text-dark-400">{{ t('smsQuery.noContent') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="queried && !error" class="rounded-2xl border border-gray-200/60 bg-white/80 p-12 text-center shadow-sm backdrop-blur-sm dark:border-dark-700/60 dark:bg-dark-800/80">
          <svg class="mx-auto mb-4 h-12 w-12 text-gray-300 dark:text-dark-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" />
          </svg>
          <p class="text-gray-500 dark:text-dark-400">{{ t('smsQuery.notFound') }}</p>
        </div>

        <!-- Error State -->
        <div v-if="error" class="rounded-2xl border border-red-200/60 bg-red-50/80 p-6 text-center backdrop-blur-sm dark:border-red-900/40 dark:bg-red-900/10">
          <svg class="mx-auto mb-3 h-10 w-10 text-red-400 dark:text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
          </svg>
          <p class="text-sm font-medium text-red-600 dark:text-red-400">{{ error }}</p>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-6 dark:border-dark-800/50">
      <div class="mx-auto max-w-6xl text-center text-sm text-gray-400 dark:text-dark-500">
        {{ t('smsQuery.footer') }}
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { refreshOrder, type SmsOrderResponse } from '@/api/smsQuery'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

const isDark = ref(document.documentElement.classList.contains('dark'))

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

// Query state
const orderNo = ref('')
const loading = ref(false)
const result = ref<SmsOrderResponse | null>(null)
const error = ref('')
const queried = ref(false)

async function doQuery() {
  const no = orderNo.value.trim()
  if (!no) return

  loading.value = true
  error.value = ''
  queried.value = true

  try {
    result.value = await refreshOrder(no)
  } catch (err: any) {
    if (err?.status === 404 || err?.code === 'SMS_ORDER_NOT_FOUND') {
      result.value = null
    } else if (err?.code === 'SMS_ORDER_FETCH_FAILED') {
      error.value = t('smsQuery.fetchFailed')
    } else {
      error.value = err?.message || t('smsQuery.queryError')
    }
  } finally {
    loading.value = false
  }
}

function canRefresh(status: string): boolean {
  return status === 'created' || status === 'pending'
}

function statusClass(status: string): string {
  switch (status) {
    case 'received':
      return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    case 'pending':
      return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'created':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
    case 'failed':
      return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
    case 'expired':
      return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400'
    default:
      return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400'
  }
}

function formatTime(iso: string): string {
  if (!iso) return ''
  return new Date(iso).toLocaleString()
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
})
</script>
