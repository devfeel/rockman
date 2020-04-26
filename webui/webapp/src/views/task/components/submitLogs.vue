<template>
  <div>
    <div class="tb">
        <el-table :data="dataSource.PageData" border fit  style="width: 100%">
            <el-table-column prop="TaskID" label="任务编码" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="NodeID" label="节点编码" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="NodeEndPoint" label="服务器信息" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="IsSuccess" label="执行状态" width="180" :show-overflow-tooltip="true">
                <template slot-scope="scope">
                    <span v-if="!scope.row.IsSuccess">成功</span>
                    <span v-if="scope.row.IsSuccess">失败</span>
                </template>
            </el-table-column>
            <el-table-column prop="FailureType" label="失败类型" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="FailureCause" label="失败原因" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="CreateTime" :formatter="formatDate" label="创建时间" :show-overflow-tooltip="true"></el-table-column>
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
  </div>
</template>
<script>
  import Minix from '@/common/tableminix.js';
  import { dealDate } from '@/common/utils.js';
  import { getTaskSubmitList } from '@/api/logs.js';
  export default {
    mixins: [Minix],
    data() {
      return {

        loading: false
      }
    },
    props: {
      TaskID: null,
      loadData: false
    },
    mounted() {
      this.init();
    },
    methods: {
      init() {
        this.onPageChange(this.queryParam)
      },
      onPageChange(param) {
        this.queryParam = param;
        this.loading = true;
        this.queryParam.TaskID = this.TaskID;
        getTaskSubmitList(param).then(res => {
          if (res.RetCode === 0) {
            this.dataSource = res.Message;
          }
          this.loading = false;
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
      formatDate(row, column) {
        // 获取单元格数据
        let data = row[column.property]
        if (!data) {
            return ''
        }
        return dealDate(data)
      },
      onRefresh() {
        this.init();
      }
    }
  }
</script>
<style lang="less" scoped>
.page{
  padding: 16px 16px;
  text-align: right;
}
</style>
