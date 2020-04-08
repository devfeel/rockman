<template>
    <div>
        <i-table :columns="columns" :loading="loading" :max-height="maxHeight" :data="dataSource.PageData" class="v-table" border>
        </i-table>
        <div style="margin: 10px;overflow: hidden">
            <div style="float: right;">
                <Page show-total show-elevator show-sizer
                :total="totalCount"
                :current="current"
                :page-size-opts="sizes"
                @on-change="onPageChange"
                @on-page-size-change="onSizeChange"></Page>
            </div>
        </div>
    </div>
</template>
<script>
    export default {
        props: {
            dataSource: {
                // 表数据源,配置了url就不用传这个参数了
                type: Object,
                default: () => {}
            },
            queryParam: {
                type: Object,
                default: function () {
                    return {}
                }
            },
            columns: {
                type: Array,
                default: () => []
                // [ {
                //   field: "columnType",
                //   title: "数据类型",
                //   width: 120,
                //   hidden:false,
                //   edit: { type: "text", status: false, data: [], key: "" }
                // }] //列的的数据格式edit格式： type类型(text,date,datetime,select,switch),status是否默认为编辑状态
                // data如果是select这里data应该有数据源，如果没有数据请设置key字典编号
            },
            url: {
                type: String,
                default: ''
            },
            loading: {
                type: Boolean,
                default: false
            },
            paginationHide: {
              type: Boolean,
              default: true
            }
        },
        data () {
            return {
                colDefs: {},
                totalCount: 0,
                current: 1,
                maxHeight: 0,
                sizes: [10, 20, 30, 50, 100]
            }
        },
        watch: {
            dataSource: {
                handler: function () {
                    this.colDefs = this.dataSource.ColDefs;
                    this.totalCount = this.dataSource.TotalCount;
                },
                deep: true
            }
        },
        mounted() {
        },
        methods: {
            onSizeChange(val) {
                this.queryParam.pageSize = val;
                this.$emit('onPageChange', this.queryParam);
            },
            onPageChange(val) {
                this.queryParam.pageIndex = val;
                this.$emit('onPageChange', this.queryParam);
            }
        }
    }
</script>
<style scoped>

</style>
