<template>
    <div >
      <tableH icon="md-apps" text="日志列表">
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
      <tableC  id="table" :loading="loading"  :columns="columns" :dataSource="dataSource" :queryParam="queryParam"
        @onPageChange="onPageChange" ref="table"></tableC>
    </div>
</template>
<script>
import Minix from '@/common/tableminix.js';
import { dealDate } from '@/common/utils.js';
import tableC from '@/components/table/table.vue';
import tableH from '@/components/table/table-header.vue';
import { getLogList } from '@/api/logs.js';
export default {
  components: { tableC, tableH },
  mixins: [Minix],
  data() {
    return {
      columns: [
        {
          title: '任务Id',
          key: 'TaskId'
        }, {
          title: 'NodeId',
          key: 'NodeId'
        }, {
          title: '服务器',
          key: 'NodeEndPoint'
        }, {
          title: '是否执行成功',
          key: 'IsSuccess'
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
        }, {
          title: '操作',
          key: 'action',
          align: 'center',
          render: (h, params) => {
            return h('div', [
              h('Button', {
                props: {
                  type: 'warning',
                  size: 'small'

                },
                style: {
                  marginRight: '5px'
                },
                on: {
                  click: () => {

                  }
                }

              }, '查看')
            ]);
          }
        }
      ],
      model: false,
      loading: false,
      closable: false,
      footerHide: false

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
      this.queryParam = param;
      if (!param.params) param.params = {};
      this.loading = true;
      // this.queryParam.NodeId = '0be88880b0d945d3b4d55d75d4da0213'
      getLogList(param).then(res => {
        if (res.code === 200) {
          this.dataSource = res.data;
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
