<template>
  <div>
    <el-space wrap>
      <el-select v-model="year" style="width: 120px" @change="load">
        <el-option v-for="y in years" :key="y" :label="y + ' 年'" :value="y" />
      </el-select>
      <el-button type="primary" @click="load">刷新</el-button>
    </el-space>
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div class="metric">{{ s.completed_drills_year ?? '-' }}</div>
          <div class="label">本年已完成演练次数</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div class="metric">{{ s.plans_met_annual_minimum ?? 0 }} / {{ s.plan_total ?? 0 }}</div>
          <div class="label">预案满足每年至少 1 次演练</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div class="metric">{{ (s.dept_coverage_percent ?? 0).toFixed(1) }}%</div>
          <div class="label">演练覆盖部门比例（本年 / 基准）</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div class="metric">{{ (s.rectification_done_percent ?? 0).toFixed(1) }}%</div>
          <div class="label">整改完成率（{{ s.rectification_done ?? 0 }}/{{ s.rectification_total ?? 0 }}）</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="8">
        <el-card shadow="hover">
          <div class="metric">{{ s.dept_distinct_year ?? 0 }}</div>
          <div class="label">本年参演组织部门种类数（去重）</div>
        </el-card>
      </el-col>
    </el-row>
    <el-alert
      v-if="(s.plan_ids_below_minimum?.length ?? 0) > 0"
      type="warning"
      title="以下预案本年尚未完成演练，请关注"
      style="margin-top: 16px"
      show-icon
    >
      预案 ID：{{ s.plan_ids_below_minimum.join('、') }}
    </el-alert>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import * as api from '../api'

const y0 = new Date().getFullYear()
const years = computed(() => [y0 + 1, y0, y0 - 1, y0 - 2])
const year = ref(y0)
const s = ref({})

async function load() {
  const { data } = await api.fetchCompliance(year.value)
  s.value = data
}

onMounted(load)
</script>

<style scoped>
.metric {
  font-size: 26px;
  font-weight: 700;
}
.label {
  color: var(--el-text-color-secondary);
  margin-top: 6px;
  font-size: 13px;
}
</style>
