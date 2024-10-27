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
                  <div
                    @click="$router.push('/server')"
                    class="cursor-pointer pb-2 px-4 hover:text-indigo-500 hover:border-b-2 hover:border-indigo-500 duration-200 ease-linear hover:font-semibold"
                  >
                    Controller
                  </div>
                </li>
                <li class="mt-8 mb-8 lg:mt-0 lg:mb-0">
                  <div
                    @click="$router.push('/dashboard')"
                    class="cursor-pointer pb-2 px-4 hover:text-indigo-500 hover:border-b-2 hover:border-indigo-500 duration-200 ease-linear hover:font-semibold"
                  >
                    Monitor
                  </div>
                </li>
                <li class="mt-4 lg:mt-0">
                  <div
                    @click="$router.push('/mainLog')"
                    class="cursor-pointer pb-2 px-4 hover:text-indigo-500 hover:border-b-2 hover:border-indigo-500 duration-200 ease-linear hover:font-semibold"
                  >
                    Logger
                  </div>
                </li>
                <li class="mt-8 lg:mt-0">
                  <div
                    @click="$router.push('/devops')"
                    class="cursor-pointer pb-2 px-4 hover:text-indigo-500 hover:border-b-2 hover:border-indigo-500 duration-200 ease-linear hover:font-semibold"
                  >
                    Devops
                  </div>
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
