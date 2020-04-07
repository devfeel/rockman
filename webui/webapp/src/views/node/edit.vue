<template>
    <Form :model="formItem" :label-width="80">
        <FormItem label="Node名称">
            <Input v-model="formItem.Name" placeholder="Node名称"></Input>
        </FormItem>
        <FormItem label="Node服务器IP">
            <Input v-model="formItem.ServerIp" placeholder="Node服务器IP"></Input>
        </FormItem>
        <FormItem label="备注">
            <Input v-model="formItem.Remark" type="textarea" :autosize="{minRows: 2,maxRows: 5}" placeholder="备注"></Input>
        </FormItem>
        <FormItem>
            <Button type="primary">Submit</Button>
            <Button style="margin-left: 8px">Cancel</Button>
        </FormItem>
    </Form>
</template>
<script>
    export default {
        data () {
            return {
                formData: {
                    Id: '',
                    Name: '',
                    ServerIp: '',
                    Remark: ''
                }
            }
        }
    }
</script>
