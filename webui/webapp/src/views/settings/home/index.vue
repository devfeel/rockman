<template>
  <div class="content-main">
    <tableH icon="md-apps" text="节点列表">
        <div slot="content"></div>
        <slot>
        <div style="text-align: right;float: right;">
            <div class="search" >
                <!-- <Input v-model="queryParam.Name" placeholder="Node名称" style="width:160px;"/> -->
            </div>
            <div class="btn">
                <!-- <Input v-model="queryParam.Name" placeholder="Node名称" style="width:160px;"/> -->
                <!-- <i-button type="info" icon="md-add" @click="onAdd">新建Node</i-button> -->
                <!-- <i-button type="info" icon="md-refresh" @click="onRefresh(false)">刷新</i-button> -->
            </div>
        </div>
        </slot>
    </tableH>
    <i-table ref="table" :columns="columns" :loading="loading" :data="NodeLists" border>
    </i-table>
  </div>
</template>
<script>
import Minix from '@/common/tableminix.js';
import tableC from '@/components/table/table.vue';
import tableH from '@/components/table/table-header.vue';
import { getNodeList, nodeSave, getNodeOnce } from '@/api/node.js';
export default {
    components: { tableC, tableH },
    mixins: [Minix],
    data() {
        return {
            columns: [
                {
                    title: 'Node编码',
                    key: 'NodeId'
                }, {
                    title: 'Node名称',
                    key: 'NodeName'
                }, {
                    title: '是否Leader',
                    key: 'isLeader',
                    render: (h, params) => {
                            const row = params.row;
                            let str = ''
                            if (row.isLeader) {
                                str += 'Leader'
                            }
                            return h('Span', str);
                    }
                }, {
                    title: '状态',
                    key: 'Status',
                    render: (h, params) => {
                            const row = params.row;
                            let str = ''
                            if (row.isLeader) {
                                str += 'Leader'
                            }
                            return h('Span', str);
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
                                            this.$Message.success('暂未实现!');
                                        }
                                    }

                                }, '下线&上线'),
                                h('Button', {
                                    props: {
                                        size: 'small'
                                    },
                                    on: {
                                        click: () => {
                                            this.$Message.success('暂未实现!');
                                        }
                                    }
                                }, '删除')
                            ]);
                        }
                }
            ],
            model: false,
            loading: false,
            modelMessage: 'Node编辑',
            closable: false,
            footerHide: false,
            NodeLists: [],
            formData: {
                Name: '',
                ServerIp: '',
                Remark: ''
            },
            ruleValidate: {
                Name: [{ required: true, message: 'Node名称必填', trigger: 'blur' }],
                ServerIp: [{ required: true, message: '服务器IP必填', trigger: 'blur' }]
            }
        }
    },
    mounted() {
      this.init();
    },
    methods: {
        init() {
          this.loading = true;
          getNodeList().then(res => {
            if (res.RetCode === 0) {
                this.NodeLists = res.data;
            } else {
                this.$Message.error(res.RetMsg);
            }
            this.loading = false;
          })
        },
        onAdd() {
            this.$refs['formValidate'].resetFields();
            this.model = true;
        },
        onSave(name) {
            this.$refs[name].validate((valid) => {
                if (valid) {
                    nodeSave(this.formData).then(res => {
                        if (res.code === 200) {
                            this.$Message.success('保存成功');
                            this.init();
                            this.model = false;
                        } else {
                            this.$Message.error(res.msg);
                        }
                    })
                } else {
                    return false;
                }
            })
        },
        onEdit() {

        },
        onRowEdit(row) {
            getNodeOnce({Id: row.Id}).then(res => {
                if (res.code === 200) {
                    this.formData = res.data;
                    this.model = true;
                } else {
                    this.$Message.error(res.msg);
                }
            })
        },
        onRowDelete(row) {
            this.$Message.success('暂未实现!');
        },
        tiggerAction() {

        },
        onRefresh() {
            this.init();
        }
    }

}
</script>
<style scoped>
.content-main{
    position: absolute;
    left: 0;
    right: 0;
    padding: 10px 30px 10px 30px;
    height: calc(100%);
    /* background:#f5f7f9; */
    /* height: calc(100% - 104px); */
}

</style>
