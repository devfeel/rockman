<template>
    <div >
      <Tabs>
        <TabPane label="Node日志" name="NodeTraceLog">
          <nodeTraceLog></nodeTraceLog>
        </TabPane>
        <TabPane label="任务提交日志" name="TaskSubmitLog">
          <taskSubmitLog></taskSubmitLog>
        </TabPane>
        <TabPane label="任务执行日志" name="TaskExecLog">
          <taskExecLog></taskExecLog>
        </TabPane>
        <TabPane label="任务状态日志" name="TaskStateLog">
          <taskStateLog></taskStateLog>
        </TabPane>
      </Tabs>
    </div>
</template>
<script>
  import nodeTraceLog from './components/nodeTraceLog.vue';
  import taskExecLog from './components/taskExecLog.vue';
  import taskStateLog from './components/taskStateLog.vue';
  import taskSubmitLog from './components/taskSubmitLog.vue';
export default {
  components: { nodeTraceLog, taskExecLog, taskStateLog, taskSubmitLog },
  data() {
    return {
    }
  },
  mounted() {
    this.init();
  },
  methods: {
    init() {
      // this.onPageChange(this.queryParam)
    }
  }

}
</script>
<style lang="less" scoped>
.search {
}

.btn {
  color: #999;
  height: 40px;
  font-size: 28px;
  font-weight: bold;
}
</style>
