export const useApi = () => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    const apiUrl = process.server 
        ? config.apiInternal 
        : config.public.apiBase

    const fetch = async <T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> => {
        const url = `${apiUrl}${endpoint}`
        
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
            ...options.headers,
        }

        // Add auth token if available
        if (authStore.token) {
            headers['Authorization'] = `Bearer ${authStore.token}`
        }

        try {
            const response = await $fetch<T>(url, {
                ...options,
                headers,
            })
            return response
        } catch (error: any) {
            // Handle error response - $fetch throws error with data property
            console.error('API Error:', error)
            console.error('API Error URL:', url)
            console.error('API Error Options:', options)
            
            // Re-throw with proper format
            if (error.data) {
                throw {
                    data: error.data,
                    message: error.message || 'Request failed',
                    status: error.status || error.statusCode || 500
                }
            }
            
            // Network error or other issues
            const errorMessage = error.message || 'Network error or server unavailable'
            console.error('Network error:', errorMessage)
            throw {
                data: { error: errorMessage },
                message: errorMessage,
                status: 0
            }
        }
    }

    const post = async <T>(endpoint: string, data: any): Promise<T> => {
        return fetch<T>(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        })
    }

    const put = async <T>(endpoint: string, data: any): Promise<T> => {
        return fetch<T>(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data),
        })
    }

    const del = async <T>(endpoint: string): Promise<T> => {
        return fetch<T>(endpoint, {
            method: 'DELETE',
        })
    }

    return {
        apiUrl,
        fetch,
        post,
        put,
        delete: del,
        del,
    }
}

