<script lang="ts">
	import { onMount } from 'svelte';
	import { getList } from '$lib/api';
	import { getUserId } from '$lib/auth';
	import type { Expense } from '$lib/types';

	let expenses = $state<Expense[]>([]);
	let error = $state('');
	let loading = $state(true);

	onMount(async () => {
		const userId = getUserId();
		if (!userId) return;
		try {
			const res = await getList<Expense>(`/api/expenses?user_id=${userId}`);
			expenses = res.data ?? [];
		} catch {
			error = 'Failed to load expenses';
		} finally {
			loading = false;
		}
	});

	function formatAmount(amount: number): string {
		return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
	}
</script>

<div class="header">
	<h1>Expenses</h1>
	<a href="/expenses/new" class="btn">+ Add</a>
</div>

{#if loading}
	<p class="muted">Loading...</p>
{:else if error}
	<p class="error">{error}</p>
{:else if expenses.length === 0}
	<p class="muted">No expenses yet. <a href="/expenses/new">Add one.</a></p>
{:else}
	<table>
		<thead>
			<tr>
				<th>Description</th>
				<th>Date</th>
				<th class="right">Amount</th>
			</tr>
		</thead>
		<tbody>
			{#each expenses as expense}
				<tr>
					<td>{expense.description}</td>
					<td>{expense.date}</td>
					<td class="right">{formatAmount(expense.amount)}</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/if}

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}
	.btn {
		padding: 0.5rem 1rem;
		background: #111827;
		color: white;
		text-decoration: none;
		border-radius: 4px;
		font-size: 0.875rem;
	}
	table {
		width: 100%;
		border-collapse: collapse;
	}
	th,
	td {
		text-align: left;
		padding: 0.75rem;
		border-bottom: 1px solid #e5e7eb;
	}
	th {
		font-weight: 600;
		font-size: 0.875rem;
		color: #6b7280;
	}
	.right {
		text-align: right;
	}
	.error {
		color: #dc2626;
	}
	.muted {
		color: #6b7280;
	}
</style>
