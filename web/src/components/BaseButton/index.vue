<template>
  <button
    :type="type"
    @click="$emit('click', $event)"
    :class="buttonClasses"
    :style="variant === 'outline' ? outlineStyle : undefined"
    :disabled="disabled || loading"
  >
    <Loader2 v-if="loading" :size="iconSize" class="mr-1.5 flex-shrink-0 animate-spin" />
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'

interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost' | 'outline'
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
  loading: false,
})

defineEmits<{
  click: [event: MouseEvent]
}>()

const iconSize = computed(() => {
  const sizes = { xs: 12, sm: 14, md: 16, lg: 18 }
  return sizes[props.size]
})

const outlineStyle = computed(() => ({
  color: 'var(--color-text-primary)',
  borderColor: 'var(--color-border-primary)',
}))

const buttonClasses = computed(() => {
  const base = props.fullWidth ? 'w-full' : ''
  const disabled = props.disabled ? 'opacity-50 cursor-not-allowed' : 'active:scale-[0.98]'

  const sizes = {
    xs: 'px-2 py-0.5 text-xs rounded',
    sm: 'px-2.5 py-1 text-xs rounded',
    md: 'px-3 py-1.5 text-sm rounded-md',
    lg: 'px-4 py-2 text-base rounded-md',
  }

  const variants = {
    primary: 'bg-[var(--color-primary)] hover:bg-primary-hover text-primary-text shadow-sm',
    secondary: 'bg-bg-tertiary text-text-primary hover:bg-bg-active',
    danger: 'bg-error text-white hover:bg-error-hover shadow-sm',
    success: 'bg-success text-white hover:bg-success-hover shadow-sm',
    ghost: 'bg-transparent text-text-secondary hover:bg-bg-hover border border-border-primary',
    outline: 'bg-transparent border-2 hover:bg-bg-hover',
  }

  return `${base} ${disabled} ${sizes[props.size]} ${variants[props.variant]} font-medium transition-all duration-200 focus:outline-none whitespace-nowrap inline-flex items-center justify-center shrink-0`
})
</script>
