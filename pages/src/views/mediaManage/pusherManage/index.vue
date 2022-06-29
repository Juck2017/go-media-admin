<template>
  <div class="app-container clientDiv">
    <my-card title="客户端管理">
      <div class="filter-container">
        <el-input
          v-model="listQuery.name"
          placeholder="请输入名称"
          prefix-icon="el-icon-search"
          style="width: 200px;margin-right: 10px;"
          class="filter-item"
          clearable
        />
        <el-button v-waves class="filter-item" type="primary" @click="handleFilter">
          查询
        </el-button>
        <el-button v-waves class="filter-item" type="primary" @click="resetFilter">
          重置
        </el-button>
        <el-button class="filter-item" style="margin-left: 10px;" type="primary" @click="addCamera">
          新增
        </el-button>
      </div>
      <table-list
        :data="pusherList"
        :columns="columns"
        :list-loading="listLoading"
        class="dataTable"
        :total="total"
        :page-size="listQuery.pageSize"
        :show-index="true"
        @currentChange="currentChange"
      />
    </my-card>
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogFormVisible"
      :before-close="beforeClose"
    >
      <el-form
        ref="pusherForm"
        :rules="rules"
        :model="pusherForm"
        label-position="left"
        label-width="100px"
        size="small"
      >
        <el-row>
          <el-col :span="12">
            <el-form-item label="设备区域" prop="AreaId">
              <el-select v-model="pusherForm.AreaId" class="filter-item" placeholder="请选择设备区域" clearable>
                <el-option
                  v-for="item in dicMap['video_area']"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="取流地址" prop="streamAddress">
              <el-input v-model="pusherForm.streamAddress" placeholder="请输入取流地址" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="是否NVR地址" prop="isNvrAddress">
              <el-select v-model="pusherForm.isNvrAddress" class="filter-item" placeholder="请选择是否NVR" clearable>
                <el-option
                  v-for="item in dicMap['is_nvr_address']"
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
        <el-button @click="cancelSave">
          取消
        </el-button>
        <el-button type="primary" @click="handleSave">
          确认
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import MyCard from '@/components/MyCard'
import tableList from '@/components/Table/TableList'
import { getPusherList } from '@/api/media'
import waves from '@/directive/waves' // waves directive
import { mapGetters } from 'vuex'

export default {
  name: 'PusherManage',
  components: {
    MyCard,
    tableList
  },
  directives: { waves },
  data() {
    return {
      tableKey: 0,
      pusherList: [],
      total: 0,
      listLoading: true,
      listQuery: {
        pageSize: 3,
        currPage: 1,
        name: ''
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '新增'
      },
      pusherForm: {
        id: '',
        clientId: '',
        accessTime: '',
        status: '',
        name: ''
      },
      columns: [
        {
          text: '客户端ID',
          value: 'clientId'
        },
        {
          text: '客户端名称',
          value: 'name'
        },
        {
          text: '客户端当前状态',
          value: 'status'
        },
        {
          text: '接入系统时间',
          value: 'accessTime'
        },
        {
          text: '操作',
          type: 'iconButton',
          width: 180,
          list: this.operButton
        }
      ],
      rules: {
        areaId: [{ required: true, message: '请选择区域', trigger: 'blur' }],
        streamAddress: [{ required: true, message: '取流地址不能为空', trigger: 'blur' }],
        isNvrAddress: [{ required: true, message: '请选择是否NVR地址', trigger: 'blur' }]
      }
    }
  },
  computed: {
    ...mapGetters(['dicMap'])
  },
  mounted() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      if (!this.listQuery.areaId) {
        delete this.listQuery.areaId
      }
      getPusherList(this.listQuery).then(res => {
        if (res.success) {
          console.log(res.result.list)
          this.pusherList = res.result.list
          this.total = res.result.total
        }
        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.currPage = 1
      this.getList()
    },
    // 重置过滤
    resetFilter() {
      this.listQuery.name = ''
      this.getList()
    },
    // 处理改变
    currentChange(val) {
      this.listQuery.currPage = val
      this.getList()
    },
    // 新增弹窗
    addCamera() {
      this.dialogFormVisible = true
      this.dialogStatus = 'create'
    },
    // 保存操作
    handleSave() {
      this.$refs['pusherForm'].validate((valid) => {
        if (valid) {
          console.log('confire success', this.pusherForm)
        }
      })
    },
    // 关闭之前
    beforeClose() {
      this.cancelSave()
    },
    // 取消操作
    cancelSave() {
      this.dialogFormVisible = false
      this.$refs['pusherForm'].clearValidate()
      this.$refs['pusherForm'].resetFields()
    },
    // 表格操作按鈕
    operButton() {
      return [
        { class: 'icon-xiugai', value: '编辑', click: this.updateCamera, privilege: 'B_camera_edit' },
        { class: 'icon-shanchu1', value: '删除', click: this.deleteCamera, privilege: 'B_camera_del' }
      ]
    },
    // 编辑流媒体
    updateStream(val) {

    },
    // 删除流媒体
    deleteCamera(val) {

    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.clientDiv {
    .filter-container {
      .filter-item{
        margin-right: 10px;
      }
    }
  }
</style>
