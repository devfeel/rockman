<template>
    <div>
        <div :style="{background: '#fff',padding:'16px'}">
            <b>基本信息</b>
            <Form label-position="left" :label-width="150" :model="taskForm">
                    <Row>
                        <Col span="5">
                            <FormItem label="任务编码："><span v-text="taskForm.TaskID"></span></FormItem>
                        </Col>
                        <Col span="5">
                            <FormItem label="任务类型：">
                                <span v-text="taskForm.TargetType"></span>
                            </FormItem>
                        </Col>
                    </Row>
                    <Row>
                        <Col span="5">
                            <FormItem label="任务执行类型："><span v-text="taskForm.TaskType"></span></FormItem>
                        </Col>
                        <Col span="5">
                            <FormItem label="cron表达式："><span v-text="taskForm.Express"></span></FormItem>
                        </Col>
                    </Row>
                    <Row>
                        <Col span="5">
                            <FormItem label="任务状态："><span v-if="taskForm.IsRun">执行中</span><span v-if="!taskForm.IsRun">已就绪</span></FormItem>
                        </Col>
                        <Col span="5">

                        </Col>
                    </Row>
                    <Row>
                        <Col span="10">
                            <FormItem label="备注："><span v-text="taskForm.Remark"></span></FormItem>
                        </Col>
                    </Row>
            </Form>
        </div>
        <Card >
            <Tabs @on-click="onTabClick" v-model="tabName">
                <TabPane label="执行统计" name="statistics">
                    <statistics :data="taskForm" :loadData="loadData"></statistics>
                </TabPane>
                <TabPane label="执行日志" name="execLogs">
                    <logs :data="taskForm" :loadData="loadExecLogData"></logs>
                </TabPane>
                <TabPane label="状态日志" name="stateLogs">
                    <stateLogs :data="taskForm" :loadData="loadStateLogData"></stateLogs>
                </TabPane>
            </Tabs>
        </Card>
    </div>
</template>
<script>
import logs from './components/logs.vue';
import stateLogs from './components/stateLogs.vue';
import statistics from './components/statistics.vue';
import { getUrlParam } from '@/common/utils.js';
import { getTaskOnce } from '@/api/task.js';
export default {
    components: { logs, statistics, stateLogs },
    data() {
        return {
            taskForm: {},
            tabName: 'statistics',
            loadData: false,
            loadExecLogData: false,
            loadStateLogData: false
        }
    },
     mounted() {
      this.init();
    },
    watch: {
        '$route': function (to, from) {
            // 执行数据更新查询
            this.init()
        }
    },
    methods: {
        init() {
            if (getUrlParam('id')) {
                getTaskOnce({ID: getUrlParam('id')}).then(res => {
                    if (res.RetCode === 0) {
                        this.taskForm = res.Message;
                    } else {
                        this.$Message.error(res.RetMsg);
                    }
                })
            }
            this.tabName = 'statistics';
        },
        onTabClick(name) {
            switch (name) {
                case 'execLogs':
                    this.loadExecLogData = true;
                    break;
                case 'stateLogs':
                    this.loadStateLogData = true;
                    break;
                default:
                    this.loadExecLogData = false;
                    this.loadStateLogData = false;
            }
        }
    }
}
</script>
<style lang="less">
.height-main{
    // position: absolute;
    // height: calc(100%);
    /* height: calc(100% - 104px); */
}
body .ivu-modal .ivu-select-dropdown{
 position: fixed !important;
}
</style>
