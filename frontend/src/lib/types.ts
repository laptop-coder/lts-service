// Response DTOs from the backend
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
  position: InstitutionAdministratorPosition;
}

export interface Staff {
  userId: string;
  position: StaffPosition;
}

export interface Parent {
  userId: string;
  students: Student[];
}

export interface Student {
  userId: string;
  parents: Parent[];
  studentGroup: StudentGroup;
}

export interface Teacher {
  userId: string;
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
  teacherId: string | null;
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
  groupAdvisorId: string | null;
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

export interface Conversation {
  id: string;
  createdAt: string;
  post: Post;
  messages: Message[];
  otherUser: User;
}

export interface ConversationListItem {
  id: string;
  updatedAt: string;
  postID: string;
  postName: string;
  unreadCount: number;
  lastMessage?: string;
  otherUser: User;
}

export interface Message {
  id: string;
  createdAt: string;
  updatedAt: string;
  senderID: string;
  content: string;
  isRead: boolean;
}

// Responses
export interface LoginResponse {
  user: User;
}
