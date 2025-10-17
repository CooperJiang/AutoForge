<template>
  <div
    :class="[
      'flex items-center justify-between',
      bordered && 'border-t border-border-primary',
      compact ? 'px-3 py-1.5' : 'px-4 py-2.5',
    ]"
  >
    <!-- 左侧信息 -->
    <div class="flex items-center gap-4">
      <!-- 总数显示 -->
      <div v-if="showTotal" class="text-xs" style="color: var(--color-text-secondary)">
        共
        <span class="font-semibold" style="color: var(--color-text-primary)">{{ total }}</span>
        条
        <span v-if="showRange" class="ml-2"> (第 {{ rangeStart }}-{{ rangeEnd }} 条) </span>
      </div>

      <!-- 每页条数选择器 -->
      <div v-if="showSizeChanger" class="flex items-center gap-2">
        <span class="text-xs" style="color: var(--color-text-secondary)">每页</span>
        <select
          :value="pageSize"
          @change="handleSizeChange"
          :class="[
            'h-[28px] px-2 text-xs rounded border focus:outline-none',
            'bg-bg-elevated border-border-primary',
          ]"
          style="
            color: var(--color-text-primary);
            background-color: var(--color-bg-elevated);
            border-color: var(--color-border-primary);
          "
        >
          <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }} 条</option>
        </select>
      </div>
    </div>

    <!-- 右侧分页控制 -->
    <div class="flex items-center gap-2">
      <!-- 上一页按钮 -->
      <button
        :disabled="current === 1"
        @click="handlePageChange(current - 1)"
        class="px-3 py-1 text-xs font-medium border rounded transition-all"
        :style="{
          color: current === 1 ? 'var(--color-text-disabled)' : 'var(--color-text-primary)',
          backgroundColor: 'var(--color-bg-elevated)',
          borderColor: 'var(--color-border-primary)',
          cursor: current === 1 ? 'not-allowed' : 'pointer',
          opacity: current === 1 ? '0.5' : '1',
        }"
        @mouseenter="
          (e) =>
            !props.current || props.current !== 1
              ? ((e.target as HTMLElement).style.borderColor = 'var(--color-primary)')
              : null
        "
        @mouseleave="
          (e) => ((e.target as HTMLElement).style.borderColor = 'var(--color-border-primary)')
        "
      >
        上一页
      </button>

      <!-- 页码按钮组 -->
      <div class="flex gap-1">
        <button
          v-for="page in visiblePages"
          :key="page"
          @click="page !== '...' && handlePageChange(page as number)"
          :disabled="page === '...'"
          :class="[
            'min-w-[32px] h-[28px] px-2 text-xs font-medium border rounded transition-colors flex items-center justify-center',
            page === current && 'bg-primary text-white border-primary',
            page === '...' && 'border-transparent cursor-default',
            page !== current && page !== '...' && 'hover:border-primary',
          ]"
          :style="
            page === current
              ? {
                  backgroundColor: 'var(--color-primary)',
                  color: 'var(--color-primary-text)',
                  borderColor: 'var(--color-primary)',
                }
              : page === '...'
                ? {
                    color: 'var(--color-text-tertiary)',
                    borderColor: 'transparent',
                  }
                : {
                    color: 'var(--color-text-primary)',
                    backgroundColor: 'var(--color-bg-elevated)',
                    borderColor: 'var(--color-border-primary)',
                  }
          "
        >
          {{ page }}
        </button>
      </div>

      <!-- 下一页按钮 -->
      <button
        :disabled="current === totalPages"
        @click="handlePageChange(current + 1)"
        class="px-3 py-1 text-xs font-medium border rounded transition-all"
        :style="{
          color:
            current === totalPages ? 'var(--color-text-disabled)' : 'var(--color-text-primary)',
          backgroundColor: 'var(--color-bg-elevated)',
          borderColor: 'var(--color-border-primary)',
          cursor: current === totalPages ? 'not-allowed' : 'pointer',
          opacity: current === totalPages ? '0.5' : '1',
        }"
        @mouseenter="
          (e) =>
            props.current < totalPages
              ? ((e.target as HTMLElement).style.borderColor = 'var(--color-primary)')
              : null
        "
        @mouseleave="
          (e) => ((e.target as HTMLElement).style.borderColor = 'var(--color-border-primary)')
        "
      >
        下一页
      </button>

      <!-- 总页数显示 -->
      <div v-if="showPageCount" class="text-xs ml-1" style="color: var(--color-text-secondary)">
        共 {{ totalPages }} 页
      </div>

      <!-- 快速跳转 -->
      <div v-if="showJumper" class="flex items-center gap-2 ml-2">
        <span class="text-xs" style="color: var(--color-text-secondary)">跳至</span>
        <input
          v-model.number="jumpPage"
          type="number"
          min="1"
          :max="totalPages"
          @keyup.enter="handleJump"
          @blur="handleJump"
          :class="[
            'w-14 h-[28px] px-2 text-xs text-center border rounded focus:outline-none focus:border-primary',
          ]"
          style="
            color: var(--color-text-primary);
            background-color: var(--color-bg-elevated);
            border-color: var(--color-border-primary);
          "
          placeholder="页码"
        />
        <span class="text-xs" style="color: var(--color-text-secondary)">页</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

