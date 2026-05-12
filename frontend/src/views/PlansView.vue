<template>
  <div>
    <el-button type="primary" @click="openCreatePlan">新建预案</el-button>
    <el-table :data="plans" v-loading="loading" style="width: 100%; margin-top: 12px" row-key="id">
      <el-table-column prop="id" label="ID" width="72" />
      <el-table-column prop="name" label="预案名称" min-width="160" />
      <el-table-column prop="plan_type" label="类型" width="120" />
      <el-table-column prop="scenario" label="适用场景" show-overflow-tooltip min-width="160" />
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openVersions(row)">版本</el-button>
          <el-button link @click="openEditPlan(row)">编辑</el-button>
          <el-button link type="danger" @click="rmPlan(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dlgPlan" :title="planForm.id ? '编辑预案' : '新建预案'" width="520px">
      <el-form label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="planForm.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="planForm.plan_type" style="width: 100%">
            <el-option v-for="t in planTypes" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="适用场景">
          <el-input v-model="planForm.scenario" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlgPlan = false">取消</el-button>
        <el-button type="primary" @click="savePlan">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dlgVer" :title="'版本：' + (activePlan?.name || '')" width="900px" top="4vh">
      <el-button type="primary" size="small" @click="openNewVersion">新增修订版本</el-button>
      <el-table :data="versions" style="width: 100%; margin-top: 10px" max-height="420">
        <el-table-column prop="version_no" label="版本" width="90" />
        <el-table-column prop="approval_status" label="状态" width="100" />
        <el-table-column prop="is_current" label="当前生效" width="90">
          <template #default="{ row }">
            <el-tag v-if="row.is_current" type="success">是</el-tag>
            <span v-else>否</span>
          </template>
        </el-table-column>
        <el-table-column prop="preparer" label="编制人" width="100" />
        <el-table-column prop="approver" label="审批人" width="100" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editVersion(row)">编辑</el-button>
            <el-button link v-if="row.approval_status === 'draft'" @click="doSubmit(row)">提交审批</el-button>
            <el-button link v-if="row.approval_status === 'pending'" type="success" @click="doApprove(row)">批准</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="dlgEditV" :title="'编辑版本 ' + (verForm.version_no || '')" width="720px">
      <el-form label-width="100px">
        <el-form-item label="版本号">
          <el-input v-model="verForm.version_no" :disabled="!!verForm.id" />
        </el-form-item>
        <el-form-item label="修订记录">
          <el-input v-model="verForm.revision_record" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="编制人">
          <el-input v-model="verForm.preparer" />
        </el-form-item>
        <el-form-item label="正文 Markdown">
          <el-input v-model="verForm.content_md" type="textarea" :rows="10" />
        </el-form-item>
        <el-form-item label="预览">
          <div class="markdown-body" v-html="mdHtml(verForm.content_md)"></div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlgEditV = false">关闭</el-button>
        <el-button type="primary" @click="saveVersion">保存内容</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { marked } from 'marked'
import * as api from '../api'

const planTypes = ['火灾', '危化品泄漏', '地震', '停电', '自然灾害']
const plans = ref([])
const loading = ref(false)
const dlgPlan = ref(false)
const planForm = ref({ name: '', plan_type: '火灾', scenario: '' })

const dlgVer = ref(false)
const activePlan = ref(null)
const versions = ref([])

const dlgEditV = ref(false)
const verForm = ref({})

function mdHtml(src) {
  try {
    return marked.parse(src || '')
  } catch {
    return ''
  }
}

async function loadPlans() {
  loading.value = true
  try {
    const { data } = await api.fetchPlans()
    plans.value = data
  } finally {
    loading.value = false
  }
}

function openCreatePlan() {
  planForm.value = { name: '', plan_type: '火灾', scenario: '' }
  dlgPlan.value = true
}

function openEditPlan(row) {
  planForm.value = { id: row.id, name: row.name, plan_type: row.plan_type, scenario: row.scenario }
  dlgPlan.value = true
}

async function savePlan() {
  const f = planForm.value
  if (!f.name?.trim()) return ElMessage.warning('请填写名称')
  if (f.id) {
    await api.patchPlan(f.id, { name: f.name, plan_type: f.plan_type, scenario: f.scenario })
    ElMessage.success('已更新')
  } else {
    await api.createPlan({ name: f.name, plan_type: f.plan_type, scenario: f.scenario })
    ElMessage.success('已创建')
  }
  dlgPlan.value = false
  loadPlans()
}

async function rmPlan(row) {
  await ElMessageBox.confirm('确认删除该预案及全部版本、关联演练？', '提示', { type: 'warning' })
  await api.deletePlan(row.id)
  ElMessage.success('已删除')
  loadPlans()
}

async function openVersions(row) {
  activePlan.value = row
  dlgVer.value = true
  const { data } = await api.fetchVersions(row.id)
  versions.value = data
}

function openNewVersion() {
  verForm.value = {
    id: null,
    version_no: 'v' + (versions.value.length + 1) + '.0',
    revision_record: '',
    content_md: '# 预案正文\n\n',
    preparer: '',
    emergency_plan_id: activePlan.value.id,
  }
  dlgEditV.value = true
}

function editVersion(row) {
  verForm.value = { ...row }
  dlgEditV.value = true
}

async function saveVersion() {
  const f = verForm.value
  if (!f.version_no?.trim()) return ElMessage.warning('版本号必填')
  if (f.id) {
    await api.patchVersion(f.id, {
      version_no: f.version_no,
      revision_record: f.revision_record,
      content_md: f.content_md,
      preparer: f.preparer,
    })
    ElMessage.success('已保存')
  } else {
    await api.createVersion(activePlan.value.id, {
      version_no: f.version_no,
      revision_record: f.revision_record,
      content_md: f.content_md,
      preparer: f.preparer,
    })
    ElMessage.success('已创建草稿')
  }
  dlgEditV.value = false
  openVersions(activePlan.value)
}

async function doSubmit(row) {
  await api.submitVersion(row.id, {})
  ElMessage.success('已提交审批')
  openVersions(activePlan.value)
}

async function doApprove(row) {
  const { value } = await ElMessageBox.prompt('审批人姓名', '批准生效', {
    inputValue: '安全负责人',
    confirmButtonText: '批准',
  })
  await api.approveVersion(row.id, { approver: value })
  ElMessage.success('已批准为当前生效版本')
  dlgEditV.value = false
  openVersions(activePlan.value)
}

onMounted(loadPlans)
</script>
