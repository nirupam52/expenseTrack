import { goto } from '$app/navigation';
import { clearToken, getToken } from './auth';
import type { ApiList, ApiResponse } from './types';

async function apiFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
	const token = getToken();
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(init.headers as Record<string, string> | undefined),
	};
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const res = await fetch(path, { ...init, headers });

	if (res.status === 401) {
		clearToken();
		goto('/login');
		throw new Error('unauthorized');
	}

	if (!res.ok) {
		throw new Error(`${res.status} ${res.statusText}`);
	}

	return res.json() as Promise<T>;
}

export function get<T>(path: string): Promise<ApiResponse<T>> {
	return apiFetch<ApiResponse<T>>(path);
}

export function getList<T>(path: string): Promise<ApiList<T>> {
	return apiFetch<ApiList<T>>(path);
}

export function post<T>(path: string, body: unknown): Promise<ApiResponse<T>> {
	return apiFetch<ApiResponse<T>>(path, {
		method: 'POST',
		body: JSON.stringify(body),
	});
}