interface Props {
  /** 当前页码 */
  current: number
  /** 每页条数 */
  pageSize: number
  /** 总记录数 */
  total: number
  /** 每页条数选项 */
  pageSizeOptions?: number[]
  /** 是否显示总数 */
  showTotal?: boolean
  /** 是否显示范围 (如: 显示第 1-10 条) */
  showRange?: boolean
  /** 是否显示每页条数选择器 */
  showSizeChanger?: boolean
  /** 是否显示总页数 */
  showPageCount?: boolean
  /** 是否显示边框 */
  bordered?: boolean
  /** 紧凑模式 */
  compact?: boolean
  /** 是否显示页面跳转 */
  showJumper?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  pageSizeOptions: () => [10, 20, 30, 50, 100],
  showTotal: true,
  showRange: false,
  showSizeChanger: false,
  showPageCount: true,
  bordered: true,
  compact: false,
  showJumper: false,
})

const emit = defineEmits<{
  'update:current': [page: number]
  'update:pageSize': [size: number]
  change: [page: number]
  sizeChange: [size: number]
}>()

// 跳转页码输入
const jumpPage = ref<number>(props.current)

// 监听 current 变化，同步到 jumpPage
watch(
  () => props.current,
  (newVal) => {
    jumpPage.value = newVal
  }
)

// 总页数
const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

// 显示范围
const rangeStart = computed(() => (props.current - 1) * props.pageSize + 1)
const rangeEnd = computed(() => Math.min(props.current * props.pageSize, props.total))

// 可见页码
const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const current = props.current
  const total = totalPages.value

  if (total <= 7) {
    // 总页数小于等于 7，全部显示
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // 总页数大于 7，显示部分页码
    pages.push(1)

    if (current > 3) {
      pages.push('...')
    }

    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < total - 2) {
      pages.push('...')
    }

    pages.push(total)
  }

  return pages
})

// 页码变化处理
const handlePageChange = (page: number) => {
  if (page === props.current || page < 1 || page > totalPages.value) return
  emit('update:current', page)
  emit('change', page)
}

// 每页条数变化处理
const handleSizeChange = (e: Event) => {
  const size = Number((e.target as HTMLSelectElement).value)
  if (size === props.pageSize) return

  emit('update:pageSize', size)
  emit('sizeChange', size)

  // 切换每页条数后，如果当前页超出范围，跳转到第一页
  const newTotalPages = Math.ceil(props.total / size)
  if (props.current > newTotalPages && newTotalPages > 0) {
    handlePageChange(1)
  }
}

// 跳转处理
const handleJump = () => {
  const page = jumpPage.value
  if (page && page >= 1 && page <= totalPages.value && page !== props.current) {
    handlePageChange(page)
  } else {
    jumpPage.value = props.current
  }
}
</script>
