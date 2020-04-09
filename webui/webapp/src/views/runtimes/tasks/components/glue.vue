<template>
<div>
    <div slot="header" class="glue-header" >
        <div class="glue-title">
            WebIDE<span>{{defaultOption.mode}}脚本编辑</span>
        </div>
        <div class="glue-btn">
            <Button type="primary" @click="onSave()">保 存</Button>
            <Button>取 消</Button>
        </div>
    </div>
    <div class="glue-conext">
        <codemirror
            ref="mycode"
            v-model="shellConfigForm.Script"
            :options="defaultOption">
        </codemirror>
    </div>
</div>
</template>
<script>
import { codemirror } from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/mode/shell/shell'
import 'codemirror/theme/duotone-light.css'
import { taskSave } from '@/api/task.js';
export default {
    components: { codemirror },
    data() {
        return {
            defaultOption: {
                tabSize: 2,
                styleActiveLine: true,
                mode: 'shell',
                theme: 'duotone-light',
                lineNumbers: true,
                line: true,
                addModeClass: false,
                lineWrapping: true // 是否强制换行
            },
            shellConfigForm: {
                Type: '',
                Script: ''
            }
        }
    },
    props: {
      data: {}
    },
    mounted() {
        // this.init();
    },
    watch: {
            data(newVal, oldVal) {
                // 执行数据更新查询
                this.init();
            }
    },
    methods: {
        init() {
            this.shellConfigForm = JSON.parse(this.data.TargetConfig);
        },
        onSave() {
             this.data.TargetConfig = JSON.stringify(this.shellConfigForm);
            taskSave(this.data).then(res => {
                if (res.RetCode === 0) {
                    this.$Message.success('保存成功');
                } else {
                    this.$Message.error(res.RetMsg);
                }
            })
        }
    }
}
</script>
<style lang="less">
.ivu-modal-body{
    padding: 0px;
}
.glue-header{
    padding: 10px;
    background: rgb(81, 90, 110);
    float: left;
    width: calc(100%);
    height: 55px;
}
.glue-title{
    float:left;
    color: white;
    vertical-align: middle;
    text-align: left;
    font-size: 18px;
    line-height: 35px;
}
.glue-title span{
    padding-left: 20px;
    font-size: 14px;
}

.glue-btn{
    float:right;
    text-align: right;
    padding-right: 50px;
}
.glue-conext{
    clear: both;
}
.CodeMirror {
    border: 1px solid #eee;
    height: auto;
}

.CodeMirror-scroll {
    min-height: 500px;
    height: auto;
    overflow-y: hidden;
    overflow-x: auto;
}
</style>
