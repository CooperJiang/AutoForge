<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">RSS 多源聚合配置</h3>

    <!-- 订阅源列表 -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <label class="block text-xs font-medium text-text-secondary">
          订阅源列表 <span class="text-error">*</span>
        </label>
        <button
          @click="addSource"
          class="flex items-center gap-1 px-2 py-1 text-xs font-medium text-primary hover:text-primary-hover bg-primary-light hover:bg-primary-light/80 rounded transition-colors"
        >
          <span class="text-lg leading-none">+</span>
          添加订阅源
        </button>
      </div>

      <!-- 订阅源卡片列表 -->
      <div v-if="localConfig.sources.length === 0" class="text-xs text-text-tertiary text-center py-4 border border-dashed border-border-primary rounded-lg">
        暂无订阅源，请点击上方按钮添加
      </div>

      <div v-else class="space-y-2">
        <div
          v-for="(source, index) in localConfig.sources"
          :key="index"
          class="border border-border-primary rounded-lg p-3 bg-bg-hover"
        >
          <div class="flex items-start gap-2">
            <div class="flex-shrink-0 w-6 h-6 rounded-full bg-primary text-white text-xs font-semibold flex items-center justify-center mt-1">
              {{ index + 1 }}
            </div>
            <div class="flex-1 space-y-2">
              <div>
                <label class="block text-xs font-medium text-text-secondary mb-1"> RSS 地址 </label>
                <BaseInput
                  v-model="source.url"
                  placeholder="https://example.com/feed"
                  class="text-xs"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-text-secondary mb-1">
                  关键词过滤（可选）
                </label>
                <BaseInput
                  v-model="source.keywords"
                  placeholder="AI, 科技, 自动化"
                  class="text-xs"
                />
                <p class="mt-0.5 text-xs text-text-tertiary">多个关键词用逗号分隔</p>
              </div>
            </div>
            <button
              @click="removeSource(index)"
              class="flex-shrink-0 w-6 h-6 rounded hover:bg-error-light text-text-tertiary hover:text-error transition-colors flex items-center justify-center"
              title="删除"
            >
              <span class="text-lg leading-none">×</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 全局配置 -->
    <div class="border-t border-border-primary pt-4 space-y-3">
      <h4 class="text-xs font-semibold text-text-primary">全局配置</h4>

      <!-- 最大条目数 -->
      <div>
        <label class="block text-xs font-medium text-text-secondary mb-1.5">
          最大条目数（总计）
        </label>
        <BaseInput
          v-model.number="localConfig.max_items"
          type="number"
          min="1"
          max="200"
          placeholder="20"
        />
        <p class="mt-1 text-xs text-text-tertiary">所有订阅源汇总后的最大文章数量（1-200）</p>
      </div>

      <!-- 时间范围 -->
      <div>
        <label class="block text-xs font-medium text-text-secondary mb-1.5">
          时间范围（小时）
        </label>
        <BaseInput
          v-model.number="localConfig.hours_ago"
          type="number"
          min="0"
          max="720"
          placeholder="0"
        />
        <p class="mt-1 text-xs text-text-tertiary">
          只获取最近 N 小时内的文章，0 表示不限制（最大 720 小时/30 天）
        </p>
      </div>

      <!-- 去重规则 -->
      <div>
        <label class="block text-xs font-medium text-text-secondary mb-1.5"> 去重规则 </label>
        <BaseSelect
          v-model="localConfig.dedup_by"
          :options="[
            { label: '按链接去重', value: 'link' },
            { label: '按标题去重', value: 'title' },
          ]"
        />
        <p class="mt-1 text-xs text-text-tertiary">避免重复文章出现在结果中</p>
      </div>

      <!-- 排序方式 -->
      <div>
        <label class="block text-xs font-medium text-text-secondary mb-1.5"> 排序方式 </label>
        <BaseSelect
          v-model="localConfig.sort_by"
          :options="[
            { label: '按发布时间（最新优先）', value: 'time' },
            { label: '按订阅源顺序', value: 'source' },
          ]"
        />
        <p class="mt-1 text-xs text-text-tertiary">控制文章在列表中的顺序</p>
      </div>
    </div>

    <!-- 可用 RSS 源列表 -->
    <details class="border border-border-primary rounded-md p-3 bg-bg-hover">
      <summary class="text-xs font-medium text-text-secondary cursor-pointer">
        📚 国内可用 RSS 源（点击展开）
      </summary>
      <div class="mt-2 space-y-2">
        <div class="text-xs font-semibold text-text-primary">科技资讯</div>
        <div class="space-y-0.5 text-xs text-text-tertiary">
          <div>• 36氪：https://36kr.com/feed</div>
          <div>• 少数派：https://sspai.com/feed</div>
          <div>• IT之家：https://www.ithome.com/rss</div>
          <div>• 爱范儿：https://www.ifanr.com/feed</div>
        </div>

        <div class="text-xs font-semibold text-text-primary mt-2">开发者</div>
        <div class="space-y-0.5 text-xs text-text-tertiary">
          <div>• 阮一峰：https://www.ruanyifeng.com/blog/atom.xml</div>
          <div>• 掘金前端：https://rsshub.app/juejin/category/frontend</div>
          <div>• 云风博客：https://blog.codingnow.com/atom.xml</div>
        </div>

        <div class="text-xs font-semibold text-text-primary mt-2">财经商业</div>
        <div class="space-y-0.5 text-xs text-text-tertiary">
          <div>• 虎嗅：https://www.huxiu.com/rss/0.xml</div>
          <div>• 钛媒体：https://www.tmtpost.com/rss.xml</div>
        </div>

        <div class="text-xs font-semibold text-text-primary mt-2">综合媒体</div>
        <div class="space-y-0.5 text-xs text-text-tertiary">
          <div>• 新浪科技：http://rss.sina.com.cn/tech/rollnews.xml</div>
          <div>• 澎湃新闻：https://www.thepaper.cn/rss</div>
        </div>

        <div class="mt-2 p-2 bg-warning-light border border-warning rounded text-xs text-warning-text">
          💡 提示：以上源已验证可用，如遇访问问题请检查网络或稍后重试
        </div>
      </div>
    </details>

    <!-- 输出说明 -->
    <div class="bg-info-light border border-info rounded-md p-3">
      <div class="text-xs font-semibold text-info-text mb-1.5">✨ 输出字段（双重结构）</div>

      <div class="text-xs font-semibold text-info-text mb-1 mt-2">📦 合并数据（统一处理）</div>
      <div class="text-xs text-info-text space-y-0.5 ml-2">
        <div>• <code v-pre>{{nodes.xxx.items}}</code> - 所有文章合并数组</div>
        <div>• <code v-pre>{{nodes.xxx.items[0].source}}</code> - 文章来源标注</div>
        <div>• <code v-pre>{{nodes.xxx.total}}</code> - 文章总数</div>
      </div>

      <div class="text-xs font-semibold text-info-text mb-1 mt-2">📂 分组数据（按源区分）</div>
      <div class="text-xs text-info-text space-y-0.5 ml-2">
        <div>• <code v-pre>{{nodes.xxx.sources_with_items}}</code> - 按订阅源分组</div>
        <div>• <code v-pre>{{nodes.xxx.sources_with_items[0].feed_title}}</code> - 源标题</div>
        <div>• <code v-pre>{{nodes.xxx.sources_with_items[0].items}}</code> - 该源的文章</div>
      </div>

      <div class="mt-2 p-2 bg-bg-elevated rounded text-xs text-info-text">
        💡 根据需求选择：想统一处理用 <code v-pre>items</code>，想区分来源用 <code v-pre>sources_with_items</code>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface RSSSource {
  url: string
  keywords: string
}

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:config', value: Record<string, any>): void
}>()

