import { ref, onMounted } from 'vue'

/**
 * 主题类型
 */
export type Theme = 'light' | 'dark' | 'auto'

/**
 * 主题存储 Key
 */
const THEME_STORAGE_KEY = 'cooper-theme'

/**
 * 当前主题（用户选择）
 */
const currentTheme = ref<Theme>('auto')

/**
 * 实际应用的主题（解析后的）
 */
const appliedTheme = ref<'light' | 'dark'>('light')

/**
 * 系统主题媒体查询
 */
let systemThemeQuery: MediaQueryList | null = null

/**
 * 从系统获取主题偏好
 */
function getSystemTheme(): 'light' | 'dark' {
  if (typeof window === 'undefined') return 'light'

  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

/**
 * 解析主题（将 auto 解析为实际主题）
 */
function resolveTheme(theme: Theme): 'light' | 'dark' {
  if (theme === 'auto') {
    return getSystemTheme()
  }
  return theme
}

/**
 * 应用主题到 DOM
 */
function applyTheme(theme: 'light' | 'dark') {
  const html = document.documentElement

  html.removeAttribute('data-theme')

  html.setAttribute('data-theme', theme)

  appliedTheme.value = theme
}

/**
 * 保存主题到 localStorage
 */
function saveTheme(theme: Theme) {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, theme)
  } catch {
    console.error('保存主题失败:', error)
  }
}

/**
 * 从 localStorage 读取主题
 */
function loadTheme(): Theme {
  try {
    const saved = localStorage.getItem(THEME_STORAGE_KEY)
    if (saved === 'light' || saved === 'dark' || saved === 'auto') {
      return saved
    }
  } catch {
    console.error('读取主题失败:', error)
  }
  return 'light'
}

/**
 * 切换主题
 */
function setTheme(theme: Theme) {
  currentTheme.value = theme
  saveTheme(theme)

  const resolved = resolveTheme(theme)
  applyTheme(resolved)
}

/**
 * 切换到下一个主题
 * 循环顺序：light -> dark -> light
 */
function toggleTheme() {
  const themes: Theme[] = ['light', 'dark']
  const currentIndex = themes.indexOf(currentTheme.value)
  const nextIndex = (currentIndex + 1) % themes.length
  setTheme(themes[nextIndex])
}

/**
 * 监听系统主题变化
 */
function watchSystemTheme() {
  if (typeof window === 'undefined') return

  systemThemeQuery = window.matchMedia('(prefers-color-scheme: dark)')

  const handler = (e: MediaQueryListEvent) => {
    if (currentTheme.value === 'auto') {
      const newTheme = e.matches ? 'dark' : 'light'
      applyTheme(newTheme)
    }
  }

  if (systemThemeQuery.addEventListener) {
    systemThemeQuery.addEventListener('change', handler)
  } else if (systemThemeQuery.addListener) {
    systemThemeQuery.addListener(handler)
  }
}

/**
 * 初始化主题
 */
function initTheme() {
  const saved = loadTheme()
  currentTheme.value = saved

  const resolved = resolveTheme(saved)
  applyTheme(resolved)

  watchSystemTheme()
}

/**
 * 主题管理 Composable
 */
export function useTheme() {
  onMounted(() => {
    if (!document.documentElement.hasAttribute('data-theme')) {
      initTheme()
    }
  })

  return {
    /**
     * 当前主题（用户选择）
     */
    currentTheme,

    /**
     * 实际应用的主题
     */
    appliedTheme,

    /**
     * 是否为暗色模式
     */
    isDark: () => appliedTheme.value === 'dark',

    /**
     * 设置主题
     */
    setTheme,

    /**
     * 切换主题
     */
    toggleTheme,

    /**
     * 获取系统主题
     */
    getSystemTheme,
  }
}

if (typeof window !== 'undefined') {
  initTheme()
}
