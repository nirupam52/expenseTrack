<script lang="ts">
	import { goto } from '$app/navigation';
	import { post } from '$lib/api';
	import { getUserId } from '$lib/auth';
	import type { Expense } from '$lib/types';

	let description = $state('');
	let amount = $state('');
	let date = $state(new Date().toISOString().split('T')[0]);
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		loading = true;
		const userId = getUserId();
		if (!userId) {
			error = 'Not authenticated';
			loading = false;
			return;
		}
		try {
			const res = await post<Expense>('/api/expenses', {
				paid_by: userId,
				description,
				amount: parseFloat(amount),
				date,
			});
			if (!res.success) {
				error = res.error ?? 'Failed to create expense';
				return;
			}
			goto('/');
		} catch {
			error = 'Failed to create expense';
		} finally {
			loading = false;
		}
	}
</script>

<div class="header">
	<a href="/" class="back">← Back</a>
	<h1>New expense</h1>
</div>

<form onsubmit={handleSubmit}>
	<label>
		Description
		<input type="text" bind:value={description} required />
	</label>
	<label>
		Amount
		<input type="number" bind:value={amount} min="0.01" step="0.01" required />
	</label>
	<label>
		Date
		<input type="date" bind:value={date} required />
	</label>
	{#if error}<p class="error">{error}</p>{/if}
	<button type="submit" disabled={loading}>
		{loading ? 'Saving...' : 'Add expense'}
	</button>
</form>

<style>
	.header {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1.5rem;
	}
	.back {
		color: #6b7280;
		text-decoration: none;
		font-size: 0.875rem;
	}
	h1 {
		margin: 0;
	}
	form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		max-width: 400px;
	}
</style>
