<template>
  <div class="content">
    <div class="tb">
      <el-table
        :data="tableData"
        border
        style="width: 100%">
        <el-table-column type="index" width="50"></el-table-column>
        <el-table-column prop="name" label="名称" ></el-table-column>
        <el-table-column prop="ipaddress" label="地址" ></el-table-column>
        <el-table-column prop="state" label="状态" width="180"></el-table-column>
        <el-table-column prop="taskNum" label="运行任务" width="180"></el-table-column>
      </el-table>
    </div>
  </div>
</template>
<script>
import { getNodeList } from '@/api/node.js';
export default {
  data() {
    return {
      tableData: []
    }
  },
  activated() {
    // this.onInit();
  },
  methods: {
    onInit() {
      getNodeList().then(res => {
        if (res.RetCode === 0) {
          for (var index = 0; index < res.Message.length; index++) {
            var row = res.Message[index];
            console.log(row)
          }
          // this.tableData = res.Message;
        } else {
          this.$message.warning(res.RetMsg);
        }
      })
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
}
</style>
