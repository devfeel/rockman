<template>
  <div class="content">
  <div class="tb">
    <div class="tool">
      <el-button type="primary" @click="onAdd">添加</el-button>
    </div>
    <el-table :data="dataSource.PageData" border fit style="width: 100%">
      <el-table-column prop="TaskID" label="任务编码" >
        <template slot-scope="scope">
          <el-button type="text" @click="onRowClick(scope.row)" style="text-decoration:underline;">
                    {{ scope.row.TaskID }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="TargetType" label="任务类型" width="180"></el-table-column>
      <el-table-column prop="Express" label="间隔(Cron)" ></el-table-column>
      <el-table-column prop="IsRun" label="执行状态" width="180">
        <template slot-scope="scope">
          <span v-if="!scope.row.IsRun">未执行</span>
          <span v-if="scope.row.IsRun">已执行</span>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="操作" width="180">
        <template slot-scope="scope">
          <el-button @click="onGLUEClick(scope.row)" v-if="isGLUE(scope.row)" type="text" size="small">GLUE</el-button>
          <el-button @click="onRowDelete(scope.row)" type="text" size="small">删除</el-button>
          <el-button @click="onLogsClick(scope.row)" type="text" size="small">日志</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="page">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="queryParam.PageIndex"
        :page-sizes="[1,10, 30, 50, 100]"
        :page-size="queryParam.PageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="dataSource.TotalCount">
      </el-pagination>
    </div>
  </div>
  <el-dialog fullscreen :visible.sync="dialogVisible" :destroy-on-close="true" :show-close="false">
    <glue :data="dataForm" @close="onCloseDialog"></glue>
  </el-dialog>
  </div>
</template>
<script>
import Minix from '@/common/tableminix.js';
import { getTaskOnce, getTaskList, taskDelete } from '@/api/task.js';
import glue from './components/glue.vue';
export default {
  components: { glue },
  mixins: [Minix],
  data() {
    return {
      dataSource: [],
      dataForm: {},
      dialogVisible: false
    };
  },
  activated() {
    this.onInit();
  },
  methods: {
    onInit() {
      this.onPageChange(this.queryParam)
    },
    onPageChange(param) {
      getTaskList(param).then(res => {
        if (res.RetCode === 0) {
          this.dataSource = res.Message;
        }
      })
    },
    handleSizeChange(val) {
      this.queryParam.PageSize = val;
      this.onPageChange(this.queryParam)
    },
    handleCurrentChange(val) {
      this.queryParam.PageIndex = val;
      this.onPageChange(this.queryParam)
    },
    onGLUEClick(row) {
      getTaskOnce({ID: row.ID}).then(res => {
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
            this.dialogVisible = true;
        } else {
            this.$message.error(res.RetMsg);
        }
      })
    },
    onCloseDialog() {
      this.dialogVisible = false;
      // this.dataForm = {};
      this.onPageChange(this.queryParam)
    },
    onRowClick(row) {
      this.$router.push({path: '/static/task/detail', query: {id: row.ID}})
    },
    onAdd() {
      this.$router.push({path: '/static/task/detail'})
    },
    onRowDelete(row) {
            this.$confirm('是否确认删除任务?', '提示', {
                        confirmButtonText: '确定',
                        cancelButtonText: '取消',
                        type: 'warning'
                }).then(() => {
                    taskDelete({ID: row.ID}).then(res => {
                        if (res.RetCode === 0) {
                            this.$message.success('删除成功!');
                            this.onInit();
                        } else {
                            this.$Message.error(res.RetMsg);
                        }
                    })
                });
    },
    isGLUE(row) {
      if (row.TargetType === 'shell') {
        var shellConfigForm = JSON.parse(row.TargetConfig);
        if (shellConfigForm.Type === 'script') {
          return true;
        }
      }
      return false;
    },
    onLogsClick(row) {
      this.$router.push({path: '/static/task/logdetail', query: {id: row.ID}})
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
.tb{
  margin: 10px 20px;
  background-color:#fff;
  padding:0px 16px 0px 16px;
}
.tool{
  //height: 45px;
  text-align: right;
  padding: 10px 0px 10px 0px;
}
.page{
  padding: 16px 16px;
  text-align: right;
}
</style>
