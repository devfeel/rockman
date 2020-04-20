<template>
    <div >
      <Tabs @on-click="onTabClick" v-model="tabName">
        <TabPane label="Node日志" name="NodeTraceLog">
          <nodeTraceLog :loadData="loadNodeTraceLogData"></nodeTraceLog>
        </TabPane>
        <TabPane label="任务提交日志" name="TaskSubmitLog">
          <taskSubmitLog :loadData="loadTaskSubmitLogData"></taskSubmitLog>
        </TabPane>
        <TabPane label="任务执行日志" name="TaskExecLog">
          <taskExecLog :loadData="loadTaskExecLogData"></taskExecLog>
        </TabPane>
        <TabPane label="任务状态日志" name="TaskStateLog">
          <taskStateLog :loadData="loadTaskStateLogData"></taskStateLog>
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
      tabName: 'NodeTraceLog',
      loadNodeTraceLogData: false,
      loadTaskSubmitLogData: false,
      loadTaskExecLogData: false,
      loadTaskStateLogData: false
    }
  },
  mounted() {
    this.init();
  },
  methods: {
    init() {
      // this.onPageChange(this.queryParam)
    },
    onTabClick(name) {
      switch (name) {
        case 'NodeTraceLog':
          this.loadNodeTraceLogData = true;
          break;
        case 'TaskSubmitLog':
          this.loadTaskSubmitLogData = true;
          break;
        case 'TaskExecLog':
          this.loadTaskExecLogData = true;
          break;
        case 'TaskStateLog':
          this.loadTaskStateLogData = true;
          break;
        default:
          this.loadNodeTraceLogData = false;
          this.loadTaskSubmitLogData = false;
          this.loadTaskExecLogData = false;
          this.loadTaskStateLogData = false;
      }
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
