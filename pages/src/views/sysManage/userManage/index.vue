<!--人员管理-->
<template>
  <div class="app-container userDiv">
    <my-card title="用户管理">
      <div class="filter-container">
        <el-input v-model="listQuery.name" class="filter-item" style="width: 200px" placeholder="姓名" clearable />
        <el-input v-model="listQuery.mobile" class="filter-item" style="width: 200px" placeholder="手机号" clearable />
        <div class="filter-item" style="width:200px">
          <treeselect
            v-model="listQuery.orgId"
            :options="orgList"
            :normalizer="normalizer"
            placeholder="请选择组织"
            no-children-text="无选择"
          />
        </div>
        <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="search">查询</el-button>
        <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="resetFilter">重置</el-button>
        <el-button v-waves class="filter-item addButton" type="primary" icon="el-icon-circle-plus-outline" @click="addUser">添加
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
        @currentChange="currentChange"
      />
    </my-card>

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="userForm" :rules="rules" :model="userForm" label-position="center" size="small" label-width="100px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="userForm.username" :disabled="userForm.id !== ''" type="text" class="filter-item" placeholder="请输入用户名" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="userForm.name" type="text" class="filter-item" placeholder="请输入姓名" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="dialogStatus == 'create'">
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input v-model="userForm.password" type="password" class="filter-item" placeholder="请输入密码" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirePassword">
              <el-input v-model="userForm.confirePassword" type="password" class="filter-item" placeholder="请再次输入密码" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="组织" prop="orgId">
              <treeselect
                v-model="userForm.orgId"
                :options="orgList"
                :normalizer="normalizer"
                placeholder="请选择组织"
                no-children-text="无选择"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="mobile">
              <el-input v-model="userForm.mobile" type="text" class="filter-item" placeholder="请输入手机号" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row />
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="cancelSaveUser">取消</el-button>
        <el-button v-waves type="primary" :loading="loading" @click="saveUser">提交</el-button>
      </div>
    </el-dialog>
    <el-dialog title="配置权限" :visible.sync="authFormVisible">
      <choose-auth :data="authList" :org-privilges="userPrivilges" @changeAuth="isPrivilegesChange = true" />
      <div slot="footer" class="dialog-footer">
        <el-button @click="authFormVisible = false">取消</el-button>
        <el-button v-waves type="primary" :loading="loading" :disabled="!isPrivilegesChange" @click="configPrivilges">提交</el-button>
      </div>
    </el-dialog>

    <el-dialog title="配置角色" :visible.sync="roleFormVisible">
      <el-form
        ref="roleForm"
        :model="roleForm"
        label-position="center"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="角色选择" prop="roleIds">
              <el-select v-model="roleForm.roleIds" multiple placeholder="请选择角色">
                <el-option
                  v-for="item in roleList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="roleFormVisible = false">取消</el-button>
        <el-button v-waves type="primary" :loading="loading" @click="configRoles">提交</el-button>
      </div>
    </el-dialog>

    <el-dialog title="密码重置" :visible.sync="passwordFormVisible">
      <el-form
        ref="passwordForm"
        :rules="passwordRules"
        :model="passwordForm"
        label-position="center"
        size="small"
        label-width="100px"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="重置密码" prop="password">
              <el-input v-model="passwordForm.password" type="password" class="filter-item" placeholder="请输入重置密码" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="passwordFormVisible = false">取消</el-button>
        <el-button v-waves type="primary" :loading="loading" @click="passwordResetSubimt">提交</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import tableList from '@/components/Table/TableList'
import MyCard from '@/components/MyCard'
import waves from '@/directive/waves'
import chooseAuth from '../orgManage/components/chooseAuth'
import {
  getUserList,
  deleteUser,
  saveUser,
  privilegesList,
  configUserPrivileges,
  enabledUser,
  resetPassword,
  getRoleList,
  configRoles
} from '@/api/sysManage/user'
import { getOrgList, getOrgPrivileges, getOrgRole } from '@/api/sysManage/organization'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'
import { validateMobile } from '@/utils/validate'
// import { parseTime } from '@/utils'