const localConfig = ref({
  sources: (props.config.sources || []) as RSSSource[],
  max_items: props.config.max_items ?? 20,
  hours_ago: props.config.hours_ago ?? 0,
  dedup_by: props.config.dedup_by || 'link',
  sort_by: props.config.sort_by || 'time',
})

// 如果初始没有订阅源，添加一个默认的
if (localConfig.value.sources.length === 0) {
  localConfig.value.sources.push({ url: '', keywords: '' })
}

// 添加订阅源
const addSource = () => {
  localConfig.value.sources.push({ url: '', keywords: '' })
}

// 删除订阅源
const removeSource = (index: number) => {
  if (localConfig.value.sources.length > 1) {
    localConfig.value.sources.splice(index, 1)
  }
}

// 防抖旗标防止递归
const updatingFromProps = ref(false)

// 子改父
watch(
  localConfig,
  (v) => {
    if (!updatingFromProps.value) {
      emit('update:config', { ...v })
    }
  },
  { deep: true }
)

// 父改子
watch(
  () => props.config,
  (cfg) => {
    updatingFromProps.value = true
    localConfig.value = {
      sources: (cfg?.sources || [{ url: '', keywords: '' }]) as RSSSource[],
      max_items: cfg?.max_items ?? 20,
      hours_ago: cfg?.hours_ago ?? 0,
      dedup_by: cfg?.dedup_by || 'link',
      sort_by: cfg?.sort_by || 'time',
    }
    setTimeout(() => {
      updatingFromProps.value = false
    }, 0)
  },
  { deep: true }
)
</script>
