<!--组织管理-->
<template>
  <div class="app-container orgDiv">
    <my-card title="组织管理">
      <el-button v-waves class="filter-item addButton" type="primary" icon="el-icon-circle-plus-outline" @click="addOrg">添加</el-button>
      <tree-table :data="data" :columns="columns" :list-loading="listLoading" border />
    </my-card>
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="orgForm" :rules="rules" :model="orgForm" label-position="center" size="small" label-width="100px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="组织名称" prop="name">
              <el-input v-model="orgForm.name" type="text" class="filter-item" placeholder="请输入组织名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组织编码" prop="code">
              <el-input v-model="orgForm.code" type="text" class="filter-item" placeholder="请输入组织编码" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组织类型" prop="type">
              <el-select v-model="orgForm.type" placeholder="请选择组织类型">
                <el-option
                  v-for="item in dicMap['org_type']"
                  :key="item.code"
                  :label="item.name"
                  :value="item.code"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取消</el-button>
        <el-button v-waves type="primary" :loading="loading" @click="saveOrg">提交</el-button>
      </div>
    </el-dialog>
    <el-dialog title="配置权限" :visible.sync="authFormVisible">
      <choose-auth :data="authList" :org-privilges="orgPrivilges" />
      <div slot="footer" class="dialog-footer">
        <el-button @click="authFormVisible = false">取消</el-button>
        <el-button v-waves type="primary" :loading="loading" @click="configPrivilges">提交</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import TreeTable from '@/components/TreeTable'
