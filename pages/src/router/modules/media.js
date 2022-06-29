import Layout from '@/layout'

const mediaRouter = {
  path: '/mediaManage',
  component: Layout,
  alwaysShow: true,
  meta: { title: '流媒体管理', icon: 'el-icon-video-camera', privilege: 'P_camera' },
  redirect: '/mediaManage',
  children: [
    {
      path: 'channelManage',
      component: () => import('@/views/mediaManage/channelManage'),
      name: 'ChannelManage',
      meta: { title: '通道管理', icon: 'el-icon-video-camera', privilege: 'P_camera_video' }
    },
    {
      path: 'pusherManage',
      component: () => import('@/views/mediaManage/pusherManage'),
      name: 'PushertManage',
      meta: { title: '客户端管理', icon: 'el-icon-video-camera', privilege: 'P_camera_video' }
    },
    {
      path: 'splitScreen',
      component: () => import('@/views/mediaManage/splitScreen'),
      name: 'SplitScreen',
      meta: { title: '分屏播放', icon: 'el-icon-video-camera', privilege: 'P_camera_video' }
    }
  ]
}
export default mediaRouter
