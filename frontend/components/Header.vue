<template>
	<div id="Header">
		<div class="left">
			<router-link to="/" tag="img" src="https://cdn.tcole.me/themodestland.png" alt="Timothy Cole Logo"></router-link>
			<router-link to="/" tag="h4" >Modest Land Subscriber Perks</router-link>
		</div>
		<div class="right">
			<div class="login" @click="loginPopup" v-if="$parent.$parent.user === null">Login with Twitch</div>
			<div class="login" @click="logout" v-if="$parent.$parent.user !== null">Logout</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'Header',
	data () {
		return {
			loginPopup: () => {
				const endpoint = `https://api.twitch.tv/kraken/oauth2/authorize`;
				const scopes = [
					"user_subscriptions",
				].join(" ");

				const url = `${endpoint}?scope=${scopes}&client_id=${twitchClientID}&redirect_uri=${twitchClientRedirect}&response_type=code&force_verify=true`;
				const popup = window.open(url, "_blank", "toolbar=no,scrollbars=yes,resizable=yes,width=440,height=580");
				setInterval(() => { 
					if (!popup.closed) return;
					location.reload();
				}, 250);
			},
			logout: async () => {
				let vm = this.$parent.$parent;
				await fetch(`/logout`, {
					method: "DELETE",
					headers: {
						"Authorization": `Session ${vm.getCookie("modestguard")}==`
					}
				});
				location.reload();
			}
		}
	}
}
</script>

<style lang="scss">
	$white: #dfebf5;
	$whiteish: #aabccb;

	#Header {
		width: 100%;
		padding: 15px 0;
		display: flex;

		-webkit-touch-callout: none;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
		user-select: none;

		* { vertical-align: middle; }
	}

	.left {
		order: 1;
		margin-left: 15px;
		flex: 2 0px;

		h4 {
			margin: 0;
			display: inline;
			font-size: 1.65em;
			font-weight: 100;
			padding-left: 10px;

			@media (max-width: 700px) {
				& { display: none; }
			}
		}

		img {
			height: 40px;
			padding-top: 3px;
		}
	}

	.right {
		order: 2;
		color: $whiteish;
		font-weight: 300;
		margin-right: 15px;

		ul {
			list-style: none;
			display: inline-block;
			margin: 0;

			li {
				display: inline-block;
				margin-left: 40px;
				&:nth-child(1) { margin-left: 0; }

				a {
					color: $whiteish;
					text-decoration: none;
					font-weight: 300;
					padding: 0 .25px;
					line-height: 40px;

					&.router-link-exact-active, &:hover {
						color: $white;
						font-weight: 400;
						padding: 0px;
					}
				}
			}
		}
	}

	.login {
		display: inline-block;
		margin-left: 35px;
		border: 1px solid $white;
		color: $white;
		padding: 10px 20px;
		border-radius: 7px;

		&:hover {
			cursor: pointer;
			border-color: transparent;
			transition: background 0.5s ease;
			background: rgb(48,139,205);
			background: linear-gradient(133deg, rgba(48,139,205,1) 0%, rgba(81,181,156,1) 100%); 
			color: $white;
		}
	}
</style>
