<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-button type="primary" @click="createGroup">创建属性</el-button>
      </el-form-item>
    </el-form>
    <el-table border :data="servantList">
      <el-table-column type="index" label="序号" width="180"></el-table-column>
      <el-table-column prop="tagName" label="标签名"></el-table-column>
      <el-table-column prop="tagEnglishName" label="标签英文名"></el-table-column>
      <el-table-column prop="creatTime" label="creatTime"></el-table-column>
      <el-table-column label="操作">
        <template #default="scoped">
          <el-button
            @click="updateGroup(scoped.row)"
            type="text"
            style="color: var(--sgrid-primary-choose-color)"
            >修改</el-button
          >
          <el-button @click="deleteGroup(scoped.row)" type="text">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="editDialogVisible" title="编辑服务组信息" width="50%">
      <el-form :model="servant" label-width="100px">
        <el-form-item label="tagName">
          <el-input v-model="servant.tagName"></el-input>
        </el-form-item>
        <el-form-item label="tagEnglishName">
          <el-input v-model="servant.tagEnglishName"></el-input>
        </el-form-item>
        <el-form-item label="userId">
          <el-input v-model="servant.userId" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="creatTime">
          <el-input v-model="servant.creatTime" :disabled="true"></el-input>
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
import { useUserStore } from '@/stores/counter'
import { ElMessage, ElMessageBox } from 'element-plus'
import _ from 'lodash'
import { onMounted, ref } from 'vue'
const userStore = useUserStore()
const servantList = ref<Array<any>>([])
async function getGroupList() {
  const servantsResp = await api.getGroup(userStore.userInfo.id)
  servantList.value = servantsResp.data
  console.log('servantResp', servantsResp)
}
onMounted(async () => {
  await getGroupList()
})
const editDialogVisible = ref(false)
const servant = ref({
  id: 0,
  tagName: '',
  tagEnglishName: '',
  userId: 0,
  creatTime: ''
})
function updateGroup(row: any) {
  editDialogVisible.value = true
  servant.value = _.cloneDeep(row)
  console.log('row', row)
}

function createGroup() {
  reset()
  editDialogVisible.value = true
}

function reset() {
  servant.value.id = 0
  servant.value.tagName = ''
  servant.value.tagEnglishName = ''
  servant.value.userId = 0
}
async function deleteGroup(row: typeof servant.value) {
  ElMessageBox.confirm('确认删除?', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      const resp = await api.deleteGroup(row.id)
      await getGroupList()
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
  await api.saveGroup(servant.value)
  await getGroupList()
  editDialogVisible.value = false
}
</script>

<style scoped></style>
