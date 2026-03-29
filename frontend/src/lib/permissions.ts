import {useAuth} from './auth'

export function usePermissions () {
  const auth = useAuth()

  const hasPermission = (permission: string): boolean => {
    const user = auth.user()
    if (!user) return false
    return user.roles.some(role => role.permissions.some(p => p.name === permission))
  }

  const hasAnyPermission = (...perms: string[]): boolean => {
    return perms.some(p => hasPermission(p))
  }

  const hasAllPermissions = (...perms: string[]): boolean => {
    return perms.every(p => hasPermission(p))
  }

  const hasRole = (roleName: string): boolean => {
    const user = auth.user()
    if (!user) return false
    return user.roles.some(r => r.name === roleName)
  }
  
  return {hasPermission, hasAnyPermission, hasAllPermissions, hasRole}
}



