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
            <img v-if="siteLogo" :src="siteLogo" alt="Logo" class="h-full w-full object-contain" />
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
      <div class="mx-auto max-w-6xl">
        <!-- Title -->
        <div class="mb-10 text-center">
          <h1 class="mb-3 text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">
            {{ t('models.title') }}
          </h1>
          <p class="text-gray-500 dark:text-dark-400">
            {{ t('models.subtitle') }}
          </p>
        </div>

        <!-- Filter Tabs -->
        <div class="mb-8 flex flex-wrap justify-center gap-2">
          <button
            v-for="tier in tiers"
            :key="tier.key"
            @click="activeTier = tier.key"
            class="rounded-full px-4 py-1.5 text-sm font-medium transition-all"
            :class="activeTier === tier.key
              ? 'bg-primary-500 text-white shadow-md shadow-primary-500/30'
              : 'bg-white text-gray-600 hover:bg-gray-100 dark:bg-dark-800 dark:text-dark-300 dark:hover:bg-dark-700'"
          >
            {{ tier.label }}
          </button>
        </div>

        <!-- Model Cards Grid -->
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="model in filteredModels"
            :key="model.id"
            class="group relative overflow-hidden rounded-2xl border border-gray-200/60 bg-white/80 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-lg dark:border-dark-700/60 dark:bg-dark-800/80"
          >
            <!-- Tier Badge -->
            <div class="mb-4 flex items-center justify-between">
              <span
                class="rounded-full px-3 py-1 text-xs font-semibold"
                :class="tierBadgeClass(model.tier)"
              >
                {{ model.tier }}
              </span>
              <span v-if="model.latest" class="rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/30 dark:text-green-400">
                {{ t('models.latest') }}
              </span>
            </div>

            <!-- Model Name -->
            <h3 class="mb-1 text-lg font-bold text-gray-900 dark:text-white">
              {{ model.displayName }}
            </h3>
            <p class="mb-4 text-sm text-gray-500 dark:text-dark-400">
              {{ model.id }}
            </p>

            <!-- Context Window -->
            <div class="mb-4 rounded-xl bg-gray-50 p-3 dark:bg-dark-700/50">
              <div class="mb-1 text-xs font-medium text-gray-400 dark:text-dark-500">
                {{ t('models.contextWindow') }}
              </div>
              <div class="flex items-baseline gap-1">
                <span class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ formatTokens(model.maxInputTokens) }}
                </span>
                <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('models.input') }}</span>
                <span class="mx-1 text-gray-300 dark:text-dark-600">/</span>
                <span class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ formatTokens(model.maxOutputTokens) }}
                </span>
                <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('models.output') }}</span>
              </div>
            </div>

            <!-- Pricing -->
            <div class="mb-4 space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-500 dark:text-dark-400">{{ t('models.inputPrice') }}</span>
                <span class="font-semibold text-gray-900 dark:text-white">
                  ${{ formatPrice(model.inputCostPerToken) }} / MTok
                </span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-500 dark:text-dark-400">{{ t('models.outputPrice') }}</span>
                <span class="font-semibold text-gray-900 dark:text-white">
                  ${{ formatPrice(model.outputCostPerToken) }} / MTok
                </span>
              </div>
              <div v-if="model.cacheReadCost" class="flex items-center justify-between text-sm">
                <span class="text-gray-500 dark:text-dark-400">{{ t('models.cacheReadPrice') }}</span>
                <span class="font-semibold text-gray-900 dark:text-white">
                  ${{ formatPrice(model.cacheReadCost) }} / MTok
                </span>
              </div>
            </div>

            <!-- Capability Tags -->
            <div class="flex flex-wrap gap-1.5">
              <span
                v-for="cap in model.capabilities"
                :key="cap"
                class="rounded-md bg-primary-50 px-2 py-0.5 text-xs font-medium text-primary-700 dark:bg-primary-900/20 dark:text-primary-400"
              >
                {{ t(`models.caps.${cap}`) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-6 dark:border-dark-800/50">
      <div class="mx-auto max-w-6xl text-center text-sm text-gray-400 dark:text-dark-500">
        {{ t('models.pricingNote') }}
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

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '')
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

onMounted(() => {
  initTheme()
  authStore.checkAuth()
})

// --- Tier filter ---
const activeTier = ref('all')
const tiers = computed(() => [
  { key: 'all', label: t('models.filterAll') },
  { key: 'Opus', label: 'Opus' },
  { key: 'Sonnet', label: 'Sonnet' },
  { key: 'Haiku', label: 'Haiku' },
])

interface ClaudeModel {
  id: string
  displayName: string
  tier: 'Opus' | 'Sonnet' | 'Haiku'
  latest: boolean
  maxInputTokens: number
  maxOutputTokens: number
  inputCostPerToken: number
  outputCostPerToken: number
  cacheReadCost: number
  capabilities: string[]
}

const models: ClaudeModel[] = [
  {
    id: 'claude-opus-4-6',
    displayName: 'Claude Opus 4.6',
    tier: 'Opus',
    latest: true,
    maxInputTokens: 1000000,
    maxOutputTokens: 128000,
    inputCostPerToken: 5e-06,
    outputCostPerToken: 2.5e-05,
    cacheReadCost: 5e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse', 'reasoning', 'webSearch'],
  },
  {
    id: 'claude-sonnet-4-6',
    displayName: 'Claude Sonnet 4.6',
    tier: 'Sonnet',
    latest: true,
    maxInputTokens: 1000000,
    maxOutputTokens: 128000,
    inputCostPerToken: 2e-06,
    outputCostPerToken: 1e-05,
    cacheReadCost: 2e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse', 'reasoning', 'webSearch'],
  },
  {
    id: 'claude-haiku-4-5',
    displayName: 'Claude Haiku 4.5',
    tier: 'Haiku',
    latest: true,
    maxInputTokens: 200000,
    maxOutputTokens: 8192,
    inputCostPerToken: 1e-06,
    outputCostPerToken: 5e-06,
    cacheReadCost: 1e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf'],
  },
  {
    id: 'claude-opus-4-5',
    displayName: 'Claude Opus 4.5',
    tier: 'Opus',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 32000,
    inputCostPerToken: 1.5e-05,
    outputCostPerToken: 7.5e-05,
    cacheReadCost: 1.5e-06,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse', 'reasoning'],
  },
  {
    id: 'claude-sonnet-4-5',
    displayName: 'Claude Sonnet 4.5',
    tier: 'Sonnet',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 64000,
    inputCostPerToken: 3e-06,
    outputCostPerToken: 1.5e-05,
    cacheReadCost: 3e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse', 'reasoning', 'webSearch'],
  },
  {
    id: 'claude-opus-4-1',
    displayName: 'Claude Opus 4.1',
    tier: 'Opus',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 32000,
    inputCostPerToken: 1.5e-05,
    outputCostPerToken: 7.5e-05,
    cacheReadCost: 1.5e-06,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse', 'reasoning'],
  },
  {
    id: 'claude-3-7-sonnet-latest',
    displayName: 'Claude 3.7 Sonnet',
    tier: 'Sonnet',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 64000,
    inputCostPerToken: 3e-06,
    outputCostPerToken: 1.5e-05,
    cacheReadCost: 3e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse', 'reasoning'],
  },
  {
    id: 'claude-3-5-sonnet-latest',
    displayName: 'Claude 3.5 Sonnet',
    tier: 'Sonnet',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 8192,
    inputCostPerToken: 3e-06,
    outputCostPerToken: 1.5e-05,
    cacheReadCost: 3e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'computerUse'],
  },
  {
    id: 'claude-3-5-haiku-latest',
    displayName: 'Claude 3.5 Haiku',
    tier: 'Haiku',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 8192,
    inputCostPerToken: 1e-06,
    outputCostPerToken: 5e-06,
    cacheReadCost: 1e-07,
    capabilities: ['vision', 'toolUse', 'caching', 'pdf', 'webSearch'],
  },
  {
    id: 'claude-3-opus-latest',
    displayName: 'Claude 3 Opus',
    tier: 'Opus',
    latest: false,
    maxInputTokens: 200000,
    maxOutputTokens: 4096,
    inputCostPerToken: 1.5e-05,
    outputCostPerToken: 7.5e-05,
    cacheReadCost: 1.5e-06,
    capabilities: ['vision', 'toolUse', 'caching'],
  },
]

const filteredModels = computed(() => {
  if (activeTier.value === 'all') return models
  return models.filter(m => m.tier === activeTier.value)
})

function tierBadgeClass(tier: string): string {
  switch (tier) {
    case 'Opus':
      return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400'
    case 'Sonnet':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
    case 'Haiku':
      return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
    default:
      return 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-dark-300'
  }
}

function formatTokens(tokens: number): string {
  if (tokens >= 1000000) return `${tokens / 1000000}M`
  if (tokens >= 1000) return `${tokens / 1000}K`
  return String(tokens)
}

function formatPrice(costPerToken: number): string {
  const perMillion = costPerToken * 1000000
  return perMillion % 1 === 0 ? perMillion.toFixed(0) : perMillion.toFixed(2)
}
</script>
