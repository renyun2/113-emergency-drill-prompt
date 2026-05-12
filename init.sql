-- 应急预案与演练记录 — 初始化表结构与演示数据（10份预案 + 近两年演练）

SET client_encoding TO 'UTF8';

CREATE TABLE emergency_plans (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  plan_type VARCHAR(40) NOT NULL,
  scenario TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE plan_versions (
  id BIGSERIAL PRIMARY KEY,
  emergency_plan_id BIGINT NOT NULL REFERENCES emergency_plans(id) ON DELETE CASCADE,
  version_no VARCHAR(32) NOT NULL,
  revision_record TEXT,
  content_md TEXT NOT NULL DEFAULT '',
  published_date DATE,
  preparer VARCHAR(120),
  prepared_date DATE,
  approver VARCHAR(120),
  approved_date DATE,
  approval_status VARCHAR(20) NOT NULL DEFAULT 'draft',
  is_current BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (emergency_plan_id, version_no)
);
CREATE INDEX ix_plan_versions_plan ON plan_versions(emergency_plan_id);

CREATE TABLE drills (
  id BIGSERIAL PRIMARY KEY,
  emergency_plan_id BIGINT NOT NULL REFERENCES emergency_plans(id) ON DELETE CASCADE,
  drill_kind VARCHAR(32) NOT NULL,
  scheduled_date DATE NOT NULL,
  location VARCHAR(300),
  org_dept VARCHAR(200),
  participant_scope VARCHAR(500),
  objectives TEXT,
  notify_depts TEXT,
  drill_status VARCHAR(24) NOT NULL DEFAULT 'planned',
  actual_participants INT,
  duration_minutes INT,
  process_description TEXT,
  problem_list TEXT,
  evaluation VARCHAR(20),
  photo_paths TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ix_drills_plan_date ON drills(emergency_plan_id, scheduled_date);
CREATE INDEX ix_drills_scheduled ON drills(scheduled_date);

CREATE TABLE drill_issues (
  id BIGSERIAL PRIMARY KEY,
  drill_id BIGINT NOT NULL REFERENCES drills(id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rectifications (
  id BIGSERIAL PRIMARY KEY,
  drill_issue_id BIGINT NOT NULL REFERENCES drill_issues(id) ON DELETE CASCADE,
  responsible_person VARCHAR(120),
  corrective_measure TEXT,
  due_date DATE,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  completed_at DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ---- 10 份预案 ----
INSERT INTO emergency_plans (id, name, plan_type, scenario) VALUES
(1, '厂区火灾专项应急预案', '火灾', '生产车间、仓库、办公楼'),
(2, '危化品泄漏应急处置预案', '危化品泄漏', '储罐区、装卸区、实验室'),
(3, '地震灾害应对预案', '地震', '全厂区及高层作业面'),
(4, '大面积停电应急保障预案', '停电', '数据中心、冷链、关键产线'),
(5, '台风暴雨自然灾害预案', '自然灾害', '厂区排水、户外设备、物流'),
(6, '办公楼宇火灾疏散预案', '火灾', '行政楼、研发中心'),
(7, '实验室易燃挥发物泄漏预案', '危化品泄漏', '质检与研发实验室'),
(8, '受限空间作业中毒窒息预案', '危化品泄漏', '罐体、地下管网维修'),
(9, '电力设备火灾联动预案', '火灾', '配电室、变压器及电缆沟'),
(10, '洪涝与内涝响应预案', '自然灾害', '低洼仓库与停车场');

SELECT setval(pg_get_serial_sequence('emergency_plans', 'id'), 10);

-- 各预案版本（含历史存档与当前生效）
INSERT INTO plan_versions (emergency_plan_id, version_no, revision_record, content_md, published_date, preparer, prepared_date, approver, approved_date, approval_status, is_current) VALUES
(1, 'v1.0', '首版发布', E'# 厂区火灾专项预案\n\n## 1. 目的\n快速控制初火、组织疏散。\n\n## 2. 响应分级\n- 一级：局部烟雾\n- 二级：明火蔓延', '2023-06-01', '张伟', '2023-05-28', '李总工', '2023-06-01', 'approved', false),
(1, 'v2.0', '增加充电区域消防要点', E'# 厂区火灾（v2）\n\n新增：电动车集中充电区巡查与断电流程。', '2024-03-15', '王芳', '2024-03-10', '李总工', '2024-03-15', 'approved', true),
(2, 'v1.0', '首版', E'# 危化品泄漏\n\n隔离、收容、洗消流程。', '2023-08-01', '刘强', '2023-07-25', '陈主任', '2023-08-01', 'approved', false),
(2, 'v1.1', '修订卸料区监测', E'# 危化品泄漏 v1.1\n\n强化卸料区气体检测频次。', '2025-01-10', '刘强', '2025-01-05', '陈主任', '2025-01-10', 'approved', true),
(3, 'v1.0', '首版', E'# 地震应对\n\n停产排查、结构安全评估。', '2022-11-20', '赵敏', '2022-11-15', '李总工', '2022-11-20', 'approved', true),
(4, 'v1.0', '首版', E'# 停电保障\n\nUPS切换、柴油发电机启停。', '2024-05-01', '周杰', '2024-04-28', '设备部长', '2024-05-01', 'approved', true),
(5, 'v1.0', '首版', E'# 台风暴雨\n\n门窗加固、成品防雨、人员避险。', '2024-06-01', '孙丽', '2024-05-28', '安环经理', '2024-06-01', 'approved', true),
(6, 'v1.0', '首版', E'# 办公楼火灾\n\n疏散楼梯、集合点。', '2023-09-01', '钱勇', '2023-08-28', '行政总监', '2023-09-01', 'approved', true),
(7, 'v1.0', '首版', E'# 实验室泄漏\n\n通风橱、紧急喷淋。', '2024-02-01', '吴静', '2024-01-28', '实验室主任', '2024-02-01', 'approved', true),
(8, 'v1.0', '首版', E'# 受限空间\n\n气体检测、监护与救援。', '2024-07-01', '郑华', '2024-06-28', '安全主管', '2024-07-01', 'approved', true),
(9, 'v1.0', '首版', E'# 电气火灾\n\n断电、CO2/干粉选用。', '2023-12-01', '马超', '2023-11-28', '电气主任', '2023-12-01', 'approved', true),
(10, 'v1.0', '首版', E'# 洪涝内涝\n\n沙袋、排水泵、库存转移。', '2024-04-20', '林涛', '2024-04-15', '后勤经理', '2024-04-20', 'approved', true);

-- ---- 近两年（2024、2025）演练计划与记录 ----
INSERT INTO drills (emergency_plan_id, drill_kind, scheduled_date, location, org_dept, participant_scope, objectives, notify_depts, drill_status, actual_participants, duration_minutes, process_description, problem_list, evaluation, photo_paths, created_at, updated_at) VALUES
(1, '桌面推演', '2024-02-20', '会议室A', '安环部', '车间主任、班组长', '熟悉报警与疏散口令', '生产部、行政部', 'completed', 18, 90, '按脚本推演报警、集合、清点。', '部分疏散图标识褪色', '良好', '/uploads/drill/2024/fire_desk1.jpg', '2024-02-01', '2024-02-20'),
(1, '实战演练', '2024-09-10', '一号车间外广场', '安环部+消防志愿队', '当班员工约120人', '灭火器实操与疏散', '全厂部门', 'completed', 125, 120, '分区疏散，微型消防站出水演示。', '个别通道临时堆放物', '优秀', '/uploads/drill/2024/fire_live1.jpg,/uploads/drill/2024/fire_live2.jpg', '2024-08-01', '2024-09-10'),
(2, '联合演练', '2024-04-15', '储罐区模拟点', '安环+设备+辖区应急', '操作工、维修、医护', '管线泄漏收容', '医院、消防站', 'completed', 45, 150, '模拟法兰泄漏，围堤收容。', '应急泵启动略慢', '良好', '/uploads/drill/2024/hazmat1.jpg', '2024-03-20', '2024-04-15'),
(2, '桌面推演', '2025-01-20', '应急指挥中心', '安环部', '班组长以上', '新版卸料程序桌面验证', '储运部', 'completed', 22, 60, '逐步核对检测点位与汇报链。', NULL, '优秀', NULL, '2025-01-05', '2025-01-20'),
(3, '桌面推演', '2024-05-08', '会议室B', '综合办', '各部门安全员', '地震后停产与排查', '设备、IT', 'completed', 15, 75, '角色扮演信息上报。', 'IT灾备切换说明待更新', '一般', NULL, '2024-04-10', '2024-05-08'),
(4, '实战演练', '2024-11-05', '数据中心楼', '信息部+设备部', '运维、值班电工', 'UPS与油机切换', '采购（燃油）', 'completed', 12, 100, '完成双路市电掉电切换测试。', '燃油记录表未签字', '良好', '/uploads/drill/2024/power1.jpg', '2024-10-15', '2024-11-05'),
(5, '桌面推演', '2024-07-12', '会议室A', '安环部', '物流、仓储、保安', '暴雨前检查清单', '后勤', 'completed', 20, 70, '逐项确认排水与沙袋。', NULL, '优秀', NULL, '2024-06-20', '2024-07-12'),
(5, '实战演练', '2025-08-20', '厂区南门及低洼区', '后勤部', '保安、仓库、应急队', '强降水内涝抢排', '安环、生产', 'planned', NULL, NULL, NULL, NULL, NULL, NULL, '2025-07-01', '2025-07-01'),
(6, '实战演练', '2024-10-18', '研发楼', '行政部', '楼内全体', '火灾疏散与集合清点', '物业', 'completed', 80, 45, '上午第二节课时段演练。', '两部电梯未贴封条提示', '良好', '/uploads/drill/2024/office1.jpg', '2024-09-25', '2024-10-18'),
(7, '桌面推演', '2024-03-22', '实验室会议室', '质检部', '实验员', '泄漏报警与洗消', '安环', 'completed', 14, 50, '回顾SOP与PPE。', NULL, '优秀', NULL, '2024-03-01', '2024-03-22'),
(8, '实战演练', '2025-03-10', '地下泵房（模拟）', '设备部', '维修班组+监护', '准入与救援程序', '安环、医疗', 'completed', 16, 110, '全程四合一气体检测。', '救援绳扣具一套老化', '一般', '/uploads/drill/2025/confined1.jpg', '2025-02-01', '2025-03-10'),
(9, '桌面推演', '2024-12-05', '配电值班室', '设备部', '电工、班长', '电气火灾先断电原则', '安环', 'completed', 10, 55, '案例剖析与分工。', NULL, '良好', NULL, '2024-11-20', '2024-12-05'),
(10, '实战演练', '2025-05-15', '三号仓库外', '后勤部', '仓库、保安', '沙袋筑堤与排水', '安环', 'completed', 28, 95, '完成50米挡水带构筑。', '潜水泵备用数量不足', '良好', '/uploads/drill/2025/flood1.jpg', '2025-04-20', '2025-05-15'),
(3, '实战演练', '2025-09-01', '总装车间', '安环部', '当班员工', '震后设备停机与排查', '设备、质量', 'planned', NULL, NULL, NULL, NULL, NULL, NULL, '2025-08-01', '2025-08-01'),

(1, '桌面推演', '2026-03-01', '应急指挥中心', '安环部', '车间安全员', '复训火灾报警与疏散联动', '生产部', 'completed', 24, 75, '逐项核对广播与门禁释放逻辑。', NULL, '良好', NULL, '2026-02-15', '2026-03-01'),
(4, '桌面推演', '2026-04-18', '机房监控室', '信息部', '运维团队', 'UPS 告警与油机启停桌面推演', '设备部', 'completed', 14, 68, '推演市电闪断场景。', NULL, '优秀', NULL, '2026-04-01', '2026-04-18');

SELECT setval(pg_get_serial_sequence('drills', 'id'), (SELECT COALESCE(MAX(id),1) FROM drills));

-- 演练问题与整改闭环（drill_id 与上方 INSERT drills 顺序一致）
INSERT INTO drill_issues (drill_id, description) VALUES
(1, '部分疏散图标识褪色，员工不易辨认'),
(2, '个别通道存在临时堆放物影响疏散'),
(3, '应急泵启动响应时间略慢'),
(6, '燃油领用记录表缺少签字'),
(9, '电梯未张贴火灾禁用提示'),
(11, '救援用绳扣具一套老化需更换'),
(13, '备用潜水泵数量不足');

INSERT INTO rectifications (drill_issue_id, responsible_person, corrective_measure, due_date, status, completed_at) VALUES
(1, '行政部-钱勇', '更换全部楼层疏散图并夜光覆膜', '2024-03-15', 'done', '2024-03-10'),
(2, '生产部-各车间', '建立通道日检并纳入5S', '2024-09-25', 'done', '2024-09-22'),
(3, '设备部-周杰', '应急泵周试车并记录', '2024-05-01', 'done', '2024-04-28'),
(4, '信息部', '完善油机燃油记录双签制度', '2024-11-20', 'done', '2024-11-18'),
(5, '行政部-物业', '全部电梯口张贴火灾禁用标识', '2024-11-01', 'done', '2024-10-28'),
(6, '设备部-郑华', '更换救援装备并建立年检台账', '2025-04-01', 'done', '2025-03-28'),
(7, '后勤部-林涛', '增配2台备用潜水泵', '2025-06-30', 'pending', NULL);

SELECT setval(pg_get_serial_sequence('drill_issues', 'id'), (SELECT COALESCE(MAX(id),1) FROM drill_issues));
SELECT setval(pg_get_serial_sequence('rectifications', 'id'), (SELECT COALESCE(MAX(id),1) FROM rectifications));
SELECT setval(pg_get_serial_sequence('plan_versions', 'id'), (SELECT COALESCE(MAX(id),1) FROM plan_versions));
