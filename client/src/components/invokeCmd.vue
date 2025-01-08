<script lang="ts" setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api/server'
import { v4 as uuidv4 } from 'uuid'
import type { T_Grid } from '@/dto/dto';

// 定义类型
interface GridNode {
  ip: string
  id: string
}

interface InvokeProps {
  cmdVisible: boolean
  serverName: string
  selectionGrid: T_Grid[]
}
const props = defineProps<InvokeProps>()
const command = ref('')
const args = ref('')

// 定义emit事件
const emit = defineEmits(['CLOSE_CMD_DIALOG', 'INVOKE_CMD'])

// 执行命令
async function handleInvokeCmd() {
  if (!command.value) {
    return ElMessage.warning('请输入命令')
  }

  try {
    const body = {
      command: command.value,
      args: args.value,
      serverName: props.serverName,
      gridIds: props.selectionGrid.map(item => item.id),
      invokeId: uuidv4()
    }
    const buildArgs = `${command.value}|${JSON.stringify(args.value)}`
    ElMessageBox.confirm('确认执行该命令?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      body.gridIds.map(v=>{
        const req = {
            cmd: buildArgs,
            gridId: v,
            invokeId: body.invokeId
        }
        return req
      }).map(v=>{
        api.invokeWithCmd(v)
      })
      ElMessage.success('命令执行成功')
      emit('CLOSE_CMD_DIALOG')
    })
  } catch (error) {
    ElMessage.error(`执行失败: ${error}`)
  }
}
</script>

<template>
  <el-dialog
    :model-value="$props.cmdVisible"
    :title="`执行命令 - ${$props.serverName}`"
    @close="$emit('CLOSE_CMD_DIALOG')"
  >
    <!-- 显示选中的节点 -->
    <div v-for="item in $props.selectionGrid" :key="item.id" style="margin: 3px">
      {{ item.gridNode.ip }}
    </div>

    <!-- 命令输入区域 -->
    <div style="margin: 20px 0">
      <el-input v-model="command">
        <template #prepend>
          <span>命令</span>
        </template>
      </el-input>
    </div>

    <!-- 参数输入区域 -->
    <div style="margin: 10px 0">
      <el-input v-model="args" type="textarea" placeholder="请输入参数 (JSON形式)" rows="5">
      </el-input>
    </div>

    <!-- 按钮区域 -->
    <template #footer>
      <div style="display: flex; align-items: center; justify-content: center">
        <el-button @click="$emit('CLOSE_CMD_DIALOG')">取消</el-button>
        <el-button type="primary" @click="handleInvokeCmd">执行</el-button>
      </div>
    </template>
  </el-dialog>
</template>
