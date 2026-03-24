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
	authMiddleware func(http.Handler) http.Handler,
	requireRoles func(bool, ...string) func(http.Handler) http.Handler,
	requirePermissions func(bool, ...string) func(http.Handler) http.Handler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	studentGroupHandler *handler.StudentGroupHandler,
	roomHandler *handler.RoomHandler,
	subjectHandler *handler.SubjectHandler,
	studentHandler *handler.StudentHandler,
	teacherHandler *handler.TeacherHandler,
	parentHandler *handler.ParentHandler,
	staffHandler *handler.StaffHandler,
	institutionAdministratorHandler *handler.InstitutionAdministratorHandler,
	roleHandler *handler.RoleHandler,
	inviteHandler *handler.InviteHandler,
) {
	// Public routes (no auth required)
	mux.HandleFunc("POST /api/v1/users", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/posts/public", postHandler.GetPostsPublic)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("GET /api/v1/tokens/invite/{token}", inviteHandler.GetRoles)

	// Secure routes (auth required)

	// User
	mux.Handle("PATCH /api/v1/users/me", authMiddleware(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnProfile))))
	mux.Handle("PUT /api/v1/users/me/avatar", authMiddleware(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnAvatar))))
	mux.Handle("DELETE /api/v1/users/me/avatar", authMiddleware(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.RemoveOwnAvatar))))
	mux.Handle("GET /api/v1/users/{id}", authMiddleware(requirePermissions(false, permissions.UserReadOther)(http.HandlerFunc(userHandler.GetUserByID))))
	mux.Handle("GET /api/v1/users", authMiddleware(requirePermissions(false, permissions.UserReadAll)(http.HandlerFunc(userHandler.GetUsers))))
	mux.Handle("GET /api/v1/users/me", authMiddleware(requirePermissions(false, permissions.UserReadOwn)(http.HandlerFunc(userHandler.GetOwnUser))))
	// User roles
	mux.Handle("PUT /api/v1/users/{id}/roles", authMiddleware(requirePermissions(false, permissions.RoleAssign)(http.HandlerFunc(userHandler.AssignRoles))))
	mux.Handle("POST /api/v1/users/{id}/roles", authMiddleware(requirePermissions(false, permissions.RoleAdd)(http.HandlerFunc(userHandler.AddRoles))))
	mux.Handle("DELETE /api/v1/users/{userId}/roles/{roleId}", authMiddleware(requirePermissions(false, permissions.RoleDelete)(http.HandlerFunc(userHandler.RemoveRole))))
	mux.Handle("GET /api/v1/users/{id}/roles", authMiddleware(requirePermissions(false, permissions.RoleReadAny)(http.HandlerFunc(userHandler.GetRoles))))
	mux.Handle("GET /api/v1/users/me/roles", authMiddleware(requirePermissions(false, permissions.RoleReadOwn)(http.HandlerFunc(userHandler.GetOwnRoles))))
	// Student groups
	mux.Handle("GET /api/v1/student_groups", authMiddleware(requirePermissions(false, permissions.StudentGroupReadAny)(http.HandlerFunc(studentGroupHandler.GetStudentGroups))))
	mux.Handle("GET /api/v1/student_groups/{id}", authMiddleware(requirePermissions(false, permissions.StudentGroupReadAny)(http.HandlerFunc(studentGroupHandler.GetStudentGroupByID))))
	mux.Handle("GET /api/v1/student_groups/{id}/advisor", authMiddleware(requirePermissions(false, permissions.StudentGroupAdvisorRead)(http.HandlerFunc(studentGroupHandler.GetAdvisorByGroupID))))
	mux.Handle("DELETE /api/v1/student_groups/{id}", authMiddleware(requirePermissions(false, permissions.StudentGroupDelete)(http.HandlerFunc(studentGroupHandler.Delete))))
	mux.Handle("POST /api/v1/student_groups/{id}/advisor", authMiddleware(requirePermissions(false, permissions.StudentGroupAdvisorAssign)(http.HandlerFunc(studentGroupHandler.AssignAdvisor))))
	mux.Handle("DELETE /api/v1/student_groups/{id}/advisor", authMiddleware(requirePermissions(false, permissions.StudentGroupAdvisorUnassignAny, permissions.StudentGroupAdvisorUnassignOwn)(http.HandlerFunc(studentGroupHandler.UnassignAdvisor))))
	mux.Handle("POST /api/v1/student_groups", authMiddleware(requirePermissions(false, permissions.StudentGroupCreate)(http.HandlerFunc(studentGroupHandler.Create))))
	mux.Handle("PATCH /api/v1/student_groups/{id}", authMiddleware(requirePermissions(false, permissions.StudentGroupUpdate)(http.HandlerFunc(studentGroupHandler.Update))))
	// Auth
	mux.Handle("DELETE /api/v1/users/{id}", authMiddleware(requirePermissions(false, permissions.UserDeleteAny)(http.HandlerFunc(authHandler.DeleteAccount))))
	mux.Handle("DELETE /api/v1/users/me", authMiddleware(requirePermissions(false, permissions.UserDeleteOwn)(http.HandlerFunc(authHandler.DeleteOwnAccount))))
	mux.Handle("POST /api/v1/auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))
	// Posts
	mux.Handle("POST /api/v1/posts", authMiddleware(requirePermissions(false, permissions.PostCreate)(http.HandlerFunc(postHandler.Create))))
	mux.Handle("DELETE /api/v1/posts/{id}", authMiddleware(requirePermissions(false, permissions.PostDeleteAny, permissions.PostDeleteOwn)(http.HandlerFunc(postHandler.Delete))))
	mux.Handle("DELETE /api/v1/posts/{id}/photo", authMiddleware(requirePermissions(false, permissions.PostPhotoDeleteAny, permissions.PostPhotoDeleteOwn)(http.HandlerFunc(postHandler.RemovePhoto))))
	mux.Handle("PATCH /api/v1/posts/{id}", authMiddleware(requirePermissions(false, permissions.PostUpdateAny, permissions.PostUpdateOwn)(http.HandlerFunc(postHandler.Update))))
	mux.Handle("GET /api/v1/posts", authMiddleware(requirePermissions(false, permissions.PostReadAny)(http.HandlerFunc(postHandler.GetPosts))))
	mux.Handle("GET /api/v1/users/me/posts", authMiddleware(requirePermissions(false, permissions.PostReadOwn)(http.HandlerFunc(postHandler.GetOwnPosts))))
	mux.Handle("PATCH /api/v1/posts/{id}/verify", authMiddleware(requirePermissions(false, permissions.PostVerify)(http.HandlerFunc(postHandler.Verify))))
	mux.Handle("PATCH /api/v1/posts/{id}/return", authMiddleware(requirePermissions(false, permissions.PostMarkReturnedAny, permissions.PostMarkReturnedOwn)(http.HandlerFunc(postHandler.ReturnToOwner))))
	// Rooms
	mux.Handle("POST /api/v1/rooms", authMiddleware(requirePermissions(false, permissions.RoomCreate)(http.HandlerFunc(roomHandler.Create))))
	mux.Handle("DELETE /api/v1/rooms/{id}", authMiddleware(requirePermissions(false, permissions.RoomDelete)(http.HandlerFunc(roomHandler.Delete))))
	mux.Handle("PATCH /api/v1/rooms/{id}", authMiddleware(requirePermissions(false, permissions.RoomUpdate)(http.HandlerFunc(roomHandler.Update))))
	mux.Handle("GET /api/v1/rooms", authMiddleware(requirePermissions(false, permissions.RoomRead)(http.HandlerFunc(roomHandler.GetRooms))))
	// Subjects
	mux.Handle("POST /api/v1/subjects", authMiddleware(requirePermissions(false, permissions.SubjectCreate)(http.HandlerFunc(subjectHandler.Create))))
	mux.Handle("DELETE /api/v1/subjects/{id}", authMiddleware(requirePermissions(false, permissions.SubjectDelete)(http.HandlerFunc(subjectHandler.Delete))))
	mux.Handle("PATCH /api/v1/subjects/{id}", authMiddleware(requirePermissions(false, permissions.SubjectUpdate)(http.HandlerFunc(subjectHandler.Update))))
	mux.Handle("GET /api/v1/subjects", authMiddleware(requirePermissions(false, permissions.SubjectRead)(http.HandlerFunc(subjectHandler.GetSubjects))))
	// Students
	mux.Handle("GET /api/v1/students/{id}", authMiddleware(requirePermissions(false, permissions.StudentReadOther)(http.HandlerFunc(studentHandler.GetStudentByID))))
	mux.Handle("GET /api/v1/students/me", authMiddleware(requirePermissions(false, permissions.StudentReadOwn)(http.HandlerFunc(studentHandler.GetOwn))))
	mux.Handle("GET /api/v1/students/{id}/classroom", authMiddleware(requirePermissions(false, permissions.StudentClassroomReadAny)(http.HandlerFunc(studentHandler.GetClassroom))))
	mux.Handle("GET /api/v1/students/me/classroom", authMiddleware(requirePermissions(false, permissions.StudentClassroomReadOwn)(http.HandlerFunc(studentHandler.GetClassroomOwn))))
	mux.Handle("GET /api/v1/students/{id}/advisor", authMiddleware(requirePermissions(false, permissions.StudentAdvisorReadAny)(http.HandlerFunc(studentHandler.GetAdvisor))))
	mux.Handle("GET /api/v1/students/me/advisor", authMiddleware(requirePermissions(false, permissions.StudentAdvisorReadOwn)(http.HandlerFunc(studentHandler.GetAdvisorOwn))))
	mux.Handle("GET /api/v1/students/{id}/parents", authMiddleware(requirePermissions(false, permissions.StudentParentReadAny)(http.HandlerFunc(studentHandler.GetParents))))
	mux.Handle("GET /api/v1/students/me/parents", authMiddleware(requirePermissions(false, permissions.StudentParentReadOwn)(http.HandlerFunc(studentHandler.GetParentsOwn))))
	mux.Handle("GET /api/v1/students/me/student_group", authMiddleware(requirePermissions(false, permissions.StudentStudentGroupReadOwn)(http.HandlerFunc(studentHandler.GetStudentGroupOwn))))
	// Teacher
	mux.Handle("GET /api/v1/teachers/{id}", authMiddleware(requirePermissions(false, permissions.TeacherReadOther)(http.HandlerFunc(teacherHandler.GetTeacherByID))))
	mux.Handle("GET /api/v1/teachers/me", authMiddleware(requirePermissions(false, permissions.TeacherReadOwn)(http.HandlerFunc(teacherHandler.GetOwn))))
	mux.Handle("GET /api/v1/teachers/{id}/classroom", authMiddleware(requirePermissions(false, permissions.TeacherClassroomReadAny)(http.HandlerFunc(teacherHandler.GetClassroom))))
	mux.Handle("GET /api/v1/teachers/me/classroom", authMiddleware(requirePermissions(false, permissions.TeacherClassroomReadOwn)(http.HandlerFunc(teacherHandler.GetClassroomOwn))))
	mux.Handle("GET /api/v1/teachers/{id}/subjects", authMiddleware(requirePermissions(false, permissions.TeacherSubjectReadAny)(http.HandlerFunc(teacherHandler.GetSubjects))))
	mux.Handle("GET /api/v1/teachers/me/subjects", authMiddleware(requirePermissions(false, permissions.TeacherSubjectReadOwn)(http.HandlerFunc(teacherHandler.GetSubjectsOwn))))
	mux.Handle("PUT /api/v1/teachers/{id}/classroom", authMiddleware(requirePermissions(false, permissions.TeacherClassroomAssignAny)(http.HandlerFunc(teacherHandler.AssignClassroom))))
	mux.Handle("PUT /api/v1/teachers/me/classroom", authMiddleware(requirePermissions(false, permissions.TeacherClassroomAssignOwn)(http.HandlerFunc(teacherHandler.AssignClassroomOwn))))
	mux.Handle("DELETE /api/v1/teachers/{id}/classroom", authMiddleware(requirePermissions(false, permissions.TeacherClassroomUnassignAny)(http.HandlerFunc(teacherHandler.UnassignClassroom))))
	mux.Handle("DELETE /api/v1/teachers/me/classroom", authMiddleware(requirePermissions(false, permissions.TeacherClassroomUnassignOwn)(http.HandlerFunc(teacherHandler.UnassignClassroomOwn))))
	mux.Handle("POST /api/v1/teachers/{id}/subjects", authMiddleware(requirePermissions(false, permissions.TeacherSubjectAddAny)(http.HandlerFunc(teacherHandler.AddSubjects))))
	mux.Handle("POST /api/v1/teachers/me/subjects", authMiddleware(requirePermissions(false, permissions.TeacherSubjectAddOwn)(http.HandlerFunc(teacherHandler.AddSubjectsOwn))))
	mux.Handle("PUT /api/v1/teachers/{id}/subjects", authMiddleware(requirePermissions(false, permissions.TeacherSubjectAssignAny)(http.HandlerFunc(teacherHandler.AssignSubjects))))
	mux.Handle("PUT /api/v1/teachers/me/subjects", authMiddleware(requirePermissions(false, permissions.TeacherSubjectAssignOwn)(http.HandlerFunc(teacherHandler.AssignSubjectsOwn))))
	mux.Handle("DELETE /api/v1/teachers/{userId}/subjects/{subjectId}", authMiddleware(requirePermissions(false, permissions.TeacherSubjectUnassignAny)(http.HandlerFunc(teacherHandler.UnassignSubject))))
	mux.Handle("DELETE /api/v1/teachers/me/subjects/{id}", authMiddleware(requirePermissions(false, permissions.TeacherSubjectUnassignOwn)(http.HandlerFunc(teacherHandler.UnassignSubjectOwn))))
	mux.Handle("GET /api/v1/teachers/me/student_groups", authMiddleware(requirePermissions(false, permissions.TeacherStudentGroupReadOwn)(http.HandlerFunc(teacherHandler.GetStudentGroupsOwn))))
	// Parents
	mux.Handle("GET /api/v1/parents/{id}", authMiddleware(requirePermissions(false, permissions.ParentReadOther)(http.HandlerFunc(parentHandler.GetParentByID))))
	mux.Handle("GET /api/v1/parents/me", authMiddleware(requirePermissions(false, permissions.ParentReadOwn)(http.HandlerFunc(parentHandler.GetOwn))))
	mux.Handle("GET /api/v1/parents/{id}/students", authMiddleware(requirePermissions(false, permissions.ParentStudentReadAny)(http.HandlerFunc(parentHandler.GetStudents))))
	mux.Handle("GET /api/v1/parents/me/students", authMiddleware(requirePermissions(false, permissions.ParentStudentReadOwn)(http.HandlerFunc(parentHandler.GetStudentsOwn))))
	mux.Handle("GET /api/v1/parents/me/student_groups", authMiddleware(requirePermissions(false, permissions.ParentStudentGroupReadOwn)(http.HandlerFunc(parentHandler.GetStudentGroupsOwn))))
	mux.Handle("POST /api/v1/parents/{id}/students", authMiddleware(requirePermissions(false, permissions.ParentStudentAddAny)(http.HandlerFunc(parentHandler.AddStudents))))
	mux.Handle("POST /api/v1/parents/me/students", authMiddleware(requirePermissions(false, permissions.ParentStudentAddOwn)(http.HandlerFunc(parentHandler.AddStudentsOwn))))
	mux.Handle("DELETE /api/v1/parents/{parentId}/students/{studentId}", authMiddleware(requirePermissions(false, permissions.ParentStudentUnassignAny)(http.HandlerFunc(parentHandler.UnassignStudent))))
	mux.Handle("DELETE /api/v1/parents/me/students/{id}", authMiddleware(requirePermissions(false, permissions.ParentStudentUnassignOwn)(http.HandlerFunc(parentHandler.UnassignStudentOwn))))
	// Staff
	mux.Handle("GET /api/v1/staff/{id}", authMiddleware(requirePermissions(false, permissions.StaffReadOther)(http.HandlerFunc(staffHandler.GetStaffByID))))
	mux.Handle("GET /api/v1/staff/me", authMiddleware(requirePermissions(false, permissions.StaffReadOwn)(http.HandlerFunc(staffHandler.GetOwn))))
	mux.Handle("PUT /api/v1/staff/{id}/position", authMiddleware(requirePermissions(false, permissions.StaffPositionAssign)(http.HandlerFunc(staffHandler.AssignPosition))))
	mux.Handle("GET /api/v1/staff/{id}/position", authMiddleware(requirePermissions(false, permissions.StaffPositionRead)(http.HandlerFunc(staffHandler.GetPosition))))
	// Institution administrator
	mux.Handle("GET /api/v1/institution_administrators/{id}", authMiddleware(requirePermissions(false, permissions.InstitutionAdministratorReadOther)(http.HandlerFunc(institutionAdministratorHandler.GetInstitutionAdministratorByID))))
	mux.Handle("GET /api/v1/institution_administrators/me", authMiddleware(requirePermissions(false, permissions.InstitutionAdministratorReadOwn)(http.HandlerFunc(institutionAdministratorHandler.GetOwn))))
	mux.Handle("PUT /api/v1/institution_administrators/{id}/position", authMiddleware(requirePermissions(false, permissions.InstitutionAdministratorPositionAssign)(http.HandlerFunc(institutionAdministratorHandler.AssignPosition))))
	mux.Handle("GET /api/v1/institution_administrators/{id}/position", authMiddleware(requirePermissions(false, permissions.InstitutionAdministratorPositionRead)(http.HandlerFunc(institutionAdministratorHandler.GetPosition))))
	// Invite tokens
	mux.Handle("POST /api/v1/tokens/invite", authMiddleware(requirePermissions(false, permissions.TokenInviteAdminCreate, permissions.TokenInviteUserCreate)(http.HandlerFunc(inviteHandler.Create))))
	mux.Handle("DELETE /api/v1/tokens/invite/{token}", authMiddleware(requirePermissions(false, permissions.TokenInviteAdminDelete, permissions.TokenInviteUserDelete)(http.HandlerFunc(inviteHandler.Revoke))))
	// Roles
	// mux.Handle("GET /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.GetPermissions)))
	// mux.Handle("PUT /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.AssignPermissions)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}
