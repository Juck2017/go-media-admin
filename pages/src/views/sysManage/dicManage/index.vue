<!--字典管理-->
<template>
  <div class="app-container dicDiv">
    <my-card title="字典管理">
      <el-button v-waves icon="el-icon-circle-plus-outline" class="filter-item addButton" type="primary" @click="addDic">添加</el-button>
      <tree-table :data="data" :columns="columns" :list-loading="listLoading" border />
    </my-card>

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dicForm" :rules="rules" :model="dicForm" label-position="center" size="small" label-width="100px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="字典名称" prop="name">
              <el-input v-model="dicForm.name" type="text" class="filter-item" placeholder="请输入字典名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="字典编码" prop="code">
              <el-input v-model="dicForm.code" type="text" class="filter-item" placeholder="请输入字典编码" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取消</el-button>
        <el-button v-waves :loading="loading" type="primary" @click="saveDic">提交</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
/* eslint-disable no-undef */
import TreeTable from '@/components/TreeTable'
import MyCard from '@/components/MyCard'
import waves from '@/directive/waves'
import { getDicList, saveDic, deleteDic } from '@/api/sysManage/dictionary'
export default {
  name: 'DicManage',
  components: {
    TreeTable,
    MyCard
  },
  directives: {
    waves
  },
  data() {
    return {
      dialogFormVisible: false, // 弹出框显示判断
      dialogStatus: 'create',
      textMap: {
        update: '编辑',
        create: '新建'
      },
      authList: [], // 字典数组
      loading: false, // 按钮重复提交判断
      listLoading: false,
      rules: {
        name: [{ required: true, message: '请填写字典名称' }],
        code: [{ required: true, message: '请填写字典编码' }]
      },
      dicForm: {
        id: '',
        parentId: '',
        name: '',
        code: ''
      },
      columns: [
        {
          text: '字典名称',
          value: 'name'
        },
        {
          text: '字典编码',
          value: 'code'
        },
        {
          text: '操作',
          type: 'iconButton',
          width: 120,
          list: this.operButton
        }
      ],
      data: []
    }
  },
  mounted() {
    this.getList()
  },
  methods: {
    // 获取所有字典
    getList: function() {
      this.listLoading = true
      getDicList().then(response => {
        if (response.code === 200) {
          this.data = response.result.reduce(function(prev, cur, index, arr) {
            if (cur.parentId === 0) {
              cur.parentId = undefined
            }
            prev.push(cur)
            return prev
          }, [])
        }
        this.listLoading = false
      })
    },
    // 表格操作按鈕
    operButton(val) {
      const temp = []
      if (!val.parentId) {
        temp.push({ class: 'icon-xinzeng', value: '添加', click: this.addChild })
      }
      temp.push({ class: 'icon-xiugai', value: '修改', click: this.updateDic })
      if (!val.children || val.children.length === 0) {
        temp.push({ class: 'icon-shanchu1', value: '删除', click: this.deleteDic })
      }
      return temp
    },
    // 更新字典
    updateDic(val) {
      this.dialogStatus = 'update'
      this.dicForm.id = val.row.id
      this.dicForm.parentId = val.row.parentId
      this.dicForm.name = val.row.name
      this.dicForm.code = val.row.code
      this.dialogFormVisible = true
    },
    // 添加子节点
    addChild(val) {
      this.dicForm = {
        id: '',
        parentId: val.row.id,
        name: '',
        code: ''
      }
      this.dialogFormVisible = true
    },
    // 删除字典
    deleteDic(val) {
      this.$confirm('是否确定删除该记录?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.loading = true
        deleteDic({ id: val.row.id }).then(response => {
          if (response.code === 200) {
            this.data.splice(this.data.findIndex(item => item.id === val.row.id), 1)
            this.$message.success('删除成功')
          } else {
            this.$message.error(response.message)
          }
          this.loading = false
        })
      })
    },
    reset() { // 清空
      this.dicForm.id = ''
      this.dicForm.parentId = ''
      this.$refs['dicForm'].resetFields()
    },
    // 新增字典
    addDic() {
      this.dialogStatus = 'create'
      this.dicForm = {
        id: '',
        parentId: '',
        name: '',
        code: ''
      }
      this.dialogFormVisible = true
    },
    setDgData(item, result) {
      item.name = result.name
      item.code = result.code
    },
    // 保存字典
    saveDic() {
      this.$refs['dicForm'].validate((valid) => {
        if (valid) {
          this.loading = true
          if (this.dicForm.id === '') {
            delete this.dicForm.id
          }
          if (this.dicForm.parentId === '') {
            delete this.dicForm.parentId
          }
          saveDic(this.dicForm).then(response => {
            if (response.code === 200) {
            // 操作数结构
              response.result.children = []
              if (this.dicForm.id) {
                // 说明是更新
                const findVal = this.data.find(item => {
                  return item.id === response.result.id
                })
                this.setDgData(findVal, response.result)
              } else {
                if (response.result.parentId === 0) {
                  delete response.result.parentId
                }
                this.data.push(response.result)
              }
              this.dialogFormVisible = false
              this.reset()
              this.$message.success('保存成功')
            } else {
              this.$message.error(response.message)
            }
            this.loading = false
          })
        }
      })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
  .dicDiv {
    .addButton {
      float: right;
      margin-right: 10px;
      margin-bottom: 10px;
    }
  }
</style>
