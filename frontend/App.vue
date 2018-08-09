<template>
	<div id="app">
		<router-view class="page"></router-view>
	</div>
</template>

<script>
export default {
	name: 'Main',
	data: () => {
		return {
			user: null
		}
	},
	created: async function () {
		try {
			this.user = (await fetch(`/me`, {
				method: "GET",
				headers: {
					"Authorization": `Session ${this.getCookie("modestguard")}==`
				}
			}).then(r => r.json()))
		} catch(_) { this.user = null }
	},
	methods: {
		getCookie: name => {
			return document.cookie.split('; ').reduce((r, v) => {
				const parts = v.split('=')
				return parts[0] === name ? decodeURIComponent(parts[1]) : r
			}, '')
		}
	}
}
</script>

<style lang="scss">
	@import url('https://fonts.googleapis.com/css?family=Material+Icons|Roboto:100,300,400,500');
	html, body {
		font-family: 'Roboto', sans-serif;
		font-weight: 100;
		color: #dfebf5;
		height: 100%;
		width: 100%;
		padding: 0;
		margin: 0;
	}

	div#app {
		height: 100%;
		width: 100%;
		display: flex;
		overflow: hidden;

		div.page {
			display: block;
			flex-direction: column;
			width: 100%;
			height: 100%;
			overflow: auto;
			z-index: 1;
		}

		div.container {
			max-width: 1275px;
			margin: 0 auto;
		}
	}
</style>