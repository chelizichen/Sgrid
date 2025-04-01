<template>
    <div class="app-container">
        <el-table border :data="list" :span-method="mergeCells">
            <el-table-column label="Server" width="440">
                <template #default="scope">
                    <span class="text-black font-bold">{{ scope.row.tagEnglishName
                        }} | {{ scope.row.serverName }} </span>
                </template>
            </el-table-column>
            <el-table-column label="IP:PORT">
                <template #default="scope">
                    <span class="text-indigo-700 font-bold">{{ scope.row.ip }}:{{ scope.row.port }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="platForm" label="platForm" width="180">
                <template #default="scope">
                    <span>{{ scope.row.platForm }}({{ scope.row.isMain ? "主" : "备"
                        }})</span>
                </template>
            </el-table-column>
            <el-table-column prop="serverStatus" label="serverStatus">
                <template #default="scoped">
                    <div :class="getGridStatus(scoped.row.serverStatus)" class="cursor-pointer">
                        <div :class="getGridStatus(scoped.row.serverStatus) + '-icon'"></div>
                        {{ getGridStatus(scoped.row.serverStatus) }}
                        {{ scoped.row.pid ? ` (${scoped.row.pid})` : '' }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column prop="updateTime" label="updateTime">
                <template #default="scope">
                    <span>{{ FMT_DATE(scope.row.updateTime) }}</span>
                </template>
            </el-table-column>
            <el-table-column label="operate">
                <template #default="scope">
                    <span class="text-indigo-700 text-sm cursor-pointer" @click="handleRowClick(scope.row)">Link</span>
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>

<script setup lang="ts">
import api from '@/api/server';
import type { Item } from '@/dto/dto';
import _ from 'lodash';
import moment from 'moment';
import { onMounted, ref } from 'vue';

const list = ref([])
onMounted(() => {
    api.getServerStatusByUser({}).then(res => {
        list.value = _.sortBy(res.data, 'tagEnglishName', 'serverName')
    })
})
function FMT_DATE(t: string) {
    return moment(t || undefined).format('YYYY-MM-DD HH:mm:ss')
}
function getGridStatus(status: number | undefined) {
    const OFF = 'offline'
    const ON = 'online'
    return status == 1 ? ON : OFF
}

const emits = defineEmits(["handleOpen"]);
function handleRowClick(item: Partial<Item>) {
    console.log('item', item);
    const body = {
        id: item.servantId,
        language: item.language,
        serverName: item.serverName,
    }
    emits("handleOpen", body);
}
function mergeCells({ rowIndex, columnIndex }) {
    if (columnIndex === 0) { // 只处理第一列（Server列）
        const currentRow = list.value[rowIndex]
        const prevRow = list.value[rowIndex - 1]
        const nextRow = list.value[rowIndex + 1]
        const isSameAsPrev = prevRow && currentRow.tagEnglishName === prevRow.tagEnglishName && currentRow.serverName === prevRow.serverName
        const isSameAsNext = nextRow && currentRow.tagEnglishName === nextRow.tagEnglishName && currentRow.serverName === nextRow.serverName

        if (!isSameAsPrev && isSameAsNext) {
            let rowSpan = 1;
            for (let i = rowIndex + 1; i < list.value.length; i++) {
                const next = list.value[i]
                if (currentRow.tagEnglishName === next.tagEnglishName && currentRow.serverName === next.serverName) {
                    rowSpan++
                } else {
                    break
                }
            }
            return { rowspan: rowSpan, colspan: 1 }
        } else if (isSameAsPrev) {
            return { rowspan: 0, colspan: 0 } // 不显示当前单元格
        }
    }
    return { rowspan: 1, colspan: 1 } // 默认行为
}

</script>

<style scoped></style>