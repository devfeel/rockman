<template>
    <div class="tb">
        <div>
            <tableH icon="md-apps" text="任务列表">
                <!-- <div slot="content"></div> -->
                <slot>
                <div style="text-align: right;">
                    <i-button type="info" icon="md-add" @click="onAdd">新建任务</i-button>
                    <i-button type="info" icon="md-refresh" @click="onRefresh(false)">刷新</i-button>
                </div>
                </slot>
            </tableH>
            <tableC
                id="table"
                :columns="columns"
                :dataSource="dataSource"
                :queryParam="queryParam"
                @onPageChange="onPageChange"
                ref="table"
                :height="315">
            </tableC>
        </div>
    <div>
        <Modal v-model="model"
            v-bind:title="modelMessage"
            width="700"
            v-bind:mask-closable="closable"
            v-bind:footer-hide="footerHide"
            @on-ok="onSave"
            class-name="vertical-center-modal">
            <div class="model-content">
                <i-form ref="formValidate" :label-width="160" :rules="ruleValidate" :model="taskForm">
                     <form-item label="任务编码" prop='TaskID'>
                        <Input v-model="taskForm.TaskID" placeholder="任务编码"/>
                    </form-item>
                    <form-item label="任务类型" prop='TargetType'>
                        <Select v-model="taskForm.TargetType">
                            <Option value="http" key="http">HTTP</Option>
                            <Option value="shell" key="shell">Shell</Option>
                            <Option value="goso" key="goso">GOSO</Option>
                        </Select>
                    </form-item>
                    <form-item label="任务执行类型" prop='TaskType'>
                        <Select v-model="taskForm.TaskType">
                            <Option value="cron" key="cron">Cron</Option>
                            <Option value="loop" key="loop">Loop</Option>
                        </Select>
                    </form-item>
                    <form-item label="cron表达式" prop='Express'>
                        <Input v-model="taskForm.Express" placeholder="cron表达式,配置作业触发时间"/>
                    </form-item>
                    <form-item label="任务延迟时间" prop='DueTime'>
                        <InputNumber  :min="1" :step="1" v-model="taskForm.DueTime" ></InputNumber>
                    </form-item>
                    <div name="http" v-if="taskForm.TargetType==='http'">
                        <Divider />
                        <form-item label="请求地址">
                            <Input v-model="httpTaskInfoForm.Url" placeholder="http任务执行请求url"/>
                        </form-item>
                        <form-item label="请求方式">
                            <Select v-model="httpTaskInfoForm.Method">
                                <Option value="get" key="get">Get</Option>
                                <Option value="post" key="post">Post</Option>
                                <Option value="head" key="head">Head</Option>
                            </Select>
                        </form-item>
                        <form-item label="数据类型">
                            <Input v-model="httpTaskInfoForm.ContentType" placeholder="数据类型"/>
                        </form-item>
                        <form-item label="请求超时(s)">
                            <InputNumber :min="0" v-model="httpTaskInfoForm.Timeout"></InputNumber>
                        </form-item>
                    </div>
                    <div name="shell" v-if="taskForm.TargetType==='shell'">
                        <Divider />
                        <form-item label="类型">
                            <Select v-model="shellConfigForm.Type">
                                <Option value="Script" key="Script">Script Mode</Option>
                                <Option value="File" key="File">File Mode</Option>
                            </Select>
                        </form-item>
                        <form-item label="文件选择" v-if="shellConfigForm.Type==='File'">
                            <Upload
                                multiple
                                type="drag"
                                action="//jsonplaceholder.typicode.com/posts/">
                                <div style="padding: 20px 0">
                                    <Icon type="ios-cloud-upload" size="52" style="color: #3399ff"></Icon>
                                    <p>点击或者拖动文件上传</p>
                                </div>
                            </Upload>
                        </form-item>
                    </div>
                    <div name="goso" v-if="taskForm.TargetType==='goso'">
                        <form-item >
                            <Upload
                                multiple
                                type="drag"
                                action="//jsonplaceholder.typicode.com/posts/">
                                <div style="padding: 20px 0">
                                    <Icon type="ios-cloud-upload" size="52" style="color: #3399ff"></Icon>
                                    <p>点击或者拖动文件上传</p>
                                </div>
                            </Upload>
                        </form-item>
                    </div>
                    <form-item label="备注">
                        <Input type="textarea" maxlength="100" show-word-limit v-model="taskForm.Remark" placeholder="Remark"/>
                    </form-item>
                </i-form>
            </div>
            <div slot="footer">
                <Button @click="model=false">取消</Button>
                <Button type="primary" @click="onSave('formValidate')">确定</Button>
            </div>
        </Modal>
        <Modal  v-model="glueModel"  fullscreen :footer-hide="true" >
            <div slot="close" class="model-close"><Button type="primary" @click="onSave()">关 闭</Button></div>
            <glue :data="glueTaskForm" ></glue>
        </Modal>
    </div>
    </div>
