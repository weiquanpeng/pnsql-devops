<template>
  <div class="dashboard-container">
    <!-- 搜索区域卡片 -->
    <t-space size="16px">
      <!-- VM名称选择框 -->
      <t-select
        v-model="searchValue"
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
        v-model="dateRange"
        placeholder="请选择日期时间范围"
        enable-time-picker
        style="width: 380px"
        :presets="timePresets"
        @change="handleDateRangeChange"
      />

      <!-- 查询按钮 -->
      <t-button theme="primary" @click="handleSearch" :loading="isSearching">
        <template #icon><search-icon /></template>
        查询
      </t-button>
    </t-space>

    <!-- 图表区域卡片 -->
    <t-card class="chart-card" title="活跃会话趋势" :collapsible="false">
      <template #content>
        <div ref="chartRef" style="width: 100%; height: 400px;"></div>
      </template>
    </t-card>

    <!-- 会话详情表格 -->
    <t-card class="table-card" title="会话详情" :collapsible="false">
      <template #content>
        <t-table
          row-key="id"
          :data="selectedSessions"
          :columns="columns"
          :sort="sort"
          @sort-change="handleSortChange"
          hover
          stripe
          size="small"
        >
        </t-table>
      </template>
    </t-card>
  </div>
</template>

<script lang="tsx" setup>
import {ref, onMounted, onBeforeUnmount, watch} from 'vue';
import { SearchIcon } from 'tdesign-icons-vue-next';
import { getMysqlVmName, getPgSession } from '@/api/services/DatabaseMetaList';
import * as echarts from 'echarts';
import type { TableSort } from 'tdesign-vue-next';

// 数据定义
const defaultPreset = '最近3小时';
const isSearching = ref(false);
const searchValue = ref('');
const dateRange = ref(['', '']);
const vmNameOptions = ref<Array<{ value: string, label: string }>>([]);
const loading = ref(false);
const chartRef = ref<HTMLElement>();
const activePreset = ref<string | null>(null);
const sessionData = ref<any[]>([]);
const selectedSessions = ref<any[]>([]);
const sort = ref<TableSort>({
  sortBy: 'sql_timing',
  descending: true,
});

let chartInstance: echarts.ECharts | null = null;

// 表格列定义
const columns = ref([
  { colKey: 'id', title: 'ID', width: 80, ellipsis: true },
  { colKey: 'pid', title: 'PID', width: 90, ellipsis: true },
  { colKey: 'datname', title: '数据库', width: 90, ellipsis: true },
  { colKey: 'ip', title: 'ip', width: 90, ellipsis: true },
  { colKey: 'application_name', title: '应用名称', width: 110, ellipsis: true },
  {
    colKey: 'xact_start',
    title: '事务开始时间',
    width: 120,
    cell: (h, { row }) => formatTime(row.xact_start),
    ellipsis: true,
  },
  {
    colKey: 'open_timing',
    title: '打开时间',
    width: 180,
    sorter: true,
    ellipsis: true,
    cell: (h, { row }) => {
      const val = row.open_timing || '';
      // 去掉负号
      return val.startsWith('-') ? val.substring(1) : val;
    },
  },
  {
    colKey: 'sql_timing',
    title: 'SQL时间',
    width: 180,
    sorter: true,
    ellipsis: true,
    cell: (h, { row }) => {
      const val = row.sql_timing || '';
      return val.startsWith('-') ? val.substring(1) : val;
    },
  },
  { colKey: 'state', title: '状态', width: 150, ellipsis: true },
  { colKey: 'wait_event', title: '等待事件', width: 150, ellipsis: true },
  {
    colKey: 'query',
    title: '查询语句',
    ellipsis: true,
  },
]);

// 创建时间快捷选项
const createTimePresets = () => {
  const now = new Date();
  return {
    '最近3小时': () => {
      const start = new Date(now.getTime() - 3 * 60 * 60 * 1000);
      return [start, now];
    },
    '最近6小时': () => {
      const start = new Date(now.getTime() - 6 * 60 * 60 * 1000);
      return [start, now];
    },
    '最近12小时': () => {
      const start = new Date(now.getTime() - 12 * 60 * 60 * 1000);
      return [start, now];
    },
    '最近1天': () => {
      const start = new Date(now.getTime() - 24 * 60 * 60 * 1000);
      return [start, now];
    },
    '最近3天': () => {
      const start = new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000);
      return [start, now];
    },
  };
};

// 时间快捷选项
const timePresets = ref(createTimePresets());

// 处理日期范围变化
const handleDateRangeChange = () => {
  activePreset.value = null;
};

