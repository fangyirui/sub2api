<template>
  <div class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-semibold tracking-tight text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-3">
          <LocaleSwitcher />
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="flex-1 w-full max-w-2xl mx-auto px-6 py-12">
      <!-- Title -->
      <div class="text-center mb-10">
        <h1 class="text-3xl sm:text-4xl font-bold tracking-tight mb-3 text-gray-900 dark:text-white">
          {{ t('keyTest.title') }}
        </h1>
        <p class="text-gray-500 dark:text-dark-400 text-base max-w-md mx-auto">
          {{ t('keyTest.subtitle') }}
        </p>
      </div>

      <!-- Form -->
      <div class="rounded-2xl border border-gray-200/60 bg-white/80 p-6 shadow-sm backdrop-blur-sm dark:border-dark-700/60 dark:bg-dark-800/80 space-y-5">
        <!-- Base URL -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('keyTest.baseUrl') }}</label>
          <input
            v-model="baseUrl"
            type="text"
            :placeholder="DEFAULT_URL"
            class="w-full rounded-xl border border-gray-200 bg-white py-3 px-4 text-sm text-gray-900 transition-all placeholder:text-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder:text-dark-500"
          />
        </div>

        <!-- API Key -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('keyTest.apiKey') }}</label>
          <div class="relative">
            <input
              v-model="apiKey"
              :type="keyVisible ? 'text' : 'password'"
              :placeholder="t('keyTest.apiKeyPlaceholder')"
              class="w-full rounded-xl border border-gray-200 bg-white py-3 px-4 pr-12 text-sm text-gray-900 transition-all placeholder:text-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder:text-dark-500"
            />
            <button
              @click="keyVisible = !keyVisible"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700 dark:text-dark-500 dark:hover:text-white transition-colors"
            >
              <svg v-if="!keyVisible" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
              <svg v-else class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Model -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1.5">{{ t('keyTest.model') }}</label>
          <input
            v-model="model"
            type="text"
            :placeholder="DEFAULT_MODEL"
            class="w-full rounded-xl border border-gray-200 bg-white py-3 px-4 text-sm text-gray-900 transition-all placeholder:text-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder:text-dark-500"
          />
        </div>

        <!-- Test Button -->
        <button
          @click="runTest"
          :disabled="isTesting"
          class="w-full h-12 rounded-xl bg-primary-500 hover:bg-primary-600 text-white font-medium text-sm transition-all active:scale-[0.98] flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed"
        >
          <svg v-if="isTesting" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          {{ isTesting ? t('keyTest.testing') : t('keyTest.testBtn') }}
        </button>
      </div>

      <!-- Result -->
      <div v-if="result" class="mt-6 rounded-2xl border p-5 shadow-sm backdrop-blur-sm" :class="resultClass">
        <div class="flex items-center gap-2 mb-3">
          <svg v-if="result.success" class="w-5 h-5 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <svg v-else class="w-5 h-5 text-red-600 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span class="font-medium text-sm" :class="result.success ? 'text-green-700 dark:text-green-300' : 'text-red-700 dark:text-red-300'">
            {{ result.success ? t('keyTest.success') : t('keyTest.failed') }}
          </span>
          <span v-if="result.latency" class="ml-auto text-xs text-gray-500 dark:text-dark-400">{{ result.latency }}ms</span>
        </div>
        <div v-if="result.content" class="rounded-lg bg-white/60 dark:bg-dark-900/60 p-3 text-sm text-gray-700 dark:text-dark-300 whitespace-pre-wrap break-words max-h-60 overflow-y-auto">{{ result.content }}</div>
        <div v-if="result.error" class="rounded-lg bg-white/60 dark:bg-dark-900/60 p-3 text-sm text-red-600 dark:text-red-400 whitespace-pre-wrap break-words max-h-60 overflow-y-auto">{{ result.error }}</div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div class="mx-auto flex max-w-6xl items-center justify-center">
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const DEFAULT_URL = 'http://www.buglaoge.cc'
const DEFAULT_MODEL = 'claude-sonnet-4-6'

const { t } = useI18n()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const isDark = ref(document.documentElement.classList.contains('dark'))
const currentYear = computed(() => new Date().getFullYear())

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

const baseUrl = ref(DEFAULT_URL)
const apiKey = ref('')
const keyVisible = ref(false)
const model = ref(DEFAULT_MODEL)
const isTesting = ref(false)
const result = ref<{ success: boolean; content?: string; error?: string; latency?: number } | null>(null)

const resultClass = computed(() => {
  if (!result.value) return ''
  return result.value.success
    ? 'border-green-200/60 bg-green-50/80 dark:border-green-900/40 dark:bg-green-900/10'
    : 'border-red-200/60 bg-red-50/80 dark:border-red-900/40 dark:bg-red-900/10'
})

async function runTest() {
  if (!apiKey.value.trim()) {
    result.value = { success: false, error: t('keyTest.noKey') }
    return
  }

  isTesting.value = true
  result.value = null

  const url = (baseUrl.value.trim() || DEFAULT_URL).replace(/\/+$/, '')
  const fullUrl = url.startsWith('http') ? url : `https://${url}`
  const start = Date.now()

  try {
    const res = await fetch(`${fullUrl}/v1/chat/completions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${apiKey.value.trim()}`
      },
      body: JSON.stringify({
        model: model.value.trim() || DEFAULT_MODEL,
        messages: [{ role: 'user', content: 'Hi' }],
        max_tokens: 50,
        stream: false
      })
    })

    const latency = Date.now() - start
    const data = await res.json()

    if (res.ok && data.choices?.[0]?.message?.content) {
      result.value = { success: true, content: data.choices[0].message.content, latency }
    } else {
      result.value = { success: false, error: data.error?.message || JSON.stringify(data, null, 2), latency }
    }
  } catch (e: unknown) {
    const latency = Date.now() - start
    result.value = { success: false, error: e instanceof Error ? e.message : String(e), latency }
  } finally {
    isTesting.value = false
  }
}
</script>
