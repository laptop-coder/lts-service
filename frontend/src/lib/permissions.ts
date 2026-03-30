import { useAuth } from "./auth";

export function usePermissions() {
  const auth = useAuth();

  const hasPermission = (permission: string): boolean => {
    const user = auth.user();
    if (!user) return false;
    return user.roles.some((role) =>
      role.permissions.some((p) => p.name === permission),
    );
  };

  const hasAnyPermission = (...perms: string[]): boolean => {
    return perms.some((p) => hasPermission(p));
  };

  const hasAllPermissions = (...perms: string[]): boolean => {
    return perms.every((p) => hasPermission(p));
  };

  const hasRole = (roleName: string): boolean => {
    const user = auth.user();
    if (!user) return false;
    return user.roles.some((r) => r.name === roleName);
  };

  return { hasPermission, hasAnyPermission, hasAllPermissions, hasRole };
}

export function getPermissions() {
  // Post permissions
  const POST_CREATE = "post.create";
  const POST_READ_ANY = "post.read.any";
  const POST_READ_OWN = "post.read.own";
  const POST_UPDATE_ANY = "post.update.any";
  const POST_UPDATE_OWN = "post.update.own";
  const POST_DELETE_ANY = "post.delete.any";
  const POST_DELETE_OWN = "post.delete.own";
  const POST_PHOTO_DELETE_ANY = "post.photo.delete.any";
  const POST_PHOTO_DELETE_OWN = "post.photo.delete.own";
  const POST_VERIFY = "post.verify";
  const POST_MARK_RETURNED_ANY = "post.mark.returned.any";
  const POST_MARK_RETURNED_OWN = "post.mark.returned.own";

  // User permissions
  const USER_READ_OWN = "user.read.own";
  const USER_READ_OTHER = "user.read.other";
  const USER_READ_ALL = "user.read.all";
  const USER_UPDATE_OWN = "user.update.own";
  const USER_DELETE_ANY = "user.delete.any";
  const USER_DELETE_OWN = "user.delete.own";

  // Room permissions
  const ROOM_CREATE = "room.create";
  const ROOM_READ = "room.read";
  const ROOM_UPDATE = "room.update";
  const ROOM_DELETE = "room.delete";

  // Subject permissions
  const SUBJECT_CREATE = "subject.create";
  const SUBJECT_READ = "subject.read";
  const SUBJECT_UPDATE = "subject.update";
  const SUBJECT_DELETE = "subject.delete";

  // Student group permissions
  const STUDENT_GROUP_CREATE = "student_group.create";
  const STUDENT_GROUP_READ_ANY = "student_group.read.any";
  const STUDENT_GROUP_UPDATE = "student_group.update";
  const STUDENT_GROUP_DELETE = "student_group.delete";
  const STUDENT_GROUP_ADVISOR_ASSIGN = "student_group.advisor.assign";
  const STUDENT_GROUP_ADVISOR_UNASSIGN_ANY =
    "student_group.advisor.unassign.any";
  const STUDENT_GROUP_ADVISOR_UNASSIGN_OWN =
    "student_group.advisor.unassign.own";
  const STUDENT_GROUP_ADVISOR_READ = "student_group.advisor.read";

  // Teacher permissions
  const TEACHER_SUBJECT_READ_ANY = "teacher.subject.read.any";
  const TEACHER_SUBJECT_READ_OWN = "teacher.subject.read.own";
  const TEACHER_SUBJECT_ADD_ANY = "teacher.subject.add.any";
  const TEACHER_SUBJECT_ADD_OWN = "teacher.subject.add.own";
  const TEACHER_SUBJECT_ASSIGN_ANY = "teacher.subject.assign.any";
  const TEACHER_SUBJECT_ASSIGN_OWN = "teacher.subject.assign.own";
  const TEACHER_SUBJECT_UNASSIGN_ANY = "teacher.subject.unassign.any";
  const TEACHER_SUBJECT_UNASSIGN_OWN = "teacher.subject.unassign.own";
  const TEACHER_CLASSROOM_READ_ANY = "teacher.classroom.read.any";
  const TEACHER_CLASSROOM_READ_OWN = "teacher.classroom.read.own";
  const TEACHER_CLASSROOM_ASSIGN_ANY = "teacher.classroom.assign.any";
  const TEACHER_CLASSROOM_ASSIGN_OWN = "teacher.classroom.assign.own";
  const TEACHER_CLASSROOM_UNASSIGN_ANY = "teacher.classroom.unassign.any";
  const TEACHER_CLASSROOM_UNASSIGN_OWN = "teacher.classroom.unassign.own";
  const TEACHER_READ_OTHER = "teacher.read.other";
  const TEACHER_READ_OWN = "teacher.read.own";
  const TEACHER_STUDENT_GROUP_READ_OWN = "teacher.student_group.read.own";

  // Parent permissions
  const PARENT_STUDENT_READ_ANY = "parent.student.read.any";
  const PARENT_STUDENT_READ_OWN = "parent.student.read.own";
  const PARENT_STUDENT_ADD_ANY = "parent.student.add.any";
  const PARENT_STUDENT_ADD_OWN = "parent.student.add.own";
  const PARENT_STUDENT_UNASSIGN_ANY = "parent.student.unassign.any";
  const PARENT_STUDENT_UNASSIGN_OWN = "parent.student.unassign.own";
  const PARENT_READ_OTHER = "parent.read.other";
  const PARENT_READ_OWN = "parent.read.own";
  const PARENT_STUDENT_GROUP_READ_OWN = "parent.student_group.read.own";

  // User roles permissions
  const ROLE_ASSIGN = "role.assign";
  const ROLE_ADD = "role.add";
  const ROLE_DELETE = "role.delete";
  const ROLE_READ_ANY = "role.read.any";
  const ROLE_READ_OWN = "role.read.own";

  // Permissions to work with tokens
  const TOKEN_INVITE_ADMIN_CREATE = "token.invite.admin.create";
  const TOKEN_INVITE_USER_CREATE = "token.invite.user.create";
  const TOKEN_INVITE_ADMIN_DELETE = "token.invite.admin.delete";
  const TOKEN_INVITE_USER_DELETE = "token.invite.user.delete";

  // Student permissions
  const STUDENT_READ_OTHER = "student.read.other";
  const STUDENT_READ_OWN = "student.read.own";
  const STUDENT_CLASSROOM_READ_ANY = "student.classroom.read.any";
  const STUDENT_CLASSROOM_READ_OWN = "student.classroom.read.own";
  const STUDENT_ADVISOR_READ_ANY = "student.advisor.read.any";
  const STUDENT_ADVISOR_READ_OWN = "student.advisor.read.own";
  const STUDENT_PARENT_READ_ANY = "student.parent.read.any";
  const STUDENT_PARENT_READ_OWN = "student.parent.read.own";
  const STUDENT_STUDENT_GROUP_READ_OWN = "student.student_group.read.own";

  // Institution administrator
  const INSTITUTION_ADMINISTRATOR_READ_OTHER =
    "institution_administrator.read.other";
  const INSTITUTION_ADMINISTRATOR_READ_OWN =
    "institution_administrator.read.own";
  const INSTITUTION_ADMINISTRATOR_POSITION_ASSIGN =
    "institution_administrator.position.assign";
  const INSTITUTION_ADMINISTRATOR_POSITION_READ =
    "institution_administrator.position.read";

  // Staff
  const STAFF_READ_OTHER = "staff.read.other";
  const STAFF_READ_OWN = "staff.read.own";
  const STAFF_POSITION_ASSIGN = "staff.position.assign";
  const STAFF_POSITION_READ = "staff.position.read";

  // Position institution administrator
  const POSITION_INSTITUTION_ADMINISTRATOR_CREATE =
    "position.institution_administrator.create";
  const POSITION_INSTITUTION_ADMINISTRATOR_READ =
    "position.institution_administrator.read";
  const POSITION_INSTITUTION_ADMINISTRATOR_UPDATE =
    "position.institution_administrator.update";
  const POSITION_INSTITUTION_ADMINISTRATOR_DELETE =
    "position.institution_administrator.delete";

  // Position staff
  const POSITION_STAFF_CREATE = "position.staff.create";
  const POSITION_STAFF_READ = "position.staff.read";
  const POSITION_STAFF_UPDATE = "position.staff.update";
  const POSITION_STAFF_DELETE = "position.staff.delete";

  return {
    POST_CREATE,
    POST_READ_ANY,
    POST_READ_OWN,
    POST_UPDATE_ANY,
    POST_UPDATE_OWN,
    POST_DELETE_ANY,
    POST_DELETE_OWN,
    POST_PHOTO_DELETE_ANY,
    POST_PHOTO_DELETE_OWN,
    POST_VERIFY,
    POST_MARK_RETURNED_ANY,
    POST_MARK_RETURNED_OWN,
    USER_READ_OWN,
    USER_READ_OTHER,
    USER_READ_ALL,
    USER_UPDATE_OWN,
    USER_DELETE_ANY,
    USER_DELETE_OWN,
    ROOM_CREATE,
    ROOM_READ,
    ROOM_UPDATE,
    ROOM_DELETE,
    SUBJECT_CREATE,
    SUBJECT_READ,
    SUBJECT_UPDATE,
    SUBJECT_DELETE,
    STUDENT_GROUP_CREATE,
    STUDENT_GROUP_READ_ANY,
    STUDENT_GROUP_UPDATE,
    STUDENT_GROUP_DELETE,
    STUDENT_GROUP_ADVISOR_ASSIGN,
    STUDENT_GROUP_ADVISOR_UNASSIGN_ANY,
    STUDENT_GROUP_ADVISOR_UNASSIGN_OWN,
    STUDENT_GROUP_ADVISOR_READ,
    TEACHER_SUBJECT_READ_ANY,
    TEACHER_SUBJECT_READ_OWN,
    TEACHER_SUBJECT_ADD_ANY,
    TEACHER_SUBJECT_ADD_OWN,
    TEACHER_SUBJECT_ASSIGN_ANY,
    TEACHER_SUBJECT_ASSIGN_OWN,
    TEACHER_SUBJECT_UNASSIGN_ANY,
    TEACHER_SUBJECT_UNASSIGN_OWN,
    TEACHER_CLASSROOM_READ_ANY,
    TEACHER_CLASSROOM_READ_OWN,
    TEACHER_CLASSROOM_ASSIGN_ANY,
    TEACHER_CLASSROOM_ASSIGN_OWN,
    TEACHER_CLASSROOM_UNASSIGN_ANY,
    TEACHER_CLASSROOM_UNASSIGN_OWN,
    TEACHER_READ_OTHER,
    TEACHER_READ_OWN,
    TEACHER_STUDENT_GROUP_READ_OWN,
    PARENT_STUDENT_READ_ANY,
    PARENT_STUDENT_READ_OWN,
    PARENT_STUDENT_ADD_ANY,
    PARENT_STUDENT_ADD_OWN,
    PARENT_STUDENT_UNASSIGN_ANY,
    PARENT_STUDENT_UNASSIGN_OWN,
    PARENT_READ_OTHER,
    PARENT_READ_OWN,
    PARENT_STUDENT_GROUP_READ_OWN,
    ROLE_ASSIGN,
    ROLE_ADD,
    ROLE_DELETE,
    ROLE_READ_ANY,
    ROLE_READ_OWN,
    TOKEN_INVITE_ADMIN_CREATE,
    TOKEN_INVITE_USER_CREATE,
    TOKEN_INVITE_ADMIN_DELETE,
    TOKEN_INVITE_USER_DELETE,
    STUDENT_READ_OTHER,
    STUDENT_READ_OWN,
    STUDENT_CLASSROOM_READ_ANY,
    STUDENT_CLASSROOM_READ_OWN,
    STUDENT_ADVISOR_READ_ANY,
    STUDENT_ADVISOR_READ_OWN,
    STUDENT_PARENT_READ_ANY,
    STUDENT_PARENT_READ_OWN,
    STUDENT_STUDENT_GROUP_READ_OWN,
    INSTITUTION_ADMINISTRATOR_READ_OTHER,
    INSTITUTION_ADMINISTRATOR_READ_OWN,
    INSTITUTION_ADMINISTRATOR_POSITION_ASSIGN,
    INSTITUTION_ADMINISTRATOR_POSITION_READ,
    STAFF_READ_OTHER,
    STAFF_READ_OWN,
    STAFF_POSITION_ASSIGN,
    STAFF_POSITION_READ,
    POSITION_INSTITUTION_ADMINISTRATOR_CREATE,
    POSITION_INSTITUTION_ADMINISTRATOR_READ,
    POSITION_INSTITUTION_ADMINISTRATOR_UPDATE,
    POSITION_INSTITUTION_ADMINISTRATOR_DELETE,
    POSITION_STAFF_CREATE,
    POSITION_STAFF_READ,
    POSITION_STAFF_UPDATE,
    POSITION_STAFF_DELETE,
  };
}

export function getRoles() {
  const SUPERADMIN = "superadmin";
  const ADMIN = "admin";
  const INSTITUTION_ADMINISTRATOR = "institution_administrator";
  const STAFF = "staff";
  const TEACHER = "teacher";
  const PARENT = "parent";
  const STUDENT = "student";

  return {
    SUPERADMIN,
    ADMIN,
    INSTITUTION_ADMINISTRATOR,
    STAFF,
    TEACHER,
    PARENT,
    STUDENT,
  };
}
