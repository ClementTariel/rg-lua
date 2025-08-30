import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import Home from '../views/Home.vue'

export const mainRoutes: RouteRecordRaw[] = [
  {
    path: '/home',
    name: 'Home',
    component: Home,
  },
  // TODO: WIP
  // {
  //   path: '/ranking',
  //   name: 'Ranking',
  //   // route level code-splitting
  //   // this generates a separate chunk for this route
  //   // which is lazy-loaded when the route is visited.
  //   component: () => import('../views/Ranking.vue'),
  // },
  // TODO: WIP
  {
    path: '/matchs',
    name: 'Matchs',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/Matchs.vue'),
  },
  {
    path: '/rules',
    name: 'Rules',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/Rules.vue'),
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/AboutUs.vue'),
  },
];

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/home',
  },
  ...mainRoutes,
  {
    path: '/matchs/:id',
    name: 'Match details',
    // route level code-splitting
    // this generates a separate chunk for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/MatchDetails.vue'),
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes
})

export default router
