<template>
  <div class="content">
    <div class="log-header">
        <el-page-header @back="goBack" :content="task.TaskID">
            </el-page-header>
    </div>
    <div class="log">
        <el-tabs v-model="activeName" @tab-click="onTabClick">
            <el-tab-pane label="提交日志" name="submitLog">
                <submitLogs :TaskID="task.TaskID"></submitLogs>
            </el-tab-pane>
            <el-tab-pane label="执行日志" name="execLogs">
                <execLogs :TaskID="task.TaskID"></execLogs>
            </el-tab-pane>
            <el-tab-pane label="状态日志" name="stateLogs">状态日志</el-tab-pane>
        </el-tabs>
    </div>
  </div>
</template>
<script>
import { getTaskOnce } from '@/api/task.js';
import submitLogs from './components/submitLogs.vue';
import execLogs from './components/execLogs.vue';
export default {
    components: { submitLogs, execLogs },
    data() {
    return {
        task: '',
        activeName: 'submitLog'
    };
    },
    activated() {
        this.onInit();
    },
    methods: {
        onInit() {
            if (this.$route.query.id) {
                getTaskOnce({ID: this.$route.query.id}).then(res => {
                    if (res.RetCode === 0) {
                        this.task = res.Message;
                    } else {
                        this.$message.error(res.RetMsg);
                    }
                })
            }
        },
        goBack() {
            this.$router.push({path: '/static/task'})
        },
        onTabClick(tab, event) {

        }
    }
};
</script>
<style lang="less" scoped>
.content {
    background-color:rgb(243, 243, 243);
    position: absolute;
    top:64px;
    left: 0;
    right: 0;
    padding: 0 0px;
    margin: 0 auto;
    /* margin-top: 40px; */
    /* background: #EDF0F5; */
    height: calc(100%);
    overflow-y:scroll;
}
.log-header{
    margin: 10px 20px;
    font-size: 16px;
}
.log{
    margin: 10px 20px;
    background-color:#fff;
    padding:10px 20px 50px 20px;
}
</style>
