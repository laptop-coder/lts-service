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

export function permissions() {
  // Post permissions
  const PostCreate = "post.create";
  const PostReadAny = "post.read.any";
  const PostReadOwn = "post.read.own";
  const PostUpdateAny = "post.update.any";
  const PostUpdateOwn = "post.update.own";
  const PostDeleteAny = "post.delete.any";
  const PostDeleteOwn = "post.delete.own";
  const PostPhotoDeleteAny = "post.photo.delete.any";
  const PostPhotoDeleteOwn = "post.photo.delete.own";
  const PostVerify = "post.verify";
  const PostMarkReturnedAny = "post.mark.returned.any";
  const PostMarkReturnedOwn = "post.mark.returned.own";

  // User permissions
  const UserReadOwn = "user.read.own";
  const UserReadOther = "user.read.other";
  const UserReadAll = "user.read.all";
  const UserUpdateOwn = "user.update.own";
  const UserDeleteAny = "user.delete.any";
  const UserDeleteOwn = "user.delete.own";

  // Room permissions
  const RoomCreate = "room.create";
  const RoomRead = "room.read";
  const RoomUpdate = "room.update";
  const RoomDelete = "room.delete";

  // Subject permissions
  const SubjectCreate = "subject.create";
  const SubjectRead = "subject.read";
  const SubjectUpdate = "subject.update";
  const SubjectDelete = "subject.delete";

  // Student group permissions
  const StudentGroupCreate = "student_group.create";
  const StudentGroupReadAny = "student_group.read.any";
  const StudentGroupUpdate = "student_group.update";
  const StudentGroupDelete = "student_group.delete";
  const StudentGroupAdvisorAssign = "student_group.advisor.assign";
  const StudentGroupAdvisorUnassignAny = "student_group.advisor.unassign.any";
  const StudentGroupAdvisorUnassignOwn = "student_group.advisor.unassign.own";
  const StudentGroupAdvisorRead = "student_group.advisor.read";

  // Teacher permissions
  const TeacherSubjectReadAny = "teacher.subject.read.any";
  const TeacherSubjectReadOwn = "teacher.subject.read.own";
  const TeacherSubjectAddAny = "teacher.subject.add.any";
  const TeacherSubjectAddOwn = "teacher.subject.add.own";
  const TeacherSubjectAssignAny = "teacher.subject.assign.any";
  const TeacherSubjectAssignOwn = "teacher.subject.assign.own";
  const TeacherSubjectUnassignAny = "teacher.subject.unassign.any";
  const TeacherSubjectUnassignOwn = "teacher.subject.unassign.own";
  const TeacherClassroomReadAny = "teacher.classroom.read.any";
  const TeacherClassroomReadOwn = "teacher.classroom.read.own";
  const TeacherClassroomAssignAny = "teacher.classroom.assign.any";
  const TeacherClassroomAssignOwn = "teacher.classroom.assign.own";
  const TeacherClassroomUnassignAny = "teacher.classroom.unassign.any";
  const TeacherClassroomUnassignOwn = "teacher.classroom.unassign.own";
  const TeacherReadOther = "teacher.read.other";
  const TeacherReadOwn = "teacher.read.own";
  const TeacherStudentGroupReadOwn = "teacher.student_group.read.own";

  // Parent permissions
  const ParentStudentReadAny = "parent.student.read.any";
  const ParentStudentReadOwn = "parent.student.read.own";
  const ParentStudentAddAny = "parent.student.add.any";
  const ParentStudentAddOwn = "parent.student.add.own";
  const ParentStudentUnassignAny = "parent.student.unassign.any";
  const ParentStudentUnassignOwn = "parent.student.unassign.own";
  const ParentReadOther = "parent.read.other";
  const ParentReadOwn = "parent.read.own";
  const ParentStudentGroupReadOwn = "parent.student_group.read.own";

  // User roles permissions
  const RoleAssign = "role.assign";
  const RoleAdd = "role.add";
  const RoleDelete = "role.delete";
  const RoleReadAny = "role.read.any";
  const RoleReadOwn = "role.read.own";

  // Permissions to work with tokens
  const TokenInviteAdminCreate = "token.invite.admin.create";
  const TokenInviteUserCreate = "token.invite.user.create";
  const TokenInviteAdminDelete = "token.invite.admin.delete";
  const TokenInviteUserDelete = "token.invite.user.delete";

  // Student permissions
  const StudentReadOther = "student.read.other";
  const StudentReadOwn = "student.read.own";
  const StudentClassroomReadAny = "student.classroom.read.any";
  const StudentClassroomReadOwn = "student.classroom.read.own";
  const StudentAdvisorReadAny = "student.advisor.read.any";
  const StudentAdvisorReadOwn = "student.advisor.read.own";
  const StudentParentReadAny = "student.parent.read.any";
  const StudentParentReadOwn = "student.parent.read.own";
  const StudentStudentGroupReadOwn = "student.student_group.read.own";

  // Institution administrator
  const InstitutionAdministratorReadOther =
    "institution_administrator.read.other";
  const InstitutionAdministratorReadOwn = "institution_administrator.read.own";
  const InstitutionAdministratorPositionAssign =
    "institution_administrator.position.assign";
  const InstitutionAdministratorPositionRead =
    "institution_administrator.position.read";

  // Staff
  const StaffReadOther = "staff.read.other";
  const StaffReadOwn = "staff.read.own";
  const StaffPositionAssign = "staff.position.assign";
  const StaffPositionRead = "staff.position.read";

  // Position institution administrator
  const PositionInstitutionAdministratorCreate =
    "position.institution_administrator.create";
  const PositionInstitutionAdministratorRead =
    "position.institution_administrator.read";
  const PositionInstitutionAdministratorUpdate =
    "position.institution_administrator.update";
  const PositionInstitutionAdministratorDelete =
    "position.institution_administrator.delete";

  // Position staff
  const PositionStaffCreate = "position.staff.create";
  const PositionStaffRead = "position.staff.read";
  const PositionStaffUpdate = "position.staff.update";
  const PositionStaffDelete = "position.staff.delete";

  return {
    PostCreate,
    PostReadAny,
    PostReadOwn,
    PostUpdateAny,
    PostUpdateOwn,
    PostDeleteAny,
    PostDeleteOwn,
    PostPhotoDeleteAny,
    PostPhotoDeleteOwn,
    PostVerify,
    PostMarkReturnedAny,
    PostMarkReturnedOwn,
    UserReadOwn,
    UserReadOther,
    UserReadAll,
    UserUpdateOwn,
    UserDeleteAny,
    UserDeleteOwn,
    RoomCreate,
    RoomRead,
    RoomUpdate,
    RoomDelete,
    SubjectCreate,
    SubjectRead,
    SubjectUpdate,
    SubjectDelete,
    StudentGroupCreate,
    StudentGroupReadAny,
    StudentGroupUpdate,
    StudentGroupDelete,
    StudentGroupAdvisorAssign,
    StudentGroupAdvisorUnassignAny,
    StudentGroupAdvisorUnassignOwn,
    StudentGroupAdvisorRead,
    TeacherSubjectReadAny,
    TeacherSubjectReadOwn,
    TeacherSubjectAddAny,
    TeacherSubjectAddOwn,
    TeacherSubjectAssignAny,
    TeacherSubjectAssignOwn,
    TeacherSubjectUnassignAny,
    TeacherSubjectUnassignOwn,
    TeacherClassroomReadAny,
    TeacherClassroomReadOwn,
    TeacherClassroomAssignAny,
    TeacherClassroomAssignOwn,
    TeacherClassroomUnassignAny,
    TeacherClassroomUnassignOwn,
    TeacherReadOther,
    TeacherReadOwn,
    TeacherStudentGroupReadOwn,
    ParentStudentReadAny,
    ParentStudentReadOwn,
    ParentStudentAddAny,
    ParentStudentAddOwn,
    ParentStudentUnassignAny,
    ParentStudentUnassignOwn,
    ParentReadOther,
    ParentReadOwn,
    ParentStudentGroupReadOwn,
    RoleAssign,
    RoleAdd,
    RoleDelete,
    RoleReadAny,
    RoleReadOwn,
    TokenInviteAdminCreate,
    TokenInviteUserCreate,
    TokenInviteAdminDelete,
    TokenInviteUserDelete,
    StudentReadOther,
    StudentReadOwn,
    StudentClassroomReadAny,
    StudentClassroomReadOwn,
    StudentAdvisorReadAny,
    StudentAdvisorReadOwn,
    StudentParentReadAny,
    StudentParentReadOwn,
    StudentStudentGroupReadOwn,
    InstitutionAdministratorReadOther,
    InstitutionAdministratorReadOwn,
    InstitutionAdministratorPositionAssign,
    InstitutionAdministratorPositionRead,
    StaffReadOther,
    StaffReadOwn,
    StaffPositionAssign,
    StaffPositionRead,
    PositionInstitutionAdministratorCreate,
    PositionInstitutionAdministratorRead,
    PositionInstitutionAdministratorUpdate,
    PositionInstitutionAdministratorDelete,
    PositionStaffCreate,
    PositionStaffRead,
    PositionStaffUpdate,
    PositionStaffDelete,
  };
}
