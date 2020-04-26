<template>
  <div class="content">
    <div class="form-header">
        <el-page-header @back="goBack" content="编辑任务">
            </el-page-header>
    </div>
    <div class="form">
        <el-form ref="form" :model="dataForm" label-width="160px" size="mini" :rules="rules"  v-if="show">
            <el-form-item label="任务编码" prop="TaskID">
                <el-input v-model="dataForm.TaskID" placeholder="任务编码" maxlength="30"></el-input>
            </el-form-item>
            <el-form-item label="任务类型" prop="TargetType">
                <el-select v-model="dataForm.TargetType" placeholder="请选择">
                    <el-option key="http" label="http" value="http"></el-option>
                    <el-option key="shell" label="shell" value="shell"></el-option>
                    <el-option key="goso" label="goso" value="goso"></el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="任务执行类型：" prop="TaskType">
                <el-select v-model="dataForm.TaskType" placeholder="请选择">
                    <el-option key="cron" label="cron" value="cron"></el-option>
                    <el-option key="loop" label="loop" value="loop"></el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="cron表达式" prop="TaskID">
                <el-input v-model="dataForm.Express" placeholder="cron表达式" maxlength="11"></el-input>
            </el-form-item>
            <el-form-item label="延迟时间" prop="DueTime">
                <el-input-number v-model="dataForm.DueTime" :min="1" label="延迟时间"></el-input-number>
            </el-form-item>
            <div name="http" v-if="dataForm.TargetType==='http'">
                <el-form-item label="请求地址" prop="HttpTaskInfoForm.Url">
                    <el-input v-model="dataForm.HttpTaskInfoForm.Url" placeholder="任务执行请求url" maxlength="100"> <template slot="prepend">Http://</template></el-input>
                </el-form-item>
                <el-form-item label="请求方式" prop="HttpTaskInfoForm.Method">
                    <el-select v-model="dataForm.HttpTaskInfoForm.Method" placeholder="请选择">
                        <el-option key="get" label="get" value="get"></el-option>
                        <el-option key="post" label="post" value="post"></el-option>
                        <el-option key="head" label="head" value="head"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="请求超时(s)" prop="HttpTaskInfoForm.Timeout">
                    <el-input-number v-model="dataForm.HttpTaskInfoForm.Timeout" :min="0" label="请求超时(s)"></el-input-number>
                </el-form-item>
            </div>
            <div name="shell" v-if="dataForm.TargetType==='shell'">
                <el-form-item label="类型" prop="ShellConfigForm.Type">
                    <el-select v-model="dataForm.ShellConfigForm.Type" placeholder="请选择">
                        <el-option key="script" label="Script Mode" value="script"></el-option>
                        <el-option key="file" label="File Mode" value="file"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="文件路径" prop="ShellConfigForm.Script" v-if="dataForm.ShellConfigForm.Type==='file'">
                    <el-input v-model="dataForm.ShellConfigForm.Script" placeholder="文件路径"></el-input>
                </el-form-item>
            </div>
            <div name="goso" v-if="dataForm.TargetType==='goso'">
                <el-form-item label="文件路径" prop="GoSoConfigForm.FileName">
                    <el-input v-model="dataForm.GoSoConfigForm.FileName" placeholder="文件路径"></el-input>
                </el-form-item>
            </div>
            <el-form-item label="备注" prop="Remark">
                <el-input type="textarea" placeholder="" v-model="dataForm.Remark" maxlength="100" show-word-limit></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="onSubmitForm('form')">提交</el-button>
                <el-button @click="onResetForm('form')">重置</el-button>
            </el-form-item>
        </el-form>
    </div>
  </div>
</template>
<script>
import { taskSave, taskUpdate, getTaskOnce } from '@/api/task.js';
export default {
  data() {
    return {
      show: true,
      dataForm: {
          ID: 0,
          TargetConfig: '',
          HttpTaskInfoForm: {},
          ShellConfigForm: {},
          GoSoConfigForm: {}
      },
      rules: {
          TaskID: [{ required: true, message: '请输入任务编码', trigger: 'blur' }],
          TargetType: [{ required: true, message: '请选择任务类型', trigger: 'blur' }],
          Express: [{ required: true, message: '请输入cron表达式', trigger: 'blur' }],
          DueTime: [{ required: true, message: '请输入延迟时间', trigger: 'blur' }],
          HttpTaskInfoForm: {
              Url: [{ required: true, message: '请输入任务执行请求url', trigger: 'blur' }],
              Method: [{ required: true, message: '请选择请求方式', trigger: 'blur' }]
          },
          ShellConfigForm: {
            Type: [{ required: true, message: '请选择类型', trigger: 'blur' }],
            Script: [{ required: true, message: '请输入文件路径', trigger: 'blur' }]
          },
          GoSoConfigForm: {
            FileName: [{ required: true, message: '请输入文件路径', trigger: 'blur' }]
          }

      }
    };
  },
  activated() {
    this.onInit();
  },
  methods: {
    onInit() {
        this.show = true;
        // this.onResetForm('form');

        if (this.$route.query.id) {
            getTaskOnce({ID: this.$route.query.id}).then(res => {
                if (res.RetCode === 0) {
                    this.dataForm = res.Message;
                    if (this.dataForm.TargetType === 'http') {
                        this.dataForm.HttpTaskInfoForm = JSON.parse(this.dataForm.TargetConfig);
                    }
                    if (this.dataForm.TargetType === 'shell') {
                        this.dataForm.ShellConfigForm = JSON.parse(this.dataForm.TargetConfig);
                    }
                    if (this.dataForm.TargetType === 'goso') {
                        this.dataForm.GoSoConfigForm = JSON.parse(this.dataForm.TargetConfig);
                    }
                } else {
                    this.$message.error(res.RetMsg);
                }
            })
        }
    },
    goBack() {
        this.show = false;
        this.$router.push({path: '/static/task'});
    },
    onSubmitForm(formName) {
        this.$refs[formName].validate((valid) => {
          if (valid) {
              if (this.dataForm.TargetType === 'http') {
                this.dataForm.TargetConfig = JSON.stringify(this.dataForm.HttpTaskInfoForm);
            }
            if (this.dataForm.TargetType === 'shell') {
                this.dataForm.TargetConfig = JSON.stringify(this.dataForm.ShellConfigForm);
            }
            if (this.dataForm.TargetType === 'goso') {
                this.dataForm.TargetConfig = JSON.stringify(this.dataForm.GoSoConfigForm);
            }
            debugger;
            if (this.dataForm.ID === 0) {
                taskSave(this.dataForm).then(res => {
                    if (res.RetCode === 0) {
                        this.$message.success('保存成功');
                        this.init();
                        this.model = false;
                    } else {
                        this.$message.error(res.RetMsg);
                    }
                })
            } else {
                taskUpdate(this.dataForm).then(res => {
                    if (res.RetCode === 0) {
                        this.$message.success('保存成功');
                        this.init();
                        this.model = false;
                    } else {
                        this.$message.error(res.RetMsg);
                    }
                })
            }
          } else {
            console.log('error submit!!');
            return false;
          }
        });
    },
    onResetForm(formName) {
        this.$refs[formName].resetFields();
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
.form-header{
    margin: 10px 20px;
    font-size: 16px;
}
.form{
    margin: 10px 20px;
    background-color:#fff;
    padding:10px 180px 50px 280px;
}
.el-select,.el-input,.el-textarea {
    width: 600px;
}
</style>
