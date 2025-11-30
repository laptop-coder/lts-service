// Routes consist of three parts (all words separated by "_" and all parts
// separated by "__"):
// 1. available to whom (empty if for everyone), for example, "MODERATOR";
// 2. what (from general to particular), for example, "THING_ADD";
// 3. word "ROUTE".

// For moderator
export const MODERATOR__HOME__ROUTE =
  '/' + '1dae325e-0557-4977-97a9-9a45e7e6efe3' + '/';
export const MODERATOR__PROFILE__ROUTE = MODERATOR__HOME__ROUTE + 'profile';

// For all users
export const HOME__ROUTE = '/';
export const LOGIN_USER__ROUTE = HOME__ROUTE + 'login';
export const REGISTER_USER__ROUTE = HOME__ROUTE + 'register';
export const LOGIN_MODERATOR__ROUTE = MODERATOR__HOME__ROUTE + 'login';
export const REGISTER_MODERATOR__ROUTE = MODERATOR__HOME__ROUTE + 'register';

// For registered users
export const USER__PROFILE__ROUTE = HOME__ROUTE + 'profile';
export const USER__THING_ADD__ROUTE = HOME__ROUTE + 'thing/add';
export const USER__THING_EDIT__ROUTE = HOME__ROUTE + 'thing/edit';
export const USER__THING_STATUS__ROUTE = HOME__ROUTE + 'thing/status';

// ----------------------------------------------------------------------------

export const ASSETS_ROUTE = '/storage/assets';
export const STORAGE_ROUTE = '/storage/storage';

export const BACKEND_URL = 'http://localhost:37190';
export const SCHOOL_URL = 'https://лицей369.рф';
export const TECH_SUPPORT_URL = 'https://help.licey369.ru';

export const PASSWORD_MIN_LEN = '8';
export const PASSWORD_MAX_LEN = '32';
export const USERNAME_MIN_LEN = '1';
export const USERNAME_MAX_LEN = '16';

export enum Role {
  'user' = 'user',
  'moderator' = 'moderator',
  'none' = 'none',
}

export enum ThingType {
  'lost' = 'lost',
  'found' = 'found',
}

export enum UserProfileSection {
  'advertisements' = 'advertisements',
  'settings' = 'settings',
}
