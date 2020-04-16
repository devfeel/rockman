export default {
    data() {
        return {
            queryParam: {
                PageIndex: 1,
                PageSize: 10
                // params: {}
            },
            dataSource: {
                colDefs: {
                    BodyFieldParams: []
                },
                pageData: [],
                totalCount: 0
            },

            ignoreColumns: [],
            detailColumns: [],

            tableHeight: 250
        }
    },
    methods: {
        onDataSourceChange(ds) {
            this.dataSource = {
                ColDefs: {
                    BodyFieldParams: []
                },
                Result: [],
                TotalCount: 0
            };
            this.$nextTick(_ => {
                this.dataSource = ds;
            });
        },
        onColumnVisible(col, detail) {
            if (detail) {
                if (!col.Visible) return false;
                if (!this.detailColumns || this.detailColumns.length === 0) return true;
                return this.detailColumns.includes(col.FieldName);
            } else {
                return col.Visible;
            }
        }
    }
}
