package main

import (
	"backend/internal/handler"
	"backend/internal/permissions"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

func SetupRoutes(
	mux *http.ServeMux,
	log logger.Logger,
	authMiddleware func(bool) func(http.Handler) http.Handler,
	requireRoles func(logger.Logger, bool, ...string) func(http.Handler) http.Handler,
	requirePermissions func(logger.Logger, bool, ...string) func(http.Handler) http.Handler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	conversationHandler *handler.ConversationHandler,
	studentGroupHandler *handler.StudentGroupHandler,
	roomHandler *handler.RoomHandler,
	subjectHandler *handler.SubjectHandler,
	studentHandler *handler.StudentHandler,
	teacherHandler *handler.TeacherHandler,
	parentHandler *handler.ParentHandler,
	staffHandler *handler.StaffHandler,
	institutionAdministratorHandler *handler.InstitutionAdministratorHandler,
	inviteHandler *handler.InviteHandler,
	staffPositionHandler *handler.StaffPositionHandler,
	institutionAdministratorPositionHandler *handler.InstitutionAdministratorPositionHandler,
	documentHandler *handler.DocumentHandler,
) {
	mux.Handle("/health", authMiddleware(true)(http.HandlerFunc(healthHandler)))

	// Post
	mux.Handle("POST /api/v1/posts", authMiddleware(false)(requirePermissions(log, false, permissions.PostCreate)(http.HandlerFunc(postHandler.Create))))
	mux.Handle("GET /api/v1/posts", authMiddleware(false)(requirePermissions(log, false, permissions.PostReadAny)(http.HandlerFunc(postHandler.GetPosts))))
	mux.Handle("GET /api/v1/posts/public", authMiddleware(true)(http.HandlerFunc(postHandler.GetPostsPublic)))
	mux.Handle("GET /api/v1/users/me/posts", authMiddleware(false)(requirePermissions(log, false, permissions.PostReadOwn)(http.HandlerFunc(postHandler.GetOwnPosts))))
	mux.Handle("DELETE /api/v1/posts/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.PostDeleteAny, permissions.PostDeleteOwn)(http.HandlerFunc(postHandler.Delete))))
	mux.Handle("PATCH /api/v1/posts/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.PostUpdateAny, permissions.PostUpdateOwn)(http.HandlerFunc(postHandler.Update))))
	mux.Handle("GET /api/v1/posts/{id}", authMiddleware(true)(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("DELETE /api/v1/posts/{id}/photo", authMiddleware(false)(requirePermissions(log, false, permissions.PostPhotoDeleteAny, permissions.PostPhotoDeleteOwn)(http.HandlerFunc(postHandler.RemovePhoto))))
	mux.Handle("PUT /api/v1/posts/{id}/photo", authMiddleware(false)(requirePermissions(log, false, permissions.PostPhotoUpdateAny, permissions.PostPhotoUpdateOwn)(http.HandlerFunc(postHandler.UpdatePhoto))))
	mux.Handle("PATCH /api/v1/posts/{id}/moderation", authMiddleware(false)(requirePermissions(log, false, permissions.PostVerify)(http.HandlerFunc(postHandler.ChangeModerationStatus))))
	mux.Handle("PATCH /api/v1/posts/{id}/return", authMiddleware(false)(requirePermissions(log, false, permissions.PostMarkReturnedAny, permissions.PostMarkReturnedOwn)(http.HandlerFunc(postHandler.ReturnToOwner))))
	mux.Handle("POST /api/v1/posts/similar", authMiddleware(true)(http.HandlerFunc(postHandler.GetSimilar)))

	// Post conversation
	mux.Handle("POST /api/v1/posts/{postId}/contact", authMiddleware(false)(requirePermissions(log, false, permissions.ConversationCreate)(http.HandlerFunc(conversationHandler.CreateConversation))))
	mux.Handle("GET /api/v1/conversations", authMiddleware(false)(requirePermissions(log, false, permissions.ConversationReadOwn)(http.HandlerFunc(conversationHandler.GetMyConversations))))
	mux.Handle("GET /api/v1/conversations/{conversationId}", authMiddleware(false)(requirePermissions(log, false, permissions.ConversationReadOwn)(http.HandlerFunc(conversationHandler.GetConversation))))
	mux.Handle("GET /api/v1/conversations/unread_count", authMiddleware(false)(requirePermissions(log, false, permissions.ConversationReadOwn)(http.HandlerFunc(conversationHandler.GetTotalUnreadCount))))
	mux.Handle("POST /api/v1/conversations/{conversationId}/messages", authMiddleware(false)(requirePermissions(log, false, permissions.ConversationMessageSend)(http.HandlerFunc(conversationHandler.SendMessage))))
	mux.Handle("PATCH /api/v1/conversations/{conversationId}/messages/read", authMiddleware(false)(requirePermissions(log, false, permissions.ConversationMessageMarkAsRead)(http.HandlerFunc(conversationHandler.MarkAsRead))))

	// User
	mux.Handle("POST /api/v1/users", authMiddleware(true)(http.HandlerFunc(authHandler.Register)))
	mux.Handle("GET /api/v1/users/me", authMiddleware(false)(requirePermissions(log, false, permissions.UserReadOwn)(http.HandlerFunc(userHandler.GetOwnUser))))
	mux.Handle("GET /api/v1/users", authMiddleware(false)(requirePermissions(log, false, permissions.UserReadAll)(http.HandlerFunc(userHandler.GetUsers))))
	mux.Handle("GET /api/v1/users/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.UserReadOther)(http.HandlerFunc(userHandler.GetUserByID))))
	mux.Handle("PUT /api/v1/users/me/extensions", authMiddleware(false)(requirePermissions(log, false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.AssignExtensionsOwn))))
	mux.Handle("PATCH /api/v1/users/me", authMiddleware(false)(requirePermissions(log, false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnProfile))))
	mux.Handle("PUT /api/v1/users/me/avatar", authMiddleware(false)(requirePermissions(log, false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnAvatar))))
	mux.Handle("DELETE /api/v1/users/me/avatar", authMiddleware(false)(requirePermissions(log, false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.RemoveOwnAvatar))))
	mux.Handle("DELETE /api/v1/users/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.UserDeleteAnyAdmin, permissions.UserDeleteAnyUser)(http.HandlerFunc(authHandler.DeleteAccount))))
	mux.Handle("DELETE /api/v1/users/me", authMiddleware(false)(requirePermissions(log, false, permissions.UserDeleteOwn)(http.HandlerFunc(authHandler.DeleteOwnAccount))))

	// Room
	mux.Handle("POST /api/v1/rooms", authMiddleware(false)(requirePermissions(log, false, permissions.RoomCreate)(http.HandlerFunc(roomHandler.Create))))
	mux.Handle("GET /api/v1/rooms", authMiddleware(true)(http.HandlerFunc(roomHandler.GetRooms)))
	mux.Handle("PATCH /api/v1/rooms/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.RoomUpdate)(http.HandlerFunc(roomHandler.Update))))
	mux.Handle("DELETE /api/v1/rooms/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.RoomDelete)(http.HandlerFunc(roomHandler.Delete))))

	// Subject
	mux.Handle("POST /api/v1/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.SubjectCreate)(http.HandlerFunc(subjectHandler.Create))))
	mux.Handle("GET /api/v1/subjects", authMiddleware(true)(http.HandlerFunc(subjectHandler.GetSubjects)))
	mux.Handle("PATCH /api/v1/subjects/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.SubjectUpdate)(http.HandlerFunc(subjectHandler.Update))))
	mux.Handle("DELETE /api/v1/subjects/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.SubjectDelete)(http.HandlerFunc(subjectHandler.Delete))))

	// Student group
	mux.Handle("POST /api/v1/student_groups/{id}/advisor", authMiddleware(false)(requirePermissions(log, false, permissions.StudentGroupAdvisorAssign)(http.HandlerFunc(studentGroupHandler.AssignAdvisor))))
	mux.Handle("GET /api/v1/student_groups", authMiddleware(true)(http.HandlerFunc(studentGroupHandler.GetStudentGroups)))
	mux.Handle("GET /api/v1/student_groups/{id}", authMiddleware(true)(http.HandlerFunc(studentGroupHandler.GetStudentGroupByID)))
	mux.Handle("PATCH /api/v1/student_groups/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.StudentGroupUpdate)(http.HandlerFunc(studentGroupHandler.Update))))
	mux.Handle("DELETE /api/v1/student_groups/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.StudentGroupDelete)(http.HandlerFunc(studentGroupHandler.Delete))))
	mux.Handle("POST /api/v1/student_groups", authMiddleware(false)(requirePermissions(log, false, permissions.StudentGroupCreate)(http.HandlerFunc(studentGroupHandler.Create))))
	mux.Handle("DELETE /api/v1/student_groups/{id}/advisor", authMiddleware(false)(requirePermissions(log, false, permissions.StudentGroupAdvisorUnassignAny, permissions.StudentGroupAdvisorUnassignOwn)(http.HandlerFunc(studentGroupHandler.UnassignAdvisor))))
	mux.Handle("GET /api/v1/student_groups/{id}/advisor", authMiddleware(false)(requirePermissions(log, false, permissions.StudentGroupAdvisorRead)(http.HandlerFunc(studentGroupHandler.GetAdvisorByGroupID))))

	// Teacher
	mux.Handle("GET /api/v1/teachers/{id}/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectReadAny)(http.HandlerFunc(teacherHandler.GetSubjects))))
	mux.Handle("GET /api/v1/teachers/me/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectReadOwn)(http.HandlerFunc(teacherHandler.GetSubjectsOwn))))
	mux.Handle("POST /api/v1/teachers/{id}/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectAddAny)(http.HandlerFunc(teacherHandler.AddSubjects))))
	mux.Handle("POST /api/v1/teachers/me/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectAddOwn)(http.HandlerFunc(teacherHandler.AddSubjectsOwn))))
	mux.Handle("PUT /api/v1/teachers/{id}/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectAssignAny)(http.HandlerFunc(teacherHandler.AssignSubjects))))
	mux.Handle("PUT /api/v1/teachers/me/subjects", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectAssignOwn)(http.HandlerFunc(teacherHandler.AssignSubjectsOwn))))
	mux.Handle("DELETE /api/v1/teachers/{userId}/subjects/{subjectId}", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectUnassignAny)(http.HandlerFunc(teacherHandler.UnassignSubject))))
	mux.Handle("DELETE /api/v1/teachers/me/subjects/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherSubjectUnassignOwn)(http.HandlerFunc(teacherHandler.UnassignSubjectOwn))))
	mux.Handle("GET /api/v1/teachers/{id}/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherClassroomReadAny)(http.HandlerFunc(teacherHandler.GetClassroom))))
	mux.Handle("GET /api/v1/teachers/me/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherClassroomReadOwn)(http.HandlerFunc(teacherHandler.GetClassroomOwn))))
	mux.Handle("PUT /api/v1/teachers/{id}/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherClassroomAssignAny)(http.HandlerFunc(teacherHandler.AssignClassroom))))
	mux.Handle("PUT /api/v1/teachers/me/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherClassroomAssignOwn)(http.HandlerFunc(teacherHandler.AssignClassroomOwn))))
	mux.Handle("DELETE /api/v1/teachers/{id}/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherClassroomUnassignAny)(http.HandlerFunc(teacherHandler.UnassignClassroom))))
	mux.Handle("DELETE /api/v1/teachers/me/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherClassroomUnassignOwn)(http.HandlerFunc(teacherHandler.UnassignClassroomOwn))))
	mux.Handle("GET /api/v1/teachers/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherReadOther)(http.HandlerFunc(teacherHandler.GetTeacherByID))))
	mux.Handle("GET /api/v1/teachers/me", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherReadOwn)(http.HandlerFunc(teacherHandler.GetOwn))))
	mux.Handle("GET /api/v1/teachers/me/student_groups", authMiddleware(false)(requirePermissions(log, false, permissions.TeacherStudentGroupReadOwn)(http.HandlerFunc(teacherHandler.GetStudentGroupsOwn))))

	// Parent
	mux.Handle("GET /api/v1/parents/{id}/students", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentReadAny)(http.HandlerFunc(parentHandler.GetStudents))))
	mux.Handle("GET /api/v1/parents/me/students", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentReadOwn)(http.HandlerFunc(parentHandler.GetStudentsOwn))))
	mux.Handle("POST /api/v1/parents/{id}/students", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentAddAny)(http.HandlerFunc(parentHandler.AddStudents))))
	mux.Handle("POST /api/v1/parents/me/students", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentAddOwn)(http.HandlerFunc(parentHandler.AddStudentsOwn))))
	mux.Handle("DELETE /api/v1/parents/{parentId}/students/{studentId}", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentUnassignAny)(http.HandlerFunc(parentHandler.UnassignStudent))))
	mux.Handle("DELETE /api/v1/parents/me/students/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentUnassignOwn)(http.HandlerFunc(parentHandler.UnassignStudentOwn))))
	mux.Handle("GET /api/v1/parents/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.ParentReadOther)(http.HandlerFunc(parentHandler.GetParentByID))))
	mux.Handle("GET /api/v1/parents/me", authMiddleware(false)(requirePermissions(log, false, permissions.ParentReadOwn)(http.HandlerFunc(parentHandler.GetOwn))))
	mux.Handle("GET /api/v1/parents/me/student_groups", authMiddleware(false)(requirePermissions(log, false, permissions.ParentStudentGroupReadOwn)(http.HandlerFunc(parentHandler.GetStudentGroupsOwn))))

	// User roles
	mux.Handle("PUT /api/v1/users/{id}/roles", authMiddleware(false)(requirePermissions(log, false, permissions.RoleAdminAssign, permissions.RoleUserAssign)(http.HandlerFunc(userHandler.AssignRoles))))
	mux.Handle("PUT /api/v1/users/{id}/roles/non_admin", authMiddleware(false)(requirePermissions(log, false, permissions.RoleUserAssign)(http.HandlerFunc(userHandler.AssignNonAdminRoles))))
	mux.Handle("POST /api/v1/users/{id}/roles", authMiddleware(false)(requirePermissions(log, false, permissions.RoleAdminAdd, permissions.RoleUserAdd)(http.HandlerFunc(userHandler.AddRoles))))
	mux.Handle("DELETE /api/v1/users/{userId}/roles/{roleId}", authMiddleware(false)(requirePermissions(log, false, permissions.RoleAdminUnassign, permissions.RoleUserUnassign)(http.HandlerFunc(userHandler.RemoveRole))))
	mux.Handle("GET /api/v1/users/{id}/roles", authMiddleware(false)(requirePermissions(log, false, permissions.RoleReadAny)(http.HandlerFunc(userHandler.GetRoles))))
	mux.Handle("GET /api/v1/users/me/roles", authMiddleware(false)(requirePermissions(log, false, permissions.RoleReadOwn)(http.HandlerFunc(userHandler.GetOwnRoles))))

	// Invite token
	mux.Handle("POST /api/v1/tokens/invite", authMiddleware(false)(requirePermissions(log, false, permissions.TokenInviteAdminCreate, permissions.TokenInviteUserCreate)(http.HandlerFunc(inviteHandler.Create))))
	mux.Handle("DELETE /api/v1/tokens/invite/{token}", authMiddleware(false)(requirePermissions(log, false, permissions.TokenInviteAdminDelete, permissions.TokenInviteUserDelete)(http.HandlerFunc(inviteHandler.Revoke))))
	mux.Handle("GET /api/v1/tokens/invite/{token}/roles", authMiddleware(true)(http.HandlerFunc(inviteHandler.GetRoles)))
	mux.Handle("GET /api/v1/tokens/invite/{token}/email", authMiddleware(true)(http.HandlerFunc(inviteHandler.GetEmail)))

	// Invite request
	mux.Handle("POST /api/v1/invite/request/student", authMiddleware(true)(http.HandlerFunc(inviteHandler.MakeStudentInviteRequest)))
	mux.Handle("POST /api/v1/invite/request/parent", authMiddleware(true)(http.HandlerFunc(inviteHandler.MakeParentInviteRequest)))

	// Auth
	mux.Handle("POST /api/v1/auth/login", authMiddleware(true)(http.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /api/v1/auth/logout", authMiddleware(false)(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("POST /api/v1/auth/forgot-password", authMiddleware(true)(http.HandlerFunc(authHandler.ForgotPassword)))
	mux.Handle("POST /api/v1/auth/reset-password", authMiddleware(true)(http.HandlerFunc(authHandler.ResetPassword)))

	// Student
	mux.Handle("GET /api/v1/students/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.StudentReadOther)(http.HandlerFunc(studentHandler.GetStudentByID))))
	mux.Handle("GET /api/v1/students/me", authMiddleware(false)(requirePermissions(log, false, permissions.StudentReadOwn)(http.HandlerFunc(studentHandler.GetOwn))))
	mux.Handle("GET /api/v1/students/{id}/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.StudentClassroomReadAny)(http.HandlerFunc(studentHandler.GetClassroom))))
	mux.Handle("GET /api/v1/students/me/classroom", authMiddleware(false)(requirePermissions(log, false, permissions.StudentClassroomReadOwn)(http.HandlerFunc(studentHandler.GetClassroomOwn))))
	mux.Handle("GET /api/v1/students/{id}/advisor", authMiddleware(false)(requirePermissions(log, false, permissions.StudentAdvisorReadAny)(http.HandlerFunc(studentHandler.GetAdvisor))))
	mux.Handle("GET /api/v1/students/me/advisor", authMiddleware(false)(requirePermissions(log, false, permissions.StudentAdvisorReadOwn)(http.HandlerFunc(studentHandler.GetAdvisorOwn))))
	mux.Handle("GET /api/v1/students/{id}/parents", authMiddleware(false)(requirePermissions(log, false, permissions.StudentParentReadAny)(http.HandlerFunc(studentHandler.GetParents))))
	mux.Handle("GET /api/v1/students/me/parents", authMiddleware(false)(requirePermissions(log, false, permissions.StudentParentReadOwn)(http.HandlerFunc(studentHandler.GetParentsOwn))))
	mux.Handle("GET /api/v1/students/me/student_group", authMiddleware(false)(requirePermissions(log, false, permissions.StudentStudentGroupReadOwn)(http.HandlerFunc(studentHandler.GetStudentGroupOwn))))

	// Institution administrator
	mux.Handle("GET /api/v1/institution_administrators/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.InstitutionAdministratorReadOther)(http.HandlerFunc(institutionAdministratorHandler.GetInstitutionAdministratorByID))))
	mux.Handle("GET /api/v1/institution_administrators/me", authMiddleware(false)(requirePermissions(log, false, permissions.InstitutionAdministratorReadOwn)(http.HandlerFunc(institutionAdministratorHandler.GetOwn))))
	mux.Handle("PUT /api/v1/institution_administrators/{id}/position", authMiddleware(false)(requirePermissions(log, false, permissions.InstitutionAdministratorPositionAssign)(http.HandlerFunc(institutionAdministratorHandler.AssignPosition))))
	mux.Handle("GET /api/v1/institution_administrators/{id}/position", authMiddleware(false)(requirePermissions(log, false, permissions.InstitutionAdministratorPositionRead)(http.HandlerFunc(institutionAdministratorHandler.GetPosition))))

	// Staff
	mux.Handle("GET /api/v1/staff/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.StaffReadOther)(http.HandlerFunc(staffHandler.GetStaffByID))))
	mux.Handle("GET /api/v1/staff/me", authMiddleware(false)(requirePermissions(log, false, permissions.StaffReadOwn)(http.HandlerFunc(staffHandler.GetOwn))))
	mux.Handle("PUT /api/v1/staff/{id}/position", authMiddleware(false)(requirePermissions(log, false, permissions.StaffPositionAssign)(http.HandlerFunc(staffHandler.AssignPosition))))
	mux.Handle("GET /api/v1/staff/{id}/position", authMiddleware(false)(requirePermissions(log, false, permissions.StaffPositionRead)(http.HandlerFunc(staffHandler.GetPosition))))

	// Position institution administrator
	mux.Handle("POST /api/v1/institution_administrators/positions", authMiddleware(false)(requirePermissions(log, false, permissions.PositionInstitutionAdministratorCreate)(http.HandlerFunc(institutionAdministratorPositionHandler.Create))))
	mux.HandleFunc("GET /api/v1/institution_administrators/positions", institutionAdministratorPositionHandler.GetAll)
	mux.Handle("PATCH /api/v1/institution_administrators/positions/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.PositionInstitutionAdministratorUpdate)(http.HandlerFunc(institutionAdministratorPositionHandler.Update))))
	mux.Handle("DELETE /api/v1/institution_administrators/positions/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.PositionInstitutionAdministratorDelete)(http.HandlerFunc(institutionAdministratorPositionHandler.Delete))))

	// Position staff
	mux.Handle("POST /api/v1/staff/positions", authMiddleware(false)(requirePermissions(log, false, permissions.PositionStaffCreate)(http.HandlerFunc(staffPositionHandler.Create))))
	mux.HandleFunc("GET /api/v1/staff/positions", staffPositionHandler.GetAll)
	mux.Handle("PATCH /api/v1/staff/positions/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.PositionStaffUpdate)(http.HandlerFunc(staffPositionHandler.Update))))
	mux.Handle("DELETE /api/v1/staff/positions/{id}", authMiddleware(false)(requirePermissions(log, false, permissions.PositionStaffDelete)(http.HandlerFunc(staffPositionHandler.Delete))))

	// Documents
	mux.Handle("PUT /api/v1/documents/privacy", authMiddleware(false)(requirePermissions(log, false, permissions.DocumentPrivacyUpload)(http.HandlerFunc(documentHandler.UploadPrivacy))))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}
