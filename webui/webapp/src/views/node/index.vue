<template>
    <div >
        <div class="tb">
      <tableH icon="md-apps" text="Node列表">
        <div slot="content"></div>
        <slot>
          <div style="text-align: right;float: right;">
              <div class="search" >
                  <!-- <Input v-model="queryParam.Name" placeholder="Node名称" style="width:160px;"/> -->
              </div>
              <div class="btn">
                   <Input v-model="queryParam.Name" placeholder="Node名称" style="width:160px;"/>
                <i-button type="info" icon="md-add" @click="onAdd">新建Node</i-button>
                <i-button type="info" icon="md-refresh" @click="onRefresh(false)">刷新</i-button>
              </div>
          </div>
        </slot>
      </tableH>
      <tableC
        id="table"
        :loading="loading"
        :columns="columns"
        :dataSource="dataSource"
        :queryParam="queryParam"
        @onPageChange="onPageChange"
        ref="table"
      ></tableC>
        </div>
      <Modal v-model="model"
            v-bind:title="modelMessage"
            width="660"
            v-bind:mask-closable="closable"
            v-bind:footer-hide="footerHide"
            @on-ok="onSave('formValidate')"
            class-name="vertical-center-modal">
            <div class="model-content">
                <i-form ref="formValidate" :label-width="120" :rules="ruleValidate" :model="formData">
                    <FormItem label="Node名称" prop='Name'>
                        <Input v-model="formData.Name" placeholder="Node名称"></Input>
                    </FormItem>
                    <FormItem label="Node服务器IP" prop='ServerIp'>
                        <Input v-model="formData.ServerIp" placeholder="Node服务器IP"></Input>
                    </FormItem>
                    <FormItem label="备注" prop='Remark'>
                        <Input v-model="formData.Remark" type="textarea" :autosize="{minRows: 2,maxRows: 5}" placeholder="备注"></Input>
                    </FormItem>
                </i-form>
            </div>
            <div slot="footer">
                <Button @click="model=false">取消</Button>
                <Button type="primary" @click="onSave('formValidate')">确定</Button>
            </div>
        </Modal>
    </div>
</template>
<script>
import Minix from '@/common/tableminix.js';
import tableC from '@/components/table/table.vue';
import tableH from '@/components/table/table-header.vue';
import { getNodeList, nodeSave, getNodeOnce, nodeDelete } from '@/api/node.js';
export default {
    components: { tableC, tableH },
    mixins: [Minix],
    data() {
        return {
            columns: [
                {
                    title: 'Node服务器Host',
                    key: 'Host'
                }, {
                    title: 'Node服务器Port',
                    key: 'Port'
                }, {
                    title: '状态',
                    key: 'IsMaster',
                    render: (h, params) => {
                            const row = params.row;
                            let str = ''
                            if (row.IsMaster) {
                                str += '主节点'
                            }
                            if (row.IsWorker) {
                                str += '工作中'
                            }
                            if (row.IsOnline) {
                                str += '在线'
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
                                            this.onRowEdit(params.row);
                                        }
                                    }

                                }, '修改'),
                                h('Button', {
                                    props: {
                                        size: 'small'
                                    },
                                    on: {
                                        click: () => {
                                            this.onRowDelete(params.row);
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
          this.onPageChange(this.queryParam)
        },
        onPageChange(param) {
          this.queryParam = param;
          if (!param.params) param.params = {};
          this.loading = true;
          getNodeList(param).then(res => {
            if (res.code === 200) {
              this.dataSource = res.data;
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
            nodeDelete({Id: row.Id}).then(res => {
                if (res.code === 200) {
                    this.$Message.success('删除成功!');
                    this.init();
                } else {
                    this.$Message.error(res.msg);
                }
            })
        },
        tiggerAction() {

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
