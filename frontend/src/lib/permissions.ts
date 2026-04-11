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

export const PERMISSIONS = {
  // Post permissions
  POST_CREATE: "post.create",
  POST_READ_ANY: "post.read.any",
  POST_READ_OWN: "post.read.own",
  POST_UPDATE_ANY: "post.update.any",
  POST_UPDATE_OWN: "post.update.own",
  POST_DELETE_ANY: "post.delete.any",
  POST_DELETE_OWN: "post.delete.own",
  POST_PHOTO_DELETE_ANY: "post.photo.delete.any",
  POST_PHOTO_DELETE_OWN: "post.photo.delete.own",
  POST_PHOTO_UPDATE_ANY: "post.photo.update.any",
  POST_PHOTO_UPDATE_OWN: "post.photo.update.own",
  POST_VERIFY: "post.verify",
  POST_MARK_RETURNED_ANY: "post.mark.returned.any",
  POST_MARK_RETURNED_OWN: "post.mark.returned.own",

  // User permissions
  USER_READ_OWN: "user.read.own",
  USER_READ_OTHER: "user.read.other",
  USER_READ_ALL: "user.read.all",
  USER_UPDATE_OWN: "user.update.own",
  USER_DELETE_ANY: "user.delete.any",
  USER_DELETE_OWN: "user.delete.own",

  // Room permissions
  ROOM_CREATE: "room.create",
  ROOM_UPDATE: "room.update",
  ROOM_DELETE: "room.delete",

  // Subject permissions
  SUBJECT_CREATE: "subject.create",
  SUBJECT_UPDATE: "subject.update",
  SUBJECT_DELETE: "subject.delete",

  // Student group permissions
  STUDENT_GROUP_CREATE: "student_group.create",
  STUDENT_GROUP_UPDATE: "student_group.update",
  STUDENT_GROUP_DELETE: "student_group.delete",
  STUDENT_GROUP_ADVISOR_ASSIGN: "student_group.advisor.assign",
  STUDENT_GROUP_ADVISOR_UNASSIGN_ANY: "student_group.advisor.unassign.any",
  STUDENT_GROUP_ADVISOR_UNASSIGN_OWN: "student_group.advisor.unassign.own",
  STUDENT_GROUP_ADVISOR_READ: "student_group.advisor.read",

  // Teacher permissions
  TEACHER_SUBJECT_READ_ANY: "teacher.subject.read.any",
  TEACHER_SUBJECT_READ_OWN: "teacher.subject.read.own",
  TEACHER_SUBJECT_ADD_ANY: "teacher.subject.add.any",
  TEACHER_SUBJECT_ADD_OWN: "teacher.subject.add.own",
  TEACHER_SUBJECT_ASSIGN_ANY: "teacher.subject.assign.any",
  TEACHER_SUBJECT_ASSIGN_OWN: "teacher.subject.assign.own",
  TEACHER_SUBJECT_UNASSIGN_ANY: "teacher.subject.unassign.any",
  TEACHER_SUBJECT_UNASSIGN_OWN: "teacher.subject.unassign.own",
  TEACHER_CLASSROOM_READ_ANY: "teacher.classroom.read.any",
  TEACHER_CLASSROOM_READ_OWN: "teacher.classroom.read.own",
  TEACHER_CLASSROOM_ASSIGN_ANY: "teacher.classroom.assign.any",
  TEACHER_CLASSROOM_ASSIGN_OWN: "teacher.classroom.assign.own",
  TEACHER_CLASSROOM_UNASSIGN_ANY: "teacher.classroom.unassign.any",
  TEACHER_CLASSROOM_UNASSIGN_OWN: "teacher.classroom.unassign.own",
  TEACHER_READ_OTHER: "teacher.read.other",
  TEACHER_READ_OWN: "teacher.read.own",
  TEACHER_STUDENT_GROUP_READ_OWN: "teacher.student_group.read.own",

  // Parent permissions
  PARENT_STUDENT_READ_ANY: "parent.student.read.any",
  PARENT_STUDENT_READ_OWN: "parent.student.read.own",
  PARENT_STUDENT_ADD_ANY: "parent.student.add.any",
  PARENT_STUDENT_ADD_OWN: "parent.student.add.own",
  PARENT_STUDENT_UNASSIGN_ANY: "parent.student.unassign.any",
  PARENT_STUDENT_UNASSIGN_OWN: "parent.student.unassign.own",
  PARENT_READ_OTHER: "parent.read.other",
  PARENT_READ_OWN: "parent.read.own",
  PARENT_STUDENT_GROUP_READ_OWN: "parent.student_group.read.own",

  // User roles permissions
  ROLE_ADMIN_ASSIGN: "role.admin.assign",
  ROLE_USER_ASSIGN: "role.user.assign",
  ROLE_ADMIN_ADD: "role.admin.add",
  ROLE_USER_ADD: "role.user.add",
  ROLE_ADMIN_UNASSIGN: "role.admin.unassign",
  ROLE_USER_UNASSIGN: "role.user.unassign",
  ROLE_READ_ANY: "role.read.any",
  ROLE_READ_OWN: "role.read.own",

  // Permissions to work with tokens
  TOKEN_INVITE_ADMIN_CREATE: "token.invite.admin.create",
  TOKEN_INVITE_USER_CREATE: "token.invite.user.create",
  TOKEN_INVITE_ADMIN_DELETE: "token.invite.admin.delete",
  TOKEN_INVITE_USER_DELETE: "token.invite.user.delete",

  // Student permissions
  STUDENT_READ_OTHER: "student.read.other",
  STUDENT_READ_OWN: "student.read.own",
  STUDENT_CLASSROOM_READ_ANY: "student.classroom.read.any",
  STUDENT_CLASSROOM_READ_OWN: "student.classroom.read.own",
  STUDENT_ADVISOR_READ_ANY: "student.advisor.read.any",
  STUDENT_ADVISOR_READ_OWN: "student.advisor.read.own",
  STUDENT_PARENT_READ_ANY: "student.parent.read.any",
  STUDENT_PARENT_READ_OWN: "student.parent.read.own",
  STUDENT_STUDENT_GROUP_READ_OWN: "student.student_group.read.own",

  // Institution administrator
  INSTITUTION_ADMINISTRATOR_READ_OTHER: "institution_administrator.read.other",
  INSTITUTION_ADMINISTRATOR_READ_OWN: "institution_administrator.read.own",
  INSTITUTION_ADMINISTRATOR_POSITION_ASSIGN:
    "institution_administrator.position.assign",
  INSTITUTION_ADMINISTRATOR_POSITION_READ:
    "institution_administrator.position.read",

  // Staff
  STAFF_READ_OTHER: "staff.read.other",
  STAFF_READ_OWN: "staff.read.own",
  STAFF_POSITION_ASSIGN: "staff.position.assign",
  STAFF_POSITION_READ: "staff.position.read",

  // Position institution administrator
  POSITION_INSTITUTION_ADMINISTRATOR_CREATE:
    "position.institution_administrator.create",
  POSITION_INSTITUTION_ADMINISTRATOR_UPDATE:
    "position.institution_administrator.update",
  POSITION_INSTITUTION_ADMINISTRATOR_DELETE:
    "position.institution_administrator.delete",

  // Position staff
  POSITION_STAFF_CREATE: "position.staff.create",
  POSITION_STAFF_UPDATE: "position.staff.update",
  POSITION_STAFF_DELETE: "position.staff.delete",
};

export const ROLES = {
  SUPERADMIN: "superadmin",
  ADMIN: "admin",
  INSTITUTION_ADMINISTRATOR: "institution_administrator",
  STAFF: "staff",
  TEACHER: "teacher",
  PARENT: "parent",
  STUDENT: "student",
};

export const ROLES_TO_DISPLAY = [
  { id: 1, name: ROLES.SUPERADMIN, displayName: "Суперадминистратор" },
  { id: 2, name: ROLES.ADMIN, displayName: "Админ" },
  {
    id: 3,
    name: ROLES.INSTITUTION_ADMINISTRATOR,
    displayName: "Администрация",
  },
  { id: 4, name: ROLES.STAFF, displayName: "Сотрудник" },
  { id: 5, name: ROLES.TEACHER, displayName: "Преподаватель" },
  { id: 6, name: ROLES.PARENT, displayName: "Родитель" },
  { id: 7, name: ROLES.STUDENT, displayName: "Обучающийся" },
];
