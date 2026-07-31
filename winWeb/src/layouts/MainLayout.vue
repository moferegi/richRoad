<script setup>
import AppHeader from '@/components/AppHeader.vue'
import AppSidebar from '@/components/AppSidebar.vue'
</script>

<template>
  <div class="main-layout">
    <AppHeader />
    <div class="layout-body">
      <AppSidebar />
      <main class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.main-layout {
  min-height: 100vh;
  background: $bg-page;
}

.layout-body {
  display: flex;
  padding-top: $header-height;
  min-height: 100vh;
}

.layout-content {
  flex: 1;
  margin-left: $sidebar-width;
  min-height: calc(100vh - #{$header-height});
  padding: $spacing-xl;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
