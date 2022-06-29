import Layout from '@/layout'

const systemRouter = {
  path: '/sysManage',
  component: Layout,
  alwaysShow: true,
  meta: { title: '系统管理', icon: 'el-icon-setting', privilege: 'P_menu_system' },
  redirect: '/sysManage',
  children: [
    {
      path: 'dic',
      component: () => import('@/views/sysManage/dicManage'),
      name: 'DicManage',
      meta: { title: '字典管理', icon: 'el-icon-notebook-1', privilege: 'B_zdgl_zdwh' }
    },
    {
      path: 'privilege',
      component: () => import('@/views/sysManage/privilegeManage'),
      name: 'PrivilegeManage',
      meta: { title: '权限管理', icon: 'el-icon-s-order', privilege: 'P_qxgl' }
    },
    {
      path: 'org',
      component: () => import('@/views/sysManage/orgManage'),
      name: 'OrgManage',
      meta: { title: '组织管理', icon: 'el-icon-warning', privilege: 'P_zzgl' }
    },
    {
      path: 'role',
      component: () => import('@/views/sysManage/roleManage'),
      name: 'RoleManage',
      meta: { title: '角色管理', icon: 'el-icon-user-solid', privilege: 'P_jsgl' }
    },
    {
      path: 'user',
      component: () => import('@/views/sysManage/userManage'),
      name: 'UserManage',
      meta: { title: '用户管理', icon: 'el-icon-user', privilege: 'P_yhgl' }
    }
  ]
}
export default systemRouter
