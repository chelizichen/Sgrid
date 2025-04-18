<script lang="ts">
export default {
  name: 'grids-component'
}
</script>

<template>
  <div>
    <div class="mb-2.5 p-2.5">
      <div class="flex justify-between items-center">
        <div class="flex flex-col items-center">
          <div class="font-bold">ServerName</div>
          <div
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300"
            @click="state.uploadVisible = true">
            {{ props.serverName }}
            <span v-if="props.serverVersion"
              class="font-medium text-red-700 cursor-pointer hover:text-red-900 hover:font-bold transition-all duration-300">
              @{{ props.serverVersion }}
            </span>
          </div>
        </div>
        <div class="flex flex-col items-center">
          <el-button
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300 disabled:text-indigo-300"
            type="text" @click="releaseServer" :disabled="selectionGrid.length == 0">Release</el-button>
        </div>
        <div>
          <el-button
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300"
            type="text" @click="showConfiguration">Configuration</el-button>
        </div>
        <div>
          <el-button
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300 disabled:text-indigo-300"
            type="text" @click="restartServer" :disabled="selectionGrid.length == 0">Restart</el-button>
        </div>
        <div>
          <el-button
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300 disabled:text-indigo-300"
            type="text" @click="batchShutdown" :disabled="selectionGrid.length == 0">BatchShutDown</el-button>
        </div>
        <div>
          <el-button
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300 disabled:text-indigo-300"
            type="text" @click="openInvokeCmdDialog" :disabled="selectionGrid.length == 0">InvokeCmd</el-button>
        </div>
      </div>
    </div>
    <el-table :data="props.gridsList" style="width: 100%" border @selection-change="handleSelectionChange" stripe>
      <el-table-column type="selection" width="55" />
      <el-table-column label="Grid">
        <template #default="scoped">
          <el-button type="text"
            class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300"
            @click="toLog(scoped.row)">{{ scoped.row.gridNode.ip }}:{{ scoped.row.port }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="gridServant.serverName" label="serverName">
        <template #default="scoped">
          <template v-if="scoped.row.gridServant.preview">
            <el-button type="text" @click="toPreview(scoped.row.gridServant.preview)"
              class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300">{{
              scoped.row.gridServant.serverName }}</el-button>
          </template>
          <template v-else>
            <span>{{ scoped.row.gridServant.serverName }}</span>
          </template>
        </template>
      </el-table-column>
      <el-table-column label="Status">
        <template #default="scoped">
          <div @click="$emit('checkStatus')" :class="getGridStatus(scoped.row.status)" class="cursor-pointer">
            <div :class="getGridStatus(scoped.row.status) + '-icon'"></div>
            {{ getGridStatus(scoped.row.status) }}
            {{ scoped.row.pid ? ` (${scoped.row.pid})` : '' }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Type">
        <template #default="scoped">
          <div>{{ scoped.row.gridServant.language }}({{ scoped.row.gridServant.protocol }})</div>
        </template>
      </el-table-column>
      <el-table-column label="Operate">
        <template #default="scoped">
          <div>
            <div
              class="font-medium text-indigo-700 cursor-pointer hover:text-indigo-800 hover:font-bold transition-all duration-300"
              @click="shutDown(scoped.row)">
              Shutdown
            </div>
            <div
              class="font-medium text-red-500 cursor-pointer hover:text-text-800 hover:font-bold transition-all duration-300"
              @click="deleteGridById(scoped.row)">
              Delete
            </div>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <el-divider content-position="left">
      <el-button type="text" @click="getLogList($props.gridsList)"><el-icon
          style="font-size: large; font-weight: 900; color: gray">
          <Loading />
        </el-icon>
      </el-button>
    </el-divider>
    <el-table :data="statLogList" style="width: 100%; margin-top: 20px" border height="450" stripe>
      <el-table-column prop="id" label="id" width="180" />
      <el-table-column label="name" width="180">
        <template #default>
          <div>{{ props.serverName }}</div>
        </template>
      </el-table-column>
      <el-table-column label="Grid">
        <template #default="scoped">
          <div>{{ scoped.row.gridInfo.gridNode.ip }}:{{ scoped.row.gridInfo.port }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="pid" label="pid">
        <template #default="scoped">
          <div>{{ scoped.row.pid || '--' }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="threads" label="threads" width="180">
        <template #default="scoped">
          <div>{{ scoped.row.threads || '--' }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="isRunning" label="isRunning" width="180">
        <template #default="scoped">
          <div>{{ scoped.row.isRunning || '--' }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="createTime" />
      <el-table-column prop="stat" label="behavior" />
    </el-table>
    <uploadComponent :upload-visible="state.uploadVisible" :serverName="$props.serverName" :servantId="$props.servantId"
      :servantLanguage="$props.servantLanguage" @CLOSE_UPLOAD_DIALOG="() => (state.uploadVisible = false)">
    </uploadComponent>
    <releaseComponent :releaseVisible="state.releaseVisible" :serverName="$props.serverName" :releaseList="releaseList"
      :selectionGrid="selectionGrid" @CLOSE_RELEASE_DIALOG="() => (state.releaseVisible = false)"
      @RELEASE_SERVER_BY_ID="(val: number) => handleRelease(val) ">
    </releaseComponent>
    <servantConf :servantId="showConfigurationId" :dialogVisible="showConfigurationVisible"
      @CLOSE_RELEASE_DIALOG="() => (showConfigurationVisible = false)"></servantConf>
    <invokeCmd :cmd-visible="cmdVisible" :server-name="props.serverName" :selection-grid="selectedGrid"
      @CLOSE_CMD_DIALOG="cmdVisible = false" />
  </div>
</template>
<script lang="ts" setup>
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
import uploadComponent from './upload.vue'
import releaseComponent from './release.vue'
import servantConf from './servantConf.vue'
import invokeCmd from './invokeCmd.vue'
import moment from 'moment'
import api from '@/api/server'
import { ElNotification, ElMessageBox, ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import type { T_Grid, T_RelaseItem, T_StatLogListItem, T_Status } from '@/dto/dto'
import _ from 'lodash'

const props = defineProps<{
  gridsList: T_Grid[]
  serverName: string
  servantId: number
  servantLanguage: string
  serverVersion: number
}>()
const emits = defineEmits(['checkStatus'])
const svrInfo = computed(() => {
  const servantBaseInfo = _.first(props.gridsList)
  return {
    serverName: props.serverName,
    serverLanguage: servantBaseInfo?.gridServant.language,
    serverProtocol: servantBaseInfo?.gridServant.protocol,
    execPath: servantBaseInfo?.gridServant.execPath,
    servantId: Number(props.servantId),
  }
})

async function getLogList(gridList: T_Grid[]) {
  const resp = await Promise.all(
    gridList.map(async (v: T_Grid) => {
      const list = await api.getStatLogList({
        id: v.id
      })
      const ret = list.data.list.map((item: T_Status) => {
        item.gridInfo = v
        return item
      })
      return ret
    })
  )
  const newArr: any[] = []
  resp.forEach((v) => {
    newArr.push(...v)
  })

  statLogList.value = newArr
    .map((v) => {
      v.createTime = moment(v.createTime).format('YYYY-MM-DD HH:mm:ss')
      return v
    })
    .sort((a, b) => b.id - a.id)
  console.log('statLogList', statLogList.value);
}

const statLogList = ref<T_StatLogListItem[]>([])

onMounted(async () => {
  await checkStat(props.gridsList)
  await getLogList(props.gridsList)
})

const state = reactive({
  uploadVisible: false,
  releaseVisible: false
})

const releaseList = ref<T_RelaseItem[]>([])
async function releaseServer() {
  const data = await api.getUploadList({
    id: props.servantId
  })
  state.releaseVisible = true
  releaseList.value = data.data
}
async function handleRelease(id: number) {
  const releaseItem = releaseList.value.find((v) => v.id == id)
  console.log('selectionGrid', selectionGrid.value)
  const body = {
    filePath: releaseItem?.filePath,
    servantGrids: selectionGrid.value.map((v) => ({
      ip: v.gridNode.ip,
      port: v.port,
      gridId: v.id
    })),
    ...svrInfo.value
  }

  console.log('body', body)

  const data = await api.releaseServer(body, { releaseId: id })
  if (data.code) {
    return ElNotification.error(data.message)
  }
  ElNotification.success('发布成功')
  state.releaseVisible = false
}

async function restartServer() {
  const body = {
    filePath: '',
    servantGrids: selectionGrid.value.map((v) => ({
      ip: v.gridNode.ip,
      port: v.port,
      gridId: v.id
    })),
    ...svrInfo.value,
  }
  const data = await api.restartServer(body)
  if (data.code) {
    return ElNotification.error(data.message)
  }
  ElNotification.success('success!')
}
const selectionGrid = ref<T_Grid[]>([])
function handleSelectionChange(value: T_Grid[]) {
  selectionGrid.value = value
}

function getGridStatus(status: number | undefined) {
  const OFF = 'offline'
  const ON = 'online'
  return status == 1 ? ON : OFF
}

async function batchShutdown() {
  const body = {
    req: selectionGrid.value
      .filter((v) => v.status != 0 && v.status)
      .map((v) => ({
        pid: v.pid,
        gridId: v.id,
        host: v.gridNode.ip,
        port: v.port
      }))
  }

  const data = await api.shutdownServer(body) 
  if (data.code) {
    ElNotification.error(data.message)
  }
  ElNotification.success('关闭成功')
}

async function shutDown(v: T_Grid) {
  const body = {
    req: [
      {
        pid: v.pid,
        gridId: v.id,
        host: v.gridNode.ip,
        port: v.port
      }
    ]
  }
  const data = await api.shutdownServer(body)
  if (data.code) {
    ElNotification.error(data.message)
  }
  ElNotification.success('关闭成功')
}

const router = useRouter()
function toLog(row: T_Grid) {
  const text = router.resolve({
    path: '/logpage',
    query: {
      host: row.gridNode.ip,
      serverName: props.serverName,
      gridId: row.id
    }
  })
  window.open(text.href, '_blank')
}

function toPreview(path: string) {
  window.open(path, '_blank')
}

function checkStat(list: T_Grid[]) {
  const body = list.map((v) => {
    return {
      host: v.gridNode.ip,
      gridId: v.id
    }
  })
  api.checkStat({
    hostPids: body
  })
}

async function deleteGridById(row: T_Grid) {
  ElMessageBox.prompt('确认删除该节点？删除后不可恢复!', 'Confirm', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputPlaceholder: 'input password'
  }).then(async ({ value: password }) => {
    if (password != props.serverName) {
      return ElMessage.error(`password error`)
    }
    if (row.status) {
      return ElNotification.error('error/client :: this grid still alive')
    }
    const data = await api.deleteGrid({
      id: row.id
    })
    if (data.code) {
      return ElNotification.error(data.message)
    }
    ElNotification.success('delete success')
    emits('checkStatus')
  })
}

const showConfigurationVisible = ref(false)
const showConfigurationId = ref(0)
async function showConfiguration() {
  showConfigurationVisible.value = true
  showConfigurationId.value = props.servantId
}

const cmdVisible = ref(false)
const selectedGrid = ref<T_Grid[]>([])

function openInvokeCmdDialog() {
  cmdVisible.value = true
  selectedGrid.value = selectionGrid.value
}
</script>
