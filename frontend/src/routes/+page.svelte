<script lang="ts">
	import { onMount } from 'svelte';

	let servers: any = { servers: [0, 1] };
	let docker: any;
	let dockerData: any = { containers: [] };
	let runningDocker: number = 0;
	let dockerLogs: any = { logs: '' };
	let currentServer: string | null = null;
	let currentDockerId: string | null = null;

	onMount(async () => {
		// cleanup the subscription when the component is destroyed
		const rs = await fetch('http://localhost:8080/api/servers');
		servers = await rs.json();
		const rd = await fetch('http://localhost:8080/api/dockers');
		docker = await rd.json();
	});

	function formatLogs(logs: string) {
		const logLines = logs.split('\n');
		let result = '';

		logLines.forEach((line) => {
			const logMatch = line.match(/\[(.*?)\]/);
			const logLevel = logMatch ? logMatch[1] : '';
			const logMessage = line.replace(/\[(.*?)\]/, '').trim();

			if (logMessage.includes('Error')) {
				result += `<p class="text-red-600">${logMessage}</p>`;
			}

			if (logLevel) {
				// result += `<p class="text-${getLogLevelColor(logLevel)}-600">${logLevel}</p>`;
				result += `<p class="text-green-600">${logLevel}</p>`;
			}

			result += `<p class="text-sm">${logMessage}</p><br>`;
		});

		return result;
	}

	type Request = {
		url: string;
		body?: any;
	};

	const handlePost = async (request: Request) => {
		runningDocker = 0;
		dockerData = { containers: [] };
		dockerLogs = { logs: '' };
		document.getElementById('logs')!.innerHTML = '';
		const payload = { server: request.url, command: request.body };
		const response = await fetch('http://localhost:8080/api/send', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(payload)
		});
		dockerData = await response.json();
		if (dockerData.containers == null) dockerData.containers = [];
		runningDocker = dockerData.containers.reduce((acc: number, curr: any) => {
			if (curr.status.includes('Up')) {
				acc++;
			}
			return acc;
		}, 0);
		currentServer = request.url;
		console.log(runningDocker);
	};

	const handlePostLogs = async (request: Request) => {
		dockerLogs = null;
		const payload = { server: request.url, command: request.body };
		const response = await fetch('http://localhost:8080/api/send', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(payload)
		});
		let logs = await response.json();
		dockerLogs = formatLogs(decodeURI(logs.logs));
		document.getElementById('logs')!.innerHTML = dockerLogs;
	};

	let hour = new Date().getHours();
	let minute = new Date().getMinutes();
	setInterval(() => {
		hour = new Date().getHours();
		minute = new Date().getMinutes();
	}, 1000);

	const openMenuModal = () => {
		const modal = document.getElementById('menu-modal');

		modal!.classList.contains('hidden')
			? modal!.classList.remove('hidden')
			: modal!.classList.add('hidden');
	};
</script>

