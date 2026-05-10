<script lang="ts">
	import { goto } from '$app/navigation';
	import { post } from '$lib/api';
	import { setToken } from '$lib/auth';
	import type { LoginResponse } from '$lib/types';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		loading = true;
		try {
			const res = await post<LoginResponse>('/api/auth/login', { email, password });
			if (!res.success) {
				error = res.error ?? 'Login failed';
				return;
			}
			setToken(res.data.token, res.data.user_id);
			goto('/');
		} catch {
			error = 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="container">
	<h1>Sign in</h1>
	<form onsubmit={handleSubmit}>
		<label>
			Email
			<input type="email" bind:value={email} required />
		</label>
		<label>
			Password
			<input type="password" bind:value={password} required />
		</label>
		{#if error}<p class="error">{error}</p>{/if}
		<button type="submit" disabled={loading}>
			{loading ? 'Signing in...' : 'Sign in'}
		</button>
	</form>
	<p><a href="/register">Create account</a></p>
</div>

<style>
	.container {
		max-width: 360px;
		margin: 4rem auto;
	}
	h1 {
		margin-bottom: 1.5rem;
	}
	form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.875rem;
		font-weight: 500;
	}
	input {
		padding: 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		font-size: 1rem;
	}
	button {
		padding: 0.625rem;
		background: #111827;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		cursor: pointer;
	}
	button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
	.error {
		color: #dc2626;
		font-size: 0.875rem;
		margin: 0;
	}
	p {
		margin-top: 1rem;
		font-size: 0.875rem;
	}
</style>