// 格式化时间显示
const formatTime = (timeStr: string) => {
  if (!timeStr) return '';
  const date = new Date(timeStr);
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`;
};

onMounted(async () => {
  loading.value = true;
  try {
    const response = await getMysqlVmName();
    vmNameOptions.value = response.data.map(item => ({ value: item, label: item }));
    initChart();
  } catch (error) {
    console.error('获取 MySQL VM 名称失败:', error);
  } finally {
    loading.value = false;
    // **不需要手动调用watch**，immediate选项会自动触发回调
  }
});

// 监听vmNameOptions变化，设置默认选中第一个选项
watch(vmNameOptions, (newVal) => {
  if (newVal.length > 0) {
    // 选择第一个VM名称
    searchValue.value = newVal[0].value;
    // 设置默认时间范围（最近3小时）
    const [start, end] = timePresets.value['最近3小时']();
    dateRange.value = [start.toISOString(), end.toISOString()];
    // 触发自动搜索
    handleSearch();
  }
}, {
  immediate: true // 初始化时立即执行一次回调
});

const initChart = () => {
  if (!chartRef.value) return;

  chartInstance = echarts.init(chartRef.value);

  // 构造默认1小时内每分钟的刻度
  const now = new Date();
  const defaultXAxisData: string[] = [];
  for (let i = 59; i >= 0; i--) {
    const t = new Date(now.getTime() - i * 60 * 1000);
    defaultXAxisData.push(t.toISOString());
  }

  const defaultSeriesData = new Array(60).fill(0);

  // 记录上一个显示的日期（月份-日期）
  let lastDate = '';

  const option = {
    title: {
      text: '活跃会话趋势',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const date = formatTime(params[0].name);
        return `${date}<br/>会话数: ${params[0].value}`;
      }
    },
    xAxis: {
      type: 'category',
      data: defaultXAxisData,
      axisLabel: {
        formatter: (value: string) => {
          const date = new Date(value);
          const month = (date.getMonth() + 1).toString().padStart(2, '0');
          const day = date.getDate().toString().padStart(2, '0');
          const hour = date.getHours().toString().padStart(2, '0');
          const minute = date.getMinutes().toString().padStart(2, '0');

          const currentDate = `${month}-${day}`;
          let label = '';

          if (currentDate !== lastDate) {
            // 当天第一次出现，显示月-日 + 时间
            label = `${currentDate} ${hour}:${minute}`;
            lastDate = currentDate;
          } else {
            // 之后只显示时间
            label = `${hour}:${minute}`;
          }

          return label;
        },
        rotate: 45 // 旋转45度避免重叠
      }
    },
    yAxis: {
      type: 'value',
      name: '会话数',
      max: 50,
    },
    series: [
      {
        name: '活跃会话',
        type: 'bar',
        data: defaultSeriesData,
        itemStyle: {
          color: '#3a4de9'
        },
        emphasis: {
          itemStyle: {
            color: '#1a2bb9'
          }
        }
      }
    ],
    grid: {
      left: 40,
      right: 20,
      top: 30,
      bottom: 100,  // 关键！为 dataZoom 腾空间
    },
    dataZoom: [
      {
        type: 'slider',
        show: true,
        bottom: 5,
        height: 26,
        start: 0,
        end: 100,
        backgroundColor: '#f9fafb',
        fillerColor: 'rgba(58, 77, 233, 0.2)',
        borderColor: '#d1d5db',
        handleStyle: {
          color: '#3a4de9',
          borderColor: '#fff',
          borderWidth: 1,
          shadowBlur: 4,
          shadowColor: 'rgba(0, 0, 0, 0.2)',
          shadowOffsetX: 0,
          shadowOffsetY: 2
        },
        handleSize: 14,
        showDetail: false,
        showDataShadow: false,
        throttle: 100,
      },
      {
        type: 'inside',
        zoomOnMouseWheel: true,
        moveOnMouseWheel: true,
      }
    ],
    brush: {
      toolbox: ['zoom', 'reset']
    }
  };

  chartInstance.setOption(option);

  // 添加点击事件
  chartInstance.on('click', (params: any) => {
    const clickedTime = params.name;
    const sessions = sessionData.value.filter(session => session.captured_at === clickedTime);
    selectedSessions.value = [...sessions];

    // 默认按sql_timing降序排序
    handleSortChange({
      sortBy: 'sql_timing',
      descending: true
    });
  });

  // 响应式调整
  window.addEventListener('resize', resizeChart);
};


const updateChartData = (data: any[]) => {
  if (!chartInstance) return;

  const markLineConfig = {
    silent: true,
    lineStyle: {
      type: 'dashed',
      color: '#f5222d'  // 红色告警线
    },
    data: [{ yAxis: 30 }]
  };

  if (data.length === 0) {
    const now = new Date();
    const times: string[] = [];
    for (let i = 59; i >= 0; i--) {
      const t = new Date(now.getTime() - i * 60 * 1000);
      times.push(t.toISOString());
    }

    chartInstance.setOption({
      xAxis: {
        data: times,
        // ✅ 保留 formatter，不要覆盖
      },
      series: [{
        type: 'bar',
        data: new Array(times.length).fill(0),
        markLine: markLineConfig  // ✅ 加上告警线
      }],
      dataZoom: [{ start: 0, end: 100 }]
    });

    sessionData.value = [];
    return;
  }

  const timeMap = new Map();
  data.forEach(session => {
    const time = session.captured_at;
    timeMap.set(time, (timeMap.get(time) || 0) + 1);
  });

  const sortedTimes = Array.from(timeMap.keys()).sort();
  const counts = sortedTimes.map(time => timeMap.get(time));

  chartInstance.setOption({
    xAxis: {
      data: sortedTimes,
      // ✅ 不改 axisLabel.formatter，保留初始化的配置
    },
    series: [{
      type: 'bar',
      data: counts,
      markLine: markLineConfig  // ✅ 添加告警线
    }],
    dataZoom: [{ start: 0, end: 100 }]
  });

  sessionData.value = data;
};



// 图表响应式调整
const resizeChart = () => {
  chartInstance?.resize();
};

// 排序处理
const handleSortChange = (sortVal: TableSort) => {
  sort.value = sortVal;

  if (sortVal.sortBy) {
    selectedSessions.value.sort((a, b) => {
      const aVal = a[sortVal.sortBy as string];
      const bVal = b[sortVal.sortBy as string];

      // 处理时间格式的排序
      if (sortVal.sortBy === 'open_timing' || sortVal.sortBy === 'sql_timing') {
        const aTime = convertTimeToSeconds(aVal);
        const bTime = convertTimeToSeconds(bVal);
        return sortVal.descending ? bTime - aTime : aTime - bTime;
      }

      // 默认排序
      if (aVal > bVal) return sortVal.descending ? -1 : 1;
      if (aVal < bVal) return sortVal.descending ? 1 : -1;
      return 0;
    });
  }
};

// 将时间字符串转换为秒数用于排序
const convertTimeToSeconds = (timeStr: string) => {
  if (!timeStr) return 0;
  const parts = timeStr.split(':');
  if (parts.length !== 3) return 0;

  const hours = parseFloat(parts[0]);
  const minutes = parseFloat(parts[1]);
  const seconds = parseFloat(parts[2]);

  return hours * 3600 + minutes * 60 + seconds;
};

// 查询方法
const handleSearch = async () => {
  if (!searchValue.value || !dateRange.value[0] || !dateRange.value[1]) {
    console.warn('请选择VM名称和时间范围');
    return;
  }

  const params = {
    vmname: searchValue.value,
    startTime: formatDateTime(dateRange.value[0]),
    stopTime: formatDateTime(dateRange.value[1])
  };

  console.log('请求参数:', params);

  isSearching.value = true; // 开始 loading

  try {
    const response = await getPgSession(params);
    console.log('接口返回数据:', response.data);

    // 更新图表数据
    updateChartData(response.data);
    // 清空之前选择的会话
    selectedSessions.value = [];

  } catch (error) {
    console.error('查询失败:', error);
  } finally {
    isSearching.value = false; // 结束 loading
  }
};


// 日期格式化
const formatDateTime = (date: string | Date) => {
  if (!date) return '';
  const d = new Date(date);
  return `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`;
};

// 组件卸载前清理
onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart);
  chartInstance?.dispose();
});
</script>

<style scoped>
.dashboard-container {
  padding: 16px;
}

.search-card {
  margin-bottom: 16px;
}

.chart-card {
  margin-top: 16px;
}

.table-card {
  margin-top: 16px;
}

.query-cell {
  max-width: 300px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.query-cell pre {
  margin: 0;
  font-family: inherit;
  white-space: pre-wrap;
}

/* 调整时间选择器快捷选项样式 */
:deep(.t-date-range-picker__presets) {
  padding: 8px 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.t-date-range-picker__preset-item) {
  margin: 0;
  padding: 4px 8px;
  font-size: 12px;
  line-height: 1.5;
  border-radius: 2px;
}

/* 调整卡片样式 */
:deep(.t-card) {
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

:deep(.t-card__header) {
  padding: 12px 20px;
  background-color: #f7f8fa;
  border-bottom: 1px solid #e5e6eb;
}

:deep(.t-card__content) {
  padding: 16px 20px;
}
</style>
