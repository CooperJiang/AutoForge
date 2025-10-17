<template>
  <div class="space-y-4">
    <!-- 顶部搜索栏 -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex gap-2 items-center">
        <BaseSelect
          v-model="filterCategory"
          :options="categoryFilterOptions"
          @update:modelValue="loadTools"
          style="width: 200px"
        />
        <BaseSelect
          v-model="filterStatus"
          :options="statusFilterOptions"
          @update:modelValue="loadTools"
          style="width: 200px"
        />
        <BaseInput
          v-model="searchKeyword"
          placeholder="搜索工具..."
          style="width: 260px"
          @keyup.enter="loadTools"
        />
        <BaseButton @click="loadTools" variant="primary">
          搜索
        </BaseButton>
      </div>
      <BaseButton @click="syncToolsFromBackend" variant="secondary">
        <RefreshCw class="w-4 h-4 mr-1.5" />
        同步工具
      </BaseButton>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-16">
      <div class="inline-block w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mb-4"></div>
      <p class="text-text-secondary">加载工具列表中...</p>
    </div>

    <!-- 空状态 -->
    <div v-else-if="filteredTools.length === 0" class="text-center py-16">
      <div class="text-6xl mb-4">🔧</div>
      <p class="text-text-primary font-medium mb-2">暂无工具</p>
      <p class="text-text-secondary text-sm">请点击右上角"同步工具"按钮从后端同步工具定义</p>
    </div>

    <!-- 工具卡片网格 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
      <div
        v-for="tool in paginatedTools"
        :key="tool.id"
        class="border-2 rounded-xl overflow-hidden hover:shadow-xl transition-all shadow-md relative"
        :style="{
          backgroundColor: 'var(--color-bg-elevated)',
          borderColor: tool.enabled ? 'var(--color-border-primary)' : 'var(--color-border-secondary)',
          opacity: tool.enabled ? '1' : '0.6',
        }"
        @mouseenter="(e) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--color-primary)'"
        @mouseleave="(e) => (e.currentTarget as HTMLElement).style.borderColor = tool.enabled ? 'var(--color-border-primary)' : 'var(--color-border-secondary)'"
      >
        <!-- 工具图标区域 -->
        <div class="p-4 flex items-center gap-3 border-b-2 border-border-primary">
          <div
            class="flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center shadow-lg"
            :class="getToolIconBg(tool.tool_code)"
          >
            <component
              :is="getToolIcon(tool.tool_code)"
              class="w-6 h-6 text-white"
            />
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-base font-bold text-text-primary truncate mb-0.5">
              {{ tool.tool_name }}
            </h3>
            <p class="text-xs text-text-tertiary truncate">
              {{ tool.tool_code }}
            </p>
          </div>
          <div class="flex items-center gap-1.5 flex-shrink-0">
            <button
              @click.stop="configTool(tool)"
              class="p-1.5 rounded-lg transition-all hover:bg-info-light"
              :style="{ color: 'var(--color-info)' }"
              title="配置"
            >
              <Settings class="w-4 h-4" />
            </button>
            <button
              v-if="!tool.is_deprecated"
              @click.stop="toggleToolEnabled(tool)"
              class="p-1.5 rounded-lg transition-all"
              :class="tool.enabled ? 'hover:bg-warning-light' : 'hover:bg-success-light'"
              :style="{ color: tool.enabled ? 'var(--color-warning)' : 'var(--color-success)' }"
              :title="tool.enabled ? '禁用' : '启用'"
            >
              <Power class="w-4 h-4" />
            </button>
            <button
              v-if="tool.is_deprecated"
              @click.stop="deleteToolConfirm(tool)"
              class="p-1.5 rounded-lg transition-all hover:bg-error-light"
              :style="{ color: 'var(--color-error)' }"
              title="删除"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- 工具信息 -->
        <div class="p-4">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-1.5 flex-wrap">
              <span
                class="inline-block px-2 py-0.5 rounded-md text-xs font-medium border"
                :style="{
                  backgroundColor: 'var(--color-primary-light)',
                  color: 'var(--color-primary)',
                  borderColor: 'var(--color-primary)',
                  opacity: '0.8',
                }"
              >
                {{ getCategoryName(tool.category) }}
              </span>
              <span
                v-if="tool.version"
                class="text-xs text-text-tertiary"
              >
                v{{ tool.version }}
              </span>
            </div>
            <div class="flex items-center gap-1 flex-shrink-0">
              <span
                v-if="tool.is_deprecated"
                class="text-xs px-2 py-0.5 rounded-md font-medium"
                :style="{
                  backgroundColor: 'var(--color-error-light)',
                  color: 'var(--color-error-text)',
                }"
                title="工具已废弃"
              >
                已废弃
              </span>
              <span
                v-else-if="!tool.enabled"
                class="text-xs px-2 py-0.5 rounded-md font-medium"
                :style="{
                  backgroundColor: 'var(--color-bg-tertiary)',
                  color: 'var(--color-text-secondary)',
                }"
                title="工具已禁用"
              >
                已禁用
              </span>
              <span
                v-else-if="!tool.visible"
                class="text-xs px-2 py-0.5 rounded-md font-medium"
                :style="{
                  backgroundColor: 'var(--color-warning-light)',
                  color: 'var(--color-warning-text)',
                }"
                title="工具已隐藏"
              >
                已隐藏
              </span>
              <span
                v-else
                class="text-xs px-2 py-0.5 rounded-md font-medium"
                :style="{
                  backgroundColor: 'var(--color-success-light)',
                  color: 'var(--color-success-text)',
                }"
                title="工具已启用"
              >
                已启用
              </span>
            </div>
          </div>

          <p class="text-sm line-clamp-3 leading-relaxed text-text-secondary">
            {{ tool.description || '暂无描述' }}
          </p>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <Pagination
      v-if="filteredTools.length > 0"
      :current="currentPage"
      :page-size="pageSize"
      :total="filteredTools.length"
      :bordered="false"
      @change="handlePageChange"
    />

    <!-- 配置弹窗 -->
    <Dialog
      v-model="showConfigDialog"
      :title="`配置工具 - ${editingTool?.tool_name || ''}`"
      @confirm="saveToolConfig"
      @cancel="cancelConfig"
      confirm-text="保存"
      cancel-text="取消"
      max-width="max-w-2xl"
    >
      <div class="space-y-4" v-if="editingTool">
        <!-- 基本设置 -->
        <div class="border-2 border-border-primary rounded-lg p-4">
          <h3 class="text-sm font-semibold text-text-primary mb-3">基本设置</h3>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <label class="text-sm text-text-primary">启用状态</label>
              <BaseCheckbox v-model="configForm.enabled" :disabled="editingTool.is_deprecated" />
            </div>
            <div class="flex items-center justify-between">
              <label class="text-sm text-text-primary">对外可见</label>
              <BaseCheckbox v-model="configForm.visible" />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">排序</label>
              <BaseInput
                v-model.number="configForm.sort_order"
                type="number"
                placeholder="数字越小越靠前"
              />
            </div>
          </div>
        </div>

        <!-- 后台配置 -->
        <div class="border-2 border-border-primary rounded-lg p-4" v-if="needsBackendConfig(editingTool.tool_code)">
          <h3 class="text-sm font-semibold text-text-primary mb-3">后台配置</h3>
          <p class="text-xs text-text-secondary mb-3">这些配置用于工具的运行环境，配置后工具才能正常使用</p>

          <!-- 阿里云 OSS -->
          <div v-if="editingTool.tool_code === 'aliyun_oss'" class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Endpoint <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.endpoint"
                placeholder="例如: oss-cn-hangzhou.aliyuncs.com"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Access Key ID <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.access_key_id"
                type="password"
                placeholder="请输入 Access Key ID"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Access Key Secret <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.access_key_secret"
                type="password"
                placeholder="请输入 Access Key Secret"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Bucket <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.bucket"
                placeholder="请输入存储桶名称"
              />
            </div>
          </div>

          <!-- 腾讯云 COS -->
          <div v-else-if="editingTool.tool_code === 'tencent_cos'" class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Secret ID <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.secret_id"
                type="password"
                placeholder="请输入 Secret ID"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Secret Key <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.secret_key"
                type="password"
                placeholder="请输入 Secret Key"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Bucket <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.bucket"
                placeholder="例如: mybucket-1234567890"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Region <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.region"
                placeholder="例如: ap-guangzhou"
              />
            </div>
          </div>

          <!-- PixelPunk 图床 -->
          <div v-else-if="editingTool.tool_code === 'pixelpunk_upload'" class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                Base URL <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.base_url"
                placeholder="例如: https://api.pixelpunk.io"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                API Key <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.api_key"
                type="password"
                placeholder="请输入 API Key"
              />
            </div>
          </div>

          <!-- OpenAI Chat -->
          <div v-else-if="editingTool.tool_code === 'openai_chatgpt' || editingTool.tool_code === 'openai_image'" class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                API Key <span class="text-error">*</span>
              </label>
              <BaseInput
                v-model="configForm.config.api_key"
                type="password"
                placeholder="请输入 OpenAI API Key"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                API Base
              </label>
              <BaseInput
                v-model="configForm.config.api_base"
                placeholder="默认: https://api.openai.com"
              />
            </div>
          </div>

          <!-- 其他工具提示 -->
          <div v-else class="text-sm text-text-secondary">
            此工具暂无需要配置的后台参数
          </div>
        </div>

        <!-- 工具信息 -->
        <div class="border-2 border-border-primary rounded-lg p-4">
          <h3 class="text-sm font-semibold text-text-primary mb-3">工具信息</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-text-secondary">工具代码</span>
              <span class="text-text-primary font-mono">{{ editingTool.tool_code }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-text-secondary">分类</span>
              <span class="text-text-primary">{{ getCategoryName(editingTool.category) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-text-secondary">版本</span>
              <span class="text-text-primary">{{ editingTool.version || 'N/A' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-text-secondary">作者</span>
              <span class="text-text-primary">{{ editingTool.author || 'N/A' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-text-secondary">最后同步</span>
              <span class="text-text-primary">{{ formatDate(editingTool.last_sync_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </Dialog>

    <!-- 删除确认弹窗 -->
    <Dialog
      v-model="showDeleteDialog"
      title="确认删除"
      :message="`确定要删除工具配置 &quot;${toolToDelete?.tool_name}&quot; 吗？此操作不可恢复！`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="deleteTool"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RefreshCw, Settings, Power, Trash2 } from 'lucide-vue-next'
import * as toolConfigApi from '@/api/toolConfig'
import type { ToolConfig, ToolConfigDetail } from '@/api/toolConfig'
import { message } from '@/utils/message'
import BaseInput from '@/components/BaseInput'
import BaseButton from '@/components/BaseButton'
import BaseSelect from '@/components/BaseSelect'
import Pagination from '@/components/Pagination'
import Dialog from '@/components/Dialog'
import BaseCheckbox from '@/components/BaseCheckbox/index.vue'
import { getToolIcon, getToolIconBg } from '@/config/tools'

const loading = ref(true)
const tools = ref<ToolConfig[]>([])
const categories = ref<toolConfigApi.ToolCategory[]>([])
const searchKeyword = ref('')
const filterCategory = ref('all')
const filterStatus = ref('all')
const currentPage = ref(1)
const pageSize = ref(12)

const showConfigDialog = ref(false)
const showDeleteDialog = ref(false)
const editingTool = ref<ToolConfigDetail | null>(null)
const toolToDelete = ref<ToolConfig | null>(null)

const configForm = ref<{
  enabled: boolean
  visible: boolean
  sort_order: number
  config: Record<string, any>
}>({
  enabled: false,
  visible: true,
  sort_order: 0,
  config: {},
})

// 判断工具是否需要后台配置
const needsBackendConfig = (toolCode: string) => {
  const configTools = [
    'aliyun_oss',
    'tencent_cos',
    'pixelpunk_upload',
    'openai_chatgpt',
    'openai_image',
  ]
  return configTools.includes(toolCode)
}

// 分类选项（动态生成）
const categoryFilterOptions = computed(() => {
  const options = [{ label: '全部分类', value: 'all' }]
  categories.value.forEach((cat) => {
    options.push({
      label: cat.name,
      value: cat.code,
    })
  })
  return options
})

// 状态选项
const statusFilterOptions = [
  { label: '全部状态', value: 'all' },
  { label: '已启用', value: 'enabled' },
  { label: '已禁用', value: 'disabled' },
  { label: '已废弃', value: 'deprecated' },
]

// 过滤后的工具
const filteredTools = computed(() => {
  let result = tools.value

  // 按分类过滤
  if (filterCategory.value !== 'all') {
    result = result.filter((t) => t.category === filterCategory.value)
  }

  // 按状态过滤
  if (filterStatus.value === 'enabled') {
    result = result.filter((t) => t.enabled && !t.is_deprecated)
  } else if (filterStatus.value === 'disabled') {
    result = result.filter((t) => !t.enabled && !t.is_deprecated)
  } else if (filterStatus.value === 'deprecated') {
    result = result.filter((t) => t.is_deprecated)
  }

  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(
      (t) =>
        t.tool_name.toLowerCase().includes(keyword) ||
        t.tool_code.toLowerCase().includes(keyword) ||
        t.description?.toLowerCase().includes(keyword)
    )
  }

  // 排序：sort_order 升序，然后按名称
  return result.sort((a, b) => {
    if (a.sort_order !== b.sort_order) {
      return a.sort_order - b.sort_order
    }
    return a.tool_name.localeCompare(b.tool_name)
  })
})

// 分页后的工具
const paginatedTools = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTools.value.slice(start, end)
})

// 加载工具列表
const loadTools = async () => {
  loading.value = true
  try {
    const res = await toolConfigApi.getAllToolConfigs()
    tools.value = res.data
  } catch (error: any) {
    message.error(error.message || '获取工具列表失败')
  } finally {
    loading.value = false
  }
}

// 同步工具
const syncToolsFromBackend = async () => {
  try {
    await toolConfigApi.syncTools()
    message.success('同步成功')
    await loadTools()
  } catch (error: any) {
    message.error(error.message || '同步失败')
  }
}

// 配置工具
const configTool = async (tool: ToolConfig) => {
  try {
    const res = await toolConfigApi.getToolConfigDetail(tool.tool_code)
    editingTool.value = res.data
    configForm.value = {
      enabled: tool.enabled,
      visible: tool.visible,
      sort_order: tool.sort_order,
      config: res.data.decrypted_config || {},
    }
    showConfigDialog.value = true
  } catch (error: any) {
    message.error(error.message || '获取工具详情失败')
  }
}

// 保存工具配置
const saveToolConfig = async () => {
  if (!editingTool.value) return

  try {
    // 更新配置
    await toolConfigApi.updateToolConfig(editingTool.value.tool_code, {
      config: configForm.value.config,
    })

    // 更新设置
    await toolConfigApi.updateToolSettings(editingTool.value.tool_code, {
      enabled: configForm.value.enabled,
      visible: configForm.value.visible,
      sort_order: configForm.value.sort_order,
    })

    message.success('保存成功')
    showConfigDialog.value = false
    await loadTools()
  } catch (error: any) {
    message.error(error.message || '保存失败')
  }
}

// 取消配置
const cancelConfig = () => {
  showConfigDialog.value = false
  editingTool.value = null
}

// 切换启用状态
const toggleToolEnabled = async (tool: ToolConfig) => {
  try {
    await toolConfigApi.updateToolSettings(tool.tool_code, {
      enabled: !tool.enabled,
      visible: tool.visible,
      sort_order: tool.sort_order,
    })
    message.success(tool.enabled ? '已禁用' : '已启用')
    await loadTools()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}

// 删除工具确认
const deleteToolConfirm = (tool: ToolConfig) => {
  toolToDelete.value = tool
  showDeleteDialog.value = true
}

// 删除工具
const deleteTool = async () => {
  if (!toolToDelete.value) return

  try {
    await toolConfigApi.deleteToolConfig(toolToDelete.value.id)
    message.success('删除成功')
    showDeleteDialog.value = false
    await loadTools()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

// 分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
}

// 工具辅助函数
const getCategoryName = (category: string) => {
  const cat = categories.value.find((c) => c.code === category)
  return cat ? cat.name : category
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await toolConfigApi.getToolCategories()
    categories.value = res.data
  } catch (error: any) {
    console.error('加载分类列表失败:', error)
  }
}

onMounted(async () => {
  await loadCategories()
  loadTools()
})
</script>
