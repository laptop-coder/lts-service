// Package permissions contains all permission constants used in the application
package permissions

// Post permissions
const (
	PostCreate          = "post.create"
	PostReadAny         = "post.read.any"
	PostReadOwn         = "post.read.own"
	PostUpdateAny       = "post.update.any"
	PostUpdateOwn       = "post.update.own"
	PostDeleteAny       = "post.delete.any"
	PostDeleteOwn       = "post.delete.own"
	PostPhotoDeleteAny       = "post.photo.delete.any"
	PostPhotoDeleteOwn       = "post.photo.delete.own"
	PostVerify          = "post.verify"
	PostMarkReturnedAny = "post.mark.returned.any"
	PostMarkReturnedOwn = "post.mark.returned.own"
)

// User permissions
const (
	UserReadOwn          = "user.read.own"
	UserReadOther        = "user.read.other"
	UserReadAll          = "user.read.all"
	UserUpdateOwn        = "user.update.own"
	UserDeleteAny        = "user.delete.any"
	UserDeleteOwn        = "user.delete.own"
	UserReadOwnSubjects  = "user.read.own.subjects"
	TeacherReadClassroom = "teacher.read.classroom"
	StudentReadClassroom = "student.read.classroom"
	TeacherStudentsRead  = "teacher.students.read"
	ParentStudentsRead   = "parent.students.read"
	StudentTeacherRead   = "student.teacher.read"
	StudentParentsRead   = "student.parents.read"
)

// Room permissions
const (
	RoomCreate  = "room.create"
	RoomRead = "room.read"
	RoomUpdate  = "room.update"
	RoomDelete  = "room.delete"
)

// Subject permissions
const (
	SubjectCreate    = "subject.create"
	SubjectRead   = "subject.read"
	SubjectUpdate = "subject.update"
	SubjectDelete = "subject.delete"
)

// Student group permissions
const (
	StudentGroupCreate          = "student.group.create"
	StudentGroupReadAny         = "student.group.read.any"
	StudentGroupReadOwn         = "student.group.read.own"
	StudentGroupUpdate          = "student.group.update"
	StudentGroupDelete          = "student.group.delete"
	StudentGroupAdvisorAssign   = "student.group.advisor.assign"
	StudentGroupAdvisorUnassign = "student.group.advisor.unassign"
	StudentGroupAdvisorRead     = "student.group.advisor.read"
)

// Teacher permissions
const (
	TeacherSubjectReadAny       = "teacher.subject.read.any"
	TeacherSubjectReadOwn       = "teacher.subject.read.own"
	TeacherSubjectAssignAny     = "teacher.subject.assign.any"
	TeacherSubjectAssignOwn     = "teacher.subject.assign.own"
	TeacherSubjectUnassignAny   = "teacher.subject.unassign.any"
	TeacherSubjectUnassignOwn   = "teacher.subject.unassign.own"
	TeacherClassroomReadAny     = "teacher.classroom.read.any"
	TeacherClassroomReadOwn     = "teacher.classroom.read.own"
	TeacherClassroomAssignAny   = "teacher.classroom.assign.any"
	TeacherClassroomAssignOwn   = "teacher.classroom.assign.own"
	TeacherClassroomUnassignAny = "teacher.classroom.unassign.any"
	TeacherClassroomUnassignOwn = "teacher.classroom.unassign.own"
)

// Parent permissions
const (
	ParentStudentReadAny     = "parent.student.read.any"
	ParentStudentReadOwn     = "parent.student.read.own"
	ParentStudentAssignAny   = "parent.student.assign.any"
	ParentStudentAssignOwn   = "parent.student.assign.own"
	ParentStudentUnassignAny = "parent.student.unassign.any"
	ParentStudentUnassignOwn = "parent.student.unassign.own"
)

// User roles permissions
const (
	RoleAssign  = "role.assign"
	RoleAdd     = "role.add"
	RoleDelete  = "role.delete"
	RoleReadAny = "role.read.any"
	RoleReadOwn = "role.read.own"
)

// Permissions to work with tokens
const (
	TokenInviteAdminCreate = "token.invite.admin.create"
	TokenInviteUserCreate  = "token.invite.user.create"
	TokenInviteAdminDelete = "token.invite.admin.delete"
	TokenInviteUserDelete  = "token.invite.user.delete"
)
