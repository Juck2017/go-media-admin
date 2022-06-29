<!--权限管理-->
<template>
  <div class="app-container authorityDiv">
    <my-card title="权限管理">
      <el-button v-waves icon="el-icon-circle-plus-outline" class="filter-item addButton" type="primary" @click="addPrivilege">添加</el-button>
      <tree-table :data="data" :columns="columns" :list-loading="listLoading" border />
    </my-card>
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="authForm" :rules="rules" :model="authForm" label-position="center" size="small" label-width="100px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="权限名称" prop="name">
              <el-input v-model="authForm.name" type="text" class="filter-item" placeholder="请输入权限名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限编码" prop="code">
              <el-input v-model="authForm.code" type="text" class="filter-item" placeholder="请输入权限编码" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限类型" prop="type">
              <el-select v-model="authForm.type" placeholder="请选择权限类型">
                <el-option v-for="item in dicMap['privilege_type']" :key="item.code" :label="item.name" :value="item.code" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限描述" prop="description">
              <el-input v-model="authForm.description" type="text" class="filter-item" placeholder="请输入权限描述" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取消</el-button>
        <el-button v-waves :loading="loading" type="primary" @click="savePrivilege">提交</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import TreeTable from '@/components/TreeTable'
import waves from '@/directive/waves'
import { getPrivilegeList, deletePrivilege, savePrivilege } from '@/api/sysManage/privilege'
import { mapGetters } from 'vuex'
import MyCard from '@/components/MyCard'

export default {
  name: 'PrivilegeManage',
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
      loading: false, // 按钮重复提交判断
      listLoading: false,
      rules: {
        name: [{ required: true, message: '请填写权限名称' }],
        code: [{ required: true, message: '请填写权限编码' }],
        type: [{ required: true, message: '请选择权限类型', trigger: 'change' }]
      },
      authForm: {
        id: '',
        parentId: '',
        name: '',
        code: '',
        type: '',
        description: ''
      },
      columns: [
        {
          text: '权限名称',
          value: 'name'
        },
        {
          text: '权限编码',
          value: 'code',
          width: 200
        },
        {
          text: '类型',
          value: 'type',
          type: 'text',
          formatter: this.typeFormatter
        },
        {
          text: '描述',
          value: 'description'
        },
        {
          text: '操作',
          type: 'iconButton',
          width: 130,
          list: this.operButton
        }
      ],
      data: []
    }
  },
  computed: {
    ...mapGetters([
      'dicMap'
    ])
  },
  mounted() {
    this.getPrivilegeList()
    console.log('this.dicMap:', this.dicMap)
  },
  methods: {
    // 获取所有权限
    getPrivilegeList() {
      this.listLoading = true
      getPrivilegeList().then(response => {
        if (response.code === 200) {
          this.data = response.result.filter(item => {
            if (item.parentId === 0) {
              delete item.parentId
            }
            return item
          })
          // this.data = this.treeListUtil(response.result)
          // this.data = response.result
        }
        this.listLoading = false
      })
    },
    typeFormatter(val) {
      return val.type === 'menu' ? '菜单' : '按钮'
    },
    // 表格操作按鈕
    operButton(val) {
      if (val.children && val.children.length > 0) {
        return [
          { class: 'icon-xinzeng', value: '添加', click: this.addChild },
          { class: 'icon-xiugai', value: '编辑', click: this.updateAuth }
        ]
      } else {
        return [
          { class: 'icon-xinzeng', value: '添加', click: this.addChild },
          { class: 'icon-xiugai', value: '编辑', click: this.updateAuth },
          { class: 'icon-shanchu1', value: '删除', click: this.deleteAuth }
        ]
      }
    },
    // 更新权限
    updateAuth(val) {
      this.dialogStatus = 'update'
      this.authForm.id = val.row.id
      this.authForm.parentId = val.row.parentId
      this.authForm.name = val.row.name
      this.authForm.code = val.row.code
      this.authForm.type = val.row.type
      this.authForm.description = val.row.description
      this.authForm.definition = val.row.definition
      this.dialogFormVisible = true
    },
    // 添加子节点
    addChild(val) {
      this.authForm = {
        id: '',
        parentId: val.row.id,
        name: '',
        code: '',
        type: '',
        description: '',
        definition: ''
      }
      this.dialogFormVisible = true
    },
    // 删除权限
    deleteAuth(val) {
      this.$confirm('是否确定删除该记录?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deletePrivilege({ id: val.row.id }).then(response => {
          if (response.code === 200) {
            this.data.splice(this.data.findIndex(item => item.id === val.row.id), 1)
            this.$message.success('删除成功')
          } else {
            this.$message.error(response.message)
          }
        })
      })
    },
    reset() { // 清空
      this.authForm.id = ''
      this.authForm.parentId = ''
      this.$refs['authForm'].resetFields()
    },
    // 打开新增权限页面
    addPrivilege() {
      this.authForm = {
        id: '',
        parentId: '',
        name: '',
        code: '',
        type: '',
        description: ''
      }
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
    },
    setDgData(item, result) {
      item.name = result.name
      item.code = result.code
      item.type = result.type
      item.description = result.description
      item.definition = result.definition
    },
    // 保存权限
    savePrivilege() {
      this.$refs['authForm'].validate((valid) => {
        if (valid) {
          this.loading = true
          this.authForm.scope = 'system'
          if (!this.authForm.id) {
            delete this.authForm.id
          }
          if (!this.authForm.parentId) {
            delete this.authForm.parentId
          }
          savePrivilege(this.authForm).then(response => {
            if (response.code === 200) {
              // 操作数结构
              response.result.children = []
              // 判断是更新还是添加
              if (this.authForm.id) {
                const findVal = this.data.find(item => {
                  return item.id === response.result.id
                })
                this.setDgData(findVal, response.result)
              } else {
                // this.saveUtil(this.authForm.parentId, response.result)
                if (response.result.parentId === 0) {
                  delete response.result.parentId
                }
                this.data.push(response.result)
              }
              this.$message.success('保存成功')
              this.dialogFormVisible = false
              this.reset()
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
  .authorityDiv {
    .addButton {
      float: right;
      margin-right: 10px;
      margin-bottom: 10px;
    }
  }

</style>
