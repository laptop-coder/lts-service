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
	staffPositionHandler *handler.StaffPositionHandler,
	institutionAdministratorPositionHandler *handler.InstitutionAdministratorPositionHandler,
) {
	// Public routes (no auth required)
	// TODO: split this routes into categories (i.e. mix with secure routes)
	// TODO: use authMiddleware with allowUnauthorized = true
	mux.HandleFunc("POST /api/v1/users", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/posts/public", postHandler.GetPostsPublic)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("GET /api/v1/tokens/invite/{token}/roles", inviteHandler.GetRoles)
	mux.HandleFunc("GET /api/v1/tokens/invite/{token}/email", inviteHandler.GetEmail)
	mux.HandleFunc("POST /api/v1/invite/request/student", inviteHandler.MakeStudentInviteRequest)

	// Secure routes (auth required)

	// User
	mux.Handle("PATCH /api/v1/users/me", authMiddleware(false)(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnProfile))))
	mux.Handle("PUT /api/v1/users/me/avatar", authMiddleware(false)(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnAvatar))))
	mux.Handle("DELETE /api/v1/users/me/avatar", authMiddleware(false)(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.RemoveOwnAvatar))))
	mux.Handle("GET /api/v1/users/{id}", authMiddleware(false)(requirePermissions(false, permissions.UserReadOther)(http.HandlerFunc(userHandler.GetUserByID))))
	mux.Handle("GET /api/v1/users", authMiddleware(false)(requirePermissions(false, permissions.UserReadAll)(http.HandlerFunc(userHandler.GetUsers))))
	mux.Handle("GET /api/v1/users/me", authMiddleware(false)(requirePermissions(false, permissions.UserReadOwn)(http.HandlerFunc(userHandler.GetOwnUser))))
	// User roles
	mux.Handle("PUT /api/v1/users/{id}/roles", authMiddleware(false)(requirePermissions(false, permissions.RoleAdminAssign, permissions.RoleUserAssign)(http.HandlerFunc(userHandler.AssignRoles))))
	mux.Handle("PUT /api/v1/users/{id}/roles/non_admin", authMiddleware(false)(requirePermissions(false, permissions.RoleUserAssign)(http.HandlerFunc(userHandler.AssignNonAdminRoles))))
	mux.Handle("POST /api/v1/users/{id}/roles", authMiddleware(false)(requirePermissions(false, permissions.RoleAdminAdd, permissions.RoleUserAdd)(http.HandlerFunc(userHandler.AddRoles))))
	mux.Handle("DELETE /api/v1/users/{userId}/roles/{roleId}", authMiddleware(false)(requirePermissions(false, permissions.RoleAdminUnassign, permissions.RoleUserUnassign)(http.HandlerFunc(userHandler.RemoveRole))))
	mux.Handle("GET /api/v1/users/{id}/roles", authMiddleware(false)(requirePermissions(false, permissions.RoleReadAny)(http.HandlerFunc(userHandler.GetRoles))))
	mux.Handle("GET /api/v1/users/me/roles", authMiddleware(false)(requirePermissions(false, permissions.RoleReadOwn)(http.HandlerFunc(userHandler.GetOwnRoles))))
	// Student groups
	mux.HandleFunc("GET /api/v1/student_groups", studentGroupHandler.GetStudentGroups)
	mux.HandleFunc("GET /api/v1/student_groups/{id}", studentGroupHandler.GetStudentGroupByID)
	mux.Handle("GET /api/v1/student_groups/{id}/advisor", authMiddleware(false)(requirePermissions(false, permissions.StudentGroupAdvisorRead)(http.HandlerFunc(studentGroupHandler.GetAdvisorByGroupID))))
	mux.Handle("DELETE /api/v1/student_groups/{id}", authMiddleware(false)(requirePermissions(false, permissions.StudentGroupDelete)(http.HandlerFunc(studentGroupHandler.Delete))))
	mux.Handle("POST /api/v1/student_groups/{id}/advisor", authMiddleware(false)(requirePermissions(false, permissions.StudentGroupAdvisorAssign)(http.HandlerFunc(studentGroupHandler.AssignAdvisor))))
	mux.Handle("DELETE /api/v1/student_groups/{id}/advisor", authMiddleware(false)(requirePermissions(false, permissions.StudentGroupAdvisorUnassignAny, permissions.StudentGroupAdvisorUnassignOwn)(http.HandlerFunc(studentGroupHandler.UnassignAdvisor))))
	mux.Handle("POST /api/v1/student_groups", authMiddleware(false)(requirePermissions(false, permissions.StudentGroupCreate)(http.HandlerFunc(studentGroupHandler.Create))))
	mux.Handle("PATCH /api/v1/student_groups/{id}", authMiddleware(false)(requirePermissions(false, permissions.StudentGroupUpdate)(http.HandlerFunc(studentGroupHandler.Update))))
	// Auth
	mux.Handle("DELETE /api/v1/users/{id}", authMiddleware(false)(requirePermissions(false, permissions.UserDeleteAny)(http.HandlerFunc(authHandler.DeleteAccount))))
	mux.Handle("DELETE /api/v1/users/me", authMiddleware(false)(requirePermissions(false, permissions.UserDeleteOwn)(http.HandlerFunc(authHandler.DeleteOwnAccount))))
	mux.Handle("POST /api/v1/auth/logout", authMiddleware(false)(http.HandlerFunc(authHandler.Logout)))
	// Posts
	mux.Handle("POST /api/v1/posts", authMiddleware(false)(requirePermissions(false, permissions.PostCreate)(http.HandlerFunc(postHandler.Create))))
	mux.Handle("DELETE /api/v1/posts/{id}", authMiddleware(false)(requirePermissions(false, permissions.PostDeleteAny, permissions.PostDeleteOwn)(http.HandlerFunc(postHandler.Delete))))
	mux.Handle("DELETE /api/v1/posts/{id}/photo", authMiddleware(false)(requirePermissions(false, permissions.PostPhotoDeleteAny, permissions.PostPhotoDeleteOwn)(http.HandlerFunc(postHandler.RemovePhoto))))
	mux.Handle("PATCH /api/v1/posts/{id}", authMiddleware(false)(requirePermissions(false, permissions.PostUpdateAny, permissions.PostUpdateOwn)(http.HandlerFunc(postHandler.Update))))
	mux.Handle("GET /api/v1/posts", authMiddleware(false)(requirePermissions(false, permissions.PostReadAny)(http.HandlerFunc(postHandler.GetPosts))))
	mux.Handle("GET /api/v1/posts/{id}", authMiddleware(true)(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("GET /api/v1/users/me/posts", authMiddleware(false)(requirePermissions(false, permissions.PostReadOwn)(http.HandlerFunc(postHandler.GetOwnPosts))))
	mux.Handle("PATCH /api/v1/posts/{id}/verify", authMiddleware(false)(requirePermissions(false, permissions.PostVerify)(http.HandlerFunc(postHandler.Verify))))
	mux.Handle("PATCH /api/v1/posts/{id}/return", authMiddleware(false)(requirePermissions(false, permissions.PostMarkReturnedAny, permissions.PostMarkReturnedOwn)(http.HandlerFunc(postHandler.ReturnToOwner))))
	// Rooms
	mux.Handle("POST /api/v1/rooms", authMiddleware(false)(requirePermissions(false, permissions.RoomCreate)(http.HandlerFunc(roomHandler.Create))))
	mux.Handle("DELETE /api/v1/rooms/{id}", authMiddleware(false)(requirePermissions(false, permissions.RoomDelete)(http.HandlerFunc(roomHandler.Delete))))
	mux.Handle("PATCH /api/v1/rooms/{id}", authMiddleware(false)(requirePermissions(false, permissions.RoomUpdate)(http.HandlerFunc(roomHandler.Update))))
	mux.HandleFunc("GET /api/v1/rooms", roomHandler.GetRooms)
	// Subjects
	mux.Handle("POST /api/v1/subjects", authMiddleware(false)(requirePermissions(false, permissions.SubjectCreate)(http.HandlerFunc(subjectHandler.Create))))
	mux.Handle("DELETE /api/v1/subjects/{id}", authMiddleware(false)(requirePermissions(false, permissions.SubjectDelete)(http.HandlerFunc(subjectHandler.Delete))))
	mux.Handle("PATCH /api/v1/subjects/{id}", authMiddleware(false)(requirePermissions(false, permissions.SubjectUpdate)(http.HandlerFunc(subjectHandler.Update))))
	mux.HandleFunc("GET /api/v1/subjects", subjectHandler.GetSubjects)
	// Students
	mux.Handle("GET /api/v1/students/{id}", authMiddleware(false)(requirePermissions(false, permissions.StudentReadOther)(http.HandlerFunc(studentHandler.GetStudentByID))))
	mux.Handle("GET /api/v1/students/me", authMiddleware(false)(requirePermissions(false, permissions.StudentReadOwn)(http.HandlerFunc(studentHandler.GetOwn))))
	mux.Handle("GET /api/v1/students/{id}/classroom", authMiddleware(false)(requirePermissions(false, permissions.StudentClassroomReadAny)(http.HandlerFunc(studentHandler.GetClassroom))))
	mux.Handle("GET /api/v1/students/me/classroom", authMiddleware(false)(requirePermissions(false, permissions.StudentClassroomReadOwn)(http.HandlerFunc(studentHandler.GetClassroomOwn))))
	mux.Handle("GET /api/v1/students/{id}/advisor", authMiddleware(false)(requirePermissions(false, permissions.StudentAdvisorReadAny)(http.HandlerFunc(studentHandler.GetAdvisor))))
	mux.Handle("GET /api/v1/students/me/advisor", authMiddleware(false)(requirePermissions(false, permissions.StudentAdvisorReadOwn)(http.HandlerFunc(studentHandler.GetAdvisorOwn))))
	mux.Handle("GET /api/v1/students/{id}/parents", authMiddleware(false)(requirePermissions(false, permissions.StudentParentReadAny)(http.HandlerFunc(studentHandler.GetParents))))
	mux.Handle("GET /api/v1/students/me/parents", authMiddleware(false)(requirePermissions(false, permissions.StudentParentReadOwn)(http.HandlerFunc(studentHandler.GetParentsOwn))))
	mux.Handle("GET /api/v1/students/me/student_group", authMiddleware(false)(requirePermissions(false, permissions.StudentStudentGroupReadOwn)(http.HandlerFunc(studentHandler.GetStudentGroupOwn))))
	// Teacher
	mux.Handle("GET /api/v1/teachers/{id}", authMiddleware(false)(requirePermissions(false, permissions.TeacherReadOther)(http.HandlerFunc(teacherHandler.GetTeacherByID))))
	mux.Handle("GET /api/v1/teachers/me", authMiddleware(false)(requirePermissions(false, permissions.TeacherReadOwn)(http.HandlerFunc(teacherHandler.GetOwn))))
	mux.Handle("GET /api/v1/teachers/{id}/classroom", authMiddleware(false)(requirePermissions(false, permissions.TeacherClassroomReadAny)(http.HandlerFunc(teacherHandler.GetClassroom))))
	mux.Handle("GET /api/v1/teachers/me/classroom", authMiddleware(false)(requirePermissions(false, permissions.TeacherClassroomReadOwn)(http.HandlerFunc(teacherHandler.GetClassroomOwn))))
	mux.Handle("GET /api/v1/teachers/{id}/subjects", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectReadAny)(http.HandlerFunc(teacherHandler.GetSubjects))))
	mux.Handle("GET /api/v1/teachers/me/subjects", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectReadOwn)(http.HandlerFunc(teacherHandler.GetSubjectsOwn))))
	mux.Handle("PUT /api/v1/teachers/{id}/classroom", authMiddleware(false)(requirePermissions(false, permissions.TeacherClassroomAssignAny)(http.HandlerFunc(teacherHandler.AssignClassroom))))
	mux.Handle("PUT /api/v1/teachers/me/classroom", authMiddleware(false)(requirePermissions(false, permissions.TeacherClassroomAssignOwn)(http.HandlerFunc(teacherHandler.AssignClassroomOwn))))
	mux.Handle("DELETE /api/v1/teachers/{id}/classroom", authMiddleware(false)(requirePermissions(false, permissions.TeacherClassroomUnassignAny)(http.HandlerFunc(teacherHandler.UnassignClassroom))))
	mux.Handle("DELETE /api/v1/teachers/me/classroom", authMiddleware(false)(requirePermissions(false, permissions.TeacherClassroomUnassignOwn)(http.HandlerFunc(teacherHandler.UnassignClassroomOwn))))
	mux.Handle("POST /api/v1/teachers/{id}/subjects", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectAddAny)(http.HandlerFunc(teacherHandler.AddSubjects))))
	mux.Handle("POST /api/v1/teachers/me/subjects", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectAddOwn)(http.HandlerFunc(teacherHandler.AddSubjectsOwn))))
	mux.Handle("PUT /api/v1/teachers/{id}/subjects", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectAssignAny)(http.HandlerFunc(teacherHandler.AssignSubjects))))
	mux.Handle("PUT /api/v1/teachers/me/subjects", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectAssignOwn)(http.HandlerFunc(teacherHandler.AssignSubjectsOwn))))
	mux.Handle("DELETE /api/v1/teachers/{userId}/subjects/{subjectId}", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectUnassignAny)(http.HandlerFunc(teacherHandler.UnassignSubject))))
	mux.Handle("DELETE /api/v1/teachers/me/subjects/{id}", authMiddleware(false)(requirePermissions(false, permissions.TeacherSubjectUnassignOwn)(http.HandlerFunc(teacherHandler.UnassignSubjectOwn))))
	mux.Handle("GET /api/v1/teachers/me/student_groups", authMiddleware(false)(requirePermissions(false, permissions.TeacherStudentGroupReadOwn)(http.HandlerFunc(teacherHandler.GetStudentGroupsOwn))))
	// Parents
	mux.Handle("GET /api/v1/parents/{id}", authMiddleware(false)(requirePermissions(false, permissions.ParentReadOther)(http.HandlerFunc(parentHandler.GetParentByID))))
	mux.Handle("GET /api/v1/parents/me", authMiddleware(false)(requirePermissions(false, permissions.ParentReadOwn)(http.HandlerFunc(parentHandler.GetOwn))))
	mux.Handle("GET /api/v1/parents/{id}/students", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentReadAny)(http.HandlerFunc(parentHandler.GetStudents))))
	mux.Handle("GET /api/v1/parents/me/students", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentReadOwn)(http.HandlerFunc(parentHandler.GetStudentsOwn))))
	mux.Handle("GET /api/v1/parents/me/student_groups", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentGroupReadOwn)(http.HandlerFunc(parentHandler.GetStudentGroupsOwn))))
	mux.Handle("POST /api/v1/parents/{id}/students", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentAddAny)(http.HandlerFunc(parentHandler.AddStudents))))
	mux.Handle("POST /api/v1/parents/me/students", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentAddOwn)(http.HandlerFunc(parentHandler.AddStudentsOwn))))
	mux.Handle("DELETE /api/v1/parents/{parentId}/students/{studentId}", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentUnassignAny)(http.HandlerFunc(parentHandler.UnassignStudent))))
	mux.Handle("DELETE /api/v1/parents/me/students/{id}", authMiddleware(false)(requirePermissions(false, permissions.ParentStudentUnassignOwn)(http.HandlerFunc(parentHandler.UnassignStudentOwn))))
	// Staff
	mux.Handle("GET /api/v1/staff/{id}", authMiddleware(false)(requirePermissions(false, permissions.StaffReadOther)(http.HandlerFunc(staffHandler.GetStaffByID))))
	mux.Handle("GET /api/v1/staff/me", authMiddleware(false)(requirePermissions(false, permissions.StaffReadOwn)(http.HandlerFunc(staffHandler.GetOwn))))
	mux.Handle("PUT /api/v1/staff/{id}/position", authMiddleware(false)(requirePermissions(false, permissions.StaffPositionAssign)(http.HandlerFunc(staffHandler.AssignPosition))))
	mux.Handle("GET /api/v1/staff/{id}/position", authMiddleware(false)(requirePermissions(false, permissions.StaffPositionRead)(http.HandlerFunc(staffHandler.GetPosition))))
	// Institution administrator
	mux.Handle("GET /api/v1/institution_administrators/{id}", authMiddleware(false)(requirePermissions(false, permissions.InstitutionAdministratorReadOther)(http.HandlerFunc(institutionAdministratorHandler.GetInstitutionAdministratorByID))))
	mux.Handle("GET /api/v1/institution_administrators/me", authMiddleware(false)(requirePermissions(false, permissions.InstitutionAdministratorReadOwn)(http.HandlerFunc(institutionAdministratorHandler.GetOwn))))
	mux.Handle("PUT /api/v1/institution_administrators/{id}/position", authMiddleware(false)(requirePermissions(false, permissions.InstitutionAdministratorPositionAssign)(http.HandlerFunc(institutionAdministratorHandler.AssignPosition))))
	mux.Handle("GET /api/v1/institution_administrators/{id}/position", authMiddleware(false)(requirePermissions(false, permissions.InstitutionAdministratorPositionRead)(http.HandlerFunc(institutionAdministratorHandler.GetPosition))))
	// Invite tokens
	mux.Handle("POST /api/v1/tokens/invite", authMiddleware(false)(requirePermissions(false, permissions.TokenInviteAdminCreate, permissions.TokenInviteUserCreate)(http.HandlerFunc(inviteHandler.Create))))
	mux.Handle("DELETE /api/v1/tokens/invite/{token}", authMiddleware(false)(requirePermissions(false, permissions.TokenInviteAdminDelete, permissions.TokenInviteUserDelete)(http.HandlerFunc(inviteHandler.Revoke))))
	// Institution administrator positions
	mux.Handle("POST /api/v1/institution_administrators/positions", authMiddleware(false)(requirePermissions(false, permissions.PositionInstitutionAdministratorCreate)(http.HandlerFunc(institutionAdministratorPositionHandler.Create))))
	mux.HandleFunc("GET /api/v1/institution_administrators/positions", institutionAdministratorPositionHandler.GetAll)
	mux.Handle("PATCH /api/v1/institution_administrators/positions/{id}", authMiddleware(false)(requirePermissions(false, permissions.PositionInstitutionAdministratorUpdate)(http.HandlerFunc(institutionAdministratorPositionHandler.Update))))
	mux.Handle("DELETE /api/v1/institution_administrators/positions/{id}", authMiddleware(false)(requirePermissions(false, permissions.PositionInstitutionAdministratorDelete)(http.HandlerFunc(institutionAdministratorPositionHandler.Delete))))
	// Staff positions
	mux.Handle("POST /api/v1/staff/positions", authMiddleware(false)(requirePermissions(false, permissions.PositionStaffCreate)(http.HandlerFunc(staffPositionHandler.Create))))
	mux.HandleFunc("GET /api/v1/staff/positions", staffPositionHandler.GetAll)
	mux.Handle("PATCH /api/v1/staff/positions/{id}", authMiddleware(false)(requirePermissions(false, permissions.PositionStaffUpdate)(http.HandlerFunc(staffPositionHandler.Update))))
	mux.Handle("DELETE /api/v1/staff/positions/{id}", authMiddleware(false)(requirePermissions(false, permissions.PositionStaffDelete)(http.HandlerFunc(staffPositionHandler.Delete))))
	// Roles
	// mux.Handle("GET /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.GetPermissions)))
	// mux.Handle("PUT /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.AssignPermissions)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}
