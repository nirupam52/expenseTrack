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
	p {
		margin-top: 1rem;
		font-size: 0.875rem;
	}
</style>