<div class="flex flex-col w-full h-full p-4">
	<nav class="flex flex-row w-full justify-between pr-4 pl-4 items-center">
		<h2 on:click={openMenuModal} class="cursor-pointer text-white">Menu</h2>
		<div class="flex flex-row items-end w-40 justify-end">
			<h2 class="text-white">{hour}:{minute > 10 ? minute : '0' + minute}</h2>
			<!-- <h3>Settings</h3> -->
		</div>
	</nav>
	<div
		id="menu-modal"
		class="hidden absolute top-11 left-7  h-40 w-28 bg-zinc-700 p-2 rounded-lg text-left text-black"
	>
		<div class="flex w-full hover:bg-zinc-600 cursor-pointer p-2 transition rounded-lg">
			<p>Exit</p>
		</div>
	</div>
	<div class="flex flex-row h-full w-full p-4">
		<div class="flex flex-col h-full w-1/5">
			<div class="flex flex-col h-5/6 w-full p-2">
				{#each servers.servers as server}
					<div class="flex flex-row justify-between">
						<div class="flex flex-row border border-zinc-600 rounded-lg w-full mb-2 h-16">
							<button
								on:click={() => handlePost({ url: server, body: 'docker ps -a' })}
								class="flex w-full h-full items-center p-4 text-white justify-start rounded-lg transition ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 hover:bg-zinc-400"
							>
								<p class="mr-2">Server:</p>
								{server}
							</button>
						</div>
					</div>
				{/each}
			</div>
			<div class="flex h-12 text-white p-4   border-zinc-700 border-t">
				<h1>Resume for {currentServer == null ? servers[0] : currentServer}</h1>
			</div>
			<div class="flex flex-col h-1/5 p-4 justify-between text-white">
				<h1>
					Total Dockers:
					{dockerData.containers.length}
				</h1>
				<h1>
					Docker Runnings:
					{dockerData.containers.length - runningDocker}
				</h1>
				<h1>
					Docker Stopped:
					{dockerData.containers.length - runningDocker}
				</h1>
			</div>
		</div>
		<div class="flex w-4 bg-zinc-900 border-l border-zinc-700 rounded-l" />
		<div class="flex flex-col w-5/6 bg-zinc-900 h-full rounded-lg p-2">
			<div class="flex flex-row w-full h-3/5">
				<div class="flex flex-col w-5/6 border border-zinc-700 rounded-lg bg-zinc-900 text-white">
					<p class="flex flex-col w-full h-full overflow-y-scroll p-2" id="logs">
						{#if dockerLogs == null}
							<p class="flex flex-col w-full h-full items-center justify-center">
								Please select a docker in order to see the logs
							</p>
						{/if}
					</p>
				</div>
				<div class="flex w-4" />
				<div class="flex flex-col w-1/5 justify-between">
					<button
						class="flex text-white border border-zinc-700 rounded-lg h-16 items-center justify-center transition ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 hover:bg-zinc-400"
						on:click={() =>
							handlePostLogs({
								url: currentServer !== null ? currentServer : '',
								body: `docker start ${currentDockerId}`
							})}>Start</button
					>
					<button
						class="flex text-white border rounded-lg border-zinc-700 h-16 items-center justify-center transition ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 hover:bg-zinc-400"
						on:click={() =>
							handlePostLogs({
								url: currentServer !== null ? currentServer : '',
								body: `docker stop ${currentDockerId}`
							})}>Stop</button
					>
					<button
						class="flex text-white border rounded-lg border-zinc-700 h-16 items-center justify-center transition ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 hover:bg-zinc-400"
						on:click={() =>
							handlePostLogs({
								url: currentServer !== null ? currentServer : '',
								body: `docker restart ${currentDockerId}`
							})}>Restart</button
					>
					<button
						class="flex text-white border rounded-lg h-16 items-center border-zinc-700 justify-center transition ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 disabled:click-events-none cursor-not-allowed"
						>Rebuild</button
					>
					<button
						class="flex text-white border rounded-lg h-16 items-center justify-center border-zinc-700 transition ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 disabled:click-events-none cursor-not-allowed"
						>Access</button
					>
					<button
						class="flex text-white border rounded-lg h-16 items-center justify-center transition border-zinc-700 ease-in-out delay-10 hover:-translate-y-1 hover:scale-105 duration-300 disabled:click-events-none cursor-not-allowed"
						>Edit</button
					>
					<button
						class="flex text-white border rounded-lg h-16 items-center justify-center transition ease-in-out border-zinc-700 delay-10 hover:-translate-y-1 hover:scale-105 duration-300 hover:bg-zinc-400"
						on:click={() =>
							handlePostLogs({
								url: currentServer !== null ? currentServer : '',
								body: `docker rm ${currentDockerId}`
							})}>Remove</button
					>
				</div>
			</div>
			<div class="flex h-4" />
			<div class="flex flex-row w-full h-2/5 border-t border-zinc-700 p-4 overflow-auto ">
				<table class="table-fixed w-full text-left text-white">
					<thead>
						<tr>
							<th>ID</th>
							<th>Name</th>
							<th>Status</th>
							<th>Port</th>
							<th>Created</th>
						</tr>
					</thead>
					<tbody class="table-fixed cursor-pointer">
						{#if dockerData == null}
							{#if docker.containers == null}
								<div class="flex w-full items-center h-full">No containers found</div>
							{/if}
							{#if docker.containers != null && docker.containers == 0}
								{#each docker.containers as d}
									<tr>
										<td>{d.id}</td>
										<td>{d.names}</td>
										<td>{d.status}</td>
										<td>{d.ports}</td>
										<td>{d.created}</td>
									</tr>
								{/each}
							{/if}
						{/if}
						{#if dockerData !== null}
							{#if dockerData.containers == null}
								<tr>
									<td class="flex w-full items-center h-full">No containers found</td>
								</tr>
							{:else}
								{#each dockerData.containers as docker}
									<tr
										class="hover:bg-zinc-700 transition"
										on:click={() => {
											handlePostLogs({
												url: currentServer !== null ? currentServer : '',
												body: `docker logs -n=20 ${docker.id}`
											});
											currentDockerId = docker.id;
										}}
									>
										<td>{docker.id}</td>
										<td>{docker.names} </td>
										<td>{docker.status}</td>
										<td>{docker.ports}</td>
										<td>{docker.created}</td>
									</tr>
								{/each}
							{/if}
						{/if}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