import waves from '@/directive/waves'
import chooseAuth from './components/chooseAuth'
import { getPrivilegeList } from '@/api/sysManage/privilege'
import { getOrgList, deleteOrg, saveOrg, configPrivileges, getOrgPrivileges } from '@/api/sysManage/organization'
import { mapGetters } from 'vuex'
import MyCard from '@/components/MyCard'
export default {
  name: 'OrgManage',
  components: {
    TreeTable, chooseAuth, MyCard
  },
  directives: {
    waves
  },
  data() {
    return {
      dialogFormVisible: false, // 弹出框显示判断
      authFormVisible: false, // 配置权限显示判断
      dialogStatus: 'create',
      textMap: {
        update: '编辑',
        create: '新建'
      },
      isChoose: false,
      orgPrivilges: [], // 组织的权限
      authList: [], // 权限数组
      loading: false, // 按钮重复提交判断
      listLoading: false,
      rules: {
        name: [{ required: true, message: '请填写组织名称' }],
        code: [{ required: true, message: '请填写组织编码' }],
        type: [{ required: true, message: '请选择组织类型', trigger: 'change' }]
      },
      orgForm: {
        id: '',
        parentId: '',
        name: '',
        type: '',
        code: ''
      },
      columns: [
        {
          text: '组织名称',
          value: 'name'
        },
        {
          text: '组织编码',
          value: 'code'
        },
        {
          text: '类型',
          value: 'type',
          type: 'text',
          formatter: this.typeFormatter
        },
        {
          text: '操作',
          type: 'iconButton',
          width: 230,
          list: this.operButton
        }
      ],
      data: [],
      currentOrgId: ''
    }
  },
  computed: {
    ...mapGetters([
      'dicMap'
    ])
  },
  mounted() {
    this.getList()
    this.getAuthList()
  },
  methods: {
    // 获取所有组织
    getList() {
      this.listLoading = true
      getOrgList().then(response => {
        if (response.code === 200) {
          // this.data = this.treeListUtil(response.result)
          response.result.forEach(item => {
            if (item.parentId === 0) {
              delete item.parentId
            }
          })
          this.data = response.result
        }
        this.listLoading = false
      })
    },
    // 获取所有权限
    getAuthList() {
      getPrivilegeList({ scope: 'system' }).then(response => {
        if (response.code === 200) {
          // 去掉parentId=0的字段
          response.result.forEach(item => {
            if (item.parentId === 0) {
              delete item.parentId
            }
          })
          this.authList = response.result
        }
      })
    },
    // 获取位置
    getLocation(val) {
      this.orgForm.location = val
    },
    typeFormatter(val) {
      return val.type === 'department' ? '部门' : '公司'
    },
    // 表格操作按鈕
    operButton(val) {
      const temp = [{ class: 'icon-xinzeng', value: '添加', click: this.addChild },
        { class: 'icon-xiugai', value: '修改', click: this.updateOrg }]
      if (!val.parentId) {
        temp.push({ class: 'icon-permisssion-management', value: '配置权限', click: this.setPrivilege })
      }
      if (!val.children || val.children.length === 0) {
        temp.push({ class: 'icon-shanchu1', value: '删除', click: this.deleteOrg })
      }
      return temp
    },
    // 更新组织
    updateOrg(val) {
      this.dialogStatus = 'update'
      this.orgForm.id = val.row.id
      this.orgForm.parentId = val.row.parentId
      this.orgForm.name = val.row.name
      this.orgForm.code = val.row.code
      this.orgForm.type = val.row.type
      if (val.row.location && val.row.location.address) {
        this.isChoose = true
      } else {
        this.isChoose = false
      }
      this.orgForm.location = val.row.location
      this.dialogFormVisible = true
    },
    // 添加子节点
    addChild(val) {
      this.dialogStatus = 'create'
      this.orgForm = {
        id: '',
        parentId: val.row.id,
        name: '',
        code: '',
        type: ''
      }
      this.dialogFormVisible = true
    },
    // 配置权限
    setPrivilege(val) {
      this.currentOrgId = val.row.id
      // 获取组织权限
      getOrgPrivileges({ id: val.row.id }).then(response => {
        if (response.code === 200) {
          const tempArr = []
          if (response.result && response.result.length > 0) {
            response.result.forEach(function(item) {
              if (item.leaf) {
                tempArr.push(item.id)
              }
            })
          }
          console.log('tempArr:', tempArr)
          this.orgPrivilges = tempArr
          this.authFormVisible = true
        } else {
          this.$message.error(response.message)
        }
      })
    },
    // 配置权限
    configPrivilges() {
      this.loading = true
      configPrivileges({ orgId: this.currentOrgId, privilegeIds: this.$store.getters.checkAuth }).then(response => {
        if (response.code === 200) {
          this.$message.success(response.message)
          this.authFormVisible = false
        } else {
          this.$message.error(response.message)
        }
        this.loading = false
      })
    },
    // 删除组织
    deleteOrg(val) {
      this.$confirm('是否确定删除该记录?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.loading = true
        deleteOrg({ id: val.row.id }).then(response => {
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
      this.orgForm.id = ''
      this.orgForm.parentId = ''
      this.$refs['orgForm'].resetFields()
    },
    // 新增组织
    addOrg() {
      this.dialogStatus = 'create'
      this.orgForm = {
        id: '',
        parentId: '',
        name: '',
        type: '',
        code: ''
      }
      this.dialogFormVisible = true
    },
    setDgData(item, result) {
      item.name = result.name
      item.code = result.code
      item.type = result.type
    },
    // 保存组织
    saveOrg() {
      this.$refs['orgForm'].validate((valid) => {
        if (valid) {
          this.loading = true
          if (this.orgForm.id === '') {
            delete this.orgForm.id
          }
          if (this.orgForm.parentId === '') {
            delete this.orgForm.parentId
          }
          saveOrg(this.orgForm).then(response => {
            if (response.code === 200) {
              this.$message.success('保存成功')
              // 操作数结构
              response.result.children = []
              // 判断是更新还是添加
              if (this.orgForm.id) {
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
  .orgDiv {
    .addButton {
      float: right;
      margin-right: 10px;
      margin-bottom: 10px;
    }
  }
</style>
