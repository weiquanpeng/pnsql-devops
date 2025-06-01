import {DashboardIcon, UserIcon} from 'tdesign-icons-vue-next';
import { shallowRef } from 'vue';

import Layout from '@/layouts/index.vue';

export default [
  // {
  //   path: '/dashboard',
  //   component: Layout,
  //   redirect: '/dashboard/base',
  //   name: 'dashboard',
  //   meta: {
  //     title: {
  //       zh_CN: '大盘',
  //       en_US: 'Dashboard',
  //     },
  //     icon: shallowRef(DashboardIcon),
  //     orderNo: 0,
  //   },
  //   children: [
  //     {
  //       path: 'base',
  //       name: 'DashboardBase',
  //       component: () => import('@/pages/dashboard/base/index.vue'),
  //       meta: {
  //         title: {
  //           zh_CN: '概览仪表盘',
  //           en_US: 'Overview',
  //         },
  //       },
  //     },
  //   ],
  // },

  {
    path: '/database',
    component: Layout,
    redirect: '/database/postgresql',
    name: 'database',
    meta: {
      title: {
        zh_CN: '数据库',
        en_US: 'Database',
      },
      icon: shallowRef(DashboardIcon),
      orderNo: 1,
    },
    children: [
      {
        path: 'postgresql',
        name: 'DatabasePostgreSQL',
        redirect: '/database/postgresql/connection',
        meta: {
          title: {
            zh_CN: 'PostgreSQL',
            en_US: 'PostgreSQL',
          },
        },
        children: [
          {
            path: 'connection',
            name: 'PostgreSQLConnection',
            component: () => import('@/pages/database/postgresql/connection.vue'),
            meta: {
              title: {
                zh_CN: '会话管理',
                en_US: 'Connection',
              },
            },
          },
          {
            path: 'storage',
            name: 'PostgreSQLStorage',
            component: () => import('@/pages/database/postgresql/storage.vue'),
            meta: {
              title: {
                zh_CN: '空间管理',
                en_US: 'Storage',
              },
            },
          },
        ],
      },
      {
        path: 'mongo',
        name: 'DatabaseMongo',
        redirect: '/database/mongo/connection',
        meta: {
          title: {
            zh_CN: 'MongoDB',
            en_US: 'MongoDB',
          },
        },
        children: [
          {
            path: 'connection',
            name: 'MongoConnection',
            component: () => import('@/pages/database/mongo/index.vue'),
            meta: {
              title: {
                zh_CN: '连接管理',
                en_US: 'Connection',
              },
            },
          },
        ],
      },
    ],
  },
  {
    path: '/owner',
    component: Layout,
    name: 'owner',
    meta: {
      title: {
        zh_CN: '权限管理',
        en_US: 'User Management',
      },
      icon: shallowRef(UserIcon),
      orderNo: 2,
    },
    children: [
      {
        path: 'index',
        name: 'OwnerIndex',
        component: () => import('@/pages/owner/index.vue'),
        meta: {
          title: {
            zh_CN: '用户管理',
            en_US: 'User Management',
          },
        },
      },
    ],
  },
];
