<template>
  <el-container class="layoutRoot">
    <el-header height="56px">
      <div class="hdr">
        <span class="title">应急预案与演练记录</span>
        <span class="sub">预案 · 演练 · 整改闭环（演示）</span>
      </div>
    </el-header>
    <el-container>
      <el-aside width="210px">
        <el-menu router :default-active="$route.path">
          <el-menu-item index="/compliance">
            <el-icon><DataAnalysis /></el-icon>合规统计
          </el-menu-item>
          <el-menu-item index="/plans">
            <el-icon><Document /></el-icon>预案与版本
          </el-menu-item>
          <el-menu-item index="/drills">
            <el-icon><Flag /></el-icon>演练计划与记录
          </el-menu-item>
          <el-menu-item index="/rectifications">
            <el-icon><CircleCheck /></el-icon>问题整改
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { DataAnalysis, Document, Flag, CircleCheck } from '@element-plus/icons-vue'
</script>

<style>
html,
body,
#app {
  height: 100%;
  margin: 0;
}
.layoutRoot {
  min-height: 100%;
}
.hdr {
  display: flex;
  align-items: baseline;
  gap: 12px;
}
.title {
  font-size: 18px;
  font-weight: 600;
}
.sub {
  color: #909399;
  font-size: 13px;
}
.markdown-body {
  padding: 8px;
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  min-height: 120px;
  background: var(--el-fill-color-blank);
  text-align: left;
}
.markdown-body h1,
.markdown-body h2 {
  margin: 0.4em 0;
}
</style>
