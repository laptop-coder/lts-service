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

// Responses
export interface LoginResponse {
  user: User;
}

