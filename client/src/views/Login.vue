<script lang="ts">
export default {
  name: "login-component",
};
</script>
<script lang="ts" setup>
import API from "@/api/server";
import { ElMessage } from "element-plus";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { md5 } from "js-md5";
const router = useRouter();
const username = ref("");
const password = ref("");

async function login() {
  const data = new FormData();
  const tkn = md5(password.value);
  data.append("password", tkn);
  data.append("username", username.value);
  const ret = await API.Login(data);
  if (ret.code) {
    ElMessage.error("Please enter a valid token.");
  } else {
    router.push("/server");
    localStorage.setItem("token", tkn);
  }
}
</script>
<template>
  <div class="body">
    <div class="container">
      <div class="login-form">
        <h2>欢迎登录 SgridCloud 平台</h2>
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
  height: 100vh;
}

.container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.1);
  padding: 40px;
  max-width: 400px;
  width: 100%;
}

.login-form {
  text-align: center;
}

.login-form h2 {
  margin-bottom: 30px;
  color: #333;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 10px;
  color: #666;
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
</style>
