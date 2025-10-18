<template>
  <div class="bg-bg-tertiary/30 border border-border-secondary rounded-lg p-4">
    <div class="flex items-center justify-between mb-3">
      <div class="flex items-center gap-2">
        <Lightbulb class="w-4 h-4 text-warning" />
        <h3 class="text-sm font-semibold text-text-primary">执行计划</h3>
        <span class="text-xs text-text-tertiary">({{ plan.total_steps }} 步)</span>
      </div>
      <div class="flex items-center gap-3">
        <!-- 进度统计 -->
        <div class="flex items-center gap-2 text-xs">
          <span v-if="progress.completed > 0" class="text-success">
            ✓ {{ progress.completed }}
          </span>
          <span v-if="progress.failed > 0" class="text-error">
            ✗ {{ progress.failed }}
          </span>
          <span v-if="progress.skipped > 0" class="text-warning">
            ⊘ {{ progress.skipped }}
          </span>
          <span class="text-text-tertiary">
            {{ progress.percentage }}%
          </span>
        </div>
        <button
          @click="collapsed = !collapsed"
          class="p-1 rounded hover:bg-bg-tertiary transition-colors"
        >
          <ChevronDown
            class="w-4 h-4 text-text-secondary transition-transform"
            :class="{ 'rotate-180': !collapsed }"
          />
        </button>
      </div>
    </div>

    <!-- 进度条 -->
    <div class="mb-3 h-1.5 bg-bg-tertiary rounded-full overflow-hidden">
      <div
        class="h-full transition-all duration-300 ease-out"
        :class="{
          'bg-success': progress.failed === 0,
          'bg-error': progress.failed > 0,
        }"
        :style="{ width: `${progress.percentage}%` }"
      />
    </div>

    <div v-show="!collapsed" class="space-y-2">
      <div
        v-for="(step, index) in plan.steps"
        :key="index"
        class="flex items-start gap-3 p-3 rounded-lg transition-all duration-300"
        :class="{
          'bg-primary/5 border-l-2 border-primary': step.status === 'running',
          'bg-success/5 border-l-2 border-success': step.status === 'completed',
          'bg-error/5 border-l-2 border-error': step.status === 'failed',
          'bg-warning/5 border-l-2 border-warning': step.status === 'skipped',
          'bg-bg-tertiary/20': step.status === 'pending',
        }"
      >
        <!-- 状态图标 -->
        <div class="flex-shrink-0 mt-0.5">
          <div
            v-if="step.status === 'pending'"
            class="w-5 h-5 rounded-full border-2 border-border-primary bg-bg-primary"
          />
          <div
            v-else-if="step.status === 'running'"
            class="w-5 h-5 rounded-full border-2 border-primary bg-primary/20"
          >
            <div class="w-full h-full rounded-full border-2 border-primary border-t-transparent animate-spin" />
          </div>
          <CheckCircle2
            v-else-if="step.status === 'completed'"
            class="w-5 h-5 text-success animate-scale-in"
          />
          <XCircle
            v-else-if="step.status === 'failed'"
            class="w-5 h-5 text-error animate-shake"
          />
          <MinusCircle
            v-else-if="step.status === 'skipped'"
            class="w-5 h-5 text-warning"
          />
        </div>

        <!-- 步骤内容 -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 flex-wrap">
            <span class="text-sm font-medium text-text-primary">步骤 {{ index + 1 }}</span>
            <span
              v-if="step.tool"
              class="text-xs px-2 py-0.5 rounded bg-bg-tertiary text-text-secondary font-mono"
            >
              {{ step.tool }}
            </span>
            <span
              v-if="step.status === 'running'"
              class="text-xs text-primary animate-pulse"
            >
              执行中...
            </span>
          </div>
          <p class="text-sm text-text-secondary mt-1">{{ step.description }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Lightbulb, ChevronDown, CheckCircle2, XCircle, MinusCircle } from 'lucide-vue-next'
import type { AgentPlan } from '@/api/agent'

interface Props {
  plan: AgentPlan
}

const props = defineProps<Props>()

const collapsed = ref(false)

// 计算完成进度
const progress = computed(() => {
  const total = props.plan.steps.length
  const completed = props.plan.steps.filter(s => s.status === 'completed').length
  const failed = props.plan.steps.filter(s => s.status === 'failed').length
  const skipped = props.plan.steps.filter(s => s.status === 'skipped').length
  return {
    total,
    completed,
    failed,
    skipped,
    percentage: total > 0 ? Math.round((completed / total) * 100) : 0,
  }
})
</script>

<style scoped>
@keyframes scale-in {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-4px);
  }
  75% {
    transform: translateX(4px);
  }
}

.animate-scale-in {
  animation: scale-in 0.3s ease-out;
}

.animate-shake {
  animation: shake 0.4s ease-in-out;
}
</style>



