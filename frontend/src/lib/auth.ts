const TOKEN_KEY = 'token';
const USER_ID_KEY = 'user_id';

export function getToken(): string | null {
	return localStorage.getItem(TOKEN_KEY);
}

export function getUserId(): number | null {
	const id = localStorage.getItem(USER_ID_KEY);
	return id ? parseInt(id, 10) : null;
}

export function setToken(token: string, userId: number): void {
	localStorage.setItem(TOKEN_KEY, token);
	localStorage.setItem(USER_ID_KEY, String(userId));
}

export function clearToken(): void {
	localStorage.removeItem(TOKEN_KEY);
	localStorage.removeItem(USER_ID_KEY);
}
