// Package permissions contains all permission constants used in the application
package permissions

// Post permissions
const (
	// ability to create new posts
	PostCreate = "posts.create"
	// ability to view posts
	PostRead = "posts.read"
	// ability to update any post
	PostUpdate = "posts.update"
	// ability to update own posts
	PostUpdateOwn = "posts.update.own"
	// ability to delete any post
	PostDelete = "posts.delete"
	// ability to delete own posts
	PostDeleteOwn = "posts.delete.own"
	// ability to verify posts (moderator action)
	PostVerify = "posts.verify"
)

// User permissions
const (
	// ability to create new users
	UserCreate = "users.create"
	// ability to view users
	UserRead = "users.read"
	// ability to update any user
	UserUpdate = "users.update"
	// ability to update own profile
	UserUpdateOwn = "users.update.own"
	// ability to delete any user
	UserDelete = "users.delete"
	// ability to delete own account
	UserDeleteOwn = "users.delete.own"
)

// Room permissions
const (
	// ability to create rooms
	RoomCreate = "rooms.create"
	// ability to view rooms
	RoomRead = "rooms.read"
	// ability to update rooms
	RoomUpdate = "rooms.update"
	// ability to delete rooms
	RoomDelete = "rooms.delete"
)

// Subject permissions
const (
	// ability to create subjects
	SubjectCreate = "subjects.create"
	// ability to view subjects
	SubjectRead = "subjects.read"
	// ability to update subjects
	SubjectUpdate = "subjects.update"
	// ability to delete subjects
	SubjectDelete = "subjects.delete"
)

// StudentGroup permissions
const (
	// ability to create student groups
	GroupCreate = "groups.create"
	// ability to view student groups
	GroupRead = "groups.read"
	// ability to update student groups
	GroupUpdate = "groups.update"
	// ability to delete student groups
	GroupDelete = "groups.delete"
	// ability to assign advisor to student group
	GroupAssignAdvisor = "groups.assign.advisor"
)

// Teacher permissions
const (
	// ability to assign subjects to teacher
	TeacherAssignSubject = "teacher.assign.subject"
	// ability to assign classroom to teacher
	TeacherAssignClassroom = "teacher.assign.classroom"
)

// Student permissions
const (
	// ability to assign parent to student
	StudentAssignParent = "student.assign.parent"
)

// Parent permissions
const (
	// ability to view parent's students
	ParentViewStudents = "parent.view.students"
)

// Role permissions
const (
	// ability to assign roles to users
	RoleAssign = "roles.assign"
	// ability to create new roles
	RoleCreate = "roles.create"
	// ability to update roles
	RoleUpdate = "roles.update"
	// ability to delete roles
	RoleDelete = "roles.delete"
)
