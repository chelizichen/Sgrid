<script lang="ts">
export default {
  name: "login-component",
};
</script>
<script lang="ts" setup>
import API from "@/api/server";
import LoginPage from "@/assets/login.jpg";
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
  <div class="body">
    <div class="pic">
      <img :src="LoginPage" style="width: 70vw; height: 90vh" />
    </div>
    <div>
      <div class="container">
        <div class="login-form">
          <h2>
            欢迎登录
            <img src="@/assets/title.png" style="width: 120px; height: 60px" />
            平台
          </h2>
          <div class="input-group">
            <label for="username">用户名</label>
            <input
              type="text"
              id="username"
              placeholder="请输入用户名"
              required
              v-model="username"
            />
          </div>
          <div class="input-group">
            <label for="password">密码</label>
            <input
              type="password"
              id="password"
              placeholder="请输入密码"
              required
              v-model="password"
            />
          </div>
          <button @click="login" class="button">登录</button>
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
.body {
  font-family: Arial, sans-serif;
  background-color: #f5f5f5;
  margin: 0;
  padding: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 90vh;
  .pic {
    margin-right: 20px;
  }
}

.container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 95%;
}

.login-form {
  h2 {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.login-form h2 {
  margin-bottom: 30px;
  color: #333;
}

.input-group {
  margin-bottom: 20px;
  display: flex;
  align-items: flex-end;
}
.input-group label {
  display: block;
  margin-bottom: 10px;
  color: #666;
  width: 100px;
}

.input-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.button {
  width: 100%;
  padding: 10px;
  background-color: rgb(207, 90, 124);
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.button:hover {
  background-color: rgb(144, 90, 124);
}

.license {
  margin-top: 100px;
  width: 100%;
  .title {
    text-align: center;
    color: rgb(207, 90, 124);
    margin: 5px 0;
  }
  .h2 {
    font-size: 20px;
    font-weight: 600;
  }
  .h4 {
    font-size: 14px;
    color: #838383;
    font-weight: 300;
  }
}
</style>
