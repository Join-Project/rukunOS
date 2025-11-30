export const useAuth = () => {
  const authStore = useAuthStore()
  const user = computed(() => authStore.user)

  // Check user roles based on role name (primary) and permissions (fallback)
  // Priority: Admin > Bendahara > Sekretariat > Satpam > Warga

  const isAdmin = computed(() => {
    if (!user.value) return false
    // Check role name first (most reliable)
    if (user.value.role_name && user.value.role_name.toLowerCase() === 'admin') {
      return true
    }
    // Fallback: Check admin permissions
    return authStore.hasPermission('users.manage') ||
      authStore.hasPermission('roles.manage') ||
      authStore.hasPermission('tenant.settings')
  })

  const isBendahara = computed(() => {
    if (!user.value) return false
    // Check role name first (most reliable)
    if (user.value.role_name && user.value.role_name.toLowerCase() === 'bendahara') {
      return true
    }
    // Fallback: Check billing create/manage permissions (not just view)
    if (isAdmin.value) return false
    return authStore.hasPermission('billing.create') ||
      authStore.hasPermission('billing.payment') ||
      authStore.hasPermission('billing.update')
  })

  const isSekretariat = computed(() => {
    if (!user.value) return false
    // Check role name first (most reliable)
    if (user.value.role_name && user.value.role_name.toLowerCase() === 'sekretariat') {
      return true
    }
    // Fallback: Check communication create permissions (not just view)
    if (isAdmin.value || isBendahara.value) return false
    return authStore.hasPermission('communication.announcement.create') ||
      authStore.hasPermission('communication.announcement.update')
  })

  const isSecurity = computed(() => {
    if (!user.value) return false
    // Check role name first (most reliable)
    if (user.value.role_name && user.value.role_name.toLowerCase() === 'satpam') {
      return true
    }
    // Fallback: Check security permissions
    if (isAdmin.value || isBendahara.value || isSekretariat.value) return false
    return authStore.hasPermission('security.visitor.create') ||
      authStore.hasPermission('security.alert.respond')
  })

  const isResident = computed(() => {
    if (!user.value) return false
    // Warga is default - only basic permissions (billing.view, announcement.view, user.view)
    // Not admin, not bendahara, not sekretariat, not security
    return !isAdmin.value && !isBendahara.value && !isSekretariat.value && !isSecurity.value
  })

  // Alias for backward compatibility
  const isFinance = computed(() => isBendahara.value)

  const logout = () => {
    authStore.logout()
  }

  return {
    user,
    isAuthenticated: authStore.isAuthenticated,
    isAdmin,
    isResident,
    isFinance,
    isBendahara,
    isSecurity,
    isSecretariat: isSekretariat,
    logout,
    hasPermission: (permission: string) => {
      if (isAdmin.value) return true
      return authStore.hasPermission(permission)
    }
  }
}
