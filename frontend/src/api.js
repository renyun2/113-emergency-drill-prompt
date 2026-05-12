import axios from 'axios'

const http = axios.create({ baseURL: '/api', timeout: 60000 })

export const fetchPlans = () => http.get('/plans')
export const createPlan = (data) => http.post('/plans', data)
export const patchPlan = (id, data) => http.patch(`/plans/${id}`, data)
export const deletePlan = (id) => http.delete(`/plans/${id}`)

export const fetchVersions = (planId) => http.get(`/plans/${planId}/versions`)
export const createVersion = (planId, data) => http.post(`/plans/${planId}/versions`, data)
export const patchVersion = (vid, data) => http.patch(`/plan-versions/${vid}`, data)
export const submitVersion = (vid, body) => http.post(`/plan-versions/${vid}/submit`, body || {})
export const approveVersion = (vid, body) => http.post(`/plan-versions/${vid}/approve`, body)

export const fetchDrills = (params) => http.get('/drills', { params })
export const createDrill = (data) => http.post('/drills', data)
export const patchDrill = (id, data) => http.patch(`/drills/${id}`, data)

export const fetchIssues = () => http.get('/issues')
export const fetchDrillIssues = (drillId) => http.get(`/drills/${drillId}/issues`)
export const createDrillIssue = (drillId, data) => http.post(`/drills/${drillId}/issues`, data)
export const createRectFromIssue = (issueId, data) => http.post(`/issues/${issueId}/rectifications`, data)

export const fetchRectifications = () => http.get('/rectifications')
export const patchRectification = (id, data) => http.patch(`/rectifications/${id}`, data)

export const fetchCompliance = (year) =>
  http.get('/stats/compliance', { params: year ? { year } : {} })
