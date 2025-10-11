<template>
  <button
    :type="type"
    @click="$emit('click', $event)"
    :class="buttonClasses"
    :disabled="disabled || loading"
  >
    <Loader2 v-if="loading" :size="iconSize" class="mr-1.5 inline-block animate-spin" />
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'

interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  fullWidth?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  fullWidth: false,
  loading: false
})

defineEmits<{
  click: [event: MouseEvent]
}>()

const iconSize = computed(() => {
  const sizes = { xs: 12, sm: 14, md: 16, lg: 18 }
  return sizes[props.size]
})

const buttonClasses = computed(() => {
  const base = props.fullWidth ? 'w-full' : ''
  const disabled = props.disabled ? 'opacity-50 cursor-not-allowed' : 'active:scale-[0.98]'

  const sizes = {
    xs: 'px-2 py-0.5 text-xs rounded',
    sm: 'px-2.5 py-1 text-xs rounded',
    md: 'px-3 py-1.5 text-sm rounded-md',
    lg: 'px-4 py-2 text-base rounded-md'
  }

  const variants = {
    primary: 'bg-blue-500 text-white hover:bg-blue-600 focus:ring-blue-100 shadow-sm',
    secondary: 'bg-slate-100 text-slate-700 hover:bg-slate-200 focus:ring-slate-100',
    danger: 'bg-rose-500 text-white hover:bg-rose-600 focus:ring-rose-100 shadow-sm',
    success: 'bg-emerald-500 text-white hover:bg-emerald-600 focus:ring-emerald-100 shadow-sm',
    ghost: 'bg-transparent text-slate-600 hover:bg-slate-50 focus:ring-slate-100 border border-slate-200'
  }

  return `${base} ${disabled} ${sizes[props.size]} ${variants[props.variant]} font-medium transition-all duration-200 focus:outline-none focus:ring-2`
})
</script>