export default {
  name: 'UserManage',
  components: {
    tableList, chooseAuth,
    Treeselect,
    MyCard
  },
  directives: {
    waves
  },
  data() {
    const that = this
    // 校验确认密码,是否与原密码一致
    function validateConfirePassword(rule, confirePassword, callback) {
      if (!that.userForm.password) {
        return callback(new Error('请先输入密码'))
      }
      if (that.userForm.password !== confirePassword) {
        return callback(new Error('两次输入密码不一致'))
      }
      return callback()
    }
    // 校验用户名,判断是否存在,新增时校验,编辑时不检验
    function validateUsername(rule, username, callback) {
      if (!username) {
        return callback(new Error('请输入用户名'))
      }
      const exist = that.existUserList.find(item => {
        return item === username
      })
      if (exist && !that.userForm.id) {
        return callback(new Error('用户已存在'))
      }
      return callback()
    }
    return {
      dialogFormVisible: false, // 弹出框显示判断
      authFormVisible: false,
      passwordFormVisible: false,
      roleFormVisible: false,
      isPrivilegesChange: false, // 判断是否修改权限
      dialogStatus: 'create',
      textMap: {
        update: '编辑',
        create: '新建'
      },
      loading: false, // 按钮重复提交判断
      listLoading: false,
      rules: {
        username: [{ required: true, validator: validateUsername }],
        mobile: [{ required: true, validator: validateMobile, ruleType: 'phone' }],
        password: [{ required: true, message: '请填写密码', trigger: 'blur' }],
        confirePassword: [{ required: true, validator: validateConfirePassword }],
        name: [{ required: true, message: '请填写姓名', trigger: 'blur' }],
        orgId: [{ required: true, message: '请选择组织', trigger: 'blur' }]
      },
      userForm: {
        id: '',
        username: '',
        name: '',
        mobile: '',
        orgId: null,
        password: '',
        confirePassword: '',
        enabled: false
      },
      passwordForm: {
        password: ''
      },
      passwordRules: {
        password: [{ required: true, message: '请填写重置密码' }]
      },
      roleForm: {
        roleIds: ''
      },
      roleList: [],
      columns: [
        {
          text: '姓名',
          value: 'name'
        },
        {
          text: '手机号',
          value: 'mobile'
        },
        {
          text: '组织',
          value: 'orgName'
        },
        {
          text: '用户名',
          value: 'username'
        },
        {
          text: '创建时间',
          value: 'createdAt'
        },
        {
          text: '状态',
          type: 'switch',
          activeText: '启用',
          inactiveText: '禁用',
          activeValue: true,
          inactiveValue: false,
          model: 'enabled',
          change: this.changeStatus
        },
        {
          text: '操作',
          type: 'iconButton',
          width: 180,
          list: this.operButton
        }
      ],
      listQuery: {
        pageSize: 10,
        currPage: 1,
        name: ''
      },
      total: 0,
      list: [],
      currentUserId: '', // 当前选择的用户id
      authList: [], // 所有权限
      orgList: [], // 可选组织
      userPrivilges: [], // 用户配置的权限
      normalizer(node) {
        return {
          id: node.id,
          label: node.name,
          children: node.children
        }
      },
      existUserList: [] // 已注册用户列表
    }
  },
  mounted() {
    this.getList()
    this.getRootList()
  },
  methods: {
    // 配置角色
    configRoles() {
      this.loading = true
      configRoles({ userId: this.currentUserId, roleIds: this.roleForm.roleIds }).then(response => {
        if (response.success) {
          this.$message.success(response.message)
          this.roleFormVisible = false
          this.$refs['roleForm'].resetFields()
        } else {
          this.$message.error(response.message)
        }
        this.loading = false
      })
    },
    // 修改状态
    changeStatus(val, checked) {
      // const temp = checked ? '启用' : '禁用'
      // this.$confirm('是否确定' + temp + '该用户?', '提示', {
      //   confirmButtonText: '确定',
      //   cancelButtonText: '取消',
      //   type: 'warning'
      // }).then(() => {
      enabledUser({ id: val.row.id, enabled: checked }).then(response => {
        if (response.success) {
          this.getList()
          this.$message.success(response.message)
        } else {
          this.$message.error(response.message)
        }
      })
      // })
    },
    // 获取所有用户
    getList() {
      this.listLoading = true
      if (!this.listQuery.name) {
        delete this.listQuery.name
      }
      if (!this.listQuery.mobile) {
        delete this.listQuery.mobile
      }
      getUserList(this.listQuery).then(response => {
        if (response.success === true) {
          this.list = response.result.list
          this.total = response.result.total
          this.addExistUser()
        } else {
          this.list = []
        }
        this.listLoading = false
      })
    },
    // 获取可选组织
    getRootList() {
      getOrgList({}).then(response => {
        if (response.code === 200) {
          // 组装参数
          this.orgList = this.treeListUtil(response.result)
        }
      })
    },
    treeListUtil(data, parentId) {
      const itemArr = []
      for (let i = 0; i < data.length; i++) {
        const node = data[i]
        if (node.parentId === parentId) {
          node.children = this.treeListUtil(data, node.id)
          if (node.children.length === 0) {
            delete node.children
          }
          itemArr.push(node)
        }
      }
      return itemArr
    },
    // 清除所有搜索条件,并重新请求接口
    resetFilter() {
      this.listQuery.name = ''
      this.listQuery.mobile = ''
      delete this.listQuery.orgId
      this.listQuery.currPage = 1
      this.getList()
    },
    search() {
      this.getList()
    },
    // 表格操作按鈕
    operButton() {
      return [
        { class: 'icon-xiugai', value: '编辑', click: this.updateUser, privilege: 'B_yhgl_bj' },
        { class: 'icon-jiaoseguanli', value: '配置角色', click: this.setRole, privilege: 'B_yhgl_pzjs' },
        { class: 'icon-permisssion-management', value: '配置权限', click: this.setPrivilege, privilege: 'B_yhgl_pzqx' },
        { class: 'icon-mima', value: '密码重置', click: this.passwordReset, privilege: 'B_yhgl_mmcz' }
      ]
    },
    currentChange(val) {
      this.listQuery.currPage = val
      this.getList()
    },
    // 提交
    passwordResetSubimt() {
      this.$refs['passwordForm'].validate((valid) => {
        if (valid) {
          this.loading = true
          resetPassword({ id: this.currentUserId, password: this.passwordForm.password }).then(response => {
            if (response.success) {
              this.$message.success(response.message)
              this.$refs['passwordForm'].resetFields()
              this.passwordFormVisible = false
            } else {
              this.$message.error(response.message)
            }
            this.loading = false
          })
        }
      })
    },
    // 密码重置
    passwordReset(val) {
      this.currentUserId = val.row.id
      this.passwordFormVisible = true
    },
    // 更新用户
    updateUser(val) {
      this.dialogStatus = 'update'
      this.userForm.id = val.row.id
      this.userForm.mobile = val.row.mobile
      this.userForm.orgId = val.row.orgId
      // this.userForm.password = val.row.password
      this.userForm.username = val.row.username
      this.userForm.name = val.row.name
      this.dialogFormVisible = true
    },
    // 配置角色
    setRole(val) {
      const me = this
      this.currentUserId = val.row.id
      const a = this.getRoleListByOrgId(val.row.orgId)
      const b = this.getRoleListByUserId(val.row.id)
      const p = Promise.all([a, b])
      p.then(function(val) {
        if (val.length === 2) {
          me.roleList = val[0]
          if (val[1] && val[1].length > 0) {
            const temp = []
            val[1].forEach(function(item) {
              temp.push(item.id)
            })
            me.roleForm.roleIds = temp
          } else {
            me.roleForm.roleIds = []
          }
          me.roleFormVisible = true
        }
      })
    },
    // 根据组织获取用户角色
    getRoleListByOrgId(orgId) {
      return new Promise(function(resolve) {
        getOrgRole({ id: orgId }).then(response => {
          if (response.success === true) {
            resolve(response.result)
          }
        })
      })
    },
    // 根据用户Id用户配置的角色
    getRoleListByUserId(id) {
      return new Promise(function(resolve) {
        getRoleList({ id }).then(response => {
          if (response.success) {
            resolve(response.result)
          }
        })
      })
    },
    // 通过用户ID获取用户权限
    getUserPrivilegesList(id) {
      return new Promise(function(resolve) {
        privilegesList({ id }).then(response => {
          if (response.success) {
            resolve(response.result)
          }
        })
      })
    },
    // 通过组织ID获取组织所拥有权限
    getPrivilegesList(orgId) {
      return new Promise(function(resolve) {
        getOrgPrivileges({ id: orgId }).then(response => {
          if (response.success) {
            resolve(response.result)
          }
        })
      })
    },
    // 配置权限
    setPrivilege(val) {
      const me = this
      this.currentUserId = val.row.id
      // 获取用户可配置权限
      const a = this.getPrivilegesList(val.row.orgId)
      // 获取用户配置的权限
      const b = this.getUserPrivilegesList(val.row.id)
      const p = Promise.all([a, b])
      p.then(function(val) {
        if (val.length === 2) {
          me.authList = val[0]
          if (val[1] && val[1].length > 0) {
            const tempArr = []
            val[1].forEach(function(item) {
              if (item.leaf) {
                tempArr.push(item.id)
              }
            })
            me.userPrivilges = tempArr
          } else {
            me.userPrivilges = []
          }
          me.authFormVisible = true
        }
      })
    },
    // 配置用户权限
    configPrivilges() {
      this.loading = true
      configUserPrivileges({
        userId: this.currentUserId,
        privilegeIds: this.$store.getters.checkAuth
      }).then(response => {
        if (response.success) {
          this.$message.success(response.message)
          this.authFormVisible = false
        } else {
          this.$message.error(response.message)
        }
        this.loading = false
      })
    },
    // 禁用启用用户
    enableUser(val) {
      this.$confirm('是否确定删除该记录?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser({ id: val.row.id }).then(response => {
          if (response.success) {
            this.getList()
            this.$message.success('删除成功')
          } else {
            this.$message.error(response.message)
          }
        })
      })
    },
    reset() { // 清空
      this.userForm.id = ''
      this.$refs['userForm'].resetFields()
    },
    // 新增用户
    addUser() {
      this.dialogStatus = 'create'
      this.userForm = {
        id: '',
        username: '',
        name: '',
        mobile: '',
        orgId: null,
        password: '',
        confirePassword: '',
        enabled: true
      }
      this.dialogFormVisible = true
    },
    // 保存用户
    saveUser() {
      this.$refs['userForm'].validate((valid) => {
        if (valid) {
          delete this.userForm.confirePassword
          delete this.userForm.username
          delete this.userForm.enabled
          if (!this.userForm.id) {
            delete this.userForm.id
          } else {
            delete this.userForm.password
          }

          this.loading = true
          saveUser(this.userForm).then(response => {
            if (response.code === 200) {
              this.$message.success('保存成功')
              // 操作数结构
              this.dialogFormVisible = false
              this.reset()
              this.getList()
              this.existUserList.push(this.userForm.username)
            } else {
              this.$message.error(response.message)
            }
            this.loading = false
          })
        }
      })
    },
    // 取消保存用户
    cancelSaveUser() {
      this.$refs['userForm'].clearValidate()
      this.dialogFormVisible = false
    },
    // 添加已经存在的用户名
    addExistUser() {
      this.list.forEach(item => {
        this.existUserList.push(item.username)
      })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
  .userDiv {
    .filter-container {
      .filter-item{
        margin-right: 10px;
      }
    }
    .vue-treeselect__control{
      font-size: 14px;
    }
  }
</style>
