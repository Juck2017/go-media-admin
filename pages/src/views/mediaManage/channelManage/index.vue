<template>
  <div class="app-container">
    <my-card title="通道管理">
      <!-- 过滤 -->
      <div class="filter-container">
        <el-select v-model="listQuery.areaId" placeholder="请选择区域" clearable class="filter-item"
          style="width: 130px; margin-right: 10px;">
          <el-option v-for="item in dicMap['video_area']" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
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
      <!-- 表格数据列表 -->
      <table-list :data="channelList" :columns="columns" :list-loading="listLoading" class="dataTable" :total="total"
        :page-size="listQuery.pageSize" :show-index="true" @currentChange="currentChange" />
    </my-card>
    <!-- 新增和编辑对话框 -->
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" :before-close="beforeClose"
      :close-on-click-modal="false">
      <el-form ref="channelForm" :rules="rules" :model="channelForm" label-position="left" label-width="100px"
        size="small">
        <el-row>
          <el-col :span="12">
            <el-form-item label="设备区域" prop="AreaId">
              <el-select v-model="channelForm.AreaId" class="filter-item" placeholder="请选择设备区域" clearable>
                <el-option v-for="item in dicMap['video_area']" :key="item.id" :label="item.name" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="取流地址" prop="streamAddress">
              <el-input v-model="channelForm.streamAddress" placeholder="请输入取流地址" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="是否NVR地址" prop="isNvrAddress">
              <el-select v-model="channelForm.isNvrAddress" class="filter-item" placeholder="请选择是否NVR" clearable>
                <el-option v-for="item in dicMap['is_nvr_address']" :key="item.id" :label="item.name"
                  :value="item.id" />
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

    <!-- 播放对话框 -->
    <el-dialog id="videoDialog" title="播放" :visible.sync="dialogVideoVisible" width="680px"
      :before-close="handleVideoClose" :close-on-click-modal="false">
      <div id="canvasDiv" />
    </el-dialog>
  </div>
</template>

<script>
import MyCard from '@/components/MyCard'
import tableList from '@/components/Table/TableList'
import './js/jsmpeg.min'
import { getChannelList } from '@/api/media'
import waves from '@/directive/waves' // waves directive
import { mapGetters } from 'vuex'

export default {
  name: 'ChannelManage',
  components: {
    MyCard,
    tableList
  },
  directives: { waves },
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'info',
        deleted: 'danger'
      }
      return statusMap[status]
    },
    typeFilter(type) {
      const calendarTypeKeyValue = this.videoAreaOptions.reduce((acc, cur) => {
        acc[cur.key] = cur.display_name
        return acc
      }, {})
      return calendarTypeKeyValue[type]
    }
  },
  data() {
    return {
      tableKey: 0,
      channelList: [],
      total: 0,
      listLoading: true,
      listQuery: {
        pageSize: 3,
        currPage: 1,
        areaId: ''
      },
      dialogFormVisible: false,
      dialogVideoVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '新增'
      },
      channelForm: {
        id: '',
        AreaId: '',
        createdAt: '',
        status: '',
        streamAddress: '',
        isNvrAddress: ''
      },
      columns: [
        {
          text: '所属客户端',
          value: 'pusherName'
        },
        {
          text: '通道编码',
          value: 'code'
        },
        {
          text: '通道名称',
          value: 'name'
        },
        {
          text: '通道状态',
          value: 'status',
          type: 'Tag',
          effect: 'dark',
          getLabel: this.getLabel
        },

        {
          text: '在线数量',
          value: 'onlineNum'
        },
        {
          text: '取流地址',
          value: 'streamAddress'
        },
        {
          text: '是否NVR地址',
          value: 'isNvrAddress'
        },
        {
          text: '接入时间',
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
      getChannelList(this.listQuery).then(res => {
        if (res.success) {
          this.channelList = res.result.list
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
      this.listQuery.areaId = ''
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
      this.$refs['channelForm'].validate((valid) => {
        if (valid) {
          console.log('confire success', this.channelForm)
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
      this.$refs['channelForm'].clearValidate()
      this.$refs['channelForm'].resetFields()
    },
    // 表格操作按鈕
    operButton() {
      return [
        { class: 'icon-shipin', value: '播放', click: this.playVideo, privilege: 'B_camera_play' },
        { class: 'icon-close', value: '关闭通道', click: this.closeChannel, privilege: 'B_camera_play' },
        { class: 'icon-xiugai', value: '编辑', click: this.updateChannel, privilege: 'B_camera_edit' },
        { class: 'icon-shanchu1', value: '删除', click: this.deleteChannel, privilege: 'B_camera_del' }
      ]
    },
    // 播放地址
    playVideo(val) {
      this.dialogVideoVisible = true
      const url = 'ws://127.0.0.1:5000'
      this.$nextTick(function () {
        const canvasDiv = document.getElementById('canvasDiv')
        const canvas = document.createElement('canvas')
        canvas.style.width = '640px'
        canvas.style.height = '350px'
        canvasDiv.appendChild(canvas)
        const row = val.row
        // eslint-disable-next-line
        this.player = new JSMpeg.Player(`${url}?code=${row.code}`,
          {
            canvas,
            onPlay: function () {
              console.log('play video')
            }
          })
      })
    },
    handleVideoClose(done) {
      console.log('close video')
      this.player.destroy()
      done()
    },
    // 编辑通道
    updateChannel() {

    },
    // 删除通道
    deleteChannel() {

    },
    // 关闭通道
    closeChannel() {

    },
    getLabel(val) {
      return val.status === 1 ? '在线' : '离线'
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
#canvasDiv {
  margin: 0;
  width: 640px;
  height: 350px;
}
</style>
