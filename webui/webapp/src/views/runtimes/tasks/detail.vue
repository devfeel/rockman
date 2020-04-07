<template>
    <div>
        <div :style="{background: '#fff',padding:'16px'}">
            <b>基本信息</b>
            <Form label-position="left" :label-width="100" :model="taskForm">
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
                            <FormItem label="任务执行类型"><span v-text="taskForm.TaskType"></span></FormItem>
                        </Col>
                        <Col span="5">
                            <FormItem label="cron表达式"><span v-text="taskForm.Express"></span></FormItem>
                        </Col>
                    </Row>
                    <Row>
                        <Col span="5">
                            <FormItem label="任务状态："><span v-text="taskForm.IsRun"></span></FormItem>
                        </Col>
                        <Col span="5">

                        </Col>
                    </Row>
                    <Row>
                        <Col span="10">
                            <FormItem label="备注"><span v-text="taskForm.Remark"></span></FormItem>
                        </Col>
                    </Row>
            </Form>
        </div>
        <Card >
            <Tabs>
                <TabPane label="执行统计" >
                    <statistics :data="taskForm"></statistics>
                </TabPane>
                <TabPane label="执行日志" >
                    <logs :data="taskForm"></logs>
                </TabPane>
                <TabPane label="状态日志" >
                </TabPane>
                <TabPane label="提交日志" >
                </TabPane>
            </Tabs>
        </Card>
    </div>
</template>
<script>
import logs from './components/logs.vue';
import statistics from './components/statistics.vue';
import { getUrlParam } from '@/common/utils.js';
import { getTaskOnce } from '@/api/task.js';
export default {
    components: { logs, statistics },
    data() {
        return {
            taskForm: {}
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
                    if (res.code === 200) {
                        this.taskForm = res.data;
                    } else {
                        this.$Message.error(res.msg);
                    }
                })
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
