<template>
  <div>
    <div class="tb">
      <tableH  text="">
        <div slot="content"></div>
        <slot>
          <div style="text-align: right;float: right;">
            <div class="search">

            </div>
            <div class="btn">
              <i-button type="info" icon="md-refresh" @click="onRefresh(false)">刷新</i-button>
            </div>
          </div>
        </slot>
      </tableH>
      <tableC id="table" :loading="loading" :columns="columns" :dataSource="dataSource" :queryParam="queryParam"
              @onPageChange="onPageChange" ref="table"></tableC>
    </div>
  </div>
</template>
<script>
  import Minix from '@/common/tableminix.js';
  import { dealDate } from '@/common/utils.js';
  import tableC from '@/components/table/table.vue';
  import tableH from '@/components/table/table-header.vue';
  import { getTaskExecList } from '@/api/logs.js';
  export default {
    components: { tableC, tableH },
    mixins: [Minix],
    data() {
      return {
        columns: [
          {
            title: '任务编码',
            key: 'TaskID'
          }, {
            title: 'Node编码',
            key: 'NodeID',
            render: (h, params) => {
              return h('div', [
                h('span', {
                  style: {
                    display: 'inline-block',
                    width: '100%',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap'
                  },
                  domProps: {
                    title: params.row.NodeID
                  }
                }, params.row.NodeID)
              ]);
            }
          }, {
            title: '服务器信息',
            key: 'NodeEndPoint'
          }, {
            title: '执行开始时间',
            key: 'StartTime',
            render: (h, params) => {
              return h('div',
                dealDate(params.row.StartTime)
              )
            }
          }, {
            title: '执行结束时间',
            key: 'EndTime',
            render: (h, params) => {
              return h('div',
                dealDate(params.row.EndTime)
              )
            }
          }, {
            title: '是否执行成功',
            key: 'IsSuccess',
            render: (h, params) => {
              const row = params.row;
              if (row.IsSuccess) {
                return h('Span', '成功');
              }
              return h('Span', '失败');
            }
          }, {
            title: '执行失败类型',
            key: 'FailureType'
          }, {
            title: '执行失败原因',
            key: 'FailureCause'
          }, {
            title: '创建时间',
            key: 'CreateTime',
            render: (h, params) => {
              return h('div',
                dealDate(params.row.CreateTime)
              )
            }
          }
        ],
        loading: false
      }
    },
    props: {
      loadData: false
    },
    mounted() {
      // this.init();
    },
    watch: {
      loadData(newVal, oldVal) {
        if (newVal) {
          this.init();
        }
      }
    },
    methods: {
      init() {
        this.onPageChange(this.queryParam)
      },
      onPageChange(param) {
        this.queryParam = param;
        if (!param.params) param.params = {};
        this.loading = true;
        this.queryParam.TaskID = '';
        getTaskExecList(param).then(res => {
          if (res.RetCode === 0) {
            this.dataSource = res.Message;
          }
          this.loading = false;
        })
      },
      onRefresh() {
        this.init();
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
