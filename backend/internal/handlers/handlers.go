package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"emergency-drill/internal/models"
)

type H struct{ DB *gorm.DB }

func Register(r *gin.Engine, db *gorm.DB) {
	h := &H{DB: db}
	api := r.Group("/api")
	{
		api.GET("/plans", h.listPlans)
		api.POST("/plans", h.createPlan)
		api.GET("/plans/:id", h.getPlan)
		api.PATCH("/plans/:id", h.patchPlan)
		api.DELETE("/plans/:id", h.delPlan)

		api.GET("/plans/:id/versions", h.listVersions)
		api.POST("/plans/:id/versions", h.createVersion)
		api.PATCH("/plan-versions/:vid", h.patchVersion)
		api.POST("/plan-versions/:vid/submit", h.submitVersion)
		api.POST("/plan-versions/:vid/approve", h.approveVersion)

		api.GET("/drills", h.listDrills)
		api.POST("/drills", h.createDrill)
		api.GET("/drills/:id/issues", h.listDrillIssues)
		api.POST("/drills/:id/issues", h.createDrillIssue)
		api.GET("/drills/:id", h.getDrill)
		api.PATCH("/drills/:id", h.patchDrill)

		api.GET("/issues", h.listAllIssues)
		api.POST("/issues/:issueId/rectifications", h.createRectification)
		api.GET("/rectifications", h.listRectifications)
		api.PATCH("/rectifications/:id", h.patchRectification)

		api.GET("/stats/compliance", h.complianceStats)
	}
	r.GET("/healthz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.String(http.StatusOK, "ok")
	})
}

func parseID(c *gin.Context, key string) (uint64, bool) {
	v := c.Param(key)
	n, err := strconv.ParseUint(v, 10, 64)
	return n, err == nil && n > 0
}

// ---- plans ----

