<script lang="ts">
export default {
  name: 'aside-component'
}
</script>
<script lang="ts" setup>
import type { Item } from '@/dto/dto'
const props = defineProps<{
  serverList: any[]
}>()
const emits = defineEmits(['handleOpen'])
function handleOpen(item: Partial<Item>) {
  emits('handleOpen', item)
}
</script>

<template>
  <div>
    <el-menu active-text-color="var(--sgrid-primary-choose-color)" style="border: none">
      <template v-if="props.serverList && props.serverList.length">
        <template v-for="(parent, index) in props.serverList" :key="parent.id">
          <el-sub-menu :index="parent.id + '-' + index">
            <template #title>
              <el-icon><Grid /></el-icon>
              <span>{{ parent.tagEnglishName }}</span>
            </template>
            <el-menu-item
              v-for="(item, index) in parent.servants"
              class="app-text-center"
              :index="item"
              :key="index"
              @click="handleOpen(item)"
            >
              <el-icon>
                <TrendCharts />
              </el-icon>
              <template #title>{{ item.serverName }}</template>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </template>
      <template v-else>
        <el-empty />
      </template>
    </el-menu>
  </div>
</template>

<style scoped></style>
