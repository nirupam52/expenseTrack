<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getToken, clearToken } from '$lib/auth';
	import { onMount } from 'svelte';

	const publicRoutes = ['/login', '/register'];

	let { children } = $props();

	onMount(() => {
		if (!publicRoutes.includes(page.url.pathname) && !getToken()) {
			goto('/login');
		}
	});

	function logout() {
		clearToken();
		goto('/login');
	}
</script>

{#if !publicRoutes.includes(page.url.pathname)}
	<nav>
		<a href="/">Expenses</a>
		<button onclick={logout}>Logout</button>
	</nav>
{/if}

<main>
	{@render children()}
</main>

<style>
	nav {
		display: flex;
		gap: 1rem;
		padding: 0.875rem 1.5rem;
		border-bottom: 1px solid #e5e7eb;
		align-items: center;
	}
	nav a {
		text-decoration: none;
		color: inherit;
		font-weight: 500;
	}
	nav button {
		margin-left: auto;
		background: none;
		border: 1px solid #d1d5db;
		padding: 0.25rem 0.75rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}
	main {
		max-width: 800px;
		margin: 2rem auto;
		padding: 0 1rem;
	}
</style>