func (h *H) listPlans(c *gin.Context) {
	var plans []models.EmergencyPlan
	if err := h.DB.Order("id").Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *H) createPlan(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		PlanType string `json:"plan_type" binding:"required"`
		Scenario string `json:"scenario"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := models.EmergencyPlan{Name: body.Name, PlanType: body.PlanType, Scenario: body.Scenario}
	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *H) getPlan(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var p models.EmergencyPlan
	if err := h.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *H) patchPlan(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Name     *string `json:"name"`
		PlanType *string `json:"plan_type"`
		Scenario *string `json:"scenario"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.PlanType != nil {
		updates["plan_type"] = *body.PlanType
	}
	if body.Scenario != nil {
		updates["scenario"] = *body.Scenario
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields"})
		return
	}
	if err := h.DB.Model(&models.EmergencyPlan{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var p models.EmergencyPlan
	_ = h.DB.First(&p, id).Error
	c.JSON(http.StatusOK, p)
}

func (h *H) delPlan(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.DB.Delete(&models.EmergencyPlan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ---- versions ----

func (h *H) listVersions(c *gin.Context) {
	pid, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	var vs []models.PlanVersion
	if err := h.DB.Where("emergency_plan_id = ?", pid).Order("id DESC").Find(&vs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vs)
}

func (h *H) createVersion(c *gin.Context) {
	pid, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	var body struct {
		VersionNo      string `json:"version_no" binding:"required"`
		RevisionRecord string `json:"revision_record"`
		ContentMd      string `json:"content_md"`
		Preparer       string `json:"preparer"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v := models.PlanVersion{
		EmergencyPlanID: pid,
		VersionNo:       body.VersionNo,
		RevisionRecord:  body.RevisionRecord,
		ContentMd:       body.ContentMd,
		Preparer:        body.Preparer,
		ApprovalStatus:  "draft",
		IsCurrent:       false,
	}
	if err := h.DB.Create(&v).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *H) patchVersion(c *gin.Context) {
	vid, ok := parseID(c, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}
	var body struct {
		RevisionRecord *string    `json:"revision_record"`
		ContentMd      *string    `json:"content_md"`
		PublishedDate  *time.Time `json:"published_date"`
		Preparer       *string    `json:"preparer"`
		PreparedDate   *time.Time `json:"prepared_date"`
		Approver       *string    `json:"approver"`
		ApprovedDate   *time.Time `json:"approved_date"`
		VersionNo      *string    `json:"version_no"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if body.RevisionRecord != nil {
		updates["revision_record"] = *body.RevisionRecord
	}
	if body.ContentMd != nil {
		updates["content_md"] = *body.ContentMd
	}
	if body.PublishedDate != nil {
		updates["published_date"] = *body.PublishedDate
	}
	if body.Preparer != nil {
		updates["preparer"] = *body.Preparer
	}
	if body.PreparedDate != nil {
		updates["prepared_date"] = *body.PreparedDate
	}
	if body.Approver != nil {
		updates["approver"] = *body.Approver
	}
	if body.ApprovedDate != nil {
		updates["approved_date"] = *body.ApprovedDate
	}
	if body.VersionNo != nil {
		updates["version_no"] = *body.VersionNo
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields"})
		return
	}
	if err := h.DB.Model(&models.PlanVersion{}).Where("id = ?", vid).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var v models.PlanVersion
	_ = h.DB.First(&v, vid).Error
	c.JSON(http.StatusOK, v)
}

func (h *H) submitVersion(c *gin.Context) {
	vid, ok := parseID(c, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}
	var body struct {
		PreparedDate *time.Time `json:"prepared_date"`
	}
	_ = c.ShouldBindJSON(&body)
	now := time.Now()
	pd := now
	if body.PreparedDate != nil {
		pd = *body.PreparedDate
	}
	if err := h.DB.Model(&models.PlanVersion{}).Where("id = ?", vid).Updates(map[string]interface{}{
		"approval_status": "pending",
		"prepared_date":   pd,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var v models.PlanVersion
	_ = h.DB.First(&v, vid).Error
	c.JSON(http.StatusOK, v)
}

func (h *H) approveVersion(c *gin.Context) {
	vid, ok := parseID(c, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}
	var body struct {
		Approver     string     `json:"approver" binding:"required"`
		ApprovedDate *time.Time `json:"approved_date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ad := time.Now()
	if body.ApprovedDate != nil {
		ad = *body.ApprovedDate
	}
	var v models.PlanVersion
	if err := h.DB.First(&v, vid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.PlanVersion{}).Where("emergency_plan_id = ?", v.EmergencyPlanID).Update("is_current", false).Error; err != nil {
			return err
		}
		return tx.Model(&models.PlanVersion{}).Where("id = ?", vid).Updates(map[string]interface{}{
			"approval_status": "approved",
			"approver":        body.Approver,
			"approved_date":   ad,
			"is_current":      true,
			"published_date":  ad,
		}).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = h.DB.First(&v, vid).Error
	c.JSON(http.StatusOK, v)
}

// ---- drills ----

func (h *H) listDrills(c *gin.Context) {
	tx := h.DB.Model(&models.Drill{}).Preload("Plan")
	if pid := c.Query("plan_id"); pid != "" {
		if n, err := strconv.ParseUint(pid, 10, 64); err == nil {
			tx = tx.Where("emergency_plan_id = ?", n)
		}
	}
	if st := strings.TrimSpace(c.Query("status")); st != "" {
		tx = tx.Where("drill_status = ?", st)
	}
	if y := c.Query("year"); y != "" {
		if yi, err := strconv.Atoi(y); err == nil {
			start := time.Date(yi, 1, 1, 0, 0, 0, 0, time.Local)
			end := time.Date(yi+1, 1, 1, 0, 0, 0, 0, time.Local)
			tx = tx.Where("scheduled_date >= ? AND scheduled_date < ?", start, end)
		}
	}
	var list []models.Drill
	if err := tx.Order("scheduled_date DESC, id DESC").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *H) createDrill(c *gin.Context) {
	var body struct {
		EmergencyPlanID  uint64    `json:"emergency_plan_id" binding:"required"`
		DrillKind        string    `json:"drill_kind" binding:"required"`
		ScheduledDate    time.Time `json:"scheduled_date" binding:"required"`
		Location         string    `json:"location"`
		OrgDept          string    `json:"org_dept"`
		ParticipantScope string    `json:"participant_scope"`
		Objectives       string    `json:"objectives"`
		NotifyDepts      string    `json:"notify_depts"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := models.Drill{
		EmergencyPlanID:  body.EmergencyPlanID,
		DrillKind:        body.DrillKind,
		ScheduledDate:    body.ScheduledDate,
		Location:         body.Location,
		OrgDept:          body.OrgDept,
		ParticipantScope: body.ParticipantScope,
		Objectives:       body.Objectives,
		NotifyDepts:      body.NotifyDepts,
		DrillStatus:      "planned",
	}
	if err := h.DB.Create(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = h.DB.Preload("Plan").First(&d, d.ID).Error
	c.JSON(http.StatusCreated, d)
}

func (h *H) getDrill(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var d models.Drill
	if err := h.DB.Preload("Plan").First(&d, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *H) patchDrill(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		DrillKind          *string    `json:"drill_kind"`
		ScheduledDate      *time.Time `json:"scheduled_date"`
		Location           *string    `json:"location"`
		OrgDept            *string    `json:"org_dept"`
		ParticipantScope   *string    `json:"participant_scope"`
		Objectives         *string    `json:"objectives"`
		NotifyDepts        *string    `json:"notify_depts"`
		DrillStatus        *string    `json:"drill_status"`
		ActualParticipants *int       `json:"actual_participants"`
		DurationMinutes    *int       `json:"duration_minutes"`
		ProcessDescription *string    `json:"process_description"`
		ProblemList        *string    `json:"problem_list"`
		Evaluation         *string    `json:"evaluation"`
		PhotoPaths         *string    `json:"photo_paths"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if body.DrillKind != nil {
		updates["drill_kind"] = *body.DrillKind
	}
	if body.ScheduledDate != nil {
		updates["scheduled_date"] = *body.ScheduledDate
	}
	if body.Location != nil {
		updates["location"] = *body.Location
	}
	if body.OrgDept != nil {
		updates["org_dept"] = *body.OrgDept
	}
	if body.ParticipantScope != nil {
		updates["participant_scope"] = *body.ParticipantScope
	}
	if body.Objectives != nil {
		updates["objectives"] = *body.Objectives
	}
	if body.NotifyDepts != nil {
		updates["notify_depts"] = *body.NotifyDepts
	}
	if body.DrillStatus != nil {
		updates["drill_status"] = *body.DrillStatus
	}
	if body.ActualParticipants != nil {
		updates["actual_participants"] = *body.ActualParticipants
	}
	if body.DurationMinutes != nil {
		updates["duration_minutes"] = *body.DurationMinutes
	}
	if body.ProcessDescription != nil {
		updates["process_description"] = *body.ProcessDescription
	}
	if body.ProblemList != nil {
		updates["problem_list"] = *body.ProblemList
	}
	if body.Evaluation != nil {
		updates["evaluation"] = *body.Evaluation
	}
	if body.PhotoPaths != nil {
		updates["photo_paths"] = *body.PhotoPaths
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields"})
		return
	}
	updates["updated_at"] = time.Now()
	if err := h.DB.Model(&models.Drill{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var d models.Drill
	_ = h.DB.Preload("Plan").First(&d, id).Error
	c.JSON(http.StatusOK, d)
}

// ---- issues & rectifications ----

func (h *H) listDrillIssues(c *gin.Context) {
	did, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drill id"})
		return
	}
	var xs []models.DrillIssue
	if err := h.DB.Where("drill_id = ?", did).Order("id").Find(&xs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, xs)
}

func (h *H) createDrillIssue(c *gin.Context) {
	did, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drill id"})
		return
	}
	var body struct {
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	iss := models.DrillIssue{DrillID: did, Description: body.Description}
	if err := h.DB.Create(&iss).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, iss)
}

func (h *H) createRectification(c *gin.Context) {
	iid, ok := parseID(c, "issueId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid issue id"})
		return
	}
	if err := h.DB.First(&models.DrillIssue{}, iid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
		return
	}
	var body struct {
		ResponsiblePerson string     `json:"responsible_person"`
		CorrectiveMeasure string     `json:"corrective_measure"`
		DueDate           *time.Time `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := models.Rectification{
		DrillIssueID:      iid,
		ResponsiblePerson: body.ResponsiblePerson,
		CorrectiveMeasure: body.CorrectiveMeasure,
		DueDate:           body.DueDate,
		Status:            "pending",
	}
	if err := h.DB.Create(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, r)
}

func (h *H) listAllIssues(c *gin.Context) {
	var xs []models.DrillIssue
	if err := h.DB.Preload("Drill").Preload("Drill.Plan").Order("id DESC").Find(&xs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, xs)
}

func (h *H) listRectifications(c *gin.Context) {
	var xs []models.Rectification
	if err := h.DB.Preload("Issue").Preload("Issue.Drill").Preload("Issue.Drill.Plan").Order("id DESC").Find(&xs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, xs)
}

func (h *H) patchRectification(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		ResponsiblePerson *string    `json:"responsible_person"`
		CorrectiveMeasure *string    `json:"corrective_measure"`
		DueDate           *time.Time `json:"due_date"`
		Status            *string    `json:"status"`
		CompletedAt       *time.Time `json:"completed_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if body.ResponsiblePerson != nil {
		updates["responsible_person"] = *body.ResponsiblePerson
	}
	if body.CorrectiveMeasure != nil {
		updates["corrective_measure"] = *body.CorrectiveMeasure
	}
	if body.DueDate != nil {
		updates["due_date"] = *body.DueDate
	}
	if body.Status != nil {
		updates["status"] = *body.Status
	}
	if body.CompletedAt != nil {
		updates["completed_at"] = *body.CompletedAt
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields"})
		return
	}
	if err := h.DB.Model(&models.Rectification{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var r models.Rectification
	_ = h.DB.Preload("Issue").Preload("Issue.Drill").First(&r, id).Error
	c.JSON(http.StatusOK, r)
}

// ---- stats ----

func (h *H) complianceStats(c *gin.Context) {
	year := time.Now().Year()
	if y := c.Query("year"); y != "" {
		if yi, err := strconv.Atoi(y); err == nil {
			year = yi
		}
	}
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)

	var completedThisYear int64
	_ = h.DB.Model(&models.Drill{}).
		Where("drill_status = ? AND scheduled_date >= ? AND scheduled_date < ?", "completed", start, end).
		Count(&completedThisYear).Error

	var planTotal int64
	_ = h.DB.Model(&models.EmergencyPlan{}).Count(&planTotal).Error

	type cnt struct {
		EmergencyPlanID uint64 `gorm:"column:emergency_plan_id"`
		N               int64  `gorm:"column:n"`
	}
	var perPlan []cnt
	_ = h.DB.Model(&models.Drill{}).
		Select("emergency_plan_id, COUNT(*) as n").
		Where("drill_status = ? AND scheduled_date >= ? AND scheduled_date < ?", "completed", start, end).
		Group("emergency_plan_id").Scan(&perPlan).Error

	planMet := 0
	planShort := []uint64{}
	m := map[uint64]int64{}
	for _, r := range perPlan {
		m[r.EmergencyPlanID] = r.N
	}
	var allPlans []models.EmergencyPlan
	_ = h.DB.Select("id").Find(&allPlans).Error
	for _, p := range allPlans {
		if m[p.ID] >= 1 {
			planMet++
		} else {
			planShort = append(planShort, p.ID)
		}
	}

	var orgsYear []string
	_ = h.DB.Model(&models.Drill{}).
		Where("drill_status = ? AND scheduled_date >= ? AND scheduled_date < ?", "completed", start, end).
		Where("TRIM(org_dept) <> ?", "").
		Distinct("org_dept").
		Pluck("org_dept", &orgsYear).Error
	distYear := int64(len(orgsYear))

	var orgsAll []string
	_ = h.DB.Model(&models.Drill{}).
		Where("drill_status = ?", "completed").
		Where("TRIM(org_dept) <> ?", "").
		Distinct("org_dept").
		Pluck("org_dept", &orgsAll).Error
	distAll := int64(len(orgsAll))

	denom := distAll
	if denom < 8 {
		denom = 8
	}
	coverage := 0.0
	if denom > 0 {
		coverage = float64(distYear) * 100.0 / float64(denom)
		if coverage > 100 {
			coverage = 100
		}
	}

	var rectTotal, rectDone int64
	_ = h.DB.Model(&models.Rectification{}).Count(&rectTotal).Error
	_ = h.DB.Model(&models.Rectification{}).Where("status = ?", "done").Count(&rectDone).Error
	rectRate := 0.0
	if rectTotal > 0 {
		rectRate = float64(rectDone) * 100.0 / float64(rectTotal)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":                       year,
		"completed_drills_year":      completedThisYear,
		"plan_total":                 planTotal,
		"plans_met_annual_minimum":   planMet,
		"plan_ids_below_minimum":     planShort,
		"dept_distinct_year":         distYear,
		"dept_distinct_baseline":     denom,
		"dept_coverage_percent":      coverage,
		"rectification_total":        rectTotal,
		"rectification_done":         rectDone,
		"rectification_done_percent": rectRate,
	})
}
