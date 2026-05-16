export interface ApiResponse<T> {
	success: boolean;
	data: T;
	error?: string;
}

export interface ApiList<T> {
	success: boolean;
	data: T[];
	meta: { count: number };
}

export interface User {
	id: number;
	name: string;
	email: string;
	created_at: string;
}

export interface Expense {
	id: number;
	group_id?: number;
	paid_by: number;
	description: string;
	amount: number;
	date: string;
	created_at: string;
}

export interface LoginResponse {
	token: string;
	user_id: number;
	expires_at: string;
}
