<!--角色管理-->
<template>
  <div class="app-container roleDiv">
    <my-card title="角色管理">
      <div class="filter-container">
        <el-input v-model="listQuery.name" class="filter-item" style="width: 200px" placeholder="角色名称" clearable />
        <el-input v-model="listQuery.code" class="filter-item" style="width: 200px" placeholder="角色编码" clearable />
        <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="search">查询</el-button>
        <el-button v-waves class="filter-item" type="primary" @click="resetFilter">重置</el-button>
        <el-button class="filter-item addButton" type="primary" icon="el-icon-circle-plus-outline" @click="addRole">添加
        </el-button>
      </div>
      <table-list
        :data="list"
        :columns="columns"
        :list-loading="listLoading"
        class="dataTable"
        :total="total"
        :page-size="listQuery.pageSize"
        :show-index="true"
        :show-pagination="showPagination"
        @currentChange="currentChange"
      />
    </my-card>
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="roleForm" :rules="rules" :model="roleForm" label-position="center" size="small" label-width="100px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="角色名称" prop="name">
              <el-input v-model="roleForm.name" type="text" class="filter-item" placeholder="请输入角色名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色编码" prop="code">
              <el-input v-model="roleForm.code" type="text" class="filter-item" placeholder="请输入角色编码" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组织" prop="orgId">
              <el-select v-model="roleForm.orgId" placeholder="请选择组织">
                <el-option
                  v-for="item in orgList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色描述" prop="description">
              <el-input v-model="roleForm.description" type="text" class="filter-item" placeholder="请输入角色描述" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="loading" @click="saveRole">提交</el-button>
      </div>
    </el-dialog>

    <el-dialog title="配置权限" :visible.sync="authFormVisible">
      <choose-auth :data="authList" :org-privilges="rolePrivilges" @changeAuth="isPrivilegesChange = true" />
      <div slot="footer" class="dialog-footer">
        <el-button @click="authFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="loading" :disabled="!isPrivilegesChange" @click="configPrivilges">提交</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import TableList from '@/components/Table/TableList'
import MyCard from '@/components/MyCard'
import waves from '@/directive/waves'
import chooseAuth from '../orgManage/components/chooseAuth'
import { getRoleList, deleteRole, saveRole, configRolePrivileges, privilegesRole } from '@/api/sysManage/role'
import { getRootList, getOrgPrivileges } from '@/api/sysManage/organization'

