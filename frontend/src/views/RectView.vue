<template>
  <div>
    <el-tabs v-model="tab">
      <el-tab-pane label="整改台账" name="r">
        <el-button @click="loadRects" style="margin-bottom: 10px">刷新</el-button>
        <el-table :data="rects" style="width: 100%">
          <el-table-column prop="id" label="ID" width="64" />
          <el-table-column label="关联预案">
            <template #default="{ row }">{{ row.issue?.drill?.plan?.name || '—' }}</template>
          </el-table-column>
          <el-table-column label="问题摘要" min-width="160">
            <template #default="{ row }">{{ row.issue?.description }}</template>
          </el-table-column>
          <el-table-column prop="responsible_person" label="责任人" width="140" />
          <el-table-column prop="corrective_measure" label="整改措施" min-width="140" show-overflow-tooltip />
          <el-table-column label="期限" width="120">
            <template #default="{ row }">{{ fmt(row.due_date) }}</template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" />
          <el-table-column label="完成日期" width="120">
            <template #default="{ row }">{{ fmt(row.completed_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="edit(row)">编辑</el-button>
              <el-button link type="success" v-if="row.status !== 'done'" @click="markDone(row)">完成</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="全部问题列表" name="i">
        <el-button @click="loadIssues" style="margin-bottom: 10px">刷新</el-button>
        <el-table :data="issues" style="width: 100%">
          <el-table-column prop="id" label="#" width="64" />
          <el-table-column label="演练ID" prop="drill_id" width="84" />
          <el-table-column label="预案" min-width="120">
            <template #default="{ row }">{{ row.drill?.plan?.name }}</template>
          </el-table-column>
          <el-table-column prop="description" label="问题描述" min-width="220" />
          <el-table-column label="整改" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openRect(row)">登记整改</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="dlg" title="整改项" width="520px">
      <el-form label-width="96px">
        <el-form-item label="责任人"><el-input v-model="f.responsible_person" /></el-form-item>
        <el-form-item label="整改措施"><el-input v-model="f.corrective_measure" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="完成期限"><el-date-picker v-model="f.due_date" type="date" value-format="YYYY-MM-DD" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="f.status"><el-option label="待整改" value="pending" /><el-option label="已完成" value="done" /></el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" @click="saveRect">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import * as api from '../api'

const tab = ref('r')
const rects = ref([])
const issues = ref([])

const dlg = ref(false)
const f = ref({ id: null, issue_id: null, responsible_person: '', corrective_measure: '', due_date: '', status: 'pending' })

function fmt(v) {
  return v ? dayjs(v).format('YYYY-MM-DD') : '—'
}

async function loadRects() {
  const { data } = await api.fetchRectifications()
  rects.value = data
}

async function loadIssues() {
  const { data } = await api.fetchIssues()
  issues.value = data
}

function edit(row) {
  f.value = {
    id: row.id,
    issue_id: null,
    responsible_person: row.responsible_person || '',
    corrective_measure: row.corrective_measure || '',
    due_date: row.due_date ? dayjs(row.due_date).format('YYYY-MM-DD') : '',
    status: row.status,
    completed_at: row.completed_at,
  }
  dlg.value = true
}

async function saveRect() {
  const x = f.value
  if (!x.issue_id && !x.id) return
  if (x.issue_id) {
    await api.createRectFromIssue(x.issue_id, {
      responsible_person: x.responsible_person,
      corrective_measure: x.corrective_measure,
      due_date: x.due_date || null,
    })
    ElMessage.success('已新建整改')
  } else {
    const body = {
      responsible_person: x.responsible_person,
      corrective_measure: x.corrective_measure,
      due_date: x.due_date || null,
      status: x.status,
    }
    if (x.status === 'done' && !x.completed_at) {
      body.completed_at = dayjs().format('YYYY-MM-DD')
    }
    await api.patchRectification(x.id, body)
    ElMessage.success('已保存')
  }
  dlg.value = false
  loadRects()
  loadIssues()
}

function openRect(issue) {
  f.value = {
    id: null,
    issue_id: issue.id,
    responsible_person: '',
    corrective_measure: '',
    due_date: '',
    status: 'pending',
  }
  dlg.value = true
}

async function markDone(row) {
  await api.patchRectification(row.id, {
    status: 'done',
    completed_at: dayjs().format('YYYY-MM-DD'),
  })
  ElMessage.success('已标记完成')
  loadRects()
}

watch(tab, (v) => {
  if (v === 'r') loadRects()
  else loadIssues()
})

onMounted(() => loadRects())
</script>
