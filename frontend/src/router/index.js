import { createRouter, createWebHistory } from 'vue-router'
import ComplianceView from '../views/ComplianceView.vue'
import PlansView from '../views/PlansView.vue'
import DrillsView from '../views/DrillsView.vue'
import RectView from '../views/RectView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/compliance' },
    { path: '/compliance', component: ComplianceView },
    { path: '/plans', component: PlansView },
    { path: '/drills', component: DrillsView },
    { path: '/rectifications', component: RectView },
  ],
})

export default router
