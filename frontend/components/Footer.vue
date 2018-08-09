<template>
	<div class="footer">
		<p><i><b>Disclaimer</b> You must be either a follower or a subscriber for more than a week to gain access to perks to pervent trolls from making new accounts to grief.</i></p>
		<p>Copyright &copy; {{ start }}-{{ end }} - Timothy Cole - All Rights Reserved.</p>
	</div>
</template>

<script>
export default {
	name: "Footer",
	data () {
		return {
			start: "2017",
			end: (new Date()).getFullYear()
		}
	}
}
</script>

<style lang="scss" scoped>
	p {
		color: #aabccb;
		font-size: 0.8em;
		display: inline-block;
		width: 100%;
		text-align: center;
	}
</style>
