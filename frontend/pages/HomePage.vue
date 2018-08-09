<template>
	<div id="HomePage">
		<div class="gradient">
			<Header class="container"></Header>
			<div class="container billboard" v-if="$parent.user === null">
				<p>Login with your Twitch account on the top right to prove that you're a subscriber of one of the following channels.</p>
				<div v-for="channel of channels" v-bind:key="channel">
					<a :href="`https://twitch.tv/${channel.toLowerCase()}`" target="_blank">{{ channel }}</a>
				</div>
			</div>
			<div class="container billboard" v-if="$parent.user !== null">
				<h1>Hey {{ $parent.user.me.data[0].display_name }}!</h1>
				<div v-if="hasAccess && subbedTo && whitelist.length == 0">
					<div>
						<input @keyup.enter="whitelistAdd" v-model="username" type="text" id="mcName" placeholder="Enter your Minecraft Username to be added to the whitelist" />
						<button v-on:click="whitelistAdd">Add to Whitelist</button>
					</div>
					<p v-text="error"></p>
					<h3>THIS CAN NOT BE CHANGED!!! <small style="font-size: 0.6em; font-weight: 300; font-style: italic;">(im lazy)</small></h3>
				</div>
				<div v-if="hasAccess && subbedTo && whitelist.length != 0">
					<h2 class="name" v-text="whitelist"></h2>
					<h3>THIS CAN NOT BE CHANGED!!! <small style="font-size: 0.6em; font-weight: 300; font-style: italic;">(im lazy)</small></h3>
				</div>
				<div v-if="!hasAccess">
					<p>Thank you for your interest but you're not a subsciber of any of the following channels 😭</p>
					<div v-for="channel of channels" v-bind:key="channel">
						<a :href="`https://twitch.tv/${channel.toLowerCase()}`" target="_blank">{{ channel }}</a>
					</div>
				</div>
			</div>
			<Footer></Footer>
		</div>
	</div>
</template>

<script>
export default {
	name: 'HomePage',
	data: () => {
		return {
			channels: [
				"ModestTim",
				"Jamie254",
				"JamiePineLive",
				"Ashturbate"
			],
			whitelist: "",
			username: "",
			error: "",
		}
	},
	computed: {
		hasAccess: function () {
			var s2 = this.subbedTo()
			return s2[1] && s2[2]
		}
	},
	methods: {
		subbedTo: function () {
			var subbedTo = [];
			var timed = false;

			for (const i of this.$parent.user.checks) {
				if (i.type === "subbed" && typeof i.created_at != 'undefined') subbedTo.push(i)
				if (typeof i.created_at != 'undefined' && !timed) timed = this.weekAgo(i.created_at);
			}

			return [ subbedTo, subbedTo.length > 0, timed ]
		},
		weekAgo: (since) => {
			const sinceDate = new Date(since);
			const now = new Date();
			const days = Math.ceil(Math.abs(sinceDate.getTime() - now.getTime()) / (1000 * 3600 * 24)); 
			return days > 7;
		},
		whitelistAdd: async function () {
			if (this.username.length < 3) return this.error = `Minecraft usernames need to be at least 3 characters`;
			if (this.username.length > 16) return this.error = `Minecraft usernames can't be more than 16 characters`;
			if (!/^[a-zA-Z0-9_]*$/.test(this.username)) return this.error = `Minecraft usernames can only be alphanumeric with underscores`;
			this.error = ``;

			let body = new FormData();
			body.append("username", this.username);
			this.whitelist = await fetch(`/whitelist`, {
				method: "POST",
				body,
				headers: {
					"Authorization": `Session ${this.$parent.getCookie("modestguard")}==`
				}
			}).then(r => r.json()).catch(_ => null)

			if (this.whitelist == null) this.error = `There must of been a server error 🤔. Oops you can either try again or <a href="https://github.com/TimothyCole/subs.modest.land">fix my bad code</a>`;
			if (this.whitelist.minecraft) this.whitelist = this.whitelist.minecraft
		}
	},
	created: async function () {
		try {
			this.whitelist = (await fetch(`/whitelist`, {
				method: "GET",
				headers: {
					"Authorization": `Session ${this.$parent.getCookie("modestguard")}==`
				}
			}).then(r => r.json())).minecraft
		} catch(_) { this.whitelist = null }
	}
}
</script>

<style lang="scss" scoped>
	$white: #dfebf5;
	$whiteish: #aabccb;

	#HomePage {
		width: 100%;
		min-height: 100%;
		background-image: url('https://cdn.tcole.me/low-poly-bg-smp.jpg');
		background-position: center;
		background-repeat: no-repeat;
		background-size: cover;
	}

	.gradient {
		height: 100%;
		width: 100%;
		background: rgba(50,59,72, 0.9);
		background: linear-gradient(-55deg, rgba(46, 36, 48, 0.9) 0%, rgba(63, 49, 66, 0.9) 100%);
		display: flex;
		flex-direction: column;

		.billboard {
			padding-top: 80px;
			flex: 2 0px;
			text-align: center;

			p {
				word-break: break-all;
				text-align: center;
			}

			div {
				display: inline-block;
				margin: 0 auto;
				padding: 10px 5px;

				div {
					input {
						display: inline-block;
						outline: 0;
						background: #f2f2f23b;
						color: #ffffff;
						width: 500px;
						border: 0;
						margin: 0 0 15px;
						padding: 15px;
						box-sizing: border-box;
						font-size: 14px;
					}
					button {
						display: inline-block;
						text-transform: uppercase;
						outline: 0;
						background: #4CAF50;
						width: 175px;
						border: 0;
						padding: 15px;
						color: #FFFFFF;
						font-size: 14px;
						-webkit-transition: all 0.3 ease;
						transition: all 0.3 ease;
						cursor: pointer;
					}
				}

				h2.name {
					background: #f2f2f23b;
					color: #ffffff;
					width: 500px;
					border: 0;
					margin: 0 0 15px;
					padding: 15px;
					box-sizing: border-box;
					font-size: 17px;
					font-weight: 100;
				}

				a {
					background-color: #308ccd4f;
					text-decoration: none;
					border-radius: 4px;
					color: #eeeeee;
					padding: 7px 10px;
					font-weight: 400;
					margin: 0 5px;

					&:hover {
						background-color: #308ccd9c;
					}
				}
			}
		}
	}
</style>
