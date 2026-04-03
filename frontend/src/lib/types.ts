// Models
export interface User {
  id: string;
  createdAt: string;
  updatedAt: string;
  email: string;
  firstName: string;
  middleName?: string | null;
  lastName: string;
  hasAvatar: boolean;
  roles: Role[];
}

export interface InstitutionAdministrator {
  userId: string;
  createdAt: string;
  updatedAt: string;
  position: InstitutionAdministratorPosition;
}

export interface Staff {
  userId: string;
  createdAt: string;
  updatedAt: string;
  position: StaffPosition;
}

export interface Parent {
  userId: string;
  createdAt: string;
  updatedAt: string;
  students: Student[];
}

export interface Student {
  userId: string;
  createdAt: string;
  updatedAt: string;
  parents: Parent[];
  studentGroup: StudentGroup;
}

export interface Teacher {
  userId: string;
  createdAt: string;
  updatedAt: string;
  classroom?: Room;
  subjects: Subject[];
  studentGroups: StudentGroup[];
}

export interface Role {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
  permissions: Permission[];
}

export interface Permission {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
}

export interface Post {
  id: string;
  createdAt: string;
  updatedAt: string;
  name: string;
  description?: string;
  verified: boolean;
  thingReturnedToOwner: boolean;
  hasPhoto: boolean;
  author: User;
}

export interface Room {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
  teacherId: string;
}

export interface Subject {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
}

export interface StudentGroup {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
  groupAdvisorId: string;
}

export interface StaffPosition {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
}

export interface InstitutionAdministratorPosition {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
}

// Responses
export interface LoginResponse {
  user: User;
}
