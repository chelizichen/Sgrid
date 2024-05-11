<script lang="ts">
export default {
  name: "upload-form",
};
</script>

<template>
  <div>
    <!-- 文件上传表单 -->
    <el-dialog
      append-to-body
      :model-value="props.uploadVisible"
      width="50%"
      title="Release"
      @close="emits('CLOSE_UPLOAD_DIALOG')"
    >
      <el-form :model="uploadForm" label-width="150px">
        <el-form-item label="Server Name" required>
          <el-input :value="props.serverName" disabled></el-input>
        </el-form-item>
        <el-form-item label="Document" required>
          <el-input v-model="uploadForm.content" type="textarea" row="5"></el-input>
        </el-form-item>
        <el-form-item label="File" required>
          <el-upload
            :file-list="state.fileList"
            :show-file-list="true"
            :on-change="handleFileChange"
            :auto-upload="false"
            action="/upload/uploadServer"
          >
            <el-button slot="trigger" size="small">Choose File</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <div style="display: flex; align-items: center; justify-content: center">
          <el-button type="primary" @click="emits('CLOSE_UPLOAD_DIALOG')"
            >Close</el-button
          >
          <el-button type="success" @click="uploadFile">Upload</el-button>
        </div>
      </span>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ElLoading, ElNotification, type UploadUserFile } from "element-plus";
import { ref,reactive } from "vue";
import api from '@/api/server'
import SparkMD5 from "spark-md5";

const props = defineProps<{
  uploadVisible: boolean;
  serverName:string;
  servantId:number;
}>();

const emits = defineEmits(["CLOSE_UPLOAD_DIALOG"]);

const uploadForm = ref({
  file: null,
  content: "",
});

const state = reactive({
  fileList: <Array<UploadUserFile>>[],
  hash:'',
});
function handleFileChange(file:UploadUserFile) {
  console.log('file',file);
    if (!file.name.includes(props.serverName)) {
    ElNotification.error(`请上传正确的服务包 [ ${props.serverName} ] `)
    uploadForm.value.file = null
    state.fileList = []
  } else {
    const fileReader = new FileReader();
    const spark = new SparkMD5.ArrayBuffer();
    fileReader.onload = function(e){
      spark.append(e.target.result);
      var md5 = spark.end()
      state.hash = md5
    }
    fileReader.readAsArrayBuffer(file.raw!)
    uploadForm.value.file = file.raw as any
    state.fileList = [file]
  }
}

async function uploadFile(){
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    spinner: 'el-icon-loading',
    background: 'rgba(0, 0, 0, 0.7)'
  })
    const formData = new FormData()
    formData.append('serverName', props.serverName)
    formData.append('file', uploadForm.value.file as unknown as Blob)
    formData.append('content', uploadForm.value.content)
    formData.append('servantId', String(props.servantId))
    formData.append('fileHash', String(state.hash))
    const data = await api.uploadSgridServer(formData)
    if (data.code) {
      ElNotification({
        type: 'error',
        message: '上传失败' + data.message
      })
    } else {
      ElNotification({
        type: 'success',
        message: '上传成功'
      })
      emits('CLOSE_UPLOAD_DIALOG')
    }
    loading.close()
}
</script>
