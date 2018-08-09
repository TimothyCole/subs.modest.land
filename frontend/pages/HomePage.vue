<template>
	<div id="HomePage">
		<div class="gradient">
			<Header class="container"></Header>
			<div class="container billboard">
				<p>Login with your Twitch account on the top right to prove that you're a subscriber of one of the following channels.</p>
				<div v-for="channel of channels" v-bind:key="channel">
					<a :href="`https://twitch.tv/${channel.toLowerCase()}`" target="_blank">{{ channel }}</a>
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
			]
		}
	}
}
</script>

<style lang="scss" scoped>
	$white: #dfebf5;
	$whiteish: #aabccb;

	#HomePage {
		width: 100%;
		min-height: 100%;
		background-image: url('http://www.3sixtydisplays.co.uk/wp-content/uploads/2015/09/3d_low_poly_abstract-2560x1600.jpg');
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
