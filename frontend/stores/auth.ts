import { defineStore } from 'pinia'

export interface User {
    id: string
    email: string
    full_name: string
    phone?: string
    tenant_id?: string
    tenant_name?: string
    role_id?: string
    role_name?: string
    unit_id?: string
    permissions?: string[]
    role?: string
}

export const useAuthStore = defineStore('auth', () => {
    // Use secure only if HTTPS, not just in production
    const isSecure = typeof window !== 'undefined' && window.location.protocol === 'https:'
    
    const token = useCookie<string | null>('token', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: isSecure
    })
    
    const user = useCookie<User | null>('user', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: isSecure
    })

    const isAuthenticated = computed(() => !!token.value)
    const hasPermission = (permission: string) => {
        return user.value?.permissions?.includes(permission) || false
    }

    function setToken(newToken: string) {
        if (!newToken) {
            console.warn('Attempting to set empty token')
            return
        }
        token.value = newToken
        console.log('Token set successfully, length:', newToken.length)
    }

    function setUser(newUser: User) {
        user.value = newUser
    }

    async function fetchCurrentUser() {
        const { fetch } = useApi()
        try {
            if (!token.value) {
                console.error('No token available when fetching current user')
                return null
            }
            console.log('Fetching current user with token, length:', token.value.length)
            const currentUser = await fetch<User>('/api/me')
            // Merge dengan data user yang sudah ada untuk preserve role_name jika ada
            if (user.value && currentUser) {
                setUser({
                    ...user.value,
                    ...currentUser,
                    // Preserve role_name if it exists in current user, otherwise keep existing
                    role_name: currentUser.role_name || user.value.role_name,
                    permissions: currentUser.permissions || user.value.permissions,
                })
            } else {
                setUser(currentUser)
            }
            return currentUser
        } catch (error) {
            console.error('Failed to fetch current user:', error)
            // Don't logout on error, just log it
            // logout()
        }
    }

    function logout() {
        token.value = null
        user.value = null
        navigateTo('/login')
    }

    return {
        token,
        user,
        isAuthenticated,
        hasPermission,
        setToken,
        setUser,
        fetchCurrentUser,
        logout
    }
})