export default {
  name: 'RoleManage',
  components: {
    TableList, chooseAuth, MyCard
  },
  directives: {
    waves
  },
  data() {
    return {
      dialogFormVisible: false, // 弹出框显示判断
      authFormVisible: false,
      isPrivilegesChange: false, // 判断是否修改权限
      dialogStatus: 'create',
      textMap: {
        update: '编辑',
        create: '新建'
      },
      loading: false, // 按钮重复提交判断
      listLoading: false,
      rules: {
        name: [{ required: true, message: '请填写角色名称' }],
        code: [{ required: true, message: '请填写角色编码' }],
        orgId: [{ required: true, message: '请选择组织', trigger: 'change' }]
      },
      roleForm: {
        id: '',
        name: '',
        code: '',
        orgId: '',
        description: ''
      },
      columns: [
        {
          text: '角色名称',
          value: 'name'
        },
        {
          text: '角色编码',
          value: 'code'
        },
        {
          text: '组织',
          value: 'orgName'
        },
        {
          text: '描述',
          value: 'description'
        },
        {
          text: '操作',
          type: 'iconButton',
          width: 150,
          list: this.operButton
        }
      ],
      listQuery: {
        pageSize: 10,
        currPage: 1,
        name: '',
        code: ''
      },
      total: 0,
      list: [],
      showPagination: true,
      currentRoleId: '', // 当前选择的角色id
      authList: [], // 所有权限
      orgList: [], // 可选组织
      rolePrivilges: [] // 角色配置的权限
    }
  },
  mounted() {
    // 获取所有角色
    this.getList()
    // 获取可选组织
    this.getRootList()
  },
  methods: {
    // 获取所有角色
    getList() {
      this.listLoading = true
      getRoleList(this.listQuery).then(response => {
        if (response.code === 200) {
          this.list = response.result.list
          this.total = response.result.total
        }
        this.listLoading = false
      })
    },
    // 获取可选根组织
    getRootList() {
      getRootList().then(response => {
        if (response.code === 200) {
          this.orgList = response.result
        }
      })
    },
    // 获取所有权限
    getAuthList(orgId) {
      getOrgPrivileges({ id: orgId }).then(response => {
        if (response.code === 200) {
          response.result.forEach(item => {
            if (item.parentId === 0) {
              delete item.parentId
            }
          })
          this.authList = response.result
        }
      })
    },
    search() {
      this.listQuery.currPage = 1
      this.getList()
      this.showPagination = false
      this.$nextTick(() => {
        this.showPagination = true
      })
    },
    // 清除所有搜索条件,并重新请求接口
    resetFilter() {
      this.listQuery.code = ''
      this.listQuery.name = ''
      this.listQuery.currPage = 1
      this.getList()
    },
    // 表格操作按鈕
    operButton() {
      return [
        { class: 'icon-xiugai', value: '编辑', click: this.updateRole },
        { class: 'icon-permisssion-management', value: '配置权限', click: this.setPrivilege },
        { class: 'icon-shanchu1', value: '删除', click: this.deleteRole }]
    },
    currentChange(val) {
      this.listQuery.currPage = val
      this.getList()
    },
    // 更新角色
    updateRole(val) {
      this.dialogStatus = 'update'
      this.roleForm.id = val.row.id
      this.roleForm.name = val.row.name
      this.roleForm.code = val.row.code
      this.roleForm.orgId = val.row.orgId
      this.roleForm.description = val.row.description
      this.dialogFormVisible = true
    },
    // 配置角色
    setPrivilege(val) {
      this.currentRoleId = val.row.id
      this.getAuthList(val.row.orgId)
      // 获取角色权限
      privilegesRole({ roleId: val.row.id }).then(response => {
        if (response.code === 200) {
          const tempArr = []
          if (response.result && response.result.length > 0) {
            response.result.forEach(function(item) {
              if (item.leaf) {
                tempArr.push(item.id)
              }
            })
          }
          this.rolePrivilges = tempArr
          this.authFormVisible = true
        } else {
          this.$message.error(response.message)
        }
      })
    },
    // 配置角色权限
    configPrivilges() {
      this.loading = true
      configRolePrivileges({ roleId: this.currentRoleId, privilegeIds: this.$store.getters.checkAuth }).then(response => {
        if (response.code === 200) {
          this.$message.success(response.message)
          this.authFormVisible = false
        } else {
          this.$message.error(response.message)
        }
        this.loading = false
      })
    },
    // 删除角色
    deleteRole(val) {
      this.$confirm('是否确定删除该记录?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteRole({ id: val.row.id }).then(response => {
          if (response.code === 200) {
            this.listQuery.currPage = 1
            this.getList()
            this.$message.success('删除成功')
          } else {
            this.$message.error(response.message)
          }
        })
      })
    },
    reset() { // 清空
      this.roleForm.id = ''
      this.$refs['roleForm'].resetFields()
    },
    // 新增角色
    addRole() {
      this.dialogStatus = 'create'
      this.roleForm = {
        id: '',
        name: '',
        code: '',
        orgId: '',
        description: ''
      }
      this.dialogFormVisible = true
    },
    // 保存角色
    saveRole() {
      this.$refs['roleForm'].validate((valid) => {
        if (valid) {
          this.loading = true
          if (this.roleForm.id === '') {
            delete this.roleForm.id
          }
          saveRole(this.roleForm).then(response => {
            if (response.code === 200) {
              this.$message.success('保存成功')
              // 操作数结构
              this.dialogFormVisible = false
              this.reset()
              this.getList()
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
  .roleDiv {
    .filter-container {
      .filter-item{
        margin-right: 10px;
      }
    }
  }
</style>
