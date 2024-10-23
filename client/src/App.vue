<script setup lang="ts">
import { RouterView } from "vue-router";
import { useServersStore } from "./stores/counter";
import { onMounted, ref } from "vue";
const serverStore = useServersStore();

onMounted(() => {
  serverStore.refreshServants();
});

let open = ref(false);
</script>
<template>
  <el-container>
    <el-header class="flex items-center">
      <nav class="bg-white w-full top-0 z-20">
        <div class="items-center mx-auto lg:flex">
          <div class="flex items-center justify-start w-64">
            <a href="javascript:void(0)">
              <img src="/icon.png" class="w-10 h-10" />
            </a>
            <a href="javascript:void(0)">
              <div
                class="text-xl font-serif tracking-wider text-indigo-600 font-bold transition-all hover:text-indigo-800 ml-4"
              >
                Sgrid
              </div>
            </a>
          </div>
          <div
            class="flex-1 justify-between flex-row-reverse lg:overflow-visible lg:flex lg:pb-0 lg:pr-0 lg:h-auto"
            :class="[open ? 'h-screen pb-20 overflow-auto pr-4' : 'hidden']"
          >
            <div>
              <ul class="flex flex-col-reverse space-x-0 lg:space-x-6 lg:flex-row">
                <li class="mt-8 mb-8 lg:mt-0 lg:mb-0">
                  <a
                    @click="$router.push('/server')"
                    class="text-gray-600 hover:text-indigo-600 cursor-pointer"
                  >
                    Controller
                  </a>
                </li>
                <li class="mt-8 mb-8 lg:mt-0 lg:mb-0">
                  <a
                    @click="$router.push('/dashboard')"
                    class="text-gray-600 hover:text-indigo-600 cursor-pointer"
                  >
                    DashBoard
                  </a>
                </li>
                <li class="mt-4 lg:mt-0">
                  <a
                    @click="$router.push('/mainLog')"
                    class="py-3 px-4 text-center border cursor-pointer text-gray-600 hover:text-indigo-600 rounded-md block lg:inline lg:border-0"
                  >
                    MainLog
                  </a>
                </li>
                <li class="mt-8 lg:mt-0">
                  <a
                    @click="$router.push('/devops')"
                    class="py-3 px-4 text-center text-white cursor-pointer bg-indigo-600 hover:bg-indigo-700 rounded-md shadow block lg:inline"
                  >
                    Devops
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </nav>
    </el-header>
    <RouterView />
  </el-container>
</template>
