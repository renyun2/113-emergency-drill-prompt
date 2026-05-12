<template>
  <div>
    <el-space wrap>
      <el-select v-model="filter.plan_id" clearable placeholder="筛选预案" style="width: 200px" @change="load">
        <el-option v-for="p in plans" :key="p.id" :label="p.name" :value="p.id" />
      </el-select>
      <el-select v-model="filter.status" clearable placeholder="状态" style="width: 120px" @change="load">
        <el-option label="计划中" value="planned" />
        <el-option label="已完成" value="completed" />
      </el-select>
      <el-button type="primary" @click="openCreate">新建演练计划</el-button>
      <el-button @click="load">刷新</el-button>
    </el-space>

    <el-table :data="rows" style="width: 100%; margin-top: 12px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="64" />
      <el-table-column label="预案" min-width="140">
        <template #default="{ row }">{{ row.plan?.name || ('#' + row.emergency_plan_id) }}</template>
      </el-table-column>
      <el-table-column prop="drill_kind" label="演练类型" width="110" />
      <el-table-column label="计划日期" width="120">
        <template #default="{ row }">{{ fmtDate(row.scheduled_date) }}</template>
      </el-table-column>
      <el-table-column prop="location" label="地点" width="120" show-overflow-tooltip />
      <el-table-column prop="org_dept" label="组织部门" width="130" />
      <el-table-column prop="notify_depts" label="通知部门" min-width="120" show-overflow-tooltip />
      <el-table-column prop="drill_status" label="状态" width="90" />
      <el-table-column label="评估" width="80">
        <template #default="{ row }">{{ row.evaluation || '—' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openRecord(row)">记录</el-button>
          <el-button link @click="openEditPlan(row)">改计划</el-button>
          <el-button link type="success" @click="openIssues(row)" v-if="row.drill_status === 'completed'">
            问题
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dlg0" title="演练计划" width="560px">
      <el-form label-width="110px">
        <el-form-item label="关联预案">
          <el-select v-model="form0.emergency_plan_id" style="width: 100%">
            <el-option v-for="p in plans" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="演练类型">
          <el-select v-model="form0.drill_kind" style="width: 100%">
            <el-option v-for="k in drillKinds" :key="k" :label="k" :value="k" />
          </el-select>
        </el-form-item>
        <el-form-item label="计划日期">
          <el-date-picker v-model="form0.scheduled_date" value-format="YYYY-MM-DD" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item label="地点"><el-input v-model="form0.location" /></el-form-item>
        <el-form-item label="组织部门"><el-input v-model="form0.org_dept" /></el-form-item>
        <el-form-item label="参演范围"><el-input v-model="form0.participant_scope" /></el-form-item>
        <el-form-item label="演练目标"><el-input v-model="form0.objectives" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="提前通知"><el-input v-model="form0.notify_depts" placeholder="相关部门，逗号分隔" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg0 = false">取消</el-button>
        <el-button type="primary" @click="savePlanDrill">{{ form0.id ? '保存' : '创建' }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dlg1" title="演练完成记录" width="640px">
      <el-form label-width="120px">
        <el-form-item label="状态">
          <el-radio-group v-model="rec.drill_status">
            <el-radio label="planned">计划中</el-radio>
            <el-radio label="completed">已完成</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="实际参演人数"><el-input-number v-model="rec.actual_participants" :min="0" /></el-form-item>
        <el-form-item label="时长(分钟)"><el-input-number v-model="rec.duration_minutes" :min="0" /></el-form-item>
        <el-form-item label="过程描述"><el-input v-model="rec.process_description" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="发现问题清单"><el-input v-model="rec.problem_list" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="评估结论">
          <el-select v-model="rec.evaluation" clearable placeholder="选择" style="width: 200px">
            <el-option v-for="e in evalOpts" :key="e" :label="e" :value="e" />
          </el-select>
        </el-form-item>
        <el-form-item label="照片路径"><el-input v-model="rec.photo_paths" placeholder="可多路径逗号分隔" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg1 = false">关闭</el-button>
        <el-button type="primary" @click="saveRecord">保存记录</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dlgIssues" title="发现问题（可逐项生成整改）" width="720px">
      <el-table :data="issueRows">
        <el-table-column prop="id" label="#" width="64" />
        <el-table-column prop="description" label="描述" />
      </el-table>
      <el-divider />
      <el-form label-width="90px">
        <el-form-item label="新问题">
          <el-input v-model="newIssue" type="textarea" :rows="2" />
        </el-form-item>
        <el-button type="primary" @click="addIssue">录入问题</el-button>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import * as api from '../api'

const drillKinds = ['桌面推演', '实战演练', '联合演练']
const evalOpts = ['优秀', '良好', '一般']

const plans = ref([])
const rows = ref([])
const loading = ref(false)
const filter = reactive({ plan_id: null, status: null })

const dlg0 = ref(false)
const form0 = ref({
  emergency_plan_id: undefined,
  drill_kind: '桌面推演',
  scheduled_date: dayjs().format('YYYY-MM-DD'),
  location: '',
  org_dept: '',
  participant_scope: '',
  objectives: '',
  notify_depts: '',
})

const dlg1 = ref(false)
const rec = ref({ id: null, drill_status: 'completed' })

const dlgIssues = ref(false)
const issueDrillId = ref(null)
const issueRows = ref([])
const newIssue = ref('')

function fmtDate(d) {
  return d ? dayjs(d).format('YYYY-MM-DD') : ''
}

async function loadPlans() {
  const { data } = await api.fetchPlans()
  plans.value = data
}

async function load() {
  loading.value = true
  try {
    const params = {}
    if (filter.plan_id) params.plan_id = filter.plan_id
    if (filter.status) params.status = filter.status
    const { data } = await api.fetchDrills(params)
    rows.value = data
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form0.value = {
    id: undefined,
    emergency_plan_id: plans.value[0]?.id,
    drill_kind: '桌面推演',
    scheduled_date: dayjs().format('YYYY-MM-DD'),
    location: '',
    org_dept: '',
    participant_scope: '',
    objectives: '',
    notify_depts: '',
  }
  dlg0.value = true
}

function openEditPlan(row) {
  form0.value = {
    id: row.id,
    emergency_plan_id: row.emergency_plan_id,
    drill_kind: row.drill_kind,
    scheduled_date: fmtDate(row.scheduled_date),
    location: row.location,
    org_dept: row.org_dept,
    participant_scope: row.participant_scope,
    objectives: row.objectives,
    notify_depts: row.notify_depts,
  }
  dlg0.value = true
}

async function savePlanDrill() {
  const f = form0.value
  if (!f.emergency_plan_id) return ElMessage.warning('请选择预案')
  if (!f.drill_kind) return ElMessage.warning('请选择类型')
  if (!f.scheduled_date) return ElMessage.warning('请选择日期')

  const body = {
    emergency_plan_id: f.emergency_plan_id,
    drill_kind: f.drill_kind,
    scheduled_date: f.scheduled_date,
    location: f.location || '',
    org_dept: f.org_dept || '',
    participant_scope: f.participant_scope || '',
    objectives: f.objectives || '',
    notify_depts: f.notify_depts || '',
  }
  if (f.id) {
    await api.patchDrill(f.id, body)
    ElMessage.success('已更新计划')
  } else {
    await api.createDrill(body)
    ElMessage.success('已创建演练计划')
  }
  dlg0.value = false
  load()
}

function openRecord(row) {
  rec.value = {
    id: row.id,
    drill_status: row.drill_status,
    actual_participants: row.actual_participants ?? undefined,
    duration_minutes: row.duration_minutes ?? undefined,
    process_description: row.process_description || '',
    problem_list: row.problem_list || '',
    evaluation: row.evaluation || '',
    photo_paths: row.photo_paths || '',
  }
  dlg1.value = true
}

async function saveRecord() {
  const r = rec.value
  const body = {
    drill_status: r.drill_status,
    actual_participants: r.actual_participants ?? null,
    duration_minutes: r.duration_minutes ?? null,
    process_description: r.process_description || '',
    problem_list: r.problem_list || '',
    evaluation: r.evaluation || '',
    photo_paths: r.photo_paths || '',
  }
  await api.patchDrill(r.id, body)
  ElMessage.success('已保存演练记录')
  dlg1.value = false
  load()
}

async function openIssues(row) {
  issueDrillId.value = row.id
  dlgIssues.value = true
  const { data } = await api.fetchDrillIssues(row.id)
  issueRows.value = data
  newIssue.value = ''
}

async function addIssue() {
  const t = newIssue.value.trim()
  if (!t) return ElMessage.warning('请输入问题描述')
  await api.createDrillIssue(issueDrillId.value, { description: t })
  newIssue.value = ''
  ElMessage.success('已添加')
  const { data } = await api.fetchDrillIssues(issueDrillId.value)
  issueRows.value = data
}

onMounted(async () => {
  await loadPlans()
  await load()
})
</script>