</template>
<script>
import Minix from '@/common/tableminix.js';
import tableC from '@/components/table/table.vue';
import tableH from '@/components/table/table-header.vue';
import glue from './components/glue.vue';
import { getTaskList, taskSave, taskUpdate, getTaskOnce, taskDelete } from '@/api/task.js';
require('codemirror/mode/javascript/javascript');
export default {
    components: { tableC, tableH, glue },
    mixins: [Minix],
    data() {
        return {
            columns: [
                {
                    title: '任务编码',
                    key: 'TaskID'
                },
                {
                    title: '任务类型',
                    key: 'TargetType'
                },
                {
                    title: '间隔(Cron)',
                    key: 'Express'
                },
                {
                    title: '状态',
                    key: 'IsRun',
                    render: (h, params) => {
                            const row = params.row;
                            if (row.IsRun) {
                                return h('Span', '运行中');
                            }
                            return h('Span', '已就绪');
                        }
                },
                {
                    title: '操作',
                    slot: 'action',
                    width: 300,
                    align: 'center',
                    render: (h, params) => {
                        let row = params.row;
                        let option = h('div', []);

                        let editOptions = h('Button', {
                                    props: {
                                        type: 'warning',
                                        size: 'small'

                                    },
                                    style: {
                                        marginRight: '5px'
                                    },
                                    on: {
                                        click: () => {
                                            this.onRowEdit(row);
                                        }
                                    }
                                }, '修改');
                        option.children.push(editOptions);
                        if (row.TargetType !== 'http') {
                            let glueOptions = h('Button', {
                                            props: {
                                                type: 'warning',
                                                size: 'small'
                                            },
                                            style: {
                                                margin: '5px'
                                            },
                                            on: {
                                                click: () => {
                                                    this.onOpenGLUE(row);
                                                }
                                            }
                                        }, 'GLUE');
                            option.children.push(glueOptions)
                        }
                        let delOptions = h('Button', {
                            props: {
                                type: 'warning',
                                size: 'small'

                            },
                            style: {
                                marginRight: '5px'
                            },
                            on: {
                                click: () => {
                                    this.onRowDelete(row);
                                }
                            }
                        }, '删除');
                        option.children.push(delOptions);
                        let detailOptions = h('Button', {
                            props: {
                                type: 'warning',
                                size: 'small'

                            },
                            style: {
                                marginRight: '5px'
                            },
                            on: {
                                click: () => {
                                    this.onOpenDetail(row);
                                }
                            }
                        }, '详细');
                        option.children.push(detailOptions);
                        return option;
                    }
                }
            ],
            closable: false,
            footerHide: false,
            model: false,
            glueModel: false,
            modelMessage: '任务管理',
            taskForm: {
                ID: 0,
                TaskID: '',
                TargetType: '',
                TaskType: '',
                Express: '',
                DueTime: 0,
                TargetConfig: '',
                Remark: ''
            },
            glueTaskForm: {},
            defaultOption: {
                tabSize: 2,
                styleActiveLine: true,
                mode: 'shell',
                theme: 'monokai',
                lineNumbers: true,
                line: true,
                addModeClass: false,
                lineWrapping: true // 是否强制换行
            },
            httpTaskInfoForm: {
                Url: '',
                Method: '',
                PostBody: '',
                ContentType: '',
                Timeout: 0
            },
            shellConfigForm: {
                Type: '',
                Script: ''
            },
            goConfigForm: {
                FileName: ''
            },
            ruleValidate: {
                TaskID: [{ required: true, message: '任务编码必填', trigger: 'blur' }],
                Name: [{ required: true, message: '任务名称必填', trigger: 'blur' }],
                TargetType: [{ required: true, message: '任务类型必填', trigger: 'change' }],
                TaskType: [{ required: true, message: '任务执行类型必填', trigger: 'change' }],
                Express: [{ required: true, message: 'cron表达式必填', trigger: 'blur' }]
            }
        }
    },
    mounted() {
      this.init();
    },
    methods: {
        init() {
          this.onPageChange(this.queryParam)
        },
        onPageChange(param) {
          getTaskList(param).then(res => {
            if (res.RetCode === 0) {
              this.dataSource = res.Message;
            }
          })
        },
        onAdd() {
            // for (var key in this.taskForm) {
            //     this.taskForm[key] = '';
            // }
            // this.setFormClass(false);
            this.taskForm.TaskID = 0;
            this.taskForm.TaskID = '';
            this.taskForm.TargetType = '';
            this.taskForm.TaskType = '';
            this.taskForm.Express = '';
            this.taskForm.DueTime = 0;
            this.taskForm.Remark = '';
            this.model = true;
        },
        onRowEdit(row) {
            getTaskOnce({ID: row.ID}).then(res => {
                if (res.RetCode === 0) {
                    this.taskForm = res.Message;
                    if (this.taskForm.TargetType === 'http') {
                        this.httpTaskInfoForm = JSON.parse(this.taskForm.TargetConfig);
                    }
                    if (this.taskForm.TargetType === 'shell') {
                        this.shellConfigForm = JSON.parse(this.taskForm.TargetConfig);
                    }
                    if (this.taskForm.TargetType === 'goso') {
                        this.goConfigForm = JSON.parse(this.taskForm.TargetConfig);
                    }
                    this.model = true;
                } else {
                    this.$Message.error(res.RetMsg);
                }
            })
        },
        onRowDelete(row) {
            this.$Modal.confirm({
                title: '提示',
                content: '是否确认删除任务?',
                loading: true,
                onOk: () => {
                    taskDelete({ID: row.ID}).then(res => {
                        if (res.RetCode === 0) {
                            this.$Message.success('删除成功!');
                            this.init();
                            this.$Modal.remove();
                        } else {
                            this.$Message.error(res.RetMsg);
                        }
                    })
                }
            })
        },
        onOpenDetail(row) {
            this.$router.push({name: 'taskdetail', query: {id: row.ID}})
        },
        onShowLog(row) {
        },
        onOpenGLUE(row) {
            getTaskOnce({ID: row.ID}).then(res => {
                if (res.RetCode === 0) {
                    this.glueTaskForm = res.Message;
                    this.glueModel = true;
                } else {
                    this.$Message.error(res.RetMsg);
                }
            })
        },
        onRefresh() {
            this.init();
        },
        onSave() {
            this.$refs['formValidate'].validate((valid) => {
                if (valid) {
                    if (this.taskForm.TargetType === 'http') {
                      this.taskForm.TargetConfig = JSON.stringify(this.httpTaskInfoForm);
                    }
                    if (this.taskForm.TargetType === 'shell') {
                      this.taskForm.TargetConfig = JSON.stringify(this.shellConfigForm);
                    }
                    if (this.taskForm.TargetType === 'goso') {
                      this.taskForm.TargetConfig = JSON.stringify(this.goConfigForm);
                    }
                    if (this.taskForm.ID === 0) {
                        taskSave(this.taskForm).then(res => {
                            if (res.RetCode === 0) {
                                this.$Message.success('保存成功');
                                this.init();
                                this.model = false;
                            } else {
                                this.$Message.error(res.RetMsg);
                            }
                        })
                    } else {
                        taskUpdate(this.taskForm).then(res => {
                            if (res.RetCode === 0) {
                                this.$Message.success('保存成功');
                                this.init();
                                this.model = false;
                            } else {
                                this.$Message.error(res.RetMsg);
                            }
                        })
                    }
                } else {
                    return false;
                }
            })
        }
    }

}
</script>
<style lang="less">
.height-main{
    position: absolute;
    height: calc(100%);
    /* height: calc(100% - 104px); */
}
.model-content{
    padding-top: 5px;
    padding-left: 10px;
    padding-right: 50px;
}
.model-close{
    margin-top: 2px;;
}
</style>
