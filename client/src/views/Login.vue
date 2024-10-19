<script lang="ts">
export default {
  name: "login-component",
};
</script>
<script lang="ts" setup>
import API from "@/api/server";
import { ElNotification } from "element-plus";
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { md5 } from "js-md5";
import { constants, localGet, localSet } from "@/constant";
import { useUserStore } from "@/stores/counter";
const router = useRouter();
const userStore = useUserStore();
const username = ref("");
const password = ref("");

async function login() {
  const data = new FormData();
  const tkn = md5(password.value);
  data.append("password", tkn);
  data.append("username", username.value);
  const ret = await API.Login(data);
  if (ret.code) {
    ElNotification.error("账号或密码错误");
  } else {
    router.push("/server");
    localSet(constants.TOKEN, ret.data.token);
  }
}

onMounted(() => {
  const token = localGet(constants.TOKEN);
  if (token && userStore.userInfo.userName && userStore.userInfo.password) {
    router.push("/dashboard");
  }
});
</script>
<template>
  <div class="bg-slate-50 m-0 p-0 flex items-center justify-center min-h-screen">
    <div>
      <div class="bg-slate-150 rounded-sm p-10 shadow-2xl">
        <div>
          <div class="font-bold mb-8 text-black text-center text-lg">
            欢迎登录 SgridCloud 平台
          </div>
          <div class="mb-5 flex items-end">
            <label for="username" class="mb-3 text-gray-500 w-20">用户名</label>
            <input
              class="w-full p-2 border border-gray-300 rounded-md"
              type="text"
              id="username"
              placeholder="请输入用户名"
              required
              v-model="username"
            />
          </div>
          <div class="mb-5 flex items-end">
            <label for="password" class="mb-3 text-gray-500 w-20">密码</label>
            <input
              class="w-full p-2 border border-gray-300 rounded-md"
              type="password"
              id="password"
              placeholder="请输入密码"
              required
              v-model="password"
            />
          </div>
          <button
            @click="login"
            class="w-full p-2.5 bg-indigo-500 text-white border-none rounded-md cursor-pointer hover:bg-indigo-600"
          >
            登录
          </button>
        </div>
      </div>
      <div class="license">
        <div class="title h2">©2024 SgridCloud Devops Platform</div>
        <div class="title h4">©Powered By Golang ｜ ©Author chelizichen</div>
      </div>
    </div>
  </div>
</template>

<style lang="less">
.license {
  margin-top: 100px;
  width: 100%;
  .title {
    text-align: center;
    color: var(--sgrid-primay-color);
    margin: 5px 0;
  }
  .h2 {
    font-size: 20px;
    font-weight: 600;
  }
  .h4 {
    font-size: 14px;
    color: var(--sgrid-primay-hover-color);
    font-weight: 300;
  }
}
</style>
