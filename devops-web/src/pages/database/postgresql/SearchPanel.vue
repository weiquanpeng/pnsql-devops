<template>
  <t-space size="16px">
    <!-- VM名称选择框 -->
    <t-select
      v-model="modelValue"
      clearable
      filterable
      placeholder="请选择VM名称"
      style="width: 200px"
      :loading="loading"
    >
      <t-option
        v-for="item in vmNameOptions"
        :key="item.value"
        :value="item.value"
        :label="item.label"
      />
    </t-select>

    <!-- 时间范围选择器 -->
    <t-date-range-picker
      v-model="dateRangeProxy"
      placeholder="请选择日期时间范围"
      enable-time-picker
      style="width: 380px"
      :presets="timePresets"
      @change="handleDateRangeChange"
    />

    <!-- 查询按钮 -->
    <t-button theme="primary" @click="$emit('search')">
      <template #icon><search-icon /></template>
      查询
    </t-button>
  </t-space>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, computed } from 'vue';
import { SearchIcon } from 'tdesign-icons-vue-next';
import { getMysqlVmName } from '@/api/services/DatabaseMetaList';

const props = defineProps({
  modelValue: String,
  dateRange: Array
});
const emit = defineEmits(['update:modelValue', 'update:dateRange', 'search']);

const loading = ref(false);
const vmNameOptions = ref<Array<{ value: string, label: string }>>([]);

const timePresets = ref({
  '最近1小时': () => [new Date(Date.now() - 3600 * 1000), new Date()],
  '最近6小时': () => [new Date(Date.now() - 6 * 3600 * 1000), new Date()],
  '最近12小时': () => [new Date(Date.now() - 12 * 3600 * 1000), new Date()],
  '最近1天': () => [new Date(Date.now() - 24 * 3600 * 1000), new Date()],
  '最近3天': () => [new Date(Date.now() - 3 * 24 * 3600 * 1000), new Date()]
});

// 双向绑定代理
const dateRangeProxy = computed({
  get: () => props.dateRange,
  set: (val) => emit('update:dateRange', val)
});

onMounted(async () => {
  loading.value = true;
  try {
    const res = await getMysqlVmName();
    vmNameOptions.value = res.data.map(item => ({ value: item, label: item }));
    if (!props.modelValue && res.data.length > 0) {
      emit('update:modelValue', res.data[0]);
      const [start, end] = timePresets.value['最近1小时']();
      emit('update:dateRange', [start.toISOString(), end.toISOString()]);
    }
  } catch (e) {
    console.error('获取VM名称失败', e);
  } finally {
    loading.value = false;
  }
});

const handleDateRangeChange = () => {
  // 清除 preset 选中状态等操作如有
};

watch(() => props.modelValue, val => {
  if (!val && vmNameOptions.value.length > 0) {
    emit('update:modelValue', vmNameOptions.value[0].value);
  }
});
</script>

<style scoped>
.t-space {
  margin-bottom: 16px;
}
</style>
