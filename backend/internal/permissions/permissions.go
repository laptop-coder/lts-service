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
	PostPhotoDeleteAny  = "post.photo.delete.any"
	PostPhotoDeleteOwn  = "post.photo.delete.own"
	PostPhotoUpdateAny  = "post.photo.update.any"
	PostPhotoUpdateOwn  = "post.photo.update.own"
	PostVerify          = "post.verify"
	PostMarkReturnedAny = "post.mark.returned.any"
	PostMarkReturnedOwn = "post.mark.returned.own"
)

// Post conversation permissions
const (
	ConversationCreate            = "conversation.create"
	ConversationReadOwn           = "conversation.read.own"
	ConversationMessageSend       = "conversation.message.send"
	ConversationMessageMarkAsRead = "conversation.message.mark_as_read"
)

// User permissions
const (
	UserReadOwn        = "user.read.own"
	UserReadOther      = "user.read.other"
	UserReadAll        = "user.read.all"
	UserUpdateOwn      = "user.update.own"
	UserDeleteAnyAdmin = "user.delete.any.admin"
	UserDeleteAnyUser  = "user.delete.any.user"
	UserDeleteOwn      = "user.delete.own"
)

// Room permissions
const (
	RoomCreate = "room.create"
	RoomUpdate = "room.update"
	RoomDelete = "room.delete"
)

// Subject permissions
const (
	SubjectCreate = "subject.create"
	SubjectUpdate = "subject.update"
	SubjectDelete = "subject.delete"
)

// Student group permissions
const (
	StudentGroupCreate             = "student_group.create"
	StudentGroupUpdate             = "student_group.update"
	StudentGroupDelete             = "student_group.delete"
	StudentGroupAdvisorAssign      = "student_group.advisor.assign"
	StudentGroupAdvisorUnassignAny = "student_group.advisor.unassign.any"
	StudentGroupAdvisorUnassignOwn = "student_group.advisor.unassign.own"
	StudentGroupAdvisorRead        = "student_group.advisor.read"
)

// Teacher permissions
const (
	TeacherSubjectReadAny       = "teacher.subject.read.any"
	TeacherSubjectReadOwn       = "teacher.subject.read.own"
	TeacherSubjectAddAny        = "teacher.subject.add.any"
	TeacherSubjectAddOwn        = "teacher.subject.add.own"
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
	TeacherReadOther            = "teacher.read.other"
	TeacherReadOwn              = "teacher.read.own"
	TeacherStudentGroupReadOwn  = "teacher.student_group.read.own"
)

// Parent permissions
const (
	ParentStudentReadAny      = "parent.student.read.any"
	ParentStudentReadOwn      = "parent.student.read.own"
	ParentStudentAddAny       = "parent.student.add.any"
	ParentStudentAddOwn       = "parent.student.add.own"
	ParentStudentUnassignAny  = "parent.student.unassign.any"
	ParentStudentUnassignOwn  = "parent.student.unassign.own"
	ParentReadOther           = "parent.read.other"
	ParentReadOwn             = "parent.read.own"
	ParentStudentGroupReadOwn = "parent.student_group.read.own"
)

// User roles permissions
const (
	RoleAdminAssign   = "role.admin.assign"
	RoleUserAssign    = "role.user.assign"
	RoleAdminAdd      = "role.admin.add"
	RoleUserAdd       = "role.user.add"
	RoleAdminUnassign = "role.admin.unassign"
	RoleUserUnassign  = "role.user.unassign"
	RoleReadAny       = "role.read.any"
	RoleReadOwn       = "role.read.own"
)

// Permissions to work with tokens
const (
	TokenInviteAdminCreate = "token.invite.admin.create"
	TokenInviteUserCreate  = "token.invite.user.create"
	TokenInviteAdminDelete = "token.invite.admin.delete"
	TokenInviteUserDelete  = "token.invite.user.delete"
)

// Student permissions
const (
	StudentReadOther           = "student.read.other"
	StudentReadOwn             = "student.read.own"
	StudentClassroomReadAny    = "student.classroom.read.any"
	StudentClassroomReadOwn    = "student.classroom.read.own"
	StudentAdvisorReadAny      = "student.advisor.read.any"
	StudentAdvisorReadOwn      = "student.advisor.read.own"
	StudentParentReadAny       = "student.parent.read.any"
	StudentParentReadOwn       = "student.parent.read.own"
	StudentStudentGroupReadOwn = "student.student_group.read.own"
)

// Institution administrator
const (
	InstitutionAdministratorReadOther      = "institution_administrator.read.other"
	InstitutionAdministratorReadOwn        = "institution_administrator.read.own"
	InstitutionAdministratorPositionAssign = "institution_administrator.position.assign"
	InstitutionAdministratorPositionRead   = "institution_administrator.position.read"
)

// Staff
const (
	StaffReadOther      = "staff.read.other"
	StaffReadOwn        = "staff.read.own"
	StaffPositionAssign = "staff.position.assign"
	StaffPositionRead   = "staff.position.read"
)

// Position institution administrator
const (
	PositionInstitutionAdministratorCreate = "position.institution_administrator.create"
	PositionInstitutionAdministratorUpdate = "position.institution_administrator.update"
	PositionInstitutionAdministratorDelete = "position.institution_administrator.delete"
)

// Position staff
const (
	PositionStaffCreate = "position.staff.create"
	PositionStaffUpdate = "position.staff.update"
	PositionStaffDelete = "position.staff.delete"
)
