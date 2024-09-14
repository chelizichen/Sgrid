<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-button type="primary" @click="createProperty">创建属性</el-button>
      </el-form-item>
    </el-form>
    <el-table border :data="servantList">
      <el-table-column type="index" label="序号" width="180"></el-table-column>
      <el-table-column prop="key" label="key"></el-table-column>
      <el-table-column prop="value" label="value"></el-table-column>
      <el-table-column prop="createTime" label="createTime"></el-table-column>
      <el-table-column label="操作">
        <template #default="scoped">
          <el-button
            @click="updateProperty(scoped.row)"
            type="text"
            style="color: var(--sgrid-primary-choose-color)"
            >修改</el-button
          >
          <el-button @click="deleteProperty(scoped.row)" type="text">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="editDialogVisible" title="编辑属性信息" width="50%">
      <el-form :model="servant" label-width="100px">
        <el-form-item label="Key">
          <el-input v-model="servant.key"></el-input>
        </el-form-item>
        <el-form-item label="Value">
          <el-input v-model="servant.value"></el-input>
        </el-form-item>
        <el-form-item label="Operate">
          <el-button type="primary" @click="confirmUpdate">更新</el-button>
          <el-button @click="editDialogVisible = false">取消</el-button>
        </el-form-item>
      </el-form></el-dialog
    >
  </div>
</template>

<script setup lang="ts">
import api from '@/api/server'
import { ElMessage, ElMessageBox } from 'element-plus'
import _ from 'lodash'
import { onMounted, ref } from 'vue'

const servantList = ref<Array<any>>([])
async function getPropertyList() {
  const servantsResp = await api.getPropertys(undefined)
  servantList.value = servantsResp.data
  console.log('servantResp', servantsResp)
}
onMounted(async () => {
  await getPropertyList()
})
const editDialogVisible = ref(false)
const servant = ref<Partial<any>>({
  id: 0,
  key: '',
  value: ''
})
function updateProperty(row: any) {
  editDialogVisible.value = true
  servant.value = _.cloneDeep(row)
  console.log('row', row)
}

function createProperty() {
  reset()
  editDialogVisible.value = true
}

function reset() {
  servant.value.id = 0
  servant.value.key = ''
  servant.value.value = ''
}
async function deleteProperty(row: typeof servant.value) {
  ElMessageBox.confirm('确认删除?', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      const resp = await api.delProperty(row.id)
      await getPropertyList()
      if (resp.code) {
        return ElMessage.error({
          type: 'error',
          message: resp.message
        })
      }
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除'
      })
    })

  editDialogVisible.value = false
}

async function confirmUpdate() {
  await api.setProperty(servant.value)
  await getPropertyList()
  editDialogVisible.value = false
}
</script>

<style scoped></style>
